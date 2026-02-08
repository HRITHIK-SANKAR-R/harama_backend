package gemini

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"harama/internal/ai"
	"harama/internal/domain"
	"harama/internal/grading/profiles"
)

//go:embed prompts/*.txt
var promptsFS embed.FS

type Client struct {
    client *genai.Client
    model  *genai.GenerativeModel
}

func NewClient(apiKey string) (*Client, error) {
    ctx := context.Background()
    client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
    if err != nil {
        return nil, err
    }
    
    model := client.GenerativeModel("gemini-3-flash-preview")
    model.SetTemperature(0.2)
    model.SetTopK(40)
    model.SetTopP(0.95)
    
    return &Client{
        client: client,
        model:  model,
    }, nil
}

func (c *Client) Grade(ctx context.Context, req ai.GradingRequest) (domain.GradingResult, error) {
	// Load appropriate profile
	profile, ok := profiles.Evaluators[req.EvaluatorID]
	if !ok {
		return domain.GradingResult{}, fmt.Errorf("evaluator profile not found: %s", req.EvaluatorID)
	}

	// Load appropriate prompt template
	promptTemplate := loadPromptTemplate(req.EvaluatorID)

	// Build prompt
	prompt := buildGradingPrompt(promptTemplate, req.Answer, req.Rubric, req.Subject, req.QuestionText)

	// Create a local model instance to safely set temperature for this specific call
	model := c.client.GenerativeModel("gemini-3-flash-preview")
	model.SetTemperature(float32(profile.Temperature))

	// Prepare parts for multimodal input
	parts := []genai.Part{genai.Text(prompt)}
// ... (omitted diagrams part for now, it's already there)

	// Add diagrams if present (Phase 2)
	for _, diagramURL := range req.Answer.Diagrams {
		_ = diagramURL // Use URL to fetch or placeholder
	}

	// Call Gemini
	fmt.Printf("DEBUG: Calling Gemini for evaluator %s...\n", req.EvaluatorID)
	
	var resp *genai.GenerateContentResponse
	var err error
	maxRetries := 3

	for i := 0; i <= maxRetries; i++ {
		resp, err = model.GenerateContent(ctx, parts...)
		if err == nil {
			break
		}
		if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "quota") {
			if i < maxRetries {
				// Increased backoff: 15s, 30s, 60s
				waitTime := time.Duration(15*(1<<i)) * time.Second
				fmt.Printf("⚠️ Grading Rate limit hit. Retrying in %v... (Attempt %d/%d)\n", waitTime, i+1, maxRetries)
				time.Sleep(waitTime)
				continue
			}
		}
		return domain.GradingResult{}, fmt.Errorf("gemini API error: %w", err)
	}

	fmt.Printf("DEBUG: Gemini returned (err=%v)\n", err)

	// Parse structured response
	var result domain.GradingResult
	if err := parseResponse(resp, &result); err != nil {
		return domain.GradingResult{}, err
	}

	result.AIEvaluatorID = req.EvaluatorID
	return result, nil
}

func (c *Client) GenerateFeedback(ctx context.Context, req ai.FeedbackRequest) (string, error) {
	prompt := fmt.Sprintf(`
ROLE: Educational feedback specialist

STUDENT: %s
CURRENT GRADE: %.1f/%.1f
AI REASONING: %s

HISTORY OF OVERRIDES:
%v

TASK:
Generate personalized, encouraging, and actionable feedback for the student.
Focus on:
1. What they did well.
2. Specific areas for improvement based on current mistakes and history.
3. Actionable next steps.

TONE: Encouraging but honest.
LENGTH: 3-4 sentences.
`, req.StudentName, req.Grade.FinalScore, req.Grade.Confidence, req.Grade.Reasoning, req.History)

	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("gemini API error: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini")
	}

	part := resp.Candidates[0].Content.Parts[0]
	if text, ok := part.(genai.Text); ok {
		return string(text), nil
	}

	return "", fmt.Errorf("unexpected response part type")
}

func (c *Client) AnalyzePatterns(ctx context.Context, req ai.AnalysisRequest) (ai.AnalysisResult, error) {
	prompt := fmt.Sprintf(`
ROLE: Educational pattern analyst

QUESTION ID: %s
RUBRIC: %v
FEEDBACK EVENTS:
%v

TASK:
Analyze the feedback events where teachers have overridden AI scores.
Identify:
1. Systematic biases (e.g., AI is consistently too strict on units).
2. Common themes in teacher corrections.
3. Recommendations for improving the rubric or system prompt.

OUTPUT JSON FORMAT:
{
  "patterns": ["pattern 1", "pattern 2"],
  "common_reasons": ["reason 1", "reason 2"],
  "recommendation": "detailed recommendation"
}
`, req.QuestionID, req.Rubric, req.Events)

	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return ai.AnalysisResult{}, fmt.Errorf("gemini API error: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return ai.AnalysisResult{}, fmt.Errorf("empty response from Gemini")
	}

	part := resp.Candidates[0].Content.Parts[0]
	textPart, ok := part.(genai.Text)
	if !ok {
		return ai.AnalysisResult{}, fmt.Errorf("unexpected response part type")
	}

	var analysisResult ai.AnalysisResult
	if err := json.Unmarshal([]byte(textPart), &analysisResult); err != nil {
		// If it's not JSON, we might need a more robust parser or just return as recommendation
		return ai.AnalysisResult{
			Recommendation: string(textPart),
		}, nil
	}

	return analysisResult, nil
}

func (c *Client) RefineRubric(ctx context.Context, req ai.RefineRubricRequest) (domain.Rubric, error) {
	prompt := fmt.Sprintf(`
ROLE: Rubric Refinement Specialist

CURRENT RUBRIC:
%v

ANALYSIS OF GRADING MISTAKES:
Patterns: %v
Recommendations: %s

TASK:
Update the rubric to address the identified issues.
- Adjust point weights if criteria are consistently undervalued or overvalued.
- Clarify descriptions if they are ambiguous.
- Add or modify partial credit rules if needed.

OUTPUT:
Return ONLY the updated Rubric JSON.
`, req.CurrentRubric, req.Analysis.Patterns, req.Analysis.Recommendation)

	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return domain.Rubric{}, fmt.Errorf("gemini API error: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return domain.Rubric{}, fmt.Errorf("empty response from Gemini")
	}

	part := resp.Candidates[0].Content.Parts[0]
	textPart, ok := part.(genai.Text)
	if !ok {
		return domain.Rubric{}, fmt.Errorf("unexpected response part type")
	}

	// Clean up markdown code blocks if present
	cleanJSON := strings.TrimPrefix(string(textPart), "```json")
	cleanJSON = strings.TrimSuffix(cleanJSON, "```")
	cleanJSON = strings.TrimSpace(cleanJSON)

	var refinedRubric domain.Rubric
	if err := json.Unmarshal([]byte(cleanJSON), &refinedRubric); err != nil {
		return domain.Rubric{}, fmt.Errorf("failed to parse refined rubric: %w", err)
	}

	// Preserve ID and QuestionID from the request as the AI might hallucinate or omit them
	refinedRubric.ID = req.CurrentRubric.ID
	refinedRubric.QuestionID = req.CurrentRubric.QuestionID

	return refinedRubric, nil
}

// Stubs for helper functions
func loadPromptTemplate(evaluatorID string) string {
	data, err := promptsFS.ReadFile(filepath.Join("prompts", evaluatorID+".txt"))
	if err != nil {
		return ""
	}
	return string(data)
}

func buildGradingPrompt(evalTmpl string, answer domain.AnswerSegment, rubric domain.Rubric, subject string, questionText string) string {
	baseData, err := promptsFS.ReadFile("prompts/base_grading.txt")
	if err != nil {
		return ""
	}

	subjectProfile, _ := profiles.Subjects[strings.ToLower(subject)]

	// Execute base template
	type baseVars struct {
		QuestionText string
		RubricJSON   string
		AnswerText   string
		MaxPoints    int
	}

	rubricJSON, _ := json.MarshalIndent(rubric, "", "  ")

	baseVarsData := baseVars{
		QuestionText: questionText,
		RubricJSON:   string(rubricJSON),
		AnswerText:   answer.Text,
		MaxPoints:    10, // Default or from rubric if we added it
	}

	tmpl, _ := template.New("base").Parse(string(baseData))
	var baseBuf bytes.Buffer
	tmpl.Execute(&baseBuf, baseVarsData)

	// Execute evaluator template
	evalVars := struct {
		BasePrompt   string
		SubjectFocus string
	}{
		BasePrompt:   baseBuf.String(),
		SubjectFocus: subjectProfile.PromptBias,
	}

	eTmpl, _ := template.New("eval").Parse(evalTmpl)
	var evalBuf bytes.Buffer
	eTmpl.Execute(&evalBuf, evalVars)

	return evalBuf.String()
}

func parseResponse(resp *genai.GenerateContentResponse, result *domain.GradingResult) error {
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return fmt.Errorf("empty response")
	}

	part := resp.Candidates[0].Content.Parts[0]
	text, ok := part.(genai.Text)
	if !ok {
		return fmt.Errorf("unexpected part type")
	}

	// Clean up markdown code blocks if present
	cleanJSON := strings.TrimPrefix(string(text), "```json")
	cleanJSON = strings.TrimSuffix(cleanJSON, "```")
	cleanJSON = strings.TrimSpace(cleanJSON)

	return json.Unmarshal([]byte(cleanJSON), result)
}
