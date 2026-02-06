package service

import (
	"context"
	"harama/internal/domain"

	"github.com/google/uuid"
)

// ExamServiceInterface defines the contract for exam operations
type ExamServiceInterface interface {
	CreateExam(ctx context.Context, exam *domain.Exam) error
	GetExam(ctx context.Context, id uuid.UUID) (*domain.Exam, error)
	AddQuestion(ctx context.Context, examID uuid.UUID, question *domain.Question) error
	SetRubric(ctx context.Context, questionID uuid.UUID, rubric *domain.Rubric) error
	ListExams(ctx context.Context, tenantID uuid.UUID) ([]domain.Exam, error)
}

// Ensure ExamService implements the interface
var _ ExamServiceInterface = (*ExamService)(nil)
