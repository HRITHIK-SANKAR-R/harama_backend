package profiles

type EvaluatorProfile struct {
	ID           string
	Name         string
	SystemPrompt string
	Temperature  float64
	Perspective  string // "strict", "lenient", "balanced"
	FocusAreas   []string
}

var Evaluators = map[string]EvaluatorProfile{
	"rubric_enforcer": {
		ID:   "rubric_enforcer",
		Name: "Rubric Enforcer",
		SystemPrompt: `You are a strict grader who follows the rubric exactly.
Award full credit only when ALL criteria are explicitly met.
Do not give partial credit unless the rubric specifically allows it.
Your job is to ensure consistency and fairness.`,
		Temperature: 0.1, // Very deterministic
		Perspective: "strict",
		FocusAreas:  []string{"rubric_compliance", "completeness"},
	},

	"reasoning_validator": {
		ID:   "reasoning_validator",
		Name: "Reasoning Validator",
		SystemPrompt: `You are an educator who values logical thinking.
Reward students for correct reasoning even if execution has minor errors.
Look for conceptual understanding, not just correct final answers.
Partial credit should be generous for good reasoning with small mistakes.`,
		Temperature: 0.3, // More flexible
		Perspective: "lenient",
		FocusAreas:  []string{"logical_flow", "conceptual_understanding"},
	},

	"structural_analyzer": {
		ID:   "structural_analyzer",
		Name: "Structural Analyzer",
		SystemPrompt: `You evaluate answer structure and organization.
Check for: clear introduction, step-by-step work, labeled diagrams.
Penalize disorganized answers even if content is correct.
Reward well-structured answers with clear explanations.`,
		Temperature: 0.2,
		Perspective: "balanced",
		FocusAreas:  []string{"organization", "clarity", "presentation"},
	},
}
