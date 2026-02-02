package service

import (
	"context"
	"harama/internal/domain"
	"harama/internal/grading"
	"harama/internal/pkg/utils"
	"harama/internal/repository/postgres"

	"github.com/google/uuid"
)

type GradingService struct {
	repo          *postgres.GradeRepo
	examRepo      *postgres.ExamRepo
	subRepo       *postgres.SubmissionRepo
	gradingEngine *grading.Engine
}

func NewGradingService(repo *postgres.GradeRepo, examRepo *postgres.ExamRepo, subRepo *postgres.SubmissionRepo, engine *grading.Engine) *GradingService {
	return &GradingService{
		repo:          repo,
		examRepo:      examRepo,
		subRepo:       subRepo,
		gradingEngine: engine,
	}
}

func (s *GradingService) GradeSubmission(ctx context.Context, submissionID uuid.UUID) error {
	sub, err := s.subRepo.GetByID(ctx, submissionID)
	if err != nil {
		return err
	}

	exam, err := s.examRepo.GetByID(ctx, sub.ExamID)
	if err != nil {
		return err
	}

	for _, answer := range sub.Answers {
		// Find question for this answer
		var targetQuestion *domain.Question
		for _, q := range exam.Questions {
			if q.ID == answer.QuestionID {
				targetQuestion = &q
				break
			}
		}

		if targetQuestion == nil || targetQuestion.Rubric == nil {
			continue
		}

		finalGrade, multiEval, err := s.gradingEngine.GradeAnswer(ctx, answer, *targetQuestion.Rubric)
		if err != nil {
			return err
		}

		finalGrade.SubmissionID = submissionID
		finalGrade.QuestionID = targetQuestion.ID
		
		err = s.repo.SaveFinalGrade(ctx, finalGrade)
		if err != nil {
			return err
		}

		if multiEval.ShouldEscalate {
			escalation := &domain.EscalationCase{
				ID:             uuid.New(),
				SubmissionID:   submissionID,
				QuestionID:     targetQuestion.ID,
				AllEvaluations: multiEval.Evaluations,
				Variance:       multiEval.Variance,
				EscalatedAt:    utils.CurrentTime(),
				Status:         "pending",
			}
			err = s.repo.SaveEscalation(ctx, escalation)
			if err != nil {
				return err
			}
		}
	}

	return s.subRepo.UpdateStatus(ctx, submissionID, domain.StatusCompleted)
}

func (s *GradingService) GetGrades(ctx context.Context, submissionID uuid.UUID) ([]domain.FinalGrade, error) {
	return s.repo.GetBySubmission(ctx, submissionID)
}
