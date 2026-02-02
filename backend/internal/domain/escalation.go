package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type EscalationCase struct {
	bun.BaseModel `bun:"table:escalations,alias:esc"`

	ID             uuid.UUID       `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	SubmissionID   uuid.UUID       `bun:"submission_id,notnull,type:uuid" json:"submission_id"`
	QuestionID     uuid.UUID       `bun:"question_id,notnull,type:uuid" json:"question_id"`
	AllEvaluations []GradingResult `bun:"all_evaluations,type:jsonb" json:"all_evaluations"`
	Variance       float64         `bun:"variance,notnull" json:"variance"`
	EscalatedAt    time.Time       `bun:"escalated_at,nullzero,notnull,default:current_timestamp" json:"escalated_at"`
	AssignedTo     *uuid.UUID      `bun:"assigned_to,type:uuid" json:"assigned_to"`
	Status         string          `bun:"status,notnull,default:'pending'" json:"status"`
}
