#!/bin/bash

# API 测试脚本
# 用于测试 Category 和 Tournament 接口

BASE_URL="${1:-http://localhost:8080}"

echo "=========================================="
echo "测试 API 接口"
echo "Base URL: $BASE_URL"
echo "=========================================="
echo ""

# 测试健康检查
echo "1. 测试健康检查"
echo "GET $BASE_URL/health"
curl -s "$BASE_URL/health" | jq .
echo ""
echo ""

# 测试获取所有分类
echo "2. 测试获取所有分类 (默认分页)"
echo "GET $BASE_URL/api/categories"
curl -s "$BASE_URL/api/categories" | jq .
echo ""
echo ""

# 测试获取足球分类
echo "3. 测试获取足球分类"
echo "GET $BASE_URL/api/categories?sport_ids=sr:sport:1"
curl -s "$BASE_URL/api/categories?sport_ids=sr:sport:1" | jq .
echo ""
echo ""

# 测试分类排序 (按比赛数量降序)
echo "4. 测试分类排序 (按比赛数量降序)"
echo "GET $BASE_URL/api/categories?sort=match_count_desc"
curl -s "$BASE_URL/api/categories?sort=match_count_desc" | jq .
echo ""
echo ""

# 测试分类分页
echo "5. 测试分类分页 (第1页, 每页2条)"
echo "GET $BASE_URL/api/categories?page=1&page_size=2"
curl -s "$BASE_URL/api/categories?page=1&page_size=2" | jq .
echo ""
echo ""

# 测试获取联赛 (需要先获取一个 category_id)
echo "6. 获取第一个分类的 ID"
CATEGORY_ID=$(curl -s "$BASE_URL/api/categories?page_size=1" | jq -r '.data[0].id')
echo "Category ID: $CATEGORY_ID"
echo ""

if [ "$CATEGORY_ID" != "null" ] && [ -n "$CATEGORY_ID" ]; then
    echo "7. 测试获取联赛列表"
    echo "GET $BASE_URL/api/tournaments?category_id=$CATEGORY_ID"
    curl -s "$BASE_URL/api/tournaments?category_id=$CATEGORY_ID" | jq .
    echo ""
    echo ""

    echo "8. 测试联赛排序 (按比赛数量降序)"
    echo "GET $BASE_URL/api/tournaments?category_id=$CATEGORY_ID&sort=match_count_desc"
    curl -s "$BASE_URL/api/tournaments?category_id=$CATEGORY_ID&sort=match_count_desc" | jq .
    echo ""
    echo ""

    echo "9. 测试联赛分页 (第1页, 每页1条)"
    echo "GET $BASE_URL/api/tournaments?category_id=$CATEGORY_ID&page=1&page_size=1"
    curl -s "$BASE_URL/api/tournaments?category_id=$CATEGORY_ID&page=1&page_size=1" | jq .
    echo ""
    echo ""
else
    echo "未找到分类数据，跳过联赛测试"
    echo ""
fi

# 测试错误情况
echo "10. 测试缺少必填参数 (category_id)"
echo "GET $BASE_URL/api/tournaments"
curl -s "$BASE_URL/api/tournaments" | jq .
echo ""
echo ""

echo "11. 测试无效的 category_id"
echo "GET $BASE_URL/api/tournaments?category_id=99999"
curl -s "$BASE_URL/api/tournaments?category_id=99999" | jq .
echo ""
echo ""

echo "=========================================="
echo "测试完成"
echo "=========================================="
