#!/bin/bash

echo "🔍 HARaMA Backend Verification"
echo "================================"
echo ""

# Check compilation
echo "1️⃣ Checking compilation..."
cd backend
if go build ./cmd/api && go build ./cmd/worker && go build ./cmd/migrate; then
    echo "   ✅ All binaries compile"
else
    echo "   ❌ Compilation failed"
    exit 1
fi
echo ""

# Run tests
echo "2️⃣ Running tests..."
if go test ./internal/... -cover > /tmp/test_results.txt 2>&1; then
    echo "   ✅ All tests passing"
    echo ""
    echo "   Test Coverage:"
    grep "coverage:" /tmp/test_results.txt | grep -v "0.0%" | sed 's/^/   /'
else
    echo "   ❌ Some tests failed"
    cat /tmp/test_results.txt
    exit 1
fi
echo ""

# Check critical files
echo "3️⃣ Checking critical files..."
FILES=(
    "cmd/api/main.go"
    "cmd/worker/main.go"
    "cmd/migrate/main.go"
    "internal/grading/engine.go"
    "internal/ai/gemini/client.go"
    "internal/worker/pool.go"
    "../docker-compose.yml"
    "../Makefile"
    ".env.example"
)

for file in "${FILES[@]}"; do
    if [ -f "$file" ]; then
        echo "   ✅ $file"
    else
        echo "   ❌ Missing: $file"
    fi
done
echo ""

# Summary
echo "================================"
echo "✅ Backend Verification Complete"
echo ""
echo "Status: READY FOR FRONTEND"
echo ""
echo "Next steps:"
echo "  1. ./quickstart.sh"
echo "  2. cd backend && go run ./cmd/api"
echo "  3. Start building frontend!"
