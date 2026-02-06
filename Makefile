.PHONY: help build run test migrate-up migrate-down docker-up docker-down clean

help:
	@echo "Available commands:"
	@echo "  make build        - Build API and worker binaries"
	@echo "  make run          - Run API server locally"
	@echo "  make worker       - Run worker locally"
	@echo "  make test         - Run tests"
	@echo "  make migrate-up   - Run database migrations"
	@echo "  make migrate-down - Rollback database migrations"
	@echo "  make docker-up    - Start all services with Docker Compose"
	@echo "  make docker-down  - Stop all Docker services"
	@echo "  make clean        - Clean build artifacts"

build:
	cd backend && go build -o ../bin/api ./cmd/api
	cd backend && go build -o ../bin/worker ./cmd/worker
	cd backend && go build -o ../bin/migrate ./cmd/migrate

run:
	cd backend && go run ./cmd/api

worker:
	cd backend && go run ./cmd/worker

test:
	cd backend && go test ./... -v

migrate-up:
	cd backend && go run ./cmd/migrate -direction=up

migrate-down:
	cd backend && go run ./cmd/migrate -direction=down

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

clean:
	rm -rf bin/
	cd backend && go clean

install-deps:
	cd backend && go mod download
	cd backend && go mod tidy
