#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"

echo "=== 测试首页文档接口 ==="
HOMEPAGE_RESPONSE=$(curl -s -X GET "$BASE_URL/documents/homepage")
echo "$HOMEPAGE_RESPONSE" | jq .

# 统计返回的数据
SPACE_COUNT=$(echo "$HOMEPAGE_RESPONSE" | jq '.data.spaces | length')
echo ""
echo "=== 统计信息 ==="
echo "知识库数量: $SPACE_COUNT"

# 遍历每个知识库，统计二级知识库和文档数量
for i in $(seq 0 $((SPACE_COUNT-1))); do
    SPACE_NAME=$(echo "$HOMEPAGE_RESPONSE" | jq -r ".data.spaces[$i].name")
    SUBSPACE_COUNT=$(echo "$HOMEPAGE_RESPONSE" | jq ".data.spaces[$i].sub_spaces | length")
    echo ""
    echo "知识库 [$SPACE_NAME] 包含 $SUBSPACE_COUNT 个二级知识库"
    
    for j in $(seq 0 $((SUBSPACE_COUNT-1))); do
        SUBSPACE_NAME=$(echo "$HOMEPAGE_RESPONSE" | jq -r ".data.spaces[$i].sub_spaces[$j].name")
        DOC_COUNT=$(echo "$HOMEPAGE_RESPONSE" | jq ".data.spaces[$i].sub_spaces[$j].documents | length")
        echo "  - 二级知识库 [$SUBSPACE_NAME] 包含 $DOC_COUNT 个文档"
    done
done

echo ""
echo "✅ 测试完成"

