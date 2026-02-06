package grading

import (
	"testing"

	"harama/internal/domain"
)

func TestCalculateMean(t *testing.T) {
	e := &Engine{}
	
	tests := []struct {
		scores []float64
		want   float64
	}{
		{[]float64{5, 5, 5}, 5.0},
		{[]float64{3, 4, 5}, 4.0},
		{[]float64{10}, 10.0},
		{[]float64{}, 0.0},
	}
	
	for _, tt := range tests {
		got := e.calculateMean(tt.scores)
		if got != tt.want {
			t.Errorf("calculateMean(%v) = %v, want %v", tt.scores, got, tt.want)
		}
	}
}

func TestCalculateWeightedConsensus(t *testing.T) {
	e := &Engine{}
	
	evaluations := []domain.GradingResult{
		{Score: 8.0, Confidence: 0.9},
		{Score: 7.0, Confidence: 0.7},
		{Score: 9.0, Confidence: 0.8},
	}
	
	result := e.calculateWeightedConsensus(evaluations)
	
	// Expected: (8*0.9 + 7*0.7 + 9*0.8) / (0.9+0.7+0.8) = 8.04
	if result < 7.9 || result > 8.2 {
		t.Errorf("calculateWeightedConsensus() = %v, want ~8.0", result)
	}
}

func TestVarianceCalculation(t *testing.T) {
	calc := NewVarianceCalculator()
	
	results := []domain.GradingResult{
		{Score: 5.0},
		{Score: 5.0},
		{Score: 5.0},
	}
	
	variance := calc.Calculate(results, 5.0)
	if variance != 0.0 {
		t.Errorf("Variance of identical scores should be 0, got %v", variance)
	}
	
	// Test with variance
	results2 := []domain.GradingResult{
		{Score: 3.0},
		{Score: 5.0},
		{Score: 7.0},
	}
	
	variance2 := calc.Calculate(results2, 5.0)
	if variance2 <= 0 {
		t.Errorf("Variance should be > 0 for different scores, got %v", variance2)
	}
}

func TestConfidenceCalculation(t *testing.T) {
	calc := NewConfidenceCalculator()
	
	// High confidence case
	results := []domain.GradingResult{
		{Confidence: 0.9},
		{Confidence: 0.9},
		{Confidence: 0.9},
	}
	
	confidence := calc.Calculate(results, 0.1) // low variance
	if confidence < 0.8 {
		t.Errorf("High individual confidence + low variance should give high confidence, got %v", confidence)
	}
	
	// Low confidence case
	results2 := []domain.GradingResult{
		{Confidence: 0.5},
		{Confidence: 0.5},
		{Confidence: 0.5},
	}
	
	confidence2 := calc.Calculate(results2, 5.0) // high variance
	if confidence2 > 0.6 {
		t.Errorf("Low individual confidence + high variance should give low confidence, got %v", confidence2)
	}
}
