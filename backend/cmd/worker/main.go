package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"harama/internal/config"
	"harama/internal/repository/postgres"
	"harama/internal/worker"
)

func main() {
	cfg := config.Load()

	db, err := postgres.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create worker pool
	pool := worker.NewWorkerPool(10, 1000)

	// Start workers
	pool.Start()
	log.Println("Worker pool started with 10 workers")

	// Wait for shutdown signal
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	log.Println("Shutting down workers...")
	pool.Stop()
	log.Println("Workers stopped")
}
