#!/bin/bash

# 测试Gateway的Workflow代理功能

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

# 提取token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token // .token // empty')
if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "❌ 无法获取token，请检查IAM服务是否正常运行"
  exit 1
fi

echo "获取到token: $TOKEN"

# 提取用户ID
USER_ID=$(echo "$LOGIN_RESPONSE" | jq -r '.data.user.id // .user.id // empty')
if [ -z "$USER_ID" ] || [ "$USER_ID" = "null" ]; then
  echo "❌ 无法获取用户ID"
  exit 1
fi

echo "获取到用户ID: $USER_ID"

# 测试workflow代理 - 获取工作流列表
echo -e "\n4. 测试获取工作流列表..."
curl -s -X GET "$GATEWAY_URL/api/v1/workflow" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-User-ID: $USER_ID" | jq .

# 测试创建工作流
echo -e "\n5. 测试创建工作流..."
CREATE_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/api/v1/workflow" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-User-ID: $USER_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试审批流程",
    "description": "通过Gateway代理创建的测试流程",
    "space_id": 1,
    "steps": [
      {
        "name": "初审",
        "description": "内容初审",
        "approver_type": "role",
        "approver_value": "content_reviewer",
        "approval_strategy": "any",
        "step_order": 1,
        "timeout_hours": 24
      },
      {
        "name": "终审",
        "description": "内容终审",
        "approver_type": "role", 
        "approver_value": "space_admin",
        "approval_strategy": "any",
        "step_order": 2,
        "timeout_hours": 48
      }
    ]
  }')

echo "创建工作流响应:"
echo "$CREATE_RESPONSE" | jq .

# 提取工作流ID
WORKFLOW_ID=$(echo "$CREATE_RESPONSE" | jq -r '.data.id // .id // empty')
if [ -z "$WORKFLOW_ID" ] || [ "$WORKFLOW_ID" = "null" ]; then
  echo "❌ 无法获取工作流ID"
  exit 1
fi

echo "获取到工作流ID: $WORKFLOW_ID"

# 测试获取工作流详情
echo -e "\n6. 测试获取工作流详情..."
curl -s -X GET "$GATEWAY_URL/api/v1/workflow/$WORKFLOW_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-User-ID: $USER_ID" | jq .

# 测试启动工作流实例
echo -e "\n7. 测试启动工作流实例..."
INSTANCE_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/api/v1/workflow/$WORKFLOW_ID/instances" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-User-ID: $USER_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "测试文档审批",
    "description": "通过Gateway代理启动的测试实例",
    "initiator_id": '$USER_ID',
    "space_id": 1
  }')

echo "启动工作流实例响应:"
echo "$INSTANCE_RESPONSE" | jq .

# 提取实例ID
INSTANCE_ID=$(echo "$INSTANCE_RESPONSE" | jq -r '.data.id // .id // empty')
if [ -z "$INSTANCE_ID" ] || [ "$INSTANCE_ID" = "null" ]; then
  echo "❌ 无法获取实例ID"
  exit 1
fi

echo "获取到实例ID: $INSTANCE_ID"

# 测试获取实例列表
echo -e "\n8. 测试获取实例列表..."
curl -s -X GET "$GATEWAY_URL/api/v1/workflow/instances" \
  -H "Authorization: Bearer $TOKEN"  | jq .

# 测试获取实例详情
echo -e "\n9. 测试获取实例详情..."
curl -s -X GET "$GATEWAY_URL/api/v1/workflow/instances/$INSTANCE_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-User-ID: $USER_ID" | jq .

# 测试获取任务列表
echo -e "\n10. 测试获取任务列表..."
curl -s -X GET "$GATEWAY_URL/api/v1/workflow/tasks" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-User-ID: $USER_ID" | jq .

# 测试获取实例状态
echo -e "\n11. 测试获取实例状态..."
curl -s -X GET "$GATEWAY_URL/api/v1/workflow/instances/$INSTANCE_ID/status" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-User-ID: $USER_ID" | jq .

echo -e "\n=== Gateway Workflow代理测试完成 ==="
