package postgres

import (
	"context"
	"harama/internal/domain"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type FeedbackRepo struct {
	db *bun.DB
}

func NewFeedbackRepo(db *bun.DB) *FeedbackRepo {
	return &FeedbackRepo{db: db}
}

func (r *FeedbackRepo) SaveFeedbackEvent(ctx context.Context, event *domain.FeedbackEvent) error {
	_, err := r.db.NewInsert().Model(event).Exec(ctx)
	return err
}

func (r *FeedbackRepo) GetFeedbackByQuestion(ctx context.Context, questionID uuid.UUID) ([]domain.FeedbackEvent, error) {
	var events []domain.FeedbackEvent
	err := r.db.NewSelect().
		Model(&events).
		Where("question_id = ?", questionID).
		Order("timestamp DESC").
		Scan(ctx)
	return events, err
}

func (r *FeedbackRepo) GetFeedbackBySubmission(ctx context.Context, submissionID uuid.UUID) ([]domain.FeedbackEvent, error) {
	var events []domain.FeedbackEvent
	err := r.db.NewSelect().
		Model(&events).
		Where("submission_id = ?", submissionID).
		Order("timestamp DESC").
		Scan(ctx)
	return events, err
}
