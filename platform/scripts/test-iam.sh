#!/bin/bash

# IAM服务测试脚本

BASE_URL="http://localhost:8080/api/v1"

echo "=== IAM服务测试 ==="

# 1. 健康检查
echo "1. 健康检查..."
curl -s "$BASE_URL/../health" | jq .

echo -e "\n2. 用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "login": "admin",
    "password": "admin123"
  }')

echo "$LOGIN_RESPONSE" | jq .

# 提取token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token')

if [ "$TOKEN" != "null" ] && [ "$TOKEN" != "" ]; then
    echo -e "\n3. 获取用户列表..."
    curl -s -X GET "$BASE_URL/users" \
      -H "Authorization: Bearer $TOKEN" | jq .

    echo -e "\n4. 创建角色..."
    curl -s -X POST "$BASE_URL/roles" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "name": "admin",
        "display_name": "管理员",
        "description": "系统管理员角色"
      }' | jq .

    echo -e "\n5. 获取角色列表..."
    curl -s -X GET "$BASE_URL/roles" \
      -H "Authorization: Bearer $TOKEN" | jq .

    echo -e "\n6. 获取权限列表..."
    curl -s -X GET "$BASE_URL/permissions" \
      -H "Authorization: Bearer $TOKEN" | jq .

    echo -e "\n7. 创建用户（需要超级管理员权限）..."
    curl -s -X POST "$BASE_URL/users" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "username": "testuser",
        "phone": "13800138000",
        "email": "testuser@example.com",
        "password": "123456",
        "nickname": "测试用户",
        "department": "测试部门",
        "company": "示例公司"
      }' | jq .

    echo -e "\n8. 修改密码..."
    curl -s -X PATCH "$BASE_URL/auth/change-password" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "old_password": "admin123",
        "new_password": "newpassword123"
      }' | jq .
else
    echo "登录失败，无法继续测试"
fi

echo -e "\n=== 测试完成 ==="
