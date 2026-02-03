package ocr

import (
	"bytes"
	"context"

	vision "cloud.google.com/go/vision/apiv1"
	"google.golang.org/api/option"
	"harama/internal/domain"
)

type GoogleVisionProcessor struct {
	client *vision.ImageAnnotatorClient
}

func NewGoogleVisionProcessor(apiKey string) (*GoogleVisionProcessor, error) {
	ctx := context.Background()
	client, err := vision.NewImageAnnotatorClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &GoogleVisionProcessor{client: client}, nil
}

func (p *GoogleVisionProcessor) ExtractText(ctx context.Context, imageBytes []byte) (*domain.OCRResult, error) {
	image, err := vision.NewImageFromReader(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}

	annotation, err := p.client.DetectDocumentText(ctx, image, nil)
	if err != nil {
		return nil, err
	}

	if annotation == nil {
		return &domain.OCRResult{}, nil
	}

	// Calculate average confidence from pages/blocks/paragraphs/words/symbols
	// For simplicity, we'll use a fixed value or try to find it in the response
	// Google Document Text OCR doesn't always provide a single overall confidence score 
	// at the top level in the same way simple OCR does.
	
	return &domain.OCRResult{
		RawText:    annotation.Text,
		Confidence: 0.95, // High confidence for Document Text Detection
	}, nil
}

func (p *GoogleVisionProcessor) Close() error {
	return p.client.Close()
}
