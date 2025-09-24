#!/bin/bash

# IAM服务增强功能测试脚本

BASE_URL="http://localhost:8081/api/v1"

echo "=== IAM服务增强功能测试 ==="

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
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token')

if [ "$TOKEN" != "null" ] && [ "$TOKEN" != "" ]; then
    echo -e "\n3. 获取用户列表..."
    curl -s -X GET "$BASE_URL/users" \
      -H "Authorization: Bearer $TOKEN" | jq .

    echo -e "\n4. 获取角色列表..."
    curl -s -X GET "$BASE_URL/roles" \
      -H "Authorization: Bearer $TOKEN" | jq .

    echo -e "\n5. 获取权限列表..."
    curl -s -X GET "$BASE_URL/permissions" \
      -H "Authorization: Bearer $TOKEN" | jq .

    echo -e "\n6. 创建新空间..."
    SPACE_RESPONSE=$(curl -s -X POST "$BASE_URL/spaces" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "name": "测试空间",
        "description": "这是一个测试空间",
        "type": "project"
      }')
    
    echo "$SPACE_RESPONSE" | jq .
    
    # 提取空间ID
    SPACE_ID=$(echo "$SPACE_RESPONSE" | jq -r '.data.id')

    if [ "$SPACE_ID" != "null" ] && [ "$SPACE_ID" != "" ]; then
        echo -e "\n7. 获取空间列表..."
        curl -s -X GET "$BASE_URL/spaces" \
          -H "Authorization: Bearer $TOKEN" | jq .

        echo -e "\n8. 获取空间详情..."
        curl -s -X GET "$BASE_URL/spaces/$SPACE_ID" \
          -H "Authorization: Bearer $TOKEN" | jq .

        echo -e "\n9. 获取空间成员..."
        curl -s -X GET "$BASE_URL/spaces/$SPACE_ID/members" \
          -H "Authorization: Bearer $TOKEN" | jq .

        echo -e "\n10. 测试权限检查..."
        curl -s -X POST "$BASE_URL/permissions/check" \
          -H "Authorization: Bearer $TOKEN" \
          -H "Content-Type: application/json" \
          -d "{
            \"space_id\": $SPACE_ID,
            \"resource\": \"document\",
            \"action\": \"create\"
          }" | jq .

        echo -e "\n11. 测试跨空间权限检查..."
        curl -s -X POST "$BASE_URL/permissions/check" \
          -H "Authorization: Bearer $TOKEN" \
          -H "Content-Type: application/json" \
          -d "{
            \"space_id\": 999,
            \"resource\": \"user\",
            \"action\": \"manage\"
          }" | jq .

        echo -e "\n12. 更新空间信息..."
        curl -s -X PUT "$BASE_URL/spaces/$SPACE_ID" \
          -H "Authorization: Bearer $TOKEN" \
          -H "Content-Type: application/json" \
          -d '{
            "name": "更新后的测试空间",
            "description": "这是一个更新后的测试空间",
            "type": "department"
          }' | jq .

        echo -e "\n13. 删除测试空间..."
        curl -s -X DELETE "$BASE_URL/spaces/$SPACE_ID" \
          -H "Authorization: Bearer $TOKEN" | jq .
    fi

    echo -e "\n14. 测试角色权限分配..."
    # 先显示所有角色，然后获取第一个角色ID
    echo "14.1 获取所有角色列表..."
    ROLES_RESPONSE=$(curl -s -X GET "$BASE_URL/roles" \
      -H "Authorization: Bearer $TOKEN")
    echo "$ROLES_RESPONSE" | jq .
    
    # 获取第一个角色ID
    ROLE_ID=$(echo "$ROLES_RESPONSE" | jq -r '.data.roles[0].id')
    echo "14.2 选择的角色ID: $ROLE_ID"
    
    # 获取第一个权限ID
    PERMISSION_ID=$(curl -s -X GET "$BASE_URL/permissions" \
      -H "Authorization: Bearer $TOKEN" | jq -r '.data.permissions[0].id')
    echo "14.3 选择的权限ID: $PERMISSION_ID"

    if [ "$ROLE_ID" != "null" ] && [ "$ROLE_ID" != "" ] && [ "$PERMISSION_ID" != "null" ] && [ "$PERMISSION_ID" != "" ]; then
        echo -e "\n15. 获取角色权限列表..."
        curl -s -X GET "$BASE_URL/roles/$ROLE_ID/permissions" \
          -H "Authorization: Bearer $TOKEN" | jq .
    else
        echo "无法获取有效的角色ID或权限ID，跳过权限测试"
    fi

else
    echo "登录失败，无法继续测试"
fi

echo -e "\n=== 测试完成 ==="
