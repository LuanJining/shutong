
GATEWAY_URL="http://192.168.0.56:8080/api/v1"
IAM_URL="http://localhost:8081/api/v1"

echo "=== Gateway服务测试 ==="

# 1. Gateway健康检查
echo "1. Gateway健康检查..."
curl -s "$GATEWAY_URL/../health" | jq .

# 2. 通过Gateway登录
echo -e "\n2. 通过Gateway登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/iam/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "login": "admin",
    "password": "admin123"
  }')

echo "$LOGIN_RESPONSE" | jq .

# 提取token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token')


#创建一个user
curl -s -X POST "$GATEWAY_URL/iam/users" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "Gateway测试用户",
    "phone": "13800138001",
    "email": "zygideon@gmail.com",
    "password": "admin123",
    "nickname": "Gateway测试用户",
    "department": "Gateway测试部门",
    "company": "Gateway测试公司"
  }' 

#分配一个角色
curl -s -X POST "$GATEWAY_URL/iam/users/5/roles" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role_id": 2
  }'
  

#删除一个user
curl -s -X DELETE "$GATEWAY_URL/iam/users/2" \
  -H "Authorization: Bearer $TOKEN"