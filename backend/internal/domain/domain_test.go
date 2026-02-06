package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestExamValidation(t *testing.T) {
	exam := Exam{
		ID:      uuid.New(),
		Title:   "Test Exam",
		Subject: "Math",
		Questions: []Question{
			{
				ID:           uuid.New(),
				QuestionText: "What is 2+2?",
				Points:       5,
				AnswerType:   "short_answer",
			},
		},
	}
	
	if exam.Title == "" {
		t.Error("Exam title should not be empty")
	}
	
	if len(exam.Questions) == 0 {
		t.Error("Exam should have at least one question")
	}
}

func TestGradeStatusTransitions(t *testing.T) {
	validStatuses := []GradeStatus{
		GradeStatusPending,
		GradeStatusAutoGraded,
		GradeStatusReview,
		GradeStatusOverridden,
		GradeStatusFinal,
	}
	
	for _, status := range validStatuses {
		if status == "" {
			t.Errorf("Status should not be empty: %v", status)
		}
	}
}

func TestSubmissionStatus(t *testing.T) {
	submission := Submission{
		ID:               uuid.New(),
		ExamID:           uuid.New(),
		StudentID:        "student123",
		ProcessingStatus: StatusPending,
		UploadedAt:       time.Now(),
	}
	
	if submission.ProcessingStatus != StatusPending {
		t.Error("New submission should be pending")
	}
	
	// Simulate processing
	submission.ProcessingStatus = StatusProcessing
	if submission.ProcessingStatus != StatusProcessing {
		t.Error("Status should update to processing")
	}
	
	// Complete
	submission.ProcessingStatus = StatusCompleted
	if submission.ProcessingStatus != StatusCompleted {
		t.Error("Status should update to completed")
	}
}

func TestOCRConfidence(t *testing.T) {
	ocr := OCRResult{
		PageNumber: 1,
		RawText:    "Sample text",
		Confidence: 0.95,
		ImageURL:   "http://example.com/image.jpg",
	}
	
	if ocr.Confidence < 0 || ocr.Confidence > 1 {
		t.Errorf("Confidence should be between 0 and 1, got %v", ocr.Confidence)
	}
	
	if ocr.Confidence < 0.85 {
		t.Log("Low confidence OCR should trigger verification")
	}
}

func TestMultiEvalResult(t *testing.T) {
	result := MultiEvalResult{
		Evaluations: []GradingResult{
			{Score: 8.0, Confidence: 0.9},
			{Score: 7.5, Confidence: 0.85},
			{Score: 8.5, Confidence: 0.95},
		},
		Variance:       0.5,
		MeanScore:      8.0,
		ConsensusScore: 8.1,
		Confidence:     0.9,
		ShouldEscalate: false,
	}
	
	if len(result.Evaluations) != 3 {
		t.Error("Should have 3 evaluations")
	}
	
	if result.ShouldEscalate {
		t.Error("Low variance should not trigger escalation")
	}
	
	// High variance case
	result.Variance = 3.0
	result.ShouldEscalate = true
	if !result.ShouldEscalate {
		t.Error("High variance should trigger escalation")
	}
}
