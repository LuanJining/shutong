
# BASE_URL="http://localhost:8080/api/v1"
BASE_URL="http://182.140.132.5:30368/api/v1"
TEST_FILE="./test.txt"

# 登录管理员账户
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
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
echo "Space ID: $SPACE_ID"


# 创建一个二级知识空间
echo "=== 创建二级知识空间 ==="
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
echo "SubSpace ID: $SUBSPACE_ID"

# 创建一个知识分类
echo "=== 创建知识分类 ==="
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
echo "Class ID: $CLASS_ID"


# 创建一个用户用于文档上传
echo "=== 创建上传用户 ==="
RANDOM_USERNAME_FILE_UPLOAD="test_user_$(date +%s)"
RANDOM_PHONE_UPLOAD="138$(date +%s | tail -c 9)"
UPLOAD_USER_RESPONSE=$(curl -s -X POST "$BASE_URL/users" \
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
ADD_UPLOAD_USER_RESPONSE=$(curl -s -X POST "$BASE_URL/spaces/$SPACE_ID/members" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"user_id\": $UPLOAD_USER_ID,
    \"roles\": [\"editor\", \"approver\"]
  }")
echo "$ADD_UPLOAD_USER_RESPONSE" | jq .

echo "=== 上传用户登录 ==="
UPLOAD_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"login\": \"$RANDOM_USERNAME_FILE_UPLOAD\",
    \"password\": \"password123\"
  }")
echo "$UPLOAD_LOGIN_RESPONSE" | jq .
UPLOAD_TOKEN=$(echo "$UPLOAD_LOGIN_RESPONSE" | jq -r '.data.access_token')
echo "Upload User Token: $UPLOAD_TOKEN"

echo "=== 上传文档 ==="
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
echo "Document ID: $DOC_ID"