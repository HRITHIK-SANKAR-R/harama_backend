package service

import (
	"context"
	"harama/internal/domain"
	"harama/internal/repository/postgres"

	"github.com/google/uuid"
)

type OCRService struct {
	repo *postgres.SubmissionRepo
}

func NewOCRService(repo *postgres.SubmissionRepo) *OCRService {
	return &OCRService{repo: repo}
}

func (s *OCRService) ProcessSubmission(ctx context.Context, submissionID uuid.UUID) error {
	// 1. Mark as processing
	err := s.repo.UpdateStatus(ctx, submissionID, domain.StatusProcessing)
	if err != nil {
		return err
	}

	// 2. Mock OCR results (to be replaced with actual Google Vision/Tesseract)
	results := []domain.OCRResult{
		{
			PageNumber: 1,
			RawText:    "Sample OCR Text",
			Confidence: 0.95,
		},
	}

	return s.repo.SaveOCRResults(ctx, submissionID, results)
}
