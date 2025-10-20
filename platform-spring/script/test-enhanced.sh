#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"
TEST_FILE="./test.txt"

echo "========================================="
echo "完整业务流程测试（带异步处理等待）"
echo "========================================="

# 登录管理员账户
echo "=== 1. 管理员登录 ==="
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "login": "admin",
    "password": "admin123"
  }')
echo "$LOGIN_RESPONSE" | jq .
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token')
echo "✅ Token: $TOKEN"

# 创建知识空间
echo ""
echo "=== 2. 创建一级知识空间 ==="
SPACE_RESPONSE=$(curl -s -X POST "$BASE_URL/spaces" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "数通智汇知识库",
    "description": "数通智汇知识库",
    "type": "department"
  }')
echo "$SPACE_RESPONSE" | jq .
SPACE_ID=$(echo "$SPACE_RESPONSE" | jq -r '.data.id')
echo "✅ Space ID: $SPACE_ID"

# 创建二级知识空间
echo ""
echo "=== 3. 创建二级知识空间 ==="
SUBSPACE_RESPONSE=$(curl -s -X POST "$BASE_URL/spaces/sub-spaces" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"space_id\": $SPACE_ID,
    \"name\": \"可信数据空间相关文档\",
    \"description\": \"可信数据空间相关文档\"
  }")
echo "$SUBSPACE_RESPONSE" | jq .
SUBSPACE_ID=$(echo "$SUBSPACE_RESPONSE" | jq -r '.data.id')
echo "✅ SubSpace ID: $SUBSPACE_ID"

# 创建知识分类
echo ""
echo "=== 4. 创建知识分类 ==="
CLASS_RESPONSE=$(curl -s -X POST "$BASE_URL/spaces/classes" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"sub_space_id\": $SUBSPACE_ID,
    \"name\": \"测试分类\",
    \"description\": \"测试用知识分类\"
  }")
echo "$CLASS_RESPONSE" | jq .
CLASS_ID=$(echo "$CLASS_RESPONSE" | jq -r '.data.id')
echo "✅ Class ID: $CLASS_ID"

# 创建上传用户
echo ""
echo "=== 5. 创建上传用户 ==="
RANDOM_USERNAME="test_user_$(date +%s)"
RANDOM_PHONE="138$(date +%s | tail -c 9)"
UPLOAD_USER_RESPONSE=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"username\": \"$RANDOM_USERNAME\",
    \"phone\": \"$RANDOM_PHONE\",
    \"password\": \"password123\",
    \"nickname\": \"上传用户\",
    \"email\": \"${RANDOM_USERNAME}@test.com\"
  }")
echo "$UPLOAD_USER_RESPONSE" | jq .
UPLOAD_USER_ID=$(echo "$UPLOAD_USER_RESPONSE" | jq -r '.data.id')
echo "✅ Upload User ID: $UPLOAD_USER_ID"

# 添加用户到空间
echo ""
echo "=== 6. 添加用户到空间（editor + approver角色） ==="
ADD_USER_RESPONSE=$(curl -s -X POST "$BASE_URL/spaces/$SPACE_ID/members" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"user_id\": $UPLOAD_USER_ID,
    \"roles\": [\"editor\", \"approver\"]
  }")
echo "$ADD_USER_RESPONSE" | jq .

# 上传用户登录
echo ""
echo "=== 7. 上传用户登录 ==="
UPLOAD_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"login\": \"$RANDOM_USERNAME\",
    \"password\": \"password123\"
  }")
echo "$UPLOAD_LOGIN_RESPONSE" | jq .
UPLOAD_TOKEN=$(echo "$UPLOAD_LOGIN_RESPONSE" | jq -r '.data.access_token')
echo "✅ Upload User Token: $UPLOAD_TOKEN"

# 上传文档
echo ""
echo "=== 8. 上传文档（需要审批） ==="
UPLOAD_DOC_RESPONSE=$(curl -s -X POST "$BASE_URL/documents/upload" \
  -H "Authorization: Bearer $UPLOAD_TOKEN" \
  -F "file=@$TEST_FILE" \
  -F "file_name=test.txt" \
  -F "space_id=$SPACE_ID" \
  -F "sub_space_id=$SUBSPACE_ID" \
  -F "class_id=$CLASS_ID" \
  -F "tags=测试,文档" \
  -F "summary=招投标项目" \
  -F "department=测试部门" \
  -F "need_approval=true" \
  -F "version=v1.0.0" \
  -F "use_type=applicable")
echo "$UPLOAD_DOC_RESPONSE" | jq .
DOC_ID=$(echo "$UPLOAD_DOC_RESPONSE" | jq -r '.data.id')
echo "✅ Document ID: $DOC_ID"

# 等待异步处理完成
echo ""
echo "=== 9. 等待文档异步处理完成 ==="
MAX_WAIT=60  # 最多等待60秒
WAIT_COUNT=0
while [ $WAIT_COUNT -lt $MAX_WAIT ]; do
  DOC_INFO=$(curl -s -X GET "$BASE_URL/documents/$DOC_ID/info" \
    -H "Authorization: Bearer $UPLOAD_TOKEN")
  DOC_STATUS=$(echo "$DOC_INFO" | jq -r '.data.status')
  DOC_PROGRESS=$(echo "$DOC_INFO" | jq -r '.data.process_progress')
  
  echo "[$WAIT_COUNT秒] 文档状态: $DOC_STATUS, 进度: $DOC_PROGRESS%"
  
  # 处理完成
  if [ "$DOC_STATUS" = "pending_approval" ]; then
    echo "✅ 文档处理完成，等待审批"
    break
  fi
  
  # 处理失败
  if [ "$DOC_STATUS" = "process_failed" ] || [ "$DOC_STATUS" = "failed" ]; then
    ERROR=$(echo "$DOC_INFO" | jq -r '.data.parse_error')
    echo "❌ 文档处理失败: $ERROR"
    exit 1
  fi
  
  sleep 2
  WAIT_COUNT=$((WAIT_COUNT + 2))
done

if [ $WAIT_COUNT -ge $MAX_WAIT ]; then
  echo "❌ 等待超时，文档处理未完成"
  exit 1
fi

# 查看审批任务
echo ""
echo "=== 10. 查看待审批任务 ==="
TASKS_RESPONSE=$(curl -s -X GET "$BASE_URL/workflow/tasks?page=1&page_size=10" \
  -H "Authorization: Bearer $UPLOAD_TOKEN")
echo "$TASKS_RESPONSE" | jq .
TASK_ID=$(echo "$TASKS_RESPONSE" | jq -r '.data.items[0].id')
echo "✅ Task ID: $TASK_ID"

# 审批任务
echo ""
echo "=== 11. 审批任务（通过） ==="
APPROVE_RESPONSE=$(curl -s -X POST "$BASE_URL/workflow/tasks/approve" \
  -H "Authorization: Bearer $UPLOAD_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"task_id\": $TASK_ID,
    \"status\": \"approved\",
    \"comment\": \"审批通过\"
  }")
echo "$APPROVE_RESPONSE" | jq .

# 验证文档状态
echo ""
echo "=== 12. 验证文档状态（应为 pending_publish） ==="
DOC_INFO_RESPONSE=$(curl -s -X GET "$BASE_URL/documents/$DOC_ID/info" \
  -H "Authorization: Bearer $UPLOAD_TOKEN")
echo "$DOC_INFO_RESPONSE" | jq .
DOC_STATUS=$(echo "$DOC_INFO_RESPONSE" | jq -r '.data.status')

if [ "$DOC_STATUS" = "pending_publish" ]; then
  echo "✅ 测试成功：文档状态为 pending_publish（审批完成）"
else
  echo "❌ 测试失败：文档状态为 $DOC_STATUS，期望为 pending_publish"
  exit 1
fi

# 发布文档
echo ""
echo "=== 13. 发布文档 ==="
PUBLISH_RESPONSE=$(curl -s -X POST "$BASE_URL/documents/$DOC_ID/publish" \
  -H "Authorization: Bearer $UPLOAD_TOKEN")
echo "$PUBLISH_RESPONSE" | jq .

# 验证最终状态
echo ""
echo "=== 14. 验证最终文档状态（应为 published） ==="
FINAL_DOC_INFO=$(curl -s -X GET "$BASE_URL/documents/$DOC_ID/info" \
  -H "Authorization: Bearer $UPLOAD_TOKEN")
echo "$FINAL_DOC_INFO" | jq .
FINAL_STATUS=$(echo "$FINAL_DOC_INFO" | jq -r '.data.status')

if [ "$FINAL_STATUS" = "published" ]; then
  echo "✅✅✅ 完整流程测试成功：文档已发布 ✅✅✅"
else
  echo "❌ 测试失败：最终状态为 $FINAL_STATUS，期望为 published"
  exit 1
fi

echo ""
echo "========================================="
echo "所有测试通过！业务逻辑闭环验证成功！"
echo "========================================="

