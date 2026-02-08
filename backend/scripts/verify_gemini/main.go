package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		// Try to load from .env if not in env
		log.Println("GEMINI_API_KEY not found in environment. Checking .env file...")
		// (Simple .env parser for this script)
		content, err := os.ReadFile(".env")
		if err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "GEMINI_API_KEY=") {
					apiKey = strings.TrimSpace(strings.TrimPrefix(line, "GEMINI_API_KEY="))
					break
				}
			}
		}
	}

	if apiKey == "" {
		log.Fatal("❌ GEMINI_API_KEY is missing. Please set it in .env or export it.")
	}

	log.Printf("✅ Found API Key: %s...%s", apiKey[:4], apiKey[len(apiKey)-4:])

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("❌ Failed to create client: %v", err)
	}
	defer client.Close()

	// Test Model
	modelName := "gemini-1.5-flash" 
	// We'll try the one in the code too
	codeModelName := "gemini-3-flash-preview"

	log.Printf("Testing model: %s...", modelName)
	model := client.GenerativeModel(modelName)
	resp, err := model.GenerateContent(ctx, genai.Text("Hello, can you hear me?"))
	if err != nil {
		log.Printf("⚠️ Model %s failed: %v", modelName, err)
	} else {
		log.Printf("✅ Model %s is working! Response: %v", modelName, resp.Candidates[0].Content.Parts[0])
	}

	log.Printf("Testing code's model: %s...", codeModelName)
	model2 := client.GenerativeModel(codeModelName)
	_, err = model2.GenerateContent(ctx, genai.Text("Hello?"))
	if err != nil {
		log.Printf("❌ Code's model %s failed: %v", codeModelName, err)
		log.Println("👉 Recommendation: Change the model name in the code.")
	} else {
		log.Printf("✅ Code's model %s is working!", codeModelName)
	}
}
