package ocr

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"harama/internal/domain"
)

type GeminiOCRProcessor struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiOCRProcessor(apiKey string) (*GeminiOCRProcessor, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	
	// Use Gemini 3 Flash for faster/cheaper OCR, or Pro for accuracy
		model := client.GenerativeModel("gemini-3-flash-preview")
	 
	model.SetTemperature(0.1) // Low temperature for deterministic transcription
	
	return &GeminiOCRProcessor{
		client: client,
		model:  model,
	}, nil
}

func (p *GeminiOCRProcessor) ExtractText(ctx context.Context, fileBytes []byte, mimeType string) (*domain.OCRResult, error) {
	// Construct multimodal prompt
	prompt := genai.Text("Transcribe the handwritten or printed text in this exam page exactly as it appears. Do not correct spelling. Return only the transcribed text.")
	
	var imgData genai.Part
	
	// Handle different mime types if necessary, though Gemini handles standard images/pdfs
	if mimeType == "application/pdf" {
		imgData = genai.Blob{MIMEType: mimeType, Data: fileBytes}
	} else {
		// Default to generic image if unsure, or specific type
		if mimeType == "" {
			mimeType = "image/png"
		}
		imgData = genai.ImageData(mimeType, fileBytes)
	}

	var resp *genai.GenerateContentResponse
	var err error

	// Retry loop for 429 Rate Limits
	maxRetries := 3
	for i := 0; i <= maxRetries; i++ {
		resp, err = p.model.GenerateContent(ctx, prompt, imgData)
		if err == nil {
			break
		}
		
		if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "quota") {
			if i < maxRetries {
				// Increased backoff: 15s, 30s, 60s to handle free tier limits
				waitTime := time.Duration(15*(1<<i)) * time.Second
				fmt.Printf("⚠️ OCR Rate limit hit. Retrying in %v... (Attempt %d/%d)\n", waitTime, i+1, maxRetries)
				time.Sleep(waitTime)
				continue
			}
		}
		return nil, fmt.Errorf("gemini ocr error: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from Gemini OCR")
	}

	part := resp.Candidates[0].Content.Parts[0]
	text, ok := part.(genai.Text)
	if !ok {
		return nil, fmt.Errorf("unexpected response type from Gemini OCR")
	}

	return &domain.OCRResult{
		RawText:    strings.TrimSpace(string(text)),
		Confidence: 0.90, // Gemini doesn't give token-level confidence easily in standard response, defaulting
	}, nil
}

func (p *GeminiOCRProcessor) Close() error {
	return p.client.Close()
}
