package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"harama/internal/config"
	"harama/internal/domain"
	"harama/internal/repository/postgres"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/uptrace/bun"
)

const (
	BaseURL = "http://127.0.0.1:8081/api/v1"
)

type Tenant struct {
	bun.BaseModel `bun:"table:tenants"`
	ID            uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	Name          string    `bun:"name"`
}

func main() {
	filePath := flag.String("file", "", "Path to the student's answer image (Required)")
	questionText := flag.String("question", "Calculate the sum of 5 and 7.", "The question text")
	rubricText := flag.String("rubric", "Correct answer 12", "The criteria for full marks")
	flag.Parse()

	if *filePath == "" {
		log.Fatal("Please provide an image file using -file flag")
	}

	// Load Config
	cfg := config.Load()

	// 0. Connect to DB
	log.Println("0. Connecting to Database...")
	db, err := postgres.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	tenantID, err := getOrCreateTenant(db)
	if err != nil {
		log.Fatalf("Failed to get tenant: %v", err)
	}
	log.Printf("   Using Tenant ID: %s\n", tenantID)

	// 1. Upload to MinIO
	log.Println("1. Uploading Answer Sheet to MinIO...")
	objectName, err := uploadToMinio(cfg, *filePath)
	if err != nil {
		log.Fatalf("Failed to upload: %v", err)
	}
	log.Printf("   File uploaded: %s\n", objectName)

	// 2. Create Exam
	log.Println("2. Creating Exam (Question Paper)...")
	examID, err := createExam(tenantID)
	if err != nil {
		log.Fatalf("Failed to create exam: %v", err)
	}

	// 3. Add Question (and Answer Key/Rubric)
	log.Printf("3. Adding Question: %s", *questionText)
	questionID, err := addQuestion(tenantID, examID, *questionText, *rubricText)
	if err != nil {
		log.Fatalf("Failed to add question: %v", err)
	}

	// 4. Create Submission (Upload Answer Sheet)
	log.Println("4. Submitting Answer Sheet...")
	subID, err := createSubmission(tenantID, examID, objectName)
	if err != nil {
		log.Fatalf("Failed to create submission: %v", err)
	}

	// 5. Wait for OCR (Backend now auto-triggers grading after OCR)
	log.Println("5. Waiting for AI to read handwriting (OCR)...")
	
	// Poll DB for OCR results and then link it to the question
	ocrText := waitForOCRAndLink(db, subID, questionID)
	log.Printf("   OCR Phase Result: %q", ocrText)

	// 6. Poll for Results (Triggering is now automatic)
	log.Println("6. Polling for Final Grade...")
	for i := 0; i < 100; i++ {
		grades, err := getGrades(tenantID, subID)
		if err == nil && len(grades) > 0 {
			log.Println("\n==========================================")
			log.Println("   📝 GRADING RESULTS")
			log.Println("==========================================")
			log.Printf("   Student Answer (OCR): %s", ocrText)
			log.Printf("   Score: %.2f / %d", grades[0].FinalScore, grades[0].MaxScore)
			log.Printf("   AI Reasoning: %s", grades[0].Reasoning)
			log.Println("==========================================")
			return
		}
		
		// Check status to see if it failed
		if i % 10 == 0 {
			var checkSub domain.Submission
			if err := db.NewSelect().Model(&checkSub).Where("id = ?", subID).Scan(context.Background()); err == nil {
				log.Printf("   Grading... Current Status: %s", checkSub.ProcessingStatus)
				if checkSub.ProcessingStatus == "failed" {
					log.Printf("   ❌ Submission marked as FAILED in DB.")
					return
				}
			}
		}
		
		time.Sleep(3 * time.Second)
		fmt.Print(".")
	}
	
	log.Println("\n   Timeout waiting for grades.")
}

func waitForOCRAndLink(db *bun.DB, subID, questionID string) string {
	ctx := context.Background()
	var sub domain.Submission
	
	// Poll until processing_status is 'failed' (OCR done but maybe empty) or we see OCR results
	// Increased timeout to 120 seconds (60 iterations * 2s)
	for i := 0; i < 60; i++ {
		time.Sleep(2 * time.Second)
		err := db.NewSelect().Model(&sub).Where("id = ?", subID).Scan(ctx)
		if err != nil {
			continue
		}
		
		if i % 5 == 0 {
			log.Printf("   Waiting... Current Status: %s", sub.ProcessingStatus)
		}

		// Fix #2 & #4: Stop polling if OCR is done, timed out, or grading has already started
		if sub.ProcessingStatus == domain.StatusOCRDone || 
		   sub.ProcessingStatus == domain.StatusOCRTimeout || 
		   sub.ProcessingStatus == domain.StatusGrading || 
		   sub.ProcessingStatus == domain.StatusCompleted {
			
			log.Printf("   OCR Phase Finished with status: %s", sub.ProcessingStatus)
			
			// If we already have OCR results, return them
			if len(sub.OCRResults) > 0 {
				extractedText := sub.OCRResults[0].RawText
				if extractedText != "" {
					return extractedText
				}
			}
			
			if sub.ProcessingStatus == domain.StatusOCRTimeout {
				return "OCR Timeout Fallback"
			}
			
			break
		}
		
		// Note: Our backend marks OCR as "failed" if text is empty/bad, but "completed" if it works.
		if len(sub.OCRResults) > 0 {
			extractedText := sub.OCRResults[0].RawText
			if extractedText != "" {
				// Success! Link this text to the question (Simulate Segmentation)
				answer := domain.AnswerSegment{
					ID:           uuid.New(),
					SubmissionID: sub.ID,
					QuestionID:   uuid.MustParse(questionID),
					Text:         extractedText,
				}
				sub.Answers = []domain.AnswerSegment{answer}
				sub.ProcessingStatus = domain.StatusPending // Reset status so Grading Worker picks it up
				
				_, _ = db.NewUpdate().Model(&sub).Column("answers", "processing_status").Where("id = ?", subID).Exec(ctx)
				return extractedText
			}
		}
		
		// If status is failed, we force a dummy text to prove grading works
		if sub.ProcessingStatus == "failed" {
			log.Println("   (OCR Failed to read text. Using fallback text for grading test)")
			fallbackText := "This is a fallback answer because OCR failed on the image."
			
			answer := domain.AnswerSegment{
				ID:           uuid.New(),
				SubmissionID: sub.ID,
				QuestionID:   uuid.MustParse(questionID),
				Text:         fallbackText,
			}
			sub.Answers = []domain.AnswerSegment{answer}
			sub.ProcessingStatus = domain.StatusPending
			_, _ = db.NewUpdate().Model(&sub).Column("answers", "processing_status").Where("id = ?", subID).Exec(ctx)
			return fallbackText
		}
	}
	return "Timeout waiting for OCR"
}

// ... (Helper functions below remain mostly same, just simplified addQuestion) ...

func getOrCreateTenant(db *bun.DB) (string, error) {
	ctx := context.Background()
	var tenant Tenant
	err := db.NewSelect().Model(&tenant).Limit(1).Scan(ctx)
	if err != nil {
		tenant = Tenant{Name: "Test Tenant"}
		_, err := db.NewInsert().Model(&tenant).Returning("id").Exec(ctx)
		if err != nil { return "", err }
	}
	return tenant.ID.String(), nil
}

func uploadToMinio(cfg *config.Config, filePath string) (string, error) {
	ctx := context.Background()
	minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil { return "", err }

	exists, err := minioClient.BucketExists(ctx, cfg.MinioBucket)
	if err == nil && !exists {
		minioClient.MakeBucket(ctx, cfg.MinioBucket, minio.MakeBucketOptions{})
	}

	file, err := os.Open(filePath)
	if err != nil { return "", err }
	defer file.Close()
	stat, _ := file.Stat()

	objectName := fmt.Sprintf("submissions/%d_%s", time.Now().Unix(), filepath.Base(filePath))
	// Use "png" to avoid double-prefix issue we fixed
	_, err = minioClient.PutObject(ctx, cfg.MinioBucket, objectName, file, stat.Size(), minio.PutObjectOptions{ContentType: "png"})
	return objectName, err
}

func createExam(tenantID string) (string, error) {
	exam := map[string]interface{}{ "title": "Manual Test Exam", "subject": "General"}
	resp, err := postRequest(BaseURL+"/exams", tenantID, exam)
	if err != nil { return "", err }
	return resp["id"].(string), nil
}

func addQuestion(tenantID, examID, text, rubricDesc string) (string, error) {
	// 1. Create Question
	q := map[string]interface{}{
		"exam_id": examID, "question_text": text, "points": 10, "answer_type": "short_answer",
	}
	resp, err := postRequest(fmt.Sprintf("%s/exams/%s/questions", BaseURL, examID), tenantID, q)
	if err != nil { return "", err }
	questionID := resp["id"].(string)

	// 2. Set Rubric
	rubric := map[string]interface{}{
		"question_id": questionID,
		"full_credit_criteria": []map[string]interface{}{
			{"description": rubricDesc, "points": 10},
		},
		"partial_credit_rules": []interface{}{}, "common_mistakes": []interface{}{},
	}
	jsonData, _ := json.Marshal(rubric)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/questions/%s/rubric", BaseURL, questionID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", tenantID)
	client := &http.Client{}
	rResp, err := client.Do(req)
	if err != nil { return "", err }
	defer rResp.Body.Close()
	if rResp.StatusCode >= 400 { return "", fmt.Errorf("rubric error") }
	return questionID, nil
}

func createSubmission(tenantID, examID, objectName string) (string, error) {
	// Note: We do NOT send 'answers' here. We let OCR find them.
	sub := map[string]interface{}{
		"exam_id": examID, "student_id": "student_manual",
		"ocr_results": []map[string]interface{}{
			{"page_number": 1, "image_url": objectName},
		},
	}
	resp, err := postRequest(fmt.Sprintf("%s/exams/%s/submissions", BaseURL, examID), tenantID, sub)
	if err != nil { return "", err }
	return resp["id"].(string), nil
}

func triggerGrading(tenantID, subID string) error {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/submissions/%s/trigger-grading", BaseURL, subID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", tenantID)
	resp, err := client.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()
	if resp.StatusCode >= 400 { return fmt.Errorf("error triggering") }
	return nil
}

func getGrades(tenantID, subID string) ([]domain.FinalGrade, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/submissions/%s/grades", BaseURL, subID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", tenantID)
	resp, err := client.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	if resp.StatusCode != 200 { return nil, fmt.Errorf("status %d", resp.StatusCode) }
	var grades []domain.FinalGrade
	json.NewDecoder(resp.Body).Decode(&grades)
	return grades, nil
}

func postRequest(url, tenantID string, data interface{}) (map[string]interface{}, error) {
	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", tenantID)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	if resp.StatusCode >= 400 { return nil, fmt.Errorf("error %d", resp.StatusCode) }
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}
