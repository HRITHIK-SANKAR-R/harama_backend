package grading

import (
	"context"
	"fmt"
	"math"
	"sync"

	"harama/internal/ai"
	"harama/internal/domain"
	"harama/internal/pkg/utils"
)

type Engine struct {
	aiProvider     ai.Provider
	confidenceCalc *ConfidenceCalculator
	varianceCalc   *VarianceCalculator
}

func NewEngine(provider ai.Provider) *Engine {
	return &Engine{
		aiProvider:     provider,
		confidenceCalc: NewConfidenceCalculator(),
		varianceCalc:   NewVarianceCalculator(),
	}
}

func (e *Engine) GradeAnswer(ctx context.Context, answer domain.AnswerSegment, rubric domain.Rubric, subject string) (*domain.FinalGrade, *domain.MultiEvalResult, error) {
	// Multi-evaluator grading
	multiEval, err := e.multiEvaluatorGrade(ctx, answer, rubric, subject)
	if err != nil {
		return nil, nil, fmt.Errorf("multi-evaluator grading failed: %w", err)
	}

	// Build consensus grade
	finalGrade := e.buildConsensus(multiEval)

	// Check if escalation needed
	if multiEval.ShouldEscalate {
		finalGrade.Status = domain.GradeStatusReview
	} else {
		finalGrade.Status = domain.GradeStatusAutoGraded
	}

	return finalGrade, multiEval, nil
}

func (e *Engine) multiEvaluatorGrade(ctx context.Context, answer domain.AnswerSegment, rubric domain.Rubric, subject string) (*domain.MultiEvalResult, error) {
	evaluatorIDs := []string{
		"rubric_enforcer",
		"reasoning_validator",
		"structural_analyzer",
	}

	type resultTask struct {
		result domain.GradingResult
		err    error
	}
	resChan := make(chan resultTask, len(evaluatorIDs))

	var wg sync.WaitGroup
	for _, id := range evaluatorIDs {
		wg.Add(1)
		go func(evalID string) {
			defer wg.Done()
			res, err := e.aiProvider.Grade(ctx, ai.GradingRequest{
				Answer:      answer,
				Rubric:      rubric,
				EvaluatorID: evalID,
				Subject:     subject,
			})
			resChan <- resultTask{res, err}
		}(id)
	}

	wg.Wait()
	close(resChan)

	var results []domain.GradingResult
	for res := range resChan {
		if res.err != nil {
			return nil, res.err
		}
		results = append(results, res.result)
	}

	// Analyze multi-eval results
	scores := make([]float64, len(results))
	for i, r := range results {
		scores[i] = r.Score
	}

	mean := e.calculateMean(scores)
	variance := e.varianceCalc.Calculate(results, mean)
	confidence := e.confidenceCalc.Calculate(results, variance)

	// Threshold for escalation (e.g., 15% of max points)
	// Assuming max points is available in first result for now
	maxScore := 1.0
	if len(results) > 0 {
		maxScore = float64(results[0].MaxScore)
	}
	shouldEscalate := variance > (0.15 * maxScore) || confidence < 0.7

	return &domain.MultiEvalResult{
		Evaluations:    results,
		Variance:       variance,
		MeanScore:      mean,
		ConsensusScore: e.calculateWeightedConsensus(results),
		Confidence:     confidence,
		ShouldEscalate: shouldEscalate,
		Reasoning:      e.generateConsensusReasoning(results, variance, confidence),
	}, nil
}

func (e *Engine) generateConsensusReasoning(evaluations []domain.GradingResult, variance float64, confidence float64) string {
	if variance < 1.0 {
		return fmt.Sprintf("All evaluators agree (variance: %.2f). High confidence in consensus.", variance)
	}

	if confidence < 0.7 {
		return "Low overall confidence due to high variance between evaluators. Escalating for review."
	}

	return fmt.Sprintf("Moderate variance (%.2f). Consensus reached through weighted average.", variance)
}

func (e *Engine) calculateMean(scores []float64) float64 {
	if len(scores) == 0 {
		return 0
	}
	sum := 0.0
	for _, s := range scores {
		sum += s
	}
	return sum / float64(len(scores))
}

func (e *Engine) calculateWeightedConsensus(evaluations []domain.GradingResult) float64 {
	totalWeight := 0.0
	weightedSum := 0.0

	for _, eval := range evaluations {
		weight := eval.Confidence
		totalWeight += weight
		weightedSum += eval.Score * weight
	}

	if totalWeight == 0 {
		return 0
	}

	return weightedSum / totalWeight
}

func (e *Engine) buildConsensus(multiEval *domain.MultiEvalResult) *domain.FinalGrade {
	return &domain.FinalGrade{
		FinalScore: multiEval.ConsensusScore,
		AIScore:    &multiEval.ConsensusScore,
		Confidence: multiEval.Confidence,
		Reasoning:  multiEval.Reasoning,
		UpdatedAt:  utils.CurrentTime(),
	}
}

type ConfidenceCalculator struct{}

func NewConfidenceCalculator() *ConfidenceCalculator { return &ConfidenceCalculator{} }
func (c *ConfidenceCalculator) Calculate(results []domain.GradingResult, variance float64) float64 {
	if len(results) == 0 {
		return 0
	}
	avgIndividualConfidence := 0.0
	for _, res := range results {
		avgIndividualConfidence += res.Confidence
	}
	avgIndividualConfidence /= float64(len(results))

	// Variance penalty (normalized)
	variancePenalty := math.Max(0, 1.0-(variance/10.0))

	return (avgIndividualConfidence * 0.6) + (variancePenalty * 0.4)
}

type VarianceCalculator struct{}

func NewVarianceCalculator() *VarianceCalculator { return &VarianceCalculator{} }
func (v *VarianceCalculator) Calculate(results []domain.GradingResult, mean float64) float64 {
	if len(results) < 2 {
		return 0
	}
	variance := 0.0
	for _, res := range results {
		variance += math.Pow(res.Score-mean, 2)
	}
	return variance / float64(len(results))
}
