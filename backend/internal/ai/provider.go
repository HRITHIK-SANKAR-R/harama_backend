package ai

import (
    "context"
    "harama/internal/domain"
)

type Provider interface {
    Grade(ctx context.Context, req GradingRequest) (domain.GradingResult, error)
    GenerateFeedback(ctx context.Context, req FeedbackRequest) (string, error)
}

type GradingRequest struct {
    Answer      domain.AnswerSegment
    Rubric      domain.Rubric
    EvaluatorID string
}

type FeedbackRequest struct {
    Grade    domain.FinalGrade
    History  []domain.FeedbackEvent
    StudentName string
}
