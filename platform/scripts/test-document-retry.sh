#!/bin/bash

# 测试文档状态管理和重试功能

BASE_URL="http://localhost:8080/api/v1"
# BASE_URL="http://182.140.132.5:30368/api/v1"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== 文档状态管理和重试功能测试 ===${NC}"
echo ""

# 登录管理员账户
echo -e "${YELLOW}1. 登录管理员账户${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/iam/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "login": "admin",
    "password": "admin123"
  }')

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token')
if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo -e "${RED}❌ 登录失败${NC}"
    echo "$LOGIN_RESPONSE" | jq .
    exit 1
fi
echo -e "${GREEN}✅ 登录成功${NC}"
echo ""

# 上传一个测试文档（需要审批）
echo -e "${YELLOW}2. 上传测试文档（需要审批）${NC}"
UPLOAD_RESPONSE=$(curl -s -X POST "$BASE_URL/documents/upload" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@./test.txt" \
  -F "file_name=test-retry.txt" \
  -F "space_id=1" \
  -F "sub_space_id=1" \
  -F "class_id=1" \
  -F "tags=测试,重试" \
  -F "summary=测试文档状态管理" \
  -F "department=测试部门" \
  -F "need_approval=true" \
  -F "version=v1.0.0" \
  -F "use_type=applicable")

DOC_ID=$(echo "$UPLOAD_RESPONSE" | jq -r '.data.id')
if [ "$DOC_ID" = "null" ] || [ -z "$DOC_ID" ]; then
    echo -e "${RED}❌ 文档上传失败${NC}"
    echo "$UPLOAD_RESPONSE" | jq .
    exit 1
fi
echo -e "${GREEN}✅ 文档上传成功，ID: $DOC_ID${NC}"
echo ""

# 监控文档处理状态
echo -e "${YELLOW}3. 监控文档处理状态${NC}"
for i in {1..30}; do
    DOC_INFO=$(curl -s -X GET "$BASE_URL/documents/$DOC_ID/info" \
      -H "Authorization: Bearer $TOKEN")
    
    STATUS=$(echo "$DOC_INFO" | jq -r '.data.status')
    PROGRESS=$(echo "$DOC_INFO" | jq -r '.data.process_progress // 0')
    PARSE_ERROR=$(echo "$DOC_INFO" | jq -r '.data.parse_error // ""')
    RETRY_COUNT=$(echo "$DOC_INFO" | jq -r '.data.retry_count // 0')
    
    echo -e "  状态: ${BLUE}$STATUS${NC} | 进度: ${BLUE}$PROGRESS%${NC} | 重试次数: ${BLUE}$RETRY_COUNT${NC}"
    
    if [ "$PARSE_ERROR" != "" ] && [ "$PARSE_ERROR" != "null" ]; then
        echo -e "  错误: ${RED}$PARSE_ERROR${NC}"
    fi
    
    # 如果处理完成或失败，退出循环
    if [ "$STATUS" = "published" ] || [ "$STATUS" = "pending_publish" ] || [ "$STATUS" = "pending_approval" ] || [ "$STATUS" = "process_failed" ] || [ "$STATUS" = "failed" ]; then
        break
    fi
    
    sleep 2
done
echo ""

# 如果文档处理失败，测试重试功能
if [ "$STATUS" = "process_failed" ] || [ "$STATUS" = "failed" ]; then
    echo -e "${YELLOW}4. 测试重试功能${NC}"
    
    # 测试普通重试
    echo -e "  测试普通重试..."
    RETRY_RESPONSE=$(curl -s -X POST "$BASE_URL/documents/retry-process" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{
        \"document_id\": $DOC_ID,
        \"force_retry\": false
      }")
    
    echo "$RETRY_RESPONSE" | jq .
    
    RETRY_STATUS=$(echo "$RETRY_RESPONSE" | jq -r '.data.status')
    RETRY_MESSAGE=$(echo "$RETRY_RESPONSE" | jq -r '.data.message')
    NEW_RETRY_COUNT=$(echo "$RETRY_RESPONSE" | jq -r '.data.retry_count')
    
    echo -e "  重试结果: ${BLUE}$RETRY_STATUS${NC}"
    echo -e "  消息: ${BLUE}$RETRY_MESSAGE${NC}"
    echo -e "  新重试次数: ${BLUE}$NEW_RETRY_COUNT${NC}"
    echo ""
    
    # 再次监控处理状态
    echo -e "${YELLOW}5. 监控重试后的处理状态${NC}"
    for i in {1..20}; do
        DOC_INFO=$(curl -s -X GET "$BASE_URL/documents/$DOC_ID/info" \
          -H "Authorization: Bearer $TOKEN")
        
        STATUS=$(echo "$DOC_INFO" | jq -r '.data.status')
        PROGRESS=$(echo "$DOC_INFO" | jq -r '.data.process_progress // 0')
        RETRY_COUNT=$(echo "$DOC_INFO" | jq -r '.data.retry_count // 0')
        
        echo -e "  状态: ${BLUE}$STATUS${NC} | 进度: ${BLUE}$PROGRESS%${NC} | 重试次数: ${BLUE}$RETRY_COUNT${NC}"
        
        if [ "$STATUS" = "published" ] || [ "$STATUS" = "pending_publish" ] || [ "$STATUS" = "pending_approval" ] || [ "$STATUS" = "process_failed" ] || [ "$STATUS" = "failed" ]; then
            break
        fi
        
        sleep 2
    done
    echo ""
else
    echo -e "${GREEN}✅ 文档处理成功，无需测试重试功能${NC}"
    echo ""
fi

# 测试强制重试
echo -e "${YELLOW}6. 测试强制重试功能${NC}"
FORCE_RETRY_RESPONSE=$(curl -s -X POST "$BASE_URL/documents/retry-process" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"document_id\": $DOC_ID,
    \"force_retry\": true
  }")

echo "$FORCE_RETRY_RESPONSE" | jq .
echo ""

# 最终状态检查
echo -e "${YELLOW}7. 最终状态检查${NC}"
FINAL_INFO=$(curl -s -X GET "$BASE_URL/documents/$DOC_ID/info" \
  -H "Authorization: Bearer $TOKEN")

FINAL_STATUS=$(echo "$FINAL_INFO" | jq -r '.data.status')
FINAL_PROGRESS=$(echo "$FINAL_INFO" | jq -r '.data.process_progress // 0')
FINAL_RETRY_COUNT=$(echo "$FINAL_INFO" | jq -r '.data.retry_count // 0')
FINAL_VECTOR_COUNT=$(echo "$FINAL_INFO" | jq -r '.data.vector_count // 0')
FINAL_ERROR=$(echo "$FINAL_INFO" | jq -r '.data.parse_error // ""')

echo -e "  最终状态: ${BLUE}$FINAL_STATUS${NC}"
echo -e "  最终进度: ${BLUE}$FINAL_PROGRESS%${NC}"
echo -e "  总重试次数: ${BLUE}$FINAL_RETRY_COUNT${NC}"
echo -e "  向量数量: ${BLUE}$FINAL_VECTOR_COUNT${NC}"

if [ "$FINAL_ERROR" != "" ] && [ "$FINAL_ERROR" != "null" ]; then
    echo -e "  最终错误: ${RED}$FINAL_ERROR${NC}"
fi

echo ""
echo -e "${BLUE}=== 测试完成 ===${NC}"

# 总结
if [ "$FINAL_STATUS" = "published" ]; then
    echo -e "${GREEN}🎉 文档已发布！${NC}"
elif [ "$FINAL_STATUS" = "pending_publish" ]; then
    echo -e "${GREEN}✅ 文档处理成功，待发布！${NC}"
elif [ "$FINAL_STATUS" = "pending_approval" ]; then
    echo -e "${YELLOW}⏳ 文档处理完成，等待审批...${NC}"
elif [ "$FINAL_STATUS" = "processing" ] || [ "$FINAL_STATUS" = "vectorizing" ]; then
    echo -e "${YELLOW}⏳ 文档仍在处理中...${NC}"
else
    echo -e "${RED}❌ 文档处理失败${NC}"
fi

echo ""
echo -e "${BLUE}状态说明:${NC}"
echo -e "  ${BLUE}uploading${NC}       - 上传中"
echo -e "  ${BLUE}processing${NC}      - 解析中（OCR/文本提取）"
echo -e "  ${BLUE}vectorizing${NC}     - 向量化中"
echo -e "  ${YELLOW}pending_approval${NC} - 等待审批"
echo -e "  ${BLUE}pending_publish${NC}  - 待发布"
echo -e "  ${GREEN}published${NC}       - 已发布"
echo -e "  ${RED}process_failed${NC}   - 处理失败（可重试）"
echo -e "  ${RED}failed${NC}          - 失败"
