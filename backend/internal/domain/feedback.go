package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type FeedbackEvent struct {
	bun.BaseModel `bun:"table:feedback_events,alias:fe"`

	ID           uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	QuestionID   uuid.UUID `bun:"question_id,notnull,type:uuid" json:"question_id"`
	SubmissionID uuid.UUID `bun:"submission_id,notnull,type:uuid" json:"submission_id"`

	AIScore      float64   `bun:"ai_score,notnull" json:"ai_score"`
	TeacherScore float64   `bun:"teacher_score,notnull" json:"teacher_score"`
	Delta        float64   `bun:"delta,notnull" json:"delta"`

	AIReasoning   string    `bun:"ai_reasoning" json:"ai_reasoning"`
	TeacherReason string    `bun:"teacher_reason" json:"teacher_reason"`

	Timestamp time.Time `bun:"timestamp,nullzero,notnull,default:current_timestamp" json:"timestamp"`
}

type FeedbackPattern struct {
	QuestionID     uuid.UUID `json:"question_id"`
	OverrideCount  int       `json:"override_count"`
	AvgCorrection  float64   `json:"avg_correction"`
	CommonReasons  []string  `json:"common_reasons"`
	Recommendation string    `json:"recommendation"`
}
