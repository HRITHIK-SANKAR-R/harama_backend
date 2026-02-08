#!/bin/bash
set -e

echo "🚀 HARaMA Backend Quick Start"
echo ""

# Check if .env exists
if [ ! -f backend/.env ]; then
    echo "📝 Creating .env from template..."
    cp backend/.env.example backend/.env
    echo "⚠️  IMPORTANT: Edit backend/.env and add your GEMINI_API_KEY"
    echo ""
fi

# Start services
echo "🐳 Starting PostgreSQL and MinIO..."
sudo docker compose up -d postgres minio

# Wait for services
echo "⏳ Waiting for services to be ready..."
sleep 5

# Run migrations
echo "📊 Running database migrations..."
cd backend && go run ./cmd/migrate -direction=up || echo "⚠️  Migrations may have already run"

echo ""
echo "✅ Backend is ready!"
echo ""
echo "Next steps:"
echo "  1. Edit backend/.env and add GEMINI_API_KEY"
echo "  2. Run: cd backend && go run ./cmd/api"
echo "  3. API will be at http://localhost:8080"
echo ""
