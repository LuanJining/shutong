# BASE_URL="http://localhost:8080/api/v1"
BASE_URL="http://182.140.132.5:30368/api/v1"
TEST_FILE="./test2.pdf"

# 登录管理员账户
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/iam/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "login": "admin",
    "password": "admin123"
  }')

echo "$LOGIN_RESPONSE" | jq .

# 提取token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token')
echo "Token: $TOKEN"

# 创建一个知识空间
echo "=== 创建一级知识空间 ==="
SPACE_RESPONSE=$(curl -s -X POST "$BASE_URL/iam/spaces" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "测试知识空间",
    "description": "测试用知识空间",
    "type": "department"
  }')
echo "$SPACE_RESPONSE" | jq .
SPACE_ID=$(echo "$SPACE_RESPONSE" | jq -r '.data.id')
echo "Space ID: $SPACE_ID"

# 创建一个二级知识空间
echo "=== 创建二级知识空间 ==="
SUBSPACE_RESPONSE=$(curl -s -X POST "$BASE_URL/iam/spaces/subspaces" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"space_id\": $SPACE_ID,
    \"name\": \"测试二级知识空间\",
    \"description\": \"测试用二级知识空间\"
  }")
echo "$SUBSPACE_RESPONSE" | jq .
SUBSPACE_ID=$(echo "$SUBSPACE_RESPONSE" | jq -r '.data.id')
echo "SubSpace ID: $SUBSPACE_ID"

# 创建一个知识分类
echo "=== 创建知识分类 ==="
CLASS_RESPONSE=$(curl -s -X POST "$BASE_URL/iam/spaces/classes" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"sub_space_id\": $SUBSPACE_ID,
    \"name\": \"测试分类\",
    \"description\": \"测试用知识分类\"
  }")
echo "$CLASS_RESPONSE" | jq .
CLASS_ID=$(echo "$CLASS_RESPONSE" | jq -r '.data.id')
echo "Class ID: $CLASS_ID"

# 创建一个用户用于文档上传
echo "=== 创建上传用户 ==="
RANDOM_USERNAME_FILE_UPLOAD="test_user_$(date +%s)"
RANDOM_PHONE_UPLOAD="138$(date +%s | tail -c 9)"
UPLOAD_USER_RESPONSE=$(curl -s -X POST "$BASE_URL/iam/users" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"username\": \"$RANDOM_USERNAME_FILE_UPLOAD\",
    \"phone\": \"$RANDOM_PHONE_UPLOAD\",
    \"password\": \"password123\",
    \"nickname\": \"上传用户\",
    \"email\": \"${RANDOM_USERNAME_FILE_UPLOAD}@test.com\"
  }")
echo "$UPLOAD_USER_RESPONSE" | jq .
UPLOAD_USER_ID=$(echo "$UPLOAD_USER_RESPONSE" | jq -r '.data.id')
echo "Upload User ID: $UPLOAD_USER_ID"

# 将上传用户添加到空间（editor角色）
echo "=== 添加上传用户到空间 ==="
ADD_UPLOAD_USER_RESPONSE=$(curl -s -X POST "$BASE_URL/iam/spaces/$SPACE_ID/members" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"user_id\": $UPLOAD_USER_ID,
    \"role\": \"editor\"
  }")
echo "$ADD_UPLOAD_USER_RESPONSE" | jq .

# 创建一个空间的审核员用于审核文档
echo "=== 创建审核员用户 ==="
RANDOM_USERNAME_FILE_AUDIT="test_audit_$(date +%s)"
RANDOM_PHONE_AUDIT="139$(date +%s | tail -c 9)"
AUDIT_USER_RESPONSE=$(curl -s -X POST "$BASE_URL/iam/users" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"username\": \"$RANDOM_USERNAME_FILE_AUDIT\",
    \"phone\": \"$RANDOM_PHONE_AUDIT\",
    \"password\": \"password123\",
    \"nickname\": \"审核员\",
    \"email\": \"${RANDOM_USERNAME_FILE_AUDIT}@test.com\"
  }")
echo "$AUDIT_USER_RESPONSE" | jq .
AUDIT_USER_ID=$(echo "$AUDIT_USER_RESPONSE" | jq -r '.data.id')
echo "Audit User ID: $AUDIT_USER_ID"

# 将审核员添加到空间（approver角色）
echo "=== 添加审核员到空间 ==="
ADD_AUDIT_USER_RESPONSE=$(curl -s -X POST "$BASE_URL/iam/spaces/$SPACE_ID/members" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"user_id\": $AUDIT_USER_ID,
    \"role\": \"approver\"
  }")
echo "$ADD_AUDIT_USER_RESPONSE" | jq .

# 用上传用户登录
echo "=== 上传用户登录 ==="
UPLOAD_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/iam/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"login\": \"$RANDOM_USERNAME_FILE_UPLOAD\",
    \"password\": \"password123\"
  }")
echo "$UPLOAD_LOGIN_RESPONSE" | jq .
UPLOAD_TOKEN=$(echo "$UPLOAD_LOGIN_RESPONSE" | jq -r '.data.access_token')
echo "Upload User Token: $UPLOAD_TOKEN"

# 上传一个文档到创建的知识空间
echo "=== 上传文档 ==="
UPLOAD_DOC_RESPONSE=$(curl -s -X POST "$BASE_URL/kb/upload" \
  -H "Authorization: Bearer $UPLOAD_TOKEN" \
  -F "file=@$TEST_FILE" \
  -F "file_name=test2.pdf" \
  -F "space_id=$SPACE_ID" \
  -F "sub_space_id=$SUBSPACE_ID" \
  -F "class_id=$CLASS_ID" \
  -F "tags=测试,文档" \
  -F "summary=这是一个测试文档" \
  -F "department=测试部门" \
  -F "need_approval=true" \
  -F "version=v1.0.0" \
  -F "use_type=applicable")
echo "$UPLOAD_DOC_RESPONSE" | jq .
DOC_ID=$(echo "$UPLOAD_DOC_RESPONSE" | jq -r '.data.id')
echo "Document ID: $DOC_ID"

# 用审核员登录
echo "=== 审核员登录 ==="
AUDIT_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/iam/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"login\": \"$RANDOM_USERNAME_FILE_AUDIT\",
    \"password\": \"password123\"
  }")
echo "$AUDIT_LOGIN_RESPONSE" | jq .
AUDIT_TOKEN=$(echo "$AUDIT_LOGIN_RESPONSE" | jq -r '.data.access_token')
echo "Audit User Token: $AUDIT_TOKEN"

# 用审核员账号列举审批任务
echo "=== 审核员查看待审批任务 ==="
TASKS_RESPONSE=$(curl -s -X GET "$BASE_URL/workflow/tasks" \
  -H "Authorization: Bearer $AUDIT_TOKEN")
echo "$TASKS_RESPONSE" | jq .
TASK_ID=$(echo "$TASKS_RESPONSE" | jq -r '.data.items[0].id')
echo "Task ID: $TASK_ID"

# 用审核员账号审批任务
echo "=== 审核员审批任务 ==="
APPROVE_RESPONSE=$(curl -s -X POST "$BASE_URL/workflow/tasks/approve" \
  -H "Authorization: Bearer $AUDIT_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"task_id\": $TASK_ID,
    \"status\": \"approved\",
    \"comment\": \"审批通过\"
  }")
echo "$APPROVE_RESPONSE" | jq .

# 上传文档的账号查看文档状态是否为待发布（审批完成）
echo "=== 查看文档状态 ==="
DOC_INFO_RESPONSE=$(curl -s -X GET "$BASE_URL/kb/$DOC_ID/info" \
  -H "Authorization: Bearer $UPLOAD_TOKEN")
echo "$DOC_INFO_RESPONSE" | jq .
DOC_STATUS=$(echo "$DOC_INFO_RESPONSE" | jq -r '.data.status')
echo "Document Status: $DOC_STATUS"

if [ "$DOC_STATUS" = "pending_publish" ]; then
  echo "✅ 测试成功：文档状态为待发布（审批完成）"
else
  echo "❌ 测试失败：文档状态为 $DOC_STATUS，期望为 pending_publish"
fi


# curl -X POST "$BASE_URL/kb/search" \
#   -H "Authorization: Bearer $UPLOAD_TOKEN" \
#   -H "Content-Type: application/json" \
#   -d '{
#     "query": "沈个好远",
#     "limit": 10
#   }' | jq .

curl -X POST "$BASE_URL/kb/$DOC_ID/publish" \
  -H "Authorization: Bearer $UPLOAD_TOKEN"

DOC_INFO_RESPONSE=$(curl -s -X GET "$BASE_URL/kb/$DOC_ID/info" \
  -H "Authorization: Bearer $UPLOAD_TOKEN")
echo "$DOC_INFO_RESPONSE" | jq .
DOC_STATUS=$(echo "$DOC_INFO_RESPONSE" | jq -r '.data.status')
echo "Document Status: $DOC_STATUS"