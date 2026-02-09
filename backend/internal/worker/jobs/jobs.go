package jobs

import (
	"context"
	"harama/internal/service"
	"harama/internal/types"
	"github.com/google/uuid"
)

type OCRJob struct {
	SubmissionID uuid.UUID
	Service      service.OCRServiceInterface
}

func (j *OCRJob) Execute(ctx context.Context) error {
	return j.Service.ProcessSubmission(ctx, j.SubmissionID)
}

func (j *OCRJob) ID() string {
	return "ocr-" + j.SubmissionID.String()
}

// Ensure OCRJob implements types.Job
var _ types.Job = (*OCRJob)(nil)

type GradingJob struct {
	SubmissionID uuid.UUID
	Service      service.GradingServiceInterface
}

func (j *GradingJob) Execute(ctx context.Context) error {
	return j.Service.GradeSubmission(ctx, j.SubmissionID)
}

func (j *GradingJob) ID() string {
	return "grading-" + j.SubmissionID.String()
}

// Ensure GradingJob implements types.Job
var _ types.Job = (*GradingJob)(nil)
