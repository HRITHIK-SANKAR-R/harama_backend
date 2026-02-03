package profiles

type SubjectProfile struct {
	Subject    string
	FocusAreas []string
	PromptBias string
}

var Subjects = map[string]SubjectProfile{
	"mathematics": {
		Subject:    "Mathematics",
		FocusAreas: []string{"step-by-step reasoning", "calculation accuracy", "formula usage"},
		PromptBias: "Prioritize numerical accuracy and logical derivation steps.",
	},
	"science": {
		Subject:    "Science",
		FocusAreas: []string{"conceptual accuracy", "diagram interpretation", "scientific terminology"},
		PromptBias: "Focus on scientific principles and accurate representation of phenomena.",
	},
	"english": {
		Subject:    "English",
		FocusAreas: []string{"grammar", "coherence", "argument structure", "literary devices"},
		PromptBias: "Value clear expression, persuasive structure, and creative depth.",
	},
}
