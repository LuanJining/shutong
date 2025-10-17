
BASE_URL="http://localhost:8080/api/v1"
# BASE_URL="http://182.140.132.5:30368/api/v1"
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