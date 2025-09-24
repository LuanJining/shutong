#!/bin/bash

# Workflow服务测试脚本
# 测试审批流程的创建、启动、审批等功能

set -e

# 配置
BASE_URL="http://localhost:8082/api/v1/workflow"
IAM_URL="http://localhost:8081/api/v1"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Workflow服务测试开始 ===${NC}"

# 1. 健康检查
echo -e "\n1. 检查Workflow服务健康状态..."
curl -s -X GET "$BASE_URL/../health" | jq .

# 2. 先通过IAM获取token（假设IAM服务已启动）
echo -e "\n2. 通过IAM获取认证token..."
LOGIN_RESPONSE=$(curl -s -X POST "$IAM_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "login": "admin",
    "password": "admin123"
  }')

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token')
USER_ID=$(echo "$LOGIN_RESPONSE" | jq -r '.data.user.id')
if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo -e "${RED}❌ 获取token失败，请确保IAM服务已启动${NC}"
    echo "$LOGIN_RESPONSE" | jq .
    exit 1
fi

if [ "$USER_ID" = "null" ] || [ -z "$USER_ID" ]; then
    echo -e "${RED}❌ 获取用户ID失败${NC}"
    echo "$LOGIN_RESPONSE" | jq .
    exit 1
fi

echo -e "${GREEN}✅ Token获取成功，用户ID: $USER_ID${NC}"

# 3. 创建审批流程
echo -e "\n3. 创建审批流程..."
CREATE_WORKFLOW_RESPONSE=$(curl -s -X POST "$BASE_URL/workflows" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-User-ID: $USER_ID" \
  -d '{
    "name": "文档发布审批流程",
    "description": "用于文档发布的审批流程",
    "space_id": 1,
    "priority": 1,
    "steps": [
      {
        "step_name": "内容审核",
        "step_order": 1,
        "approver_type": "space_admin",
        "approver_id": 0,
        "is_required": true,
        "timeout_hours": 24
      },
      {
        "step_name": "最终审批",
        "step_order": 2,
        "approver_type": "user",
        "approver_id": 1,
        "is_required": true,
        "timeout_hours": 48
      }
    ]
  }')

echo "$CREATE_WORKFLOW_RESPONSE" | jq .

WORKFLOW_ID=$(echo "$CREATE_WORKFLOW_RESPONSE" | jq -r '.data.id')
if [ "$WORKFLOW_ID" = "null" ] || [ -z "$WORKFLOW_ID" ]; then
    echo -e "${RED}❌ 创建流程失败${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 流程创建成功，ID: $WORKFLOW_ID${NC}"

# 4. 获取流程列表
echo -e "\n4. 获取流程列表..."
curl -s -X GET "$BASE_URL/workflows?space_id=1&page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN" | jq .

# 5. 获取流程详情
echo -e "\n5. 获取流程详情..."
curl -s -X GET "$BASE_URL/workflows/$WORKFLOW_ID" \
  -H "Authorization: Bearer $TOKEN" | jq .

# 6. 启动流程实例
echo -e "\n6. 启动流程实例..."
START_INSTANCE_RESPONSE=$(curl -s -X POST "$BASE_URL/instances" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-User-ID: $USER_ID" \
  -d '{
    "workflow_id": '$WORKFLOW_ID',
    "title": "测试文档发布申请",
    "description": "这是一份测试文档，需要经过审批流程",
    "resource_type": "document",
    "resource_id": 123,
    "space_id": 1,
    "priority": "normal"
  }')

echo "$START_INSTANCE_RESPONSE" | jq .

INSTANCE_ID=$(echo "$START_INSTANCE_RESPONSE" | jq -r '.data.id')
if [ "$INSTANCE_ID" = "null" ] || [ -z "$INSTANCE_ID" ]; then
    echo -e "${RED}❌ 启动流程实例失败${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 流程实例启动成功，ID: $INSTANCE_ID${NC}"

# 7. 获取我的待办任务
echo -e "\n7. 获取我的待办任务..."
TASKS_RESPONSE=$(curl -s -X GET "$BASE_URL/tasks?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-User-ID: $USER_ID")

echo "$TASKS_RESPONSE" | jq .

TASK_ID=$(echo "$TASKS_RESPONSE" | jq -r '.data.items[0].id')
if [ "$TASK_ID" = "null" ] || [ -z "$TASK_ID" ]; then
    echo -e "${YELLOW}⚠️  没有找到待办任务，可能需要不同的用户权限${NC}"
else
    echo -e "${GREEN}✅ 找到待办任务，ID: $TASK_ID${NC}"
    
    # 8. 审批通过任务
    echo -e "\n8. 审批通过任务..."
    curl -s -X POST "$BASE_URL/tasks/$TASK_ID/approve" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -H "X-User-ID: $USER_ID" \
      -d '{
        "comment": "内容审核通过，可以发布"
      }' | jq .
    
    echo -e "${GREEN}✅ 任务审批完成${NC}"
fi

# 9. 再次获取待办任务
echo -e "\n9. 再次获取待办任务..."
curl -s -X GET "$BASE_URL/tasks?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-User-ID: $USER_ID" | jq .

echo -e "\n${GREEN}=== Workflow服务测试完成 ===${NC}"
