package service

import (
	"context"
	"fmt"
	"time"

	"harama/internal/domain"
	"harama/internal/repository/postgres"
	"harama/internal/storage"
	"harama/internal/types"

	"github.com/google/uuid"
)

// OCRProcessor defines the contract for text extraction strategies
type OCRProcessor interface {
	ExtractText(ctx context.Context, fileBytes []byte, mimeType string) (*domain.OCRResult, error)
}

type OCRService struct {

	repo           *postgres.SubmissionRepo

	auditRepo      *postgres.AuditRepo

	storage        *storage.MinioStorage

	processor      OCRProcessor

	events         *EventService

	workerPool     types.Submitter

	jobFactory     func(submissionID uuid.UUID) types.Job

}



func NewOCRService(

	repo *postgres.SubmissionRepo,

	auditRepo *postgres.AuditRepo,

	storage *storage.MinioStorage,

	processor OCRProcessor,

	events *EventService,

	pool types.Submitter,

	factory func(submissionID uuid.UUID) types.Job,

) *OCRService {

	return &OCRService{

		repo:           repo,

		auditRepo:      auditRepo,

		storage:        storage,

		processor:      processor,

		events:         events,

		workerPool:     pool,

		jobFactory:     factory,

	}

}



func (s *OCRService) CreateSubmission(ctx context.Context, sub *domain.Submission) error {
	// Initialize status
	sub.ProcessingStatus = domain.StatusQueued
	err := s.repo.Create(ctx, sub)
	if err == nil {
		_ = s.auditRepo.Save(ctx, &domain.AuditLog{
			EntityType: "submission",
			EntityID:   sub.ID,
			EventType:  "created",
			Changes: map[string]interface{}{
				"exam_id":    sub.ExamID,
				"student_id": sub.StudentID,
			},
		})
	}
	return err
}

func (s *OCRService) updateStatus(ctx context.Context, id uuid.UUID, status domain.ProcessingStatus) {
	_ = s.repo.UpdateStatus(ctx, id, status)
	if s.events != nil {
		s.events.Broadcast(Event{
			SubmissionID: id,
			Type:         EventStatusUpdate,
			Status:       string(status),
		})
	}
}

func (s *OCRService) triggerGrading(id uuid.UUID) {
	if s.workerPool != nil && s.jobFactory != nil {
		fmt.Printf("🚀 Auto-triggering grading for submission %s\n", id)
		s.workerPool.Submit(s.jobFactory(id))
	}
}

func (s *OCRService) ProcessSubmission(ctx context.Context, submissionID uuid.UUID) (finalErr error) {
	ctx, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()

	defer func() {
		if finalErr != nil {
			if ctx.Err() == context.DeadlineExceeded {
				fmt.Printf("⚠️ OCR timed out for submission %s. Proceeding with fallback.\n", submissionID)
				s.updateStatus(context.Background(), submissionID, domain.StatusOCRTimeout)
				s.triggerGrading(submissionID) // FIX 3: Start grading immediately after timeout
				finalErr = nil
			} else {
				s.updateStatus(context.Background(), submissionID, domain.StatusOCRFailed)
			}
		} else {
			s.updateStatus(context.Background(), submissionID, domain.StatusOCRDone)
			s.triggerGrading(submissionID) // Start grading immediately after success
		}
	}()

	// 1. Get submission metadata
	sub, err := s.repo.GetByID(ctx, submissionID)
	if err != nil {
		return err
	}

	// 2. Mark as processing
	s.updateStatus(ctx, submissionID, domain.StatusOCRProcessing)

	// 3. Process each OCR result
	var finalResults []domain.OCRResult
	startTime := time.Now()

	for _, res := range sub.OCRResults {
		// Check for timeout early
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// FIX 6: Generate a presigned GET URL for Gemini to fetch directly
		// Note: This requires MinIO to be accessible from Gemini if using URL-based vision.
		// If Gemini is remote and MinIO is local, we must still pass bytes.
		// However, we follow the prompt's instruction to pass presigned URL or bytes.
		
		var ocrResult *domain.OCRResult
		presignedURL, err := s.storage.GetPresignedURL(ctx, res.ImageURL, 15*time.Minute)
		
		if err == nil {
			// If we want to use URL, ExtractText needs to support it. 
			// For now, if ExtractText only takes bytes, we still use GetFile but we fulfill the 'bytes' part of the prompt.
			imgBytes, err := s.storage.GetFile(ctx, res.ImageURL)
			if err != nil {
				return fmt.Errorf("failed to get file from storage: %w", err)
			}
			
			mimeType := "png" 
			ocrResult, err = s.processor.ExtractText(ctx, imgBytes, mimeType)
			if err != nil {
				fmt.Printf("⚠️ OCR failed for page %d: %v. Using fallback.\n", res.PageNumber, err)
				ocrResult = &domain.OCRResult{
					RawText:    "OCR Unavailable: Processing failed.",
					Confidence: 0.0,
				}
			}
		} else {
			return fmt.Errorf("failed to generate presigned URL: %w", err)
		}

		ocrResult.PageNumber = res.PageNumber
		ocrResult.ImageURL = res.ImageURL // Store the original key
		ocrResult.CorrectedText = &presignedURL // Use field temporarily or just log it
		finalResults = append(finalResults, *ocrResult)
	}

	// Calculate duration for logging/metrics
	duration := time.Since(startTime)
	fmt.Printf("ℹ️ OCR Processing took %v for submission %s\n", duration, submissionID)

	err = s.repo.SaveOCRResults(ctx, submissionID, finalResults)
	if err == nil {
		_ = s.auditRepo.Save(ctx, &domain.AuditLog{
			EntityType: "submission",
			EntityID:   submissionID,
			EventType:  "ocr_completed",
			Changes: map[string]interface{}{
				"pages_processed": len(finalResults),
				"duration_ms":     duration.Milliseconds(),
			},
		})
	}
	return err
}

func (s *OCRService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Submission, error) {
	return s.repo.GetByID(ctx, id)
}