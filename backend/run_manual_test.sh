#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Starting Manual Test Environment...${NC}"

# 0. Check for Postgres
if ! pg_isready -h localhost -p 5432 > /dev/null 2>&1; then
    echo -e "${RED}❌ Error: PostgreSQL is not running on localhost:5432.${NC}"
    echo "Please start your local PostgreSQL server manually and try again."
    echo "Example (Ubuntu/Debian): sudo service postgresql start"
    echo "Example (MacOS): brew services start postgresql"
    exit 1
fi
echo -e "${GREEN}✅ PostgreSQL is running.${NC}"

# 1. Start MinIO (Background)
# Create data dir if not exists
mkdir -p ~/minio-data
echo -e "${GREEN}Starting MinIO server...${NC}"
minio server ~/minio-data --address ":9000" --console-address ":9001" > minio.log 2>&1 &
MINIO_PID=$!
echo "MinIO running (PID: $MINIO_PID)"

# 2. Wait for MinIO to be ready
echo "Waiting for MinIO..."
sleep 3

# 3. Start Backend API (Background)
echo -e "${GREEN}Starting Backend API...${NC}"
go run cmd/api/main.go > api.log 2>&1 &
API_PID=$!
echo "API running (PID: $API_PID)"

# 4. Start Worker (Background)
echo -e "${GREEN}Starting Worker...${NC}"
go run cmd/worker/main.go > worker.log 2>&1 &
WORKER_PID=$!
echo "Worker running (PID: $WORKER_PID)"

# 5. Wait for API to be ready
echo "Waiting for API..."
sleep 5

# 6. Run the Test
echo -e "${BLUE}🧪 Running Test...${NC}"
echo "----------------------------------------"
go run scripts/e2e_test/main.go -file scripts/e2e_test/test1.jpeg
TEST_EXIT_CODE=$?
echo "----------------------------------------"

# 7. Cleanup
echo -e "${BLUE}🧹 Cleaning up processes...${NC}"
kill $MINIO_PID || true
kill $API_PID || true
kill $WORKER_PID || true

if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}✅ Test Completed Successfully!${NC}"
else
    echo -e "${RED}❌ Test Failed.${NC}"
fi

exit $TEST_EXIT_CODE
