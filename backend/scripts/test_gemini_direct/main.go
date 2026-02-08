package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"harama/internal/config"
)

func main() {
	// Load config to get key
	cfg := config.Load()
	if cfg.GeminiAPIKey == "" || cfg.GeminiAPIKey == "your-gemini-api-key" {
		log.Fatal("GEMINI_API_KEY is not set or is invalid in .env")
	}

	fmt.Printf("Using API Key: %s...%s\n", cfg.GeminiAPIKey[:4], cfg.GeminiAPIKey[len(cfg.GeminiAPIKey)-4:])
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GeminiAPIKey))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	iter := client.ListModels(ctx)
	fmt.Println("Available models:")
	for {
		m, err := iter.Next()
		if err != nil {
			break
		}
		fmt.Printf(" - %s\n", m.Name)
	}

	modelName := "gemini-1.5-flash"
	fmt.Printf("Testing model: %s\n", modelName)
	model := client.GenerativeModel(modelName)

	fmt.Println("Sending request...")
	resp, err := model.GenerateContent(ctx, genai.Text("Hello, are you working?"))
	if err != nil {
		log.Fatalf("GenerateContent failed: %v", err)
	}

	if len(resp.Candidates) > 0 {
		fmt.Println("SUCCESS! Response received:")
		for _, part := range resp.Candidates[0].Content.Parts {
			fmt.Printf("%v\n", part)
		}
	} else {
		fmt.Println("Response received but no candidates.")
	}
}
