package service

import (
	"context"
	"fmt"
	"time"

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
	auditRepo     *postgres.AuditRepo
	gradingEngine *grading.Engine
	events        *EventService
}

func NewGradingService(repo *postgres.GradeRepo, examRepo *postgres.ExamRepo, subRepo *postgres.SubmissionRepo, auditRepo *postgres.AuditRepo, engine *grading.Engine, events *EventService) *GradingService {
	return &GradingService{
		repo:          repo,
		examRepo:      examRepo,
		subRepo:       subRepo,
		auditRepo:     auditRepo,
		gradingEngine: engine,
		events:        events,
	}
}

func (s *GradingService) updateStatus(ctx context.Context, id uuid.UUID, status domain.ProcessingStatus) {
	_ = s.subRepo.UpdateStatus(ctx, id, status)
	if s.events != nil {
		s.events.Broadcast(Event{
			SubmissionID: id,
			Type:         EventStatusUpdate,
			Status:       string(status),
		})
	}
}

func (s *GradingService) GradeSubmission(ctx context.Context, submissionID uuid.UUID) (finalErr error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		if finalErr != nil {
			fmt.Printf("❌ Grading failed for submission %s after %v: %v\n", submissionID, duration, finalErr)
			s.updateStatus(context.Background(), submissionID, domain.StatusFailed)
		} else {
			fmt.Printf("✅ Grading completed for submission %s in %v\n", submissionID, duration)
			s.updateStatus(context.Background(), submissionID, domain.StatusCompleted)
		}
	}()

	// 1. Mark as grading
	s.updateStatus(ctx, submissionID, domain.StatusGrading)

	sub, err := s.subRepo.GetByID(ctx, submissionID)
	if err != nil {
		return err
	}

	exam, err := s.examRepo.GetByID(ctx, sub.ExamID)
	if err != nil {
		return err
	}

	// If there are no answers, it might be because OCR failed or segmentation hasn't happened.
	// We check if we can generate a fallback answer from raw OCR text if available.
	if len(sub.Answers) == 0 && len(sub.OCRResults) > 0 {
		fmt.Printf("⚠️ No answer segments found for submission %s. Attempting fallback from raw OCR text.\n", submissionID)
		// Simple fallback: map first question to all OCR text
		if len(exam.Questions) > 0 {
			fullText := ""
			for _, res := range sub.OCRResults {
				fullText += res.RawText + "\n"
			}
			sub.Answers = []domain.AnswerSegment{
				{
					ID:           uuid.New(),
					SubmissionID: sub.ID,
					QuestionID:   exam.Questions[0].ID,
					Text:         fullText,
				},
			}
		}
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

		finalGrade, multiEval, err := s.gradingEngine.GradeAnswer(ctx, answer, *targetQuestion.Rubric, exam.Subject, targetQuestion.QuestionText)
		if err != nil {
			return err
		}

		finalGrade.SubmissionID = submissionID
		finalGrade.QuestionID = targetQuestion.ID
		
		err = s.repo.SaveFinalGrade(ctx, finalGrade)
		if err != nil {
			return err
		}

		// Log audit event for AI grading
		_ = s.auditRepo.Save(ctx, &domain.AuditLog{
			EntityType: "grade",
			EntityID:   finalGrade.ID,
			EventType:  "ai_graded",
			ActorType:  "ai",
			Changes: map[string]interface{}{
				"score":      finalGrade.FinalScore,
				"confidence": finalGrade.Confidence,
				"reasoning":  finalGrade.Reasoning,
			},
		})

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
