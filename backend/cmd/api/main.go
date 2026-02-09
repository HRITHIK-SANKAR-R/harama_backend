package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    "harama/internal/api"
    "harama/internal/config"
    "harama/internal/repository/postgres"
)
func main() {
    // Load configuration
    cfg := config.Load()
    
    // Initialize database
    db, err := postgres.Connect(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Cleanup: Reset stuck jobs
    log.Println("Cleaning up stuck jobs...")
    res, err := db.NewUpdate().
        Table("submissions").
        Set("processing_status = ?", "failed").
        Where("processing_status = ?", "processing").
        Exec(context.Background())
    if err != nil {
        log.Printf("⚠️ Failed to cleanup stuck jobs: %v", err)
    } else {
        count, _ := res.RowsAffected()
        if count > 0 {
            log.Printf("🔄 Reset %d stuck submissions to 'failed'", count)
        }
    }
    
    // Initialize router
    router, err := api.NewRouter(cfg, db)
    if err != nil {
        log.Fatalf("Failed to initialize router: %v", err)
    }
    // Start server
    srv := &http.Server{
        Addr:    ":" + cfg.Port,
        Handler: router,
    }
    
    // Graceful shutdown
    go func() {
        sigint := make(chan os.Signal, 1)
        signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
        <-sigint
        
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        
        if err := srv.Shutdown(ctx); err != nil {
            log.Printf("Server shutdown error: %v", err)
        }
    }()
    
    log.Printf("Server started on :%s", cfg.Port)
    if err := srv.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatalf("Server error: %v", err)
    }
}