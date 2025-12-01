#!/bin/bash

# MTS Service 测试运行脚本
# 该脚本用于运行所有 MTS API 测试

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 默认配置
BASE_URL="${BASE_URL:-http://localhost:8080}"
WAIT_TIME=5

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}MTS Service API 测试套件${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# 检查服务是否运行
echo -e "${YELLOW}[1/4] 检查服务状态...${NC}"
if curl -s -f "${BASE_URL}/health" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ 服务正在运行${NC}"
else
    echo -e "${RED}✗ 服务未运行或无法访问: ${BASE_URL}${NC}"
    echo -e "${YELLOW}请确保 mts-service 已启动并监听在 ${BASE_URL}${NC}"
    echo ""
    echo "启动服务的命令示例："
    echo "  cd /home/ubuntu/mts-service"
    echo "  go run cmd/server/mts_main.go"
    echo ""
    exit 1
fi

echo ""

# 运行主测试套件
echo -e "${YELLOW}[2/4] 运行 API 功能测试...${NC}"
python3 test_mts_api.py "${BASE_URL}"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ API 功能测试完成${NC}"
else
    echo -e "${RED}✗ API 功能测试失败${NC}"
fi

echo ""

# 等待一段时间
echo -e "${YELLOW}等待 ${WAIT_TIME} 秒后继续...${NC}"
sleep ${WAIT_TIME}

echo ""

# 运行异步验证测试
echo -e "${YELLOW}[3/4] 运行异步验证测试...${NC}"
python3 verify_async_results.py "${BASE_URL}"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 异步验证测试完成${NC}"
else
    echo -e "${RED}✗ 异步验证测试失败${NC}"
fi

echo ""

# 生成测试报告
echo -e "${YELLOW}[4/4] 生成测试报告...${NC}"

if [ -f "test_results.json" ]; then
    echo -e "${GREEN}✓ 测试结果: test_results.json${NC}"
    
    # 提取测试统计
    TOTAL=$(jq -r '.total_tests' test_results.json)
    PASSED=$(jq -r '.passed' test_results.json)
    FAILED=$(jq -r '.failed' test_results.json)
    PASS_RATE=$(jq -r '.pass_rate' test_results.json)
    
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}测试统计${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo -e "总计测试: ${TOTAL}"
    echo -e "通过: ${GREEN}${PASSED}${NC}"
    echo -e "失败: ${RED}${FAILED}${NC}"
    echo -e "通过率: ${PASS_RATE}"
    echo ""
fi

if [ -f "async_verification_results.json" ]; then
    echo -e "${GREEN}✓ 异步验证结果: async_verification_results.json${NC}"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}测试完成!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "查看详细结果:"
echo "  cat test_results.json | jq ."
echo "  cat async_verification_results.json | jq ."
echo ""
