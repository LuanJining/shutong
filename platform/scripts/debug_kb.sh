KB_SERVICE_URL="http://localhost:8083"



curl -s -X POST "$KB_SERVICE_URL/api/v1/documents/32/chat/stream" \
  -H "Content-Type: application/json" \
  -d '{
    "question": "业务流程是啥？解决啥问题",
    "document_ids": [],
    "limit": 3
  }'