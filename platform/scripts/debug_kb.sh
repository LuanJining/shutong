KB_SERVICE_URL="http://localhost:8083"



curl -s -X POST "$KB_SERVICE_URL/api/v1/documents/32/chat/stream" \
  -H "Content-Type: application/json" \
  -d '{
    "question": "这方案水平一般吧",
    "document_ids": [],
    "limit": 3
  }'