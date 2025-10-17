KB_SERVICE_URL="http://localhost:8083"



curl -s -X POST "$KB_SERVICE_URL/api/v1/documents/1/chat/stream" \
  -H "Content-Type: application/json" \
  -d '{
    "question": "怎么审核",
    "document_ids": [],
    "limit": 3
  }'