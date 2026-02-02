package service

import (
	"context"
	"harama/internal/domain"
	"harama/internal/repository/postgres"

	"github.com/google/uuid"
)

type ExamService struct {
	repo *postgres.ExamRepo
}

func NewExamService(repo *postgres.ExamRepo) *ExamService {
	return &ExamService{repo: repo}
}

func (s *ExamService) CreateExam(ctx context.Context, exam *domain.Exam) error {
	return s.repo.Create(ctx, exam)
}

func (s *ExamService) GetExam(ctx context.Context, id uuid.UUID) (*domain.Exam, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ExamService) AddQuestion(ctx context.Context, examID uuid.UUID, question *domain.Question) error {
	question.ExamID = examID
	return s.repo.CreateQuestion(ctx, question)
}

func (s *ExamService) SetRubric(ctx context.Context, questionID uuid.UUID, rubric *domain.Rubric) error {
	rubric.QuestionID = questionID
	return s.repo.UpdateRubric(ctx, rubric)
}
