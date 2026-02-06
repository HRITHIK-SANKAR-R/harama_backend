package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"harama/internal/api"
	"harama/internal/config"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func TestHealthCheck(t *testing.T) {
	// 1. Mock DB
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	bunDB := bun.NewDB(db, pgdialect.New())

	// 2. Mock Config
	cfg := &config.Config{
		Port:           "8080",
		DatabaseURL:    "postgres://user:pass@localhost:5432/db",
		GeminiAPIKey:   "dummy_key",
		MinioEndpoint:  "localhost:9000",
		MinioAccessKey: "minio",
		MinioSecretKey: "minio123",
		MinioBucket:    "uploads",
		MinioUseSSL:    false,
	}

	// 3. Initialize Router
	router, err := api.NewRouter(cfg, bunDB)
	assert.NoError(t, err)

	// 4. Create Request
	req, _ := http.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	// 5. Serve
	router.ServeHTTP(rr, req)

	// 6. Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}
