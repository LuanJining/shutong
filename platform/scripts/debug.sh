
GATEWAY_URL="http://192.168.0.56:8080"
IAM_URL="http://localhost:8081"
WORKFLOW_URL="http://localhost:8082"

echo "=== Gateway Workflow代理测试 ==="

# 检查服务状态
echo "1. 检查服务状态..."
echo "Gateway: $GATEWAY_URL"
echo "IAM: $IAM_URL" 
echo "Workflow: $WORKFLOW_URL"

# 测试Gateway健康检查
echo -e "\n2. 测试Gateway健康检查..."
curl -s "$GATEWAY_URL/api/v1/health" | jq .

# 先登录获取token
echo -e "\n3. 登录获取token..."
LOGIN_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/api/v1/iam/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "login": "admin",
    "password": "admin123"
  }')

echo "登录响应:"
echo "$LOGIN_RESPONSE" | jq .

# 提取用户ID
USER_ID=$(echo "$LOGIN_RESPONSE" | jq -r '.data.user.id // .user.id // empty')
if [ -z "$USER_ID" ] || [ "$USER_ID" = "null" ]; then
  echo "❌ 无法获取用户ID"
  exit 1
fi

echo "获取到用户ID: $USER_ID"

# 测试获取任务列表
echo -e "\n4. 测试获取任务列表..."
curl -s -X GET "$GATEWAY_URL/api/v1/workflow/tasks?page=1&page_size=1" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-User-ID: $USER_ID" | jq .

echo -e "\n5. 测试审批任务..."
curl -s -X POST "$GATEWAY_URL/api/v1/workflow/tasks/11/approve" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-User-ID: $USER_ID" \
  -d '{
    "comment": "内容审核通过，可以发布"
  }' | jq .

echo -e "\n6. 测试获取审批流程实例详情..."
curl -s -X GET "$GATEWAY_URL/api/v1/workflow/instances/user" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-User-ID: $USER_ID" | jq .

# echo -e "\n7. 测试获取审批流程"
# curl -s -X GET "$GATEWAY_URL/api/v1/workflow/workflows/18" \
#   -H "Authorization: Bearer $TOKEN" \
#   -H "X-User-ID: $USER_ID" | jq .