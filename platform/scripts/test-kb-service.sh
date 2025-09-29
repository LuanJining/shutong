#!/bin/bash

# KB Service 测试脚本
# 测试文档上传和预览功能

set -e

# 配置
KB_SERVICE_URL="http://localhost:8083"
TEST_FILE="/Users/gideonzy/Downloads/evaluation_report2.pdf"
TEST_CONTENT="这是一个测试文档内容，用于验证KB服务的上传和预览功能。\n\n包含中文内容测试。\n\n测试时间: $(date)"

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
check_service() {
    print_info "检查KB服务是否运行..."
    if curl -s "$KB_SERVICE_URL/api/v1/health" > /dev/null; then
        print_success "KB服务正在运行"
    else
        print_error "KB服务未运行，请先启动服务: ./bin/kb_serivce"
        exit 1
    fi
}

# 测试文档上传
test_upload() {
    print_info "测试文档上传..."
    
    # 上传文档
    response=$(curl -s -X POST \
        -F "file_name=test.pdf" \
        -F "file=@$TEST_FILE" \
        -F "space_id=1" \
        -F "visibility=private" \
        -F "urgency=normal" \
        -F "tags=测试,文档" \
        -F "summary=这是一个测试文档" \
        -F "created_by=1" \
        -F "department=技术部" \
        -F "need_approval=true" \
        "$KB_SERVICE_URL/api/v1/documents/upload")
    echo "响应: $response"
    # 检查响应
    if echo "$response" | grep -q "document_id"; then
        # 提取document_id
        DOCUMENT_ID=$(echo "$response" | grep -o '"document_id":[0-9]*' | grep -o '[0-9]*')
        print_success "文档上传成功，文档ID: $DOCUMENT_ID"
        echo "响应: $response"
    else
        print_error "文档上传失败"
        echo "响应: $response"
        exit 1
    fi
}

# 测试文档预览
test_preview() {
    if [ -z "$DOCUMENT_ID" ]; then
        print_error "没有可用的文档ID进行预览测试"
        return 1
    fi
    
    print_info "测试文档预览，文档ID: $DOCUMENT_ID"
    
    # 预览文档
    response=$(curl -s "$KB_SERVICE_URL/api/v1/documents/$DOCUMENT_ID/preview")
    
    # 检查响应
    if echo "$response" | grep -q "content"; then
        print_success "文档预览成功"
        echo "预览内容:"
        echo "$response" | grep -o '"content":"[^"]*"' | sed 's/"content":"//g' | sed 's/"$//g'
    else
        print_error "文档预览失败"
        echo "响应: $response"
    fi
}

# 测试文档下载
test_download() {
    if [ -z "$DOCUMENT_ID" ]; then
        print_error "没有可用的文档ID进行下载测试"
        return 1
    fi
    
    print_info "测试文档下载，文档ID: $DOCUMENT_ID"
    
    # 下载文档
    response=$(curl -s -o "downloaded-$TEST_FILE" "$KB_SERVICE_URL/api/v1/documents/$DOCUMENT_ID/download")
    
    # 检查文件是否下载成功
    if [ -f "downloaded-$TEST_FILE" ]; then
        print_success "文档下载成功"
        echo "下载的文件内容:"
        cat "downloaded-$TEST_FILE"
    else
        print_error "文档下载失败"
    fi
}

# 清理测试文件
cleanup() {
    print_info "清理测试文件..."
    print_success "清理完成"
}

# 主函数
main() {
    echo "=========================================="
    echo "KB Service 功能测试"
    echo "=========================================="
    
    # 检查服务
    check_service
    # 测试上传
    test_upload
    
    # 等待一下让处理完成
    print_info "等待文档处理完成..."
    
    # # 测试预览
    # test_preview
    
    # # 测试下载
    # test_download
    
    # # 清理
    # cleanup
    
    echo "=========================================="
    print_success "所有测试完成！"
    echo "=========================================="
}

# 捕获中断信号进行清理
trap cleanup EXIT

# 运行主函数
main "$@"
