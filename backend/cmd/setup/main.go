package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"harama/internal/config"
	"harama/internal/domain"
)

func main() {
	cfg := config.Load()
	
	// Open connection
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DatabaseURL)))
	db := bun.NewDB(sqldb, pgdialect.New())
	defer db.Close()

	ctx := context.Background()

	// 1. Run Migrations
	fmt.Println("Running migrations...")
	if err := runMigrations(ctx, db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("Migrations complete.")

	// 2. Seed Data
	fmt.Println("Seeding demo data...")
	if err := seedDemoData(ctx, db); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}
	fmt.Println("Seeding complete.")
}

func runMigrations(ctx context.Context, db *bun.DB) error {
	// Read migration files
	files, err := os.ReadDir("migrations")
	if err != nil {
		return err
	}

	var migrationFiles []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".up.sql") {
			migrationFiles = append(migrationFiles, f.Name())
		}
	}
	sort.Strings(migrationFiles)

	for _, file := range migrationFiles {
		fmt.Printf("Applying %s...\n", file)
		content, err := os.ReadFile(filepath.Join("migrations", file))
		if err != nil {
			return err
		}

		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			// Ignore "already exists" errors for idempotency
			if !strings.Contains(err.Error(), "already exists") {
				return fmt.Errorf("failed to apply %s: %w", file, err)
			}
			fmt.Printf("Skipping %s (likely already applied)\n", file)
		}
	}
	return nil
}

func seedDemoData(ctx context.Context, db *bun.DB) error {
	tenantID := uuid.New()
	
	// 0. Create Tenant
	type Tenant struct {
		bun.BaseModel `bun:"table:tenants"`
		ID            uuid.UUID `bun:"id,pk"`
		Name          string    `bun:"name"`
		CreatedAt     time.Time `bun:"created_at"`
	}
	
	tenant := &Tenant{
		ID:        tenantID,
		Name:      "Demo School District",
		CreatedAt: time.Now(),
	}
	if _, err := db.NewInsert().Model(tenant).Exec(ctx); err != nil {
		return fmt.Errorf("failed to create tenant: %w", err)
	}

	// 1. Create Exam
	exam := &domain.Exam{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Title:     "Physics 101: Midterm",
		Subject:   "physics",
		CreatedAt: time.Now(),
	}
	if _, err := db.NewInsert().Model(exam).Exec(ctx); err != nil {
		return err
	}

	// 2. Create Question (Ambiguous Logic)
	questionID := uuid.New()
	question := &domain.Question{
		ID:           questionID,
		ExamID:       exam.ID,
		QuestionText: "Explain why a satellite does not fall to the Earth while in orbit.",
		Points:       10,
		AnswerType:   domain.AnswerTypeEssay,
	}
	if _, err := db.NewInsert().Model(question).Exec(ctx); err != nil {
		return err
	}

	// 3. Create Rubric
	rubric := &domain.Rubric{
		ID:         uuid.New(),
		QuestionID: questionID,
		FullCreditCriteria: []domain.Criterion{
			{ID: "c1", Description: "Mentions gravitational force acts as centripetal force", Points: 5},
			{ID: "c2", Description: "Mentions high tangential velocity", Points: 3},
			{ID: "c3", Description: "Explains 'falling around the Earth'", Points: 2},
		},
		PartialCreditRules: []domain.PartialCreditRule{},
		CommonMistakes: []domain.CommonMistake{
			{ID: "m1", Description: "Says there is no gravity in space", Penalty: 5},
		},
	}
	if _, err := db.NewInsert().Model(rubric).Exec(ctx); err != nil {
		return err
	}

	// 4. Create Submission (The "High Variance" Student)
	subID := uuid.New()
	submission := &domain.Submission{
		ID:               subID,
		ExamID:           exam.ID,
		StudentID:        "student_123",
		UploadedAt:       time.Now(),
		ProcessingStatus: domain.StatusPending,
		OCRResults:       []domain.OCRResult{},
		TenantID:         tenantID,
	}
	if _, err := db.NewInsert().Model(submission).Exec(ctx); err != nil {
		return err
	}

	// 5. Create Answer (Correct logic, terrible formatting)
	// This should trigger:
	// - Rubric Enforcer: Low score (missing keywords, messy)
	// - Reasoning Validator: High score (concept is right)
	answer := &domain.AnswerSegment{
		SubmissionID: subID,
		QuestionID:   questionID,
		Text:         "its falling but moving sideways fast enough that it misses the ground. like gravity pulls it down but it goes forward so it curves matching the earth.",
		PageIndices:  []int{1},
	}
	
	// We need to manually insert answer into submissions table jsonb column or a separate table depending on schema.
	// Based on domain model `Submission` struct has `Answers []AnswerSegment`.
	// In Postgres repo, we usually store this as a JSONB column on submission or separate table.
	// Looking at `001_initial_schema.up.sql` (inferred), let's assume `answers` is a JSONB column on `submissions`.
	
	submission.Answers = []domain.AnswerSegment{*answer}
	_, err := db.NewUpdate().Model(submission).Column("answers").WherePK().Exec(ctx)

	fmt.Printf("Created Submission ID: %s\n", subID)
	return err
}
