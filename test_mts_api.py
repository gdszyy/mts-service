#!/usr/bin/env python3
"""
MTS Service API 测试脚本

该脚本用于测试 mts-service 的所有 API 端点，包括：
- 健康检查
- 单注投注
- 串关投注
- 系统投注
- Banker 系统投注
- 预设系统投注
- 多重彩投注

测试结果将保存到 test_results.json 文件中。
"""

import requests
import json
import time
import uuid
from datetime import datetime
from typing import Dict, List, Any


class MTSAPITester:
    """MTS API 测试类"""

    def __init__(self, base_url: str = "http://localhost:8080"):
        self.base_url = base_url
        self.results = []
        self.session = requests.Session()
        self.session.headers.update({"Content-Type": "application/json"})

    def log(self, message: str, level: str = "INFO"):
        """打印日志"""
        timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        print(f"[{timestamp}] [{level}] {message}")

    def generate_ticket_id(self, prefix: str = "test") -> str:
        """生成唯一的 ticket ID"""
        return f"{prefix}-{uuid.uuid4().hex[:12]}-{int(time.time())}"

    def make_request(self, method: str, endpoint: str, data: Dict = None) -> Dict:
        """发送 HTTP 请求"""
        url = f"{self.base_url}{endpoint}"
        try:
            if method.upper() == "GET":
                response = self.session.get(url)
            elif method.upper() == "POST":
                response = self.session.post(url, json=data)
            else:
                raise ValueError(f"Unsupported HTTP method: {method}")

            return {
                "status_code": response.status_code,
                "response": response.json() if response.text else {},
                "success": response.status_code == 200
            }
        except requests.exceptions.ConnectionError:
            self.log(f"连接失败: {url}", "ERROR")
            return {
                "status_code": 0,
                "response": {"error": "Connection failed"},
                "success": False
            }
        except Exception as e:
            self.log(f"请求异常: {str(e)}", "ERROR")
            return {
                "status_code": 0,
                "response": {"error": str(e)},
                "success": False
            }

    def record_result(self, test_id: str, description: str, result: Dict, expected: str):
        """记录测试结果"""
        passed = self._check_expectation(result, expected)
        self.results.append({
            "test_id": test_id,
            "description": description,
            "timestamp": datetime.now().isoformat(),
            "result": result,
            "expected": expected,
            "passed": passed
        })
        status = "✓ PASSED" if passed else "✗ FAILED"
        self.log(f"{test_id}: {description} - {status}", "RESULT")

    def _check_expectation(self, result: Dict, expected: str) -> bool:
        """检查结果是否符合预期"""
        if "200 OK" in expected and result["status_code"] == 200:
            if "accepted" in expected:
                return result.get("response", {}).get("data", {}).get("content", {}).get("status") == "accepted"
            return True
        elif "400 Bad Request" in expected and result["status_code"] == 400:
            return True
        elif "rejected" in expected:
            return result.get("response", {}).get("data", {}).get("content", {}).get("status") == "rejected"
        return False

    # ==================== 测试用例 ====================

    def test_health_check(self):
        """TC-HEALTH-001: 健康检查"""
        self.log("开始测试: 健康检查", "TEST")
        result = self.make_request("GET", "/health")
        self.record_result(
            "TC-HEALTH-001",
            "健康检查",
            result,
            "200 OK"
        )

    def test_single_bet_success(self):
        """TC-SINGLE-001: 成功提交单注"""
        self.log("开始测试: 成功提交单注", "TEST")
        data = {
            "ticketId": self.generate_ticket_id("single"),
            "selection": {
                "productId": "3",
                "eventId": "sr:match:12345",
                "marketId": "1",
                "outcomeId": "1",
                "odds": "2.50",
                "specifiers": ""
            },
            "stake": {
                "type": "cash",
                "currency": "EUR",
                "amount": 10.00,
                "mode": "total"
            },
            "context": {
                "channelType": "internet",
                "language": "EN",
                "ip": "192.168.1.1"
            }
        }
        result = self.make_request("POST", "/api/bets/single", data)
        self.record_result(
            "TC-SINGLE-001",
            "成功提交单注",
            result,
            "200 OK, accepted"
        )

    def test_single_bet_missing_ticket_id(self):
        """TC-SINGLE-002: 缺少 ticketId"""
        self.log("开始测试: 缺少 ticketId", "TEST")
        data = {
            "selection": {
                "productId": "3",
                "eventId": "sr:match:12345",
                "marketId": "1",
                "outcomeId": "1",
                "odds": "2.75",            },
            "stake": {
                "type": "cash",
                "currency": "EUR",
                "amount": 10.00,
                "mode": "total"
            }
        }
        result = self.make_request("POST", "/api/bets/single", data)
        self.record_result(
            "TC-SINGLE-002",
            "缺少 ticketId",
            result,
            "400 Bad Request"
        )

    def test_single_bet_invalid_odds(self):
        """TC-SINGLE-003: 无效的赔率"""
        self.log("开始测试: 无效的赔率", "TEST")
        data = {
            "ticketId": self.generate_ticket_id("single"),
            "selection": {
                "productId": "3",
                "eventId": "sr:match:12345",
                "marketId": "1",
                "outcomeId": "1",
                "odds": "-1.0"
            },
            "stake": {
                "type": "cash",
                "currency": "EUR",
                "amount": 10.00,
                "mode": "total"
            }
        }
        result = self.make_request("POST", "/api/bets/single", data)
        self.record_result(
            "TC-SINGLE-003",
            "无效的赔率",
            result,
            "400 Bad Request"
        )

    def test_accumulator_success(self):
        """TC-ACC-001: 成功提交串关"""
        self.log("开始测试: 成功提交串关", "TEST")
        data = {
            "ticketId": self.generate_ticket_id("acc"),
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
                }
            ],
            "stake": {
                "type": "cash",
                "currency": "EUR",
                "amount": 10.00,
                "mode": "total"
            }
        }
        result = self.make_request("POST", "/api/bets/accumulator", data)
        self.record_result(
            "TC-ACC-001",
            "成功提交串关",
            result,
            "200 OK, accepted"
        )

    def test_accumulator_insufficient_selections(self):
        """TC-ACC-002: 选项数量不足"""
        self.log("开始测试: 选项数量不足", "TEST")
        data = {
            "ticketId": self.generate_ticket_id("acc"),
            "selections": [
                {
                    "productId": "3",
                    "eventId": "sr:match:12345",
                    "marketId": "1",
                    "outcomeId": "1",
                    "odds": 2.50
                }
            ],
            "stake": {
                "type": "cash",
                "currency": "EUR",
                "amount": 10.00,
                "mode": "total"
            }
        }
        result = self.make_request("POST", "/api/bets/accumulator", data)
        self.record_result(
            "TC-ACC-002",
            "选项数量不足",
            result,
            "400 Bad Request"
        )

    def test_system_bet_success(self):
        """TC-SYS-001: 成功提交系统投注"""
        self.log("开始测试: 成功提交系统投注", "TEST")
        data = {
            "ticketId": self.generate_ticket_id("sys"),
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
        }
        result = self.make_request("POST", "/api/bets/system", data)
        self.record_result(
            "TC-SYS-001",
            "成功提交系统投注",
            result,
            "200 OK, accepted"
        )

    def test_system_bet_invalid_mode(self):
        """TC-SYS-002: 无效的 stake mode"""
        self.log("开始测试: 无效的 stake mode", "TEST")
        data = {
            "ticketId": self.generate_ticket_id("sys"),
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
                    "odds": "1.80",                },
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
        }
        result = self.make_request("POST", "/api/bets/system", data)
        self.record_result(
            "TC-SYS-002",
            "无效的 stake mode",
            result,
            "400 Bad Request"
        )

    def test_banker_system_success(self):
        """TC-BANKER-001: 成功提交 Banker 系统投注"""
        self.log("开始测试: 成功提交 Banker 系统投注", "TEST")
        data = {
            "ticketId": self.generate_ticket_id("banker"),
            "bankers": [
                {
                    "productId": "3",
                    "eventId": "sr:match:12345",
                    "marketId": "1",
                    "outcomeId": "1",
                    "odds": "1.50"
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
                    "odds": "2.20",                }
            ],
            "stake": {
                "type": "cash",
                "currency": "EUR",
                "amount": 1.00,
                "mode": "unit"
            }
        }
        result = self.make_request("POST", "/api/bets/banker-system", data)
        self.record_result(
            "TC-BANKER-001",
            "成功提交 Banker 系统投注",
            result,
            "200 OK, accepted"
        )

    def test_preset_trixie_success(self):
        """TC-PRESET-001: 成功提交 Trixie 投注"""
        self.log("开始测试: 成功提交 Trixie 投注", "TEST")
        data = {
            "ticketId": self.generate_ticket_id("trixie"),
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
        }
        result = self.make_request("POST", "/api/bets/preset", data)
        self.record_result(
            "TC-PRESET-001",
            "成功提交 Trixie 投注",
            result,
            "200 OK, accepted"
        )

    def test_preset_invalid_selection_count(self):
        """TC-PRESET-002: 选项数量与类型不匹配"""
        self.log("开始测试: 选项数量与类型不匹配", "TEST")
        data = {
            "ticketId": self.generate_ticket_id("trixie"),
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
                }
            ],
            "stake": {
                "type": "cash",
                "currency": "EUR",
                "amount": 1.00,
                "mode": "unit"
            }
        }
        result = self.make_request("POST", "/api/bets/preset", data)
        self.record_result(
            "TC-PRESET-002",
            "选项数量与类型不匹配",
            result,
            "400 Bad Request"
        )

    def test_multi_bet_success(self):
        """TC-MULTI-001: 成功提交多重彩"""
        self.log("开始测试: 成功提交多重彩", "TEST")
        data = {
            "ticketId": self.generate_ticket_id("multi"),
            "bets": [
                {
                    "type": "single",
                    "selections": [
                        {
                            "productId": "3",
                            "eventId": "sr:match:12345",
                            "marketId": "1",
                            "outcomeId": "1",
                            "odds": 2.50
                        }
                    ],
                    "stake": {
                        "type": "cash",
                        "currency": "EUR",
                        "amount": 5.00,
                        "mode": "total"
                    }
                },
                {
                    "type": "accumulator",
                    "selections": [
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
                            "odds": "3.00",                        }
                    ],
                    "stake": {
                        "type": "cash",
                        "currency": "EUR",
                        "amount": 10.00,
                        "mode": "total"
                    }
                }
            ]
        }
        result = self.make_request("POST", "/api/bets/multi", data)
        self.record_result(
            "TC-MULTI-001",
            "成功提交多重彩",
            result,
            "200 OK, accepted"
        )

    def run_all_tests(self):
        """运行所有测试"""
        self.log("=" * 60, "INFO")
        self.log("开始 MTS API 测试", "INFO")
        self.log("=" * 60, "INFO")

        # 健康检查
        self.test_health_check()

        # 单注测试
        self.test_single_bet_success()
        self.test_single_bet_missing_ticket_id()
        self.test_single_bet_invalid_odds()

        # 串关测试
        self.test_accumulator_success()
        self.test_accumulator_insufficient_selections()

        # 系统投注测试
        self.test_system_bet_success()
        self.test_system_bet_invalid_mode()

        # Banker 系统投注测试
        self.test_banker_system_success()

        # 预设系统投注测试
        self.test_preset_trixie_success()
        self.test_preset_invalid_selection_count()

        # 多重彩测试
        self.test_multi_bet_success()

        self.log("=" * 60, "INFO")
        self.log("测试完成", "INFO")
        self.log("=" * 60, "INFO")

    def generate_report(self, output_file: str = "test_results.json"):
        """生成测试报告"""
        total = len(self.results)
        passed = sum(1 for r in self.results if r["passed"])
        failed = total - passed

        summary = {
            "timestamp": datetime.now().isoformat(),
            "total_tests": total,
            "passed": passed,
            "failed": failed,
            "pass_rate": f"{(passed / total * 100):.2f}%" if total > 0 else "0.00%",
            "results": self.results
        }

        with open(output_file, "w", encoding="utf-8") as f:
            json.dump(summary, f, indent=2, ensure_ascii=False)

        self.log(f"测试报告已保存到: {output_file}", "INFO")
        self.log(f"总计: {total} | 通过: {passed} | 失败: {failed} | 通过率: {summary['pass_rate']}", "SUMMARY")

        return summary


def main():
    """主函数"""
    import sys

    base_url = sys.argv[1] if len(sys.argv) > 1 else "http://localhost:8080"
    
    tester = MTSAPITester(base_url)
    tester.run_all_tests()
    tester.generate_report()


if __name__ == "__main__":
    main()
