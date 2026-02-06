package unit_test

import (
	"context"
	"testing"
	"time"

	"harama/internal/domain"
	"harama/internal/repository/postgres"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func TestAuditRepo_Save(t *testing.T) {
	// 1. Setup mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	bunDB := bun.NewDB(db, pgdialect.New())
	repo := postgres.NewAuditRepo(bunDB)

	// 2. Define test data
	ctx := context.Background()
	auditLog := &domain.AuditLog{
		ID:         uuid.New(),
		EntityType: "exam",
		EntityID:   uuid.New(),
		EventType:  "created",
		ActorType:  "system",
		Changes:    map[string]interface{}{"title": "Math 101"},
		CreatedAt:  time.Now(),
	}

	// 3. Expectation: GetLastHash
	// Use generous regex for SELECT
	mock.ExpectQuery(`SELECT .*`).
		WillReturnRows(sqlmock.NewRows([]string{"hash"})) 

	// 4. Expectation: Save (Insert)
	// Match INSERT and potential RETURNING clause. Bun uses QueryRow for RETURNING.
	mock.ExpectQuery(`INSERT INTO "audit_log" .*`).
		WillReturnRows(sqlmock.NewRows([]string{"actor_id"}).AddRow(nil))

	// 5. Execute
	err = repo.Save(ctx, auditLog)

	// 6. Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, auditLog.Hash)
	assert.NotEqual(t, "initial_seed", auditLog.Hash) // Should be hashed

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAuditRepo_GetLastHash_Existing(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	bunDB := bun.NewDB(db, pgdialect.New())
	repo := postgres.NewAuditRepo(bunDB)

	expectedHash := "previous_hash_value"
	
	// Expectation
	rows := sqlmock.NewRows([]string{"id", "entity_type", "entity_id", "event_type", "changes", "hash", "created_at"}).
		AddRow(uuid.New(), "grade", uuid.New(), "updated", "{}", expectedHash, time.Now())
	
	mock.ExpectQuery(`SELECT .*`).
		WillReturnRows(rows)

	hash, err := repo.GetLastHash(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedHash, hash)
}
