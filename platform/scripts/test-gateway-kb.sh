#!/bin/bash

# Gateway KB代理测试脚本

set -e

# 配置
GATEWAY_URL="http://192.168.0.56:8080"
KB_SERVICE_URL="http://localhost:8083"
TEST_FILE="/Users/gideonzy/Downloads/evaluation_report2.pdf"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印函数
print_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查服务是否运行
check_services() {
    print_info "检查服务是否运行..."
    
    # 检查Gateway
    if curl -s "$GATEWAY_URL/api/v1/health" > /dev/null; then
        print_success "Gateway服务正在运行"
    else
        print_error "Gateway服务未运行，请先启动: ./bin/gateway"
        exit 1
    fi
    
    # 检查KB Service
    if curl -s "$KB_SERVICE_URL/api/v1/health" > /dev/null; then
        print_success "KB Service正在运行"
    else
        print_error "KB Service未运行，请先启动: ./bin/kb_serivce"
        exit 1
    fi
}

# 测试文档上传代理
test_upload_proxy() {
    print_info "测试通过Gateway上传文档..."
    
    if [ ! -f "$TEST_FILE" ]; then
        print_error "测试文件不存在: $TEST_FILE"
        return 1
    fi
    
    # 通过Gateway上传文档
    print_info "上传文件大小: $(stat -f%z "$TEST_FILE" 2>/dev/null || stat -c%s "$TEST_FILE" 2>/dev/null || echo "unknown") 字节"
    response=$(curl -v -X POST \
        -H "Authorization: Bearer $TOKEN" \
        -F "file_name=test.pdf" \
        -F "file=@$TEST_FILE" \
        -F "space_id=14" \
        -F "visibility=private" \
        -F "urgency=normal" \
        -F "tags=测试,代理" \
        -F "summary=通过Gateway代理上传的测试文档" \
        -F "created_by=1" \
        -F "department=技术部" \
        -F "need_approval=true" \
        "$GATEWAY_URL/api/v1/kb/upload" 2>&1)
    
    # 检查响应
    if echo "$response" | grep -q "document_id"; then
        # 提取document_id
        DOCUMENT_ID=$(echo "$response" | grep -o '"document_id":[0-9]*' | grep -o '[0-9]*')
        print_success "通过Gateway上传文档成功，文档ID: $DOCUMENT_ID"
        echo "响应: $response"
    else
        print_error "通过Gateway上传文档失败"
        echo "响应: $response"
        return 1
    fi
}

# 测试文档预览代理
test_preview_proxy() {
    if [ -z "$DOCUMENT_ID" ]; then
        print_error "没有可用的文档ID进行预览测试"
        return 1
    fi
    
    print_info "测试通过Gateway预览文档，文档ID: $DOCUMENT_ID"
    
    # 通过Gateway预览文档
    response=$(curl -s "$GATEWAY_URL/api/v1/kb/$DOCUMENT_ID/preview")
    
    # 检查响应
    if [ ${#response} -gt 100 ]; then
        print_success "通过Gateway预览文档成功"
        echo "预览内容长度: ${#response} 字符"
    else
        print_error "通过Gateway预览文档失败"
        echo "响应: $response"
    fi
}

# 测试文档下载代理
test_download_proxy() {
    if [ -z "$DOCUMENT_ID" ]; then
        print_error "没有可用的文档ID进行下载测试"
        return 1
    fi
    
    print_info "测试通过Gateway下载文档，文档ID: $DOCUMENT_ID"
    
    # 通过Gateway下载文档
    response=$(curl -s -o "downloaded-via-gateway.pdf" "$GATEWAY_URL/api/v1/kb/$DOCUMENT_ID/download")
    
    # 检查文件是否下载成功
    if [ -f "downloaded-via-gateway.pdf" ]; then
        file_size=$(stat -f%z "downloaded-via-gateway.pdf" 2>/dev/null || stat -c%s "downloaded-via-gateway.pdf" 2>/dev/null || echo "unknown")
        print_success "通过Gateway下载文档成功，文件大小: $file_size 字节"
    else
        print_error "通过Gateway下载文档失败"
    fi
}

# 对比直接访问和代理访问
test_direct_vs_proxy() {
    if [ -z "$DOCUMENT_ID" ]; then
        print_error "没有可用的文档ID进行对比测试"
        return 1
    fi
    
    print_info "对比直接访问和代理访问..."
    
    # 直接访问KB Service
    direct_response=$(curl -s "$KB_SERVICE_URL/api/v1/documents/$DOCUMENT_ID/preview")
    direct_size=${#direct_response}
    
    # 通过Gateway代理访问
    proxy_response=$(curl -s "$GATEWAY_URL/api/v1/kb/$DOCUMENT_ID/preview")
    proxy_size=${#proxy_response}
    
    print_info "直接访问响应大小: $direct_size 字符"
    print_info "代理访问响应大小: $proxy_size 字符"
    
    if [ "$direct_size" -eq "$proxy_size" ]; then
        print_success "直接访问和代理访问响应一致"
    else
        print_error "直接访问和代理访问响应不一致"
    fi
}

# 清理测试文件
cleanup() {
    print_info "清理测试文件..."
    rm -f "downloaded-via-gateway.pdf"
    print_success "清理完成"
}

# 主函数
main() {
    echo "=========================================="
    echo "Gateway KB代理功能测试"
    echo "=========================================="
    
    # 检查服务
    check_services
    
    # 测试上传代理
    test_upload_proxy
    
    # 等待一下让处理完成
    print_info "等待文档处理完成..."
    sleep 3
    
    # 测试预览代理
    test_preview_proxy
    
    # 测试下载代理
    test_download_proxy
    
    # 对比测试
    test_direct_vs_proxy
    
    # 清理
    cleanup
    
    echo "=========================================="
    print_success "所有代理测试完成！"
    echo "=========================================="
}

# 捕获中断信号进行清理
trap cleanup EXIT

# 运行主函数
main "$@"
