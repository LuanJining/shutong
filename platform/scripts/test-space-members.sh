#!/bin/bash

# 空间成员管理测试脚本
# 测试添加、更新、删除空间成员功能

BASE_URL="http://localhost:8080/api/v1/iam"
# 如果直接测试 IAM 服务，使用: BASE_URL="http://localhost:8081/api/v1"

echo "========================================="
echo "空间成员管理测试"
echo "========================================="

# 1. 登录获取 token
echo -e "\n1. 登录获取 token..."
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "login": "admin",
    "password": "admin123"
  }')

echo "$LOGIN_RESPONSE" | jq .

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ 登录失败，无法获取 token"
  exit 1
fi

echo "✅ Token: $TOKEN"

# 2. 创建测试空间
echo -e "\n2. 创建测试空间..."
CREATE_SPACE_RESPONSE=$(curl -s -X POST "${BASE_URL}/spaces" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "测试空间",
    "description": "用于测试成员管理的空间",
    "type": "project"
  }')

echo "$CREATE_SPACE_RESPONSE" | jq .

SPACE_ID=$(echo "$CREATE_SPACE_RESPONSE" | jq -r '.data.id')

if [ "$SPACE_ID" == "null" ] || [ -z "$SPACE_ID" ]; then
  echo "❌ 创建空间失败"
  exit 1
fi

echo "✅ 空间ID: $SPACE_ID"

# 3. 创建测试用户
echo -e "\n3. 创建测试用户..."
CREATE_USER_RESPONSE=$(curl -s -X POST "${BASE_URL}/users" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "username": "testuser",
    "phone": "13900000001",
    "email": "test@example.com",
    "password": "test123",
    "nickname": "测试用户",
    "department": "测试部门"
  }')

echo "$CREATE_USER_RESPONSE" | jq .

USER_ID=$(echo "$CREATE_USER_RESPONSE" | jq -r '.data.id')

if [ "$USER_ID" == "null" ] || [ -z "$USER_ID" ]; then
  echo "⚠️  用户可能已存在，尝试获取现有用户..."
  # 如果用户已存在，可以从 GetUsers 获取
  USER_ID=2  # 假设是第二个用户
fi

echo "✅ 用户ID: $USER_ID"

# 4. 添加空间成员（编辑者角色）
echo -e "\n4. 添加空间成员（编辑者角色）..."
ADD_MEMBER_RESPONSE=$(curl -s -X POST "${BASE_URL}/spaces/${SPACE_ID}/members" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"user_id\": ${USER_ID},
    \"role\": \"editor\"
  }")

echo "$ADD_MEMBER_RESPONSE" | jq .

# 5. 获取空间成员列表
echo -e "\n5. 获取空间成员列表..."
GET_MEMBERS_RESPONSE=$(curl -s -X GET "${BASE_URL}/spaces/${SPACE_ID}/members" \
  -H "Authorization: Bearer $TOKEN")

echo "$GET_MEMBERS_RESPONSE" | jq .

# 6. 更新成员角色（升级为管理员）
echo -e "\n6. 更新成员角色（升级为管理员）..."
UPDATE_ROLE_RESPONSE=$(curl -s -X PUT "${BASE_URL}/spaces/${SPACE_ID}/members/${USER_ID}" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "role": "admin"
  }')

echo "$UPDATE_ROLE_RESPONSE" | jq .

# 7. 按角色查询成员
echo -e "\n7. 查询所有管理员..."
GET_ADMIN_RESPONSE=$(curl -s -X GET "${BASE_URL}/spaces/${SPACE_ID}/members/role/admin" \
  -H "Authorization: Bearer $TOKEN")

echo "$GET_ADMIN_RESPONSE" | jq .

# 8. 移除空间成员
echo -e "\n8. 移除空间成员..."
REMOVE_MEMBER_RESPONSE=$(curl -s -X DELETE "${BASE_URL}/spaces/${SPACE_ID}/members/${USER_ID}" \
  -H "Authorization: Bearer $TOKEN")

echo "$REMOVE_MEMBER_RESPONSE" | jq .

# 9. 验证成员已移除
echo -e "\n9. 验证成员已移除..."
VERIFY_RESPONSE=$(curl -s -X GET "${BASE_URL}/spaces/${SPACE_ID}/members" \
  -H "Authorization: Bearer $TOKEN")

echo "$VERIFY_RESPONSE" | jq .

echo -e "\n========================================="
echo "测试完成！"
echo "========================================="

