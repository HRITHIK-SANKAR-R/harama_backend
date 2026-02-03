package service

import (
	"context"
	"fmt"
	"harama/internal/domain"
	"harama/internal/ocr"
	"harama/internal/repository/postgres"
	"harama/internal/storage"

	"github.com/google/uuid"
)

type OCRService struct {
	repo    *postgres.SubmissionRepo
	storage *storage.MinioStorage
	vision  *ocr.GoogleVisionProcessor
}

func NewOCRService(repo *postgres.SubmissionRepo, storage *storage.MinioStorage, vision *ocr.GoogleVisionProcessor) *OCRService {
	return &OCRService{
		repo:    repo,
		storage: storage,
		vision:  vision,
	}
}

func (s *OCRService) ProcessSubmission(ctx context.Context, submissionID uuid.UUID) error {
	// 1. Get submission metadata
	sub, err := s.repo.GetByID(ctx, submissionID)
	if err != nil {
		return err
	}

	// 2. Mark as processing
	err = s.repo.UpdateStatus(ctx, submissionID, domain.StatusProcessing)
	if err != nil {
		return err
	}

	// 3. Process each OCR result (which should have ImageURL as object name for now or we derive it)
	var finalResults []domain.OCRResult
	for _, res := range sub.OCRResults {
		// Assume ImageURL is the object name in MinIO for this mock/implementation
		imgBytes, err := s.storage.GetFile(ctx, res.ImageURL)
		if err != nil {
			return fmt.Errorf("failed to get file from storage: %w", err)
		}

		ocrResult, err := s.vision.ExtractText(ctx, imgBytes)
		if err != nil {
			return fmt.Errorf("failed to extract text with vision: %w", err)
		}

		ocrResult.PageNumber = res.PageNumber
		ocrResult.ImageURL = res.ImageURL
		finalResults = append(finalResults, *ocrResult)
	}

	return s.repo.SaveOCRResults(ctx, submissionID, finalResults)
}
