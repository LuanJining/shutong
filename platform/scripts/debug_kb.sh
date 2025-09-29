KB_SERVICE_URL="http://localhost:8083"



curl -s -X POST "$KB_SERVICE_URL/api/v1/documents/14/chat/stream" \
  -H "Content-Type: application/json" \
  -d '{
    "question": "如何健康减肥",
    "document_ids": [],
    "limit": 3
  }'