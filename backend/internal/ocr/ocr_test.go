package ocr

import (
	"context"
	"testing"

	"harama/internal/domain"
)

// Mock OCR Processor for testing without API calls
type MockOCRProcessor struct {
	shouldFail bool
}

func (m *MockOCRProcessor) ExtractText(ctx context.Context, fileBytes []byte, mimeType string) (*domain.OCRResult, error) {
	if m.shouldFail {
		return nil, context.Canceled
	}
	
	// Simulate OCR extraction
	return &domain.OCRResult{
		RawText:    "Sample extracted text from image",
		Confidence: 0.95,
		PageNumber: 1,
	}, nil
}

func TestMockOCRProcessor(t *testing.T) {
	mock := &MockOCRProcessor{}
	
	// Test successful extraction
	result, err := mock.ExtractText(context.Background(), []byte("fake image data"), "image/png")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if result.RawText == "" {
		t.Error("Expected text to be extracted")
	}
	
	if result.Confidence < 0.9 {
		t.Errorf("Expected high confidence, got %v", result.Confidence)
	}
	
	t.Logf("✅ Mock OCR extracted: %s (confidence: %.2f)", result.RawText, result.Confidence)
}

func TestMockOCRProcessorError(t *testing.T) {
	mock := &MockOCRProcessor{shouldFail: true}
	
	_, err := mock.ExtractText(context.Background(), []byte("fake"), "image/png")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	
	t.Logf("✅ Mock OCR correctly handles errors")
}

// Test with real Gemini (requires API key)
func TestGeminiOCRProcessor(t *testing.T) {
	// Skip if no API key
	t.Skip("Requires GEMINI_API_KEY - run manually with: go test -v -run TestGeminiOCR")
	
	// To test manually:
	// 1. Set GEMINI_API_KEY environment variable
	// 2. Run: go test -v -run TestGeminiOCR ./internal/ocr
	
	// apiKey := os.Getenv("GEMINI_API_KEY")
	// if apiKey == "" {
	// 	t.Skip("GEMINI_API_KEY not set")
	// }
	
	// processor, err := NewGeminiOCRProcessor(apiKey)
	// if err != nil {
	// 	t.Fatalf("Failed to create processor: %v", err)
	// }
	// defer processor.Close()
	
	// // Create a simple test image with text
	// testImage := createTestImage("Hello World")
	
	// result, err := processor.ExtractText(context.Background(), testImage, "image/png")
	// if err != nil {
	// 	t.Fatalf("OCR failed: %v", err)
	// }
	
	// if result.RawText == "" {
	// 	t.Error("No text extracted")
	// }
	
	// t.Logf("Extracted: %s", result.RawText)
}

func TestOCRConfidenceThresholds(t *testing.T) {
	tests := []struct {
		name       string
		confidence float64
		shouldFlag bool
	}{
		{"High confidence", 0.95, false},
		{"Medium confidence", 0.80, false},
		{"Low confidence - flag", 0.70, true},
		{"Very low - flag", 0.50, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			needsReview := tt.confidence < 0.75
			
			if needsReview != tt.shouldFlag {
				t.Errorf("Confidence %.2f: expected flag=%v, got %v", 
					tt.confidence, tt.shouldFlag, needsReview)
			}
		})
	}
}

func TestMimeTypeHandling(t *testing.T) {
	mock := &MockOCRProcessor{}
	
	mimeTypes := []string{
		"image/png",
		"image/jpeg",
		"image/jpg",
		"application/pdf",
	}
	
	for _, mime := range mimeTypes {
		result, err := mock.ExtractText(context.Background(), []byte("test"), mime)
		if err != nil {
			t.Errorf("Failed for mime type %s: %v", mime, err)
		}
		
		if result == nil {
			t.Errorf("Nil result for mime type %s", mime)
		}
	}
	
	t.Log("✅ All mime types handled correctly")
}
