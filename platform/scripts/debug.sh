GATEWAY_URL="http://localhost:8080/api/v1"
IAM_URL="http://localhost:8081/api/v1"


LOGIN_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/iam/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "login": "admin",
    "password": "admin123"
  }')

echo "$LOGIN_RESPONSE" | jq .



TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token')


curl -s -X GET "$GATEWAY_URL/iam/roles/1/permissions" \
          -H "Authorization: Bearer $TOKEN" 