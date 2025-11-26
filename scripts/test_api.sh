#!/bin/bash

# MTS Service API Test Script
# This script tests all API endpoints with sample data

BASE_URL="${BASE_URL:-http://localhost:8080}"

echo "========================================="
echo "MTS Service API Test Script"
echo "========================================="
echo "Base URL: $BASE_URL"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Function to print test result
print_result() {
    local test_name=$1
    local status_code=$2
    local expected=$3
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if [ "$status_code" -eq "$expected" ]; then
        echo -e "${GREEN}✓ PASS${NC}: $test_name (HTTP $status_code)"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}✗ FAIL${NC}: $test_name (Expected HTTP $expected, got $status_code)"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
}

# Function to make API request
test_endpoint() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local expected_status=${5:-200}
    
    echo ""
    echo -e "${YELLOW}Testing:${NC} $name"
    echo "Endpoint: $method $endpoint"
    
    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi
    
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    echo "Response: $body" | head -c 200
    if [ ${#body} -gt 200 ]; then
        echo "..."
    fi
    echo ""
    
    print_result "$name" "$status_code" "$expected_status"
}

# Test 1: Health Check
test_endpoint "Health Check" "GET" "/health" "" 200

# Test 2: Root Endpoint (API Documentation)
test_endpoint "Root Endpoint" "GET" "/" "" 200

# Test 3: Single Bet
SINGLE_BET_DATA='{
  "ticketId": "test-single-'$(date +%s)'",
  "selection": {
    "productId": "3",
    "eventId": "sr:match:12345",
    "marketId": "1",
    "outcomeId": "1",
    "odds": 2.50
  },
  "stake": {
    "type": "cash",
    "currency": "EUR",
    "amount": 10.00,
    "mode": "total"
  }
}'
test_endpoint "Single Bet" "POST" "/api/bets/single" "$SINGLE_BET_DATA" 200

# Test 4: Accumulator Bet
ACCUMULATOR_DATA='{
  "ticketId": "test-acc-'$(date +%s)'",
  "selections": [
    {
      "productId": "3",
      "eventId": "sr:match:12345",
      "marketId": "1",
      "outcomeId": "1",
      "odds": 2.50
    },
    {
      "productId": "3",
      "eventId": "sr:match:12346",
      "marketId": "1",
      "outcomeId": "2",
      "odds": 1.80
    },
    {
      "productId": "3",
      "eventId": "sr:match:12347",
      "marketId": "1",
      "outcomeId": "1",
      "odds": 3.00
    }
  ],
  "stake": {
    "type": "cash",
    "currency": "EUR",
    "amount": 10.00,
    "mode": "total"
  }
}'
test_endpoint "Accumulator Bet" "POST" "/api/bets/accumulator" "$ACCUMULATOR_DATA" 200

# Test 5: System Bet
SYSTEM_DATA='{
  "ticketId": "test-sys-'$(date +%s)'",
  "size": [2],
  "selections": [
    {
      "productId": "3",
      "eventId": "sr:match:12345",
      "marketId": "1",
      "outcomeId": "1",
      "odds": 2.50
    },
    {
      "productId": "3",
      "eventId": "sr:match:12346",
      "marketId": "1",
      "outcomeId": "2",
      "odds": 1.80
    },
    {
      "productId": "3",
      "eventId": "sr:match:12347",
      "marketId": "1",
      "outcomeId": "1",
      "odds": 3.00
    }
  ],
  "stake": {
    "type": "cash",
    "currency": "EUR",
    "amount": 1.00,
    "mode": "unit"
  }
}'
test_endpoint "System Bet (2/3)" "POST" "/api/bets/system" "$SYSTEM_DATA" 200

# Test 6: Banker System Bet
BANKER_DATA='{
  "ticketId": "test-banker-'$(date +%s)'",
  "bankers": [
    {
      "productId": "3",
      "eventId": "sr:match:12345",
      "marketId": "1",
      "outcomeId": "1",
      "odds": 1.50
    }
  ],
  "size": [2],
  "selections": [
    {
      "productId": "3",
      "eventId": "sr:match:12346",
      "marketId": "1",
      "outcomeId": "2",
      "odds": 2.50
    },
    {
      "productId": "3",
      "eventId": "sr:match:12347",
      "marketId": "1",
      "outcomeId": "1",
      "odds": 3.00
    },
    {
      "productId": "3",
      "eventId": "sr:match:12348",
      "marketId": "1",
      "outcomeId": "2",
      "odds": 2.20
    }
  ],
  "stake": {
    "type": "cash",
    "currency": "EUR",
    "amount": 1.00,
    "mode": "unit"
  }
}'
test_endpoint "Banker System Bet" "POST" "/api/bets/banker-system" "$BANKER_DATA" 200

# Test 7: Trixie Bet
TRIXIE_DATA='{
  "ticketId": "test-trixie-'$(date +%s)'",
  "type": "trixie",
  "selections": [
    {
      "productId": "3",
      "eventId": "sr:match:12345",
      "marketId": "1",
      "outcomeId": "1",
      "odds": 2.50
    },
    {
      "productId": "3",
      "eventId": "sr:match:12346",
      "marketId": "1",
      "outcomeId": "2",
      "odds": 1.80
    },
    {
      "productId": "3",
      "eventId": "sr:match:12347",
      "marketId": "1",
      "outcomeId": "1",
      "odds": 3.00
    }
  ],
  "stake": {
    "type": "cash",
    "currency": "EUR",
    "amount": 1.00,
    "mode": "unit"
  }
}'
test_endpoint "Trixie Bet" "POST" "/api/bets/preset" "$TRIXIE_DATA" 200

# Test 8: Validation Error (missing ticketId)
INVALID_DATA='{
  "selection": {
    "productId": "3",
    "eventId": "sr:match:12345",
    "marketId": "1",
    "outcomeId": "1",
    "odds": 2.50
  },
  "stake": {
    "type": "cash",
    "currency": "EUR",
    "amount": 10.00,
    "mode": "total"
  }
}'
test_endpoint "Validation Error Test" "POST" "/api/bets/single" "$INVALID_DATA" 400

# Print summary
echo ""
echo "========================================="
echo "Test Summary"
echo "========================================="
echo "Total Tests: $TOTAL_TESTS"
echo -e "${GREEN}Passed: $PASSED_TESTS${NC}"
if [ $FAILED_TESTS -gt 0 ]; then
    echo -e "${RED}Failed: $FAILED_TESTS${NC}"
else
    echo "Failed: $FAILED_TESTS"
fi
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
fi
