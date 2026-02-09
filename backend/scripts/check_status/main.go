package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"harama/internal/config"
	"harama/internal/repository/postgres"
	"github.com/google/uuid"
)

func main() {
	cfg := config.Load()
	db, err := postgres.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	subIDStr := "8b6a6183-d1fe-4973-b284-47ad449d8577"
	if len(os.Args) > 1 {
		subIDStr = os.Args[1]
	}

	subID, err := uuid.Parse(subIDStr)
	if err != nil {
		log.Fatalf("Invalid UUID: %v", err)
	}

	type Sub struct {
		ProcessingStatus string          `bun:"processing_status"`
		OCRResults       json.RawMessage `bun:"ocr_results"`
	}

	var sub Sub
	err = db.NewSelect().Table("submissions").
		Column("processing_status", "ocr_results").
		Where("id = ?", subID).
		Scan(context.Background(), &sub)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Current Status: %s\n", sub.ProcessingStatus)
	fmt.Printf("OCR Data: %s\n", string(sub.OCRResults))
}

