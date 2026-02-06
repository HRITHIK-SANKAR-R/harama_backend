#!/bin/bash

BASE_URL="http://localhost:8080"
TENANT_ID="00000000-0000-0000-0000-000000000001"

echo "🧪 Testing HARaMA API"
echo ""

# Test 1: Health Check
echo "1️⃣ Health Check..."
curl -s $BASE_URL/health
echo -e "\n"

# Test 2: Create Exam
echo "2️⃣ Creating Exam..."
EXAM_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/exams \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d '{
    "title": "Physics Midterm",
    "subject": "Physics",
    "questions": [
      {
        "question_text": "Calculate force: F=ma, m=10kg, a=5m/s²",
        "points": 5,
        "answer_type": "short_answer"
      }
    ]
  }')
echo $EXAM_RESPONSE | jq '.' 2>/dev/null || echo $EXAM_RESPONSE
echo ""

# Extract exam ID
EXAM_ID=$(echo $EXAM_RESPONSE | jq -r '.id' 2>/dev/null)

if [ "$EXAM_ID" != "null" ] && [ -n "$EXAM_ID" ]; then
  echo "✅ Exam created: $EXAM_ID"
  
  # Test 3: Get Exam
  echo ""
  echo "3️⃣ Fetching Exam..."
  curl -s $BASE_URL/api/v1/exams/$EXAM_ID \
    -H "X-Tenant-ID: $TENANT_ID" | jq '.' 2>/dev/null
else
  echo "❌ Failed to create exam"
fi

echo ""
echo "✅ Basic API tests complete"
