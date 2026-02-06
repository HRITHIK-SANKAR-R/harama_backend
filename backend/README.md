# HARaMA Backend

AI-powered exam grading system built with Go, PostgreSQL, and Gemini 3.

## Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Gemini API Key

### Setup

1. **Clone and configure**
```bash
cd backend
cp .env.example .env
# Edit .env and add your GEMINI_API_KEY
```

2. **Start services**
```bash
make docker-up
```

3. **Run migrations**
```bash
make migrate-up
```

4. **Access**
- API: http://localhost:8080
- MinIO Console: http://localhost:9001 (minioadmin/minioadmin)

## Development

### Local Development (without Docker)

1. **Start PostgreSQL and MinIO**
```bash
docker-compose up postgres minio -d
```

2. **Run API**
```bash
make run
```

3. **Run Worker (separate terminal)**
```bash
make worker
```

### Build

```bash
make build
# Binaries in ./bin/
```

### Testing

```bash
make test
```

## API Endpoints

### Exams
- `POST /api/v1/exams` - Create exam
- `GET /api/v1/exams` - List exams
- `GET /api/v1/exams/{id}` - Get exam details
- `POST /api/v1/exams/{id}/questions` - Add question
- `PUT /api/v1/questions/{id}/rubric` - Set rubric

### Submissions
- `POST /api/v1/exams/{id}/submissions` - Upload submission
- `GET /api/v1/submissions/{id}` - Get submission
- `POST /api/v1/submissions/{id}/trigger-grading` - Trigger grading

### Grading
- `GET /api/v1/submissions/{id}/grades` - Get grades
- `POST /api/v1/submissions/{sid}/questions/{qid}/override` - Override grade

### Analytics
- `GET /api/v1/analytics/grading-trends` - Get trends
- `POST /api/v1/exams/{id}/export` - Export grades (CSV)

### Health
- `GET /health` - Health check

## Architecture

```
cmd/
├── api/      # HTTP API server
├── worker/   # Background job processor
└── migrate/  # Database migration tool

internal/
├── api/          # HTTP handlers & routing
├── domain/       # Business entities
├── service/      # Business logic
├── grading/      # AI grading engine
├── ocr/          # OCR processing
├── repository/   # Data access
├── storage/      # File storage (MinIO)
└── worker/       # Job queue & workers
```

## Configuration

Environment variables (see `.env.example`):

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | API server port | `8080` |
| `DATABASE_URL` | PostgreSQL connection string | - |
| `GEMINI_API_KEY` | Gemini API key | - |
| `MINIO_ENDPOINT` | MinIO endpoint | `localhost:9000` |
| `WORKER_COUNT` | Number of worker goroutines | `10` |

## Database Migrations

Migrations are in `migrations/` directory.

```bash
# Apply migrations
make migrate-up

# Rollback
make migrate-down
```

## Docker Commands

```bash
# Start all services
make docker-up

# Stop services
make docker-down

# View logs
make docker-logs

# Rebuild and restart
docker-compose up --build -d
```

## Troubleshooting

### Database connection failed
```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# View logs
docker-compose logs postgres
```

### MinIO connection failed
```bash
# Check MinIO status
docker-compose ps minio

# Access MinIO console
open http://localhost:9001
```

### Worker not processing jobs
```bash
# Check worker logs
docker-compose logs worker

# Restart worker
docker-compose restart worker
```

## Production Deployment

### Fly.io

1. Install Fly CLI
```bash
curl -L https://fly.io/install.sh | sh
```

2. Deploy
```bash
fly launch
fly secrets set GEMINI_API_KEY=your_key
fly deploy
```

### Environment Variables (Production)
Set these secrets in your deployment platform:
- `GEMINI_API_KEY`
- `DATABASE_URL`
- `MINIO_ACCESS_KEY`
- `MINIO_SECRET_KEY`

## License

MIT
