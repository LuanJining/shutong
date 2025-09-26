#!/bin/bash

# Gateway服务测试脚本 - 通过Gateway访问IAM服务

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

if [ "$TOKEN" != "null" ] && [ "$TOKEN" != "" ]; then
    echo -e "\n3. 通过Gateway获取用户列表..."
    curl -s -X GET "$GATEWAY_URL/iam/users" \
      -H "Authorization: Bearer $TOKEN" | jq .

    echo -e "\n4. 通过Gateway获取角色列表..."
    curl -s -X GET "$GATEWAY_URL/iam/roles" \
      -H "Authorization: Bearer $TOKEN" | jq .

    echo -e "\n5. 通过Gateway获取权限列表..."
    curl -s -X GET "$GATEWAY_URL/iam/permissions" \
      -H "Authorization: Bearer $TOKEN" | jq .

    echo -e "\n6. 通过Gateway创建新空间..."
    SPACE_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/iam/spaces" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "name": "Gateway测试空间",
        "description": "通过Gateway创建的测试空间",
        "type": "project"
      }')
    
    echo "$SPACE_RESPONSE" | jq .
    
    # 提取空间ID
    SPACE_ID=$(echo "$SPACE_RESPONSE" | jq -r '.data.id')

    if [ "$SPACE_ID" != "null" ] && [ "$SPACE_ID" != "" ]; then
        echo -e "\n7. 通过Gateway获取空间列表..."
        curl -s -X GET "$GATEWAY_URL/iam/spaces" \
          -H "Authorization: Bearer $TOKEN" | jq .

        echo -e "\n8. 通过Gateway获取空间详情..."
        curl -s -X GET "$GATEWAY_URL/iam/spaces/$SPACE_ID" \
          -H "Authorization: Bearer $TOKEN" | jq .

        echo -e "\n9. 通过Gateway获取空间成员..."
        curl -s -X GET "$GATEWAY_URL/iam/spaces/$SPACE_ID/members" \
          -H "Authorization: Bearer $TOKEN" | jq .

        echo -e "\n10. 通过Gateway测试权限检查..."
        curl -s -X POST "$GATEWAY_URL/iam/permissions/check" \
          -H "Authorization: Bearer $TOKEN" \
          -H "Content-Type: application/json" \
          -d "{
            \"space_id\": $SPACE_ID,
            \"resource\": \"document\",
            \"action\": \"create\"
          }" | jq .

        echo -e "\n11. 通过Gateway测试跨空间权限检查..."
        curl -s -X POST "$GATEWAY_URL/iam/permissions/check" \
          -H "Authorization: Bearer $TOKEN" \
          -H "Content-Type: application/json" \
          -d "{
            \"space_id\": 999,
            \"resource\": \"user\",
            \"action\": \"manage\"
          }" | jq .

        echo -e "\n12. 通过Gateway更新空间信息..."
        curl -s -X PUT "$GATEWAY_URL/iam/spaces/$SPACE_ID" \
          -H "Authorization: Bearer $TOKEN" \
          -H "Content-Type: application/json" \
          -d '{
            "name": "Gateway更新后的测试空间",
            "description": "通过Gateway更新后的测试空间",
            "type": "department"
          }' | jq .

        echo -e "\n13. 通过Gateway删除测试空间..."
        curl -s -X DELETE "$GATEWAY_URL/iam/spaces/$SPACE_ID" \
          -H "Authorization: Bearer $TOKEN" | jq .
    fi

    echo -e "\n14. 通过Gateway测试角色权限分配..."
    # 先显示所有角色，然后获取第一个角色ID
    echo "14.1 获取所有角色列表..."
    ROLES_RESPONSE=$(curl -s -X GET "$GATEWAY_URL/iam/roles" \
      -H "Authorization: Bearer $TOKEN")
    echo "$ROLES_RESPONSE" | jq .
    
    # 获取第一个角色ID
    ROLE_ID=$(echo "$ROLES_RESPONSE" | jq -r '.data.roles[0].id')
    echo "14.2 选择的角色ID: $ROLE_ID"
    
    # 获取第一个权限ID
    PERMISSION_ID=$(curl -s -X GET "$GATEWAY_URL/iam/permissions" \
      -H "Authorization: Bearer $TOKEN" | jq -r '.data.permissions[0].id')
    echo "14.3 选择的权限ID: $PERMISSION_ID"

    if [ "$ROLE_ID" != "null" ] && [ "$ROLE_ID" != "" ] && [ "$PERMISSION_ID" != "null" ] && [ "$PERMISSION_ID" != "" ]; then
        echo -e "\n15. 通过Gateway获取角色权限列表..."
        curl -s -X GET "$GATEWAY_URL/iam/roles/$ROLE_ID/permissions" \
          -H "Authorization: Bearer $TOKEN" | jq .
    else
        echo "无法获取有效的角色ID或权限ID，跳过权限测试"
    fi

    echo -e "\n16. 测试Gateway代理功能..."
    echo "16.1 测试用户管理接口..."
    curl -s -X GET "$GATEWAY_URL/iam/users/1" \
      -H "Authorization: Bearer $TOKEN" | jq .

    echo -e "\n16.2 测试角色管理接口..."
    curl -s -X GET "$GATEWAY_URL/iam/roles/1" \
      -H "Authorization: Bearer $TOKEN" | jq .

    echo -e "\n16.3 测试权限管理接口..."
    curl -s -X GET "$GATEWAY_URL/iam/permissions/1" \
      -H "Authorization: Bearer $TOKEN" | jq .

    echo -e "\n17. 测试Gateway错误处理..."
    echo "17.1 测试无效路径..."
    curl -s -X GET "$GATEWAY_URL/iam/invalid/path" \
      -H "Authorization: Bearer $TOKEN" | jq .

    echo -e "\n17.2 测试无效token..."
    curl -s -X GET "$GATEWAY_URL/iam/users" \
      -H "Authorization: Bearer invalid_token" | jq .

    echo -e "\n18. 测试Gateway性能..."
    echo "18.1 并发请求测试..."
    for i in {1..5}; do
        curl -s -X GET "$GATEWAY_URL/iam/users" \
          -H "Authorization: Bearer $TOKEN" > /dev/null &
    done
    wait
    echo "并发请求完成"

else
    echo "登录失败，无法继续测试"
fi

echo -e "\n=== Gateway测试完成 ==="
echo "Gateway URL: $GATEWAY_URL"
echo "IAM URL: $IAM_URL"
