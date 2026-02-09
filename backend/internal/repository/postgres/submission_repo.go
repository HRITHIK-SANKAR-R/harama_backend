package postgres

import (
	"context"
	"harama/internal/domain"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type SubmissionRepo struct {
	db *bun.DB
}

func NewSubmissionRepo(db *bun.DB) *SubmissionRepo {
	return &SubmissionRepo{db: db}
}

func (r *SubmissionRepo) Create(ctx context.Context, sub *domain.Submission) error {
	_, err := r.db.NewInsert().Model(sub).Exec(ctx)
	return err
}

func (r *SubmissionRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Submission, error) {
	sub := new(domain.Submission)
	err := r.db.NewSelect().
		Model(sub).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (r *SubmissionRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ProcessingStatus) error {
	// Enforce forward-only transitions in SQL to prevent race conditions or backward jumps
	// queued -> ocr_processing -> (ocr_done | ocr_failed | ocr_timeout) -> grading -> completed | failed
	_, err := r.db.NewUpdate().
		Model((*domain.Submission)(nil)).
		Set("processing_status = ?", status).
		Where("id = ?", id).
		// Only allow updating if the current status is earlier in the pipeline
		// We use a CASE expression to define the rank/order of states
		Where(`(
			CASE processing_status 
				WHEN 'queued' THEN 1 
				WHEN 'ocr_processing' THEN 2 
				WHEN 'ocr_done' THEN 3 
				WHEN 'ocr_failed' THEN 3 
				WHEN 'ocr_timeout' THEN 3 
				WHEN 'grading' THEN 4 
				WHEN 'completed' THEN 5 
				WHEN 'failed' THEN 5 
				ELSE 0 
			END
		) < (
			CASE ? 
				WHEN 'queued' THEN 1 
				WHEN 'ocr_processing' THEN 2 
				WHEN 'ocr_done' THEN 3 
				WHEN 'ocr_failed' THEN 3 
				WHEN 'ocr_timeout' THEN 3 
				WHEN 'grading' THEN 4 
				WHEN 'completed' THEN 5 
				WHEN 'failed' THEN 5 
				ELSE 0 
			END
		)`, status).
		Exec(ctx)
	return err
}

func (r *SubmissionRepo) SaveOCRResults(ctx context.Context, id uuid.UUID, results []domain.OCRResult) error {
	_, err := r.db.NewUpdate().
		Model((*domain.Submission)(nil)).
		Set("ocr_results = ?", results).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *SubmissionRepo) ListByExam(ctx context.Context, examID uuid.UUID) ([]domain.Submission, error) {
	var subs []domain.Submission
	err := r.db.NewSelect().
		Model(&subs).
		Where("exam_id = ?", examID).
		Scan(ctx)
	return subs, err
}

func (r *SubmissionRepo) ListPendingReviews(ctx context.Context, tenantID uuid.UUID) ([]domain.Submission, error) {
	var subs []domain.Submission
	// Find submissions that belong to the tenant AND have at least one grade with status 'needs_review'
	err := r.db.NewSelect().
		Model(&subs).
		Join("JOIN grades ON grades.submission_id = submission.id").
		Where("submission.tenant_id = ?", tenantID).
		Where("grades.status = ?", domain.GradeStatusReview).
		Group("submission.id").
		Order("submission.uploaded_at DESC").
		Scan(ctx)
	return subs, err
}
