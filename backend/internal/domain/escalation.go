package domain

import (
	"time"

	"github.com/google/uuid"
)

type EscalationCase struct {
	ID             uuid.UUID       `json:"id"`
	SubmissionID   uuid.UUID       `json:"submission_id"`
	QuestionID     uuid.UUID       `json:"question_id"`
	AllEvaluations []GradingResult `json:"all_evaluations"`
	Variance       float64         `json:"variance"`
	EscalatedAt    time.Time       `json:"escalated_at"`
	AssignedTo     *uuid.UUID      `json:"assigned_to"` // Teacher
	Status         string          `json:"status"`      // 'pending', 'resolved'
}
