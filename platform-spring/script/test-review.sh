#!/bin/bash

# 智能审查功能测试脚本

BASE_URL="http://localhost:8080/api/v1"
API_PREFIX="/review"
AUTH_URL="http://localhost:8080/api/v1"

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 默认管理员账号
ADMIN_LOGIN="admin"
ADMIN_PASSWORD="admin123"
TOKEN=""

# 登录获取token
login() {
    echo -e "${BLUE}=== 登录系统 ===${NC}"
    
    LOGIN_RESPONSE=$(curl -s -X POST "${AUTH_URL}/auth/login" \
        -H "Content-Type: application/json" \
        -d "{
            \"login\": \"${ADMIN_LOGIN}\",
            \"password\": \"${ADMIN_PASSWORD}\"
        }")
    
    echo "登录响应: $LOGIN_RESPONSE"
    
    # 提取token
    TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
    
    if [ -n "$TOKEN" ]; then
        echo -e "${GREEN}✓ 登录成功${NC}"
        echo "Token: ${TOKEN:0:50}..."
        echo $TOKEN > .token
        return 0
    else
        echo -e "${RED}✗ 登录失败，请检查账号密码或服务是否启动${NC}"
        return 1
    fi
}

# 创建测试文档
create_test_document() {
    echo ""
    echo "创建测试文档..."
    cat > test-document.txt << 'EOF'
第一章 总则

第一条 为了规范XX管理,根据《中华人民共和国XX法》第10条的规定,制定本办法.

第二条 本办法适用于全国范围内的XX活动。

第三条 XX工作应当遵循以下原则:
(一)依法依规;
(二)公开透明;
(三)便民高效.

第二章 管理职责

第四条 国务院XX部门负责全国XX工作的监督管理。

第五条 县级以上地方人民政府XX部门负责本行政区域内的XX管理工作。

第三章 附则

第六条 本办法自2024-01-01起施行。
EOF
    echo -e "${GREEN}✓ 测试文档创建成功${NC}"
}

# 测试1: 上传文档
test_upload() {
    echo ""
    echo -e "${YELLOW}=== 测试1: 上传文档 ===${NC}"
    
    if [ ! -f .token ]; then
        echo -e "${RED}✗ 未找到token，请先登录${NC}"
        return 1
    fi
    
    TOKEN=$(cat .token)
    
    UPLOAD_RESPONSE=$(curl -s -X POST "${BASE_URL}${API_PREFIX}/upload" \
        -H "Authorization: Bearer $TOKEN" \
        -F "file=@test-document.txt")
    
    echo "响应: $UPLOAD_RESPONSE"
    
    # 提取sessionId
    SESSION_ID=$(echo $UPLOAD_RESPONSE | grep -o '"data":"[^"]*"' | cut -d'"' -f4)
    
    if [ -n "$SESSION_ID" ]; then
        echo -e "${GREEN}✓ 上传成功，sessionId: $SESSION_ID${NC}"
        echo $SESSION_ID > .session_id
        return 0
    else
        echo -e "${RED}✗ 上传失败${NC}"
        return 1
    fi
}

# 测试2: 获取审查建议（SSE流式）
test_suggestions() {
    echo ""
    echo -e "${YELLOW}=== 测试2: 获取审查建议（SSE流式） ===${NC}"
    
    if [ ! -f .session_id ]; then
        echo -e "${RED}✗ 未找到sessionId，请先运行上传测试${NC}"
        return 1
    fi
    
    if [ ! -f .token ]; then
        echo -e "${RED}✗ 未找到token，请先登录${NC}"
        return 1
    fi
    
    SESSION_ID=$(cat .session_id)
    TOKEN=$(cat .token)
    
    echo "开始接收实时建议..."
    echo "SessionId: $SESSION_ID"
    echo ""
    
    # 使用curl -N保持连接接收SSE
    # 禁用AI建议以加快测试速度（AI建议可能超时60秒）
    curl -N "${BASE_URL}${API_PREFIX}/${SESSION_ID}/suggestions?fileName=test-document.txt&fileType=.txt&checkFormat=true&verifyReferences=true&suggestContent=false" \
        -H "Authorization: Bearer $TOKEN" \
        2>/dev/null | while IFS= read -r line; do
        
        if [[ $line == data:* ]]; then
            # 提取data部分
            data=${line#data: }
            
            # 解析并美化显示
            if [ -n "$data" ]; then
                echo -e "${GREEN}收到建议:${NC}"
                echo "$data" | python3 -m json.tool 2>/dev/null || echo "$data"
                echo ""
            fi
        elif [[ $line == *"DONE"* ]]; then
            echo -e "${GREEN}✓ 审查完成${NC}"
            break
        fi
    done
}

# 测试3: 获取审查摘要
test_summary() {
    echo ""
    echo -e "${YELLOW}=== 测试3: 获取审查摘要 ===${NC}"
    
    if [ ! -f .session_id ]; then
        echo -e "${RED}✗ 未找到sessionId，请先运行上传测试${NC}"
        return 1
    fi
    
    if [ ! -f .token ]; then
        echo -e "${RED}✗ 未找到token，请先登录${NC}"
        return 1
    fi
    
    SESSION_ID=$(cat .session_id)
    TOKEN=$(cat .token)
    
    SUMMARY_RESPONSE=$(curl -s "${BASE_URL}${API_PREFIX}/${SESSION_ID}/summary?fileName=test-document.txt&fileType=.txt" \
        -H "Authorization: Bearer $TOKEN")
    
    echo "审查摘要:"
    echo "$SUMMARY_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$SUMMARY_RESPONSE"
    
    # 提取统计信息
    TOTAL=$(echo $SUMMARY_RESPONSE | grep -o '"totalSuggestions":[0-9]*' | cut -d':' -f2)
    ERRORS=$(echo $SUMMARY_RESPONSE | grep -o '"errorCount":[0-9]*' | cut -d':' -f2)
    WARNINGS=$(echo $SUMMARY_RESPONSE | grep -o '"warningCount":[0-9]*' | cut -d':' -f2)
    INFOS=$(echo $SUMMARY_RESPONSE | grep -o '"infoCount":[0-9]*' | cut -d':' -f2)
    
    echo ""
    echo -e "${GREEN}统计信息:${NC}"
    echo "  总建议数: $TOTAL"
    echo "  错误: $ERRORS"
    echo "  警告: $WARNINGS"
    echo "  信息: $INFOS"
}

# 测试4: 完整审查（含AI建议，可能较慢）
test_full_review() {
    echo ""
    echo -e "${YELLOW}=== 测试4: 完整审查（含AI建议） ===${NC}"
    echo -e "${BLUE}注意: AI建议功能较慢，可能需要1-2分钟${NC}"
    
    if [ ! -f .session_id ]; then
        echo -e "${RED}✗ 未找到sessionId，请先运行上传测试${NC}"
        return 1
    fi
    
    if [ ! -f .token ]; then
        echo -e "${RED}✗ 未找到token，请先登录${NC}"
        return 1
    fi
    
    SESSION_ID=$(cat .session_id)
    TOKEN=$(cat .token)
    
    echo "执行完整审查（格式+引用+AI建议）..."
    echo ""
    
    # 设置超时时间
    timeout 120 curl -N "${BASE_URL}${API_PREFIX}/${SESSION_ID}/suggestions?fileName=test-document.txt&fileType=.txt&checkFormat=true&verifyReferences=true&suggestContent=true" \
        -H "Authorization: Bearer $TOKEN" \
        2>/dev/null | while IFS= read -r line; do
        
        if [[ $line == data:* ]]; then
            data=${line#data: }
            if [ -n "$data" ]; then
                echo -e "${GREEN}收到建议:${NC}"
                echo "$data" | python3 -m json.tool 2>/dev/null || echo "$data"
                echo ""
            fi
        elif [[ $line == *"DONE"* ]]; then
            echo -e "${GREEN}✓ 完整审查完成${NC}"
            break
        fi
    done
    
    if [ $? -eq 124 ]; then
        echo -e "${YELLOW}⚠ 审查超时（2分钟），这是正常的${NC}"
        echo -e "${BLUE}提示: 可以禁用AI建议加快速度（见 TIMEOUT_OPTIMIZATION.md）${NC}"
    fi
}

# 清理
cleanup() {
    echo ""
    echo "清理测试文件..."
    rm -f test-document.txt .session_id .token
    echo -e "${GREEN}✓ 清理完成${NC}"
}

# 主函数
main() {
    echo -e "${YELLOW}╔═══════════════════════════════════╗${NC}"
    echo -e "${YELLOW}║   智能审查功能测试脚本           ║${NC}"
    echo -e "${YELLOW}╚═══════════════════════════════════╝${NC}"
    echo ""
    echo -e "${BLUE}服务地址: ${BASE_URL}${NC}"
    echo -e "${BLUE}认证地址: ${AUTH_URL}${NC}"
    echo ""
    
    # 检查依赖
    if ! command -v python3 &> /dev/null; then
        echo -e "${YELLOW}警告: python3未安装，JSON美化显示将不可用${NC}"
    fi
    
    if ! command -v jq &> /dev/null; then
        echo -e "${YELLOW}警告: jq未安装，建议安装: brew install jq${NC}"
    fi
    
    # 登录获取token
    if ! login; then
        echo -e "${RED}登录失败，终止测试${NC}"
        exit 1
    fi
    
    # 创建测试文档
    create_test_document
    
    # 运行测试
    if test_upload; then
        test_suggestions  # 快速测试（无AI）
        test_summary      # 摘要统计
        
        # 询问是否运行完整测试
        echo ""
        echo -e "${BLUE}是否运行完整审查测试（含AI建议）？这可能需要1-2分钟 [y/N]${NC}"
        read -t 5 -r response || response="n"
        
        if [[ "$response" =~ ^[Yy]$ ]]; then
            test_full_review
        else
            echo -e "${YELLOW}跳过完整审查测试（已禁用AI建议以加快速度）${NC}"
            echo -e "${BLUE}如需测试AI建议，请手动运行或查看 TIMEOUT_OPTIMIZATION.md${NC}"
        fi
    else
        echo -e "${RED}上传失败，跳过后续测试${NC}"
    fi
    
    # 清理
    cleanup
    
    echo ""
    echo -e "${GREEN}╔═══════════════════════════════════╗${NC}"
    echo -e "${GREEN}║   测试完成！                      ║${NC}"
    echo -e "${GREEN}╚═══════════════════════════════════╝${NC}"
}

# 运行
main

