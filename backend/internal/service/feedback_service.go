package service

import (
	"context"
	"harama/internal/domain"
	"harama/internal/repository/postgres"

	"github.com/google/uuid"
)

type FeedbackService struct {
	repo     *postgres.FeedbackRepo
	gradeRepo *postgres.GradeRepo
}

func NewFeedbackService(repo *postgres.FeedbackRepo, gradeRepo *postgres.GradeRepo) *FeedbackService {
	return &FeedbackService{
		repo:     repo,
		gradeRepo: gradeRepo,
	}
}

func (s *FeedbackService) CaptureOverrideFeedback(ctx context.Context, submissionID uuid.UUID, questionID uuid.UUID, teacherScore float64, teacherReason string) error {
	// 1. Get the existing grade to find the AI score and reasoning
	grades, err := s.gradeRepo.GetBySubmission(ctx, submissionID)
	if err != nil {
		return err
	}

	var originalGrade *domain.FinalGrade
	for _, g := range grades {
		if g.QuestionID == questionID {
			originalGrade = &g
			break
		}
	}

	if originalGrade == nil {
		return nil // Or return error if grade should exist
	}

	aiScore := 0.0
	if originalGrade.AIScore != nil {
		aiScore = *originalGrade.AIScore
	}

	event := &domain.FeedbackEvent{
		ID:            uuid.New(),
		QuestionID:    questionID,
		SubmissionID:  submissionID,
		AIScore:       aiScore,
		TeacherScore:  teacherScore,
		Delta:         teacherScore - aiScore,
		AIReasoning:   originalGrade.Reasoning,
		TeacherReason: teacherReason,
	}

	return s.repo.SaveFeedbackEvent(ctx, event)
}

func (s *FeedbackService) GetFeedbackByQuestion(ctx context.Context, questionID uuid.UUID) ([]domain.FeedbackEvent, error) {
	return s.repo.GetFeedbackByQuestion(ctx, questionID)
}
