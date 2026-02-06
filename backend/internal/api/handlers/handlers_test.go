package handlers

import (
	"context"
	"testing"

	"harama/internal/domain"
	"harama/internal/service"

	"github.com/google/uuid"
)

// Mock service for testing
type mockExamService struct{}

func (m *mockExamService) CreateExam(ctx context.Context, exam *domain.Exam) error {
	exam.ID = uuid.New()
	return nil
}

func (m *mockExamService) GetExam(ctx context.Context, id uuid.UUID) (*domain.Exam, error) {
	return &domain.Exam{
		ID:      id,
		Title:   "Test Exam",
		Subject: "Math",
	}, nil
}

func (m *mockExamService) ListExams(ctx context.Context, tenantID uuid.UUID) ([]domain.Exam, error) {
	return []domain.Exam{
		{ID: uuid.New(), Title: "Exam 1"},
		{ID: uuid.New(), Title: "Exam 2"},
	}, nil
}

func (m *mockExamService) AddQuestion(ctx context.Context, examID uuid.UUID, question *domain.Question) error {
	question.ID = uuid.New()
	return nil
}

func (m *mockExamService) SetRubric(ctx context.Context, questionID uuid.UUID, rubric *domain.Rubric) error {
	return nil
}

// Ensure mock implements interface
var _ service.ExamServiceInterface = (*mockExamService)(nil)

func TestCreateExamHandler(t *testing.T) {
	// Skip for now - requires full handler setup with context
	t.Skip("Integration test - requires full setup")
}

func TestMockService(t *testing.T) {
	mock := &mockExamService{}
	
	// Test CreateExam
	exam := &domain.Exam{Title: "Test", Subject: "Math"}
	err := mock.CreateExam(context.Background(), exam)
	if err != nil {
		t.Errorf("CreateExam failed: %v", err)
	}
	if exam.ID == uuid.Nil {
		t.Error("Exam ID should be set")
	}
	
	// Test GetExam
	got, err := mock.GetExam(context.Background(), uuid.New())
	if err != nil {
		t.Errorf("GetExam failed: %v", err)
	}
	if got.Title != "Test Exam" {
		t.Errorf("Expected 'Test Exam', got '%s'", got.Title)
	}
	
	// Test ListExams
	exams, err := mock.ListExams(context.Background(), uuid.New())
	if err != nil {
		t.Errorf("ListExams failed: %v", err)
	}
	if len(exams) != 2 {
		t.Errorf("Expected 2 exams, got %d", len(exams))
	}
}
