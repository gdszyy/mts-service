#!/bin/bash

# 测试 /api/events 接口的 market_types 过滤功能

BASE_URL="${1:-http://localhost:8080}"

echo "=========================================="
echo "测试 /api/events 接口的 market_types 过滤"
echo "Base URL: $BASE_URL"
echo "=========================================="
echo ""

# 测试 1: 不带 market_types 参数（不返回 markets）
echo "1. 测试不带 market_types 参数（不返回 markets）"
echo "GET $BASE_URL/api/events?limit=2"
curl -s "$BASE_URL/api/events?limit=2" | jq .
echo ""
echo ""

# 测试 2: 带单个 market_type
echo "2. 测试带单个 market_type (1x2)"
echo "GET $BASE_URL/api/events?limit=2&market_types=1x2"
curl -s "$BASE_URL/api/events?limit=2&market_types=1x2" | jq .
echo ""
echo ""

# 测试 3: 带多个 market_types
echo "3. 测试带多个 market_types (1x2,handicap,totals)"
echo "GET $BASE_URL/api/events?limit=2&market_types=1x2,handicap,totals"
curl -s "$BASE_URL/api/events?limit=2&market_types=1x2,handicap,totals" | jq .
echo ""
echo ""

# 测试 4: 带空格的 market_types（测试 trim 功能）
echo "4. 测试带空格的 market_types (1x2, handicap, totals)"
echo "GET $BASE_URL/api/events?limit=2&market_types=1x2,%20handicap,%20totals"
curl -s "$BASE_URL/api/events?limit=2&market_types=1x2,%20handicap,%20totals" | jq .
echo ""
echo ""

# 测试 5: 结合其他参数
echo "5. 测试结合 status 参数"
echo "GET $BASE_URL/api/events?status=scheduled&limit=2&market_types=1x2"
curl -s "$BASE_URL/api/events?status=scheduled&limit=2&market_types=1x2" | jq .
echo ""
echo ""

# 测试 6: 检查返回的 markets 是否包含 outcomes
echo "6. 检查返回的 markets 是否包含 outcomes"
echo "GET $BASE_URL/api/events?limit=1&market_types=1x2"
curl -s "$BASE_URL/api/events?limit=1&market_types=1x2" | jq '.events[0].markets'
echo ""
echo ""

echo "=========================================="
echo "测试完成"
echo "=========================================="
