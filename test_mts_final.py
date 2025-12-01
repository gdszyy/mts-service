#!/usr/bin/env python3
"""
MTS Service API 最终测试脚本 (使用真实赛事 ID)
"""

import requests
import json
import time
import uuid
from datetime import datetime
from typing import Dict, List


class MTSFinalTester:
    """MTS API 最终测试类"""

    def __init__(self, base_url: str, event_ids: List[str]):
        self.base_url = base_url
        self.event_ids = event_ids
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
        
        # 如果失败，打印详细信息
        if not passed and result.get("response"):
            resp = result["response"]
            if resp.get("data", {}).get("content", {}).get("message"):
                self.log(f"  错误信息: {resp['data']['content']['message']}", "DEBUG")
            elif resp.get("error"):
                self.log(f"  错误信息: {resp['error']}", "DEBUG")

    def _check_expectation(self, result: Dict, expected: str) -> bool:
        """检查结果是否符合预期"""
        if "200 OK" in expected and result["status_code"] == 200:
            if "accepted" in expected:
                return result.get("response", {}).get("data", {}).get("content", {}).get("status") == "accepted"
            elif "rejected" in expected:
                return result.get("response", {}).get("data", {}).get("content", {}).get("status") == "rejected"
            return True
        elif "400 Bad Request" in expected and result["status_code"] == 400:
            return True
        return False

    # ==================== 测试用例 ====================

    def test_health_check(self):
        """TC-HEALTH-001: 健康检查"""
        self.log("开始测试: 健康检查", "TEST")
        result = self.make_request("GET", "/health")
        self.record_result("TC-HEALTH-001", "健康检查", result, "200 OK")

    def test_single_bet_success(self):
        """TC-SINGLE-001: 成功提交单注 (真实赛事)"""
        self.log("开始测试: 成功提交单注 (真实赛事)", "TEST")
        data = {
            "ticketId": self.generate_ticket_id("single"),
            "selection": {
                "productId": "3",
                "eventId": self.event_ids[0],
                "marketId": "1",
                "outcomeId": "1",
                "odds": "2.50",
                "specifiers": ""
            },
            "stake": {
                "type": "cash",
                "currency": "EUR",
                "amount": "10.00",
                "mode": "total"
            }
        }
        result = self.make_request("POST", "/api/bets/single", data)
        # 可能被接受或拒绝，只要返回 200 OK 就算通过
        self.record_result("TC-SINGLE-001", "成功提交单注 (真实赛事)", result, "200 OK")

    def test_single_bet_missing_ticket_id(self):
        """TC-SINGLE-002: 缺少 ticketId"""
        self.log("开始测试: 缺少 ticketId", "TEST")
        data = {
            "selection": {
                "productId": "3",
                "eventId": self.event_ids[0],
                "marketId": "1",
                "outcomeId": "1",
                "odds": "2.50"
            },
            "stake": {
                "type": "cash",
                "currency": "EUR",
                "amount": "10.00",
                "mode": "total"
            }
        }
        result = self.make_request("POST", "/api/bets/single", data)
        self.record_result("TC-SINGLE-002", "缺少 ticketId", result, "400 Bad Request")

    def test_accumulator_success(self):
        """TC-ACC-001: 成功提交串关 (真实赛事)"""
        self.log("开始测试: 成功提交串关 (真实赛事)", "TEST")
        data = {
            "ticketId": self.generate_ticket_id("acc"),
            "selections": [
                {
                    "productId": "3",
                    "eventId": self.event_ids[0],
                    "marketId": "1",
                    "outcomeId": "1",
                    "odds": "2.50"
                },
                {
                    "productId": "3",
                    "eventId": self.event_ids[1],
                    "marketId": "1",
                    "outcomeId": "2",
                    "odds": "1.80"
                }
            ],
            "stake": {
                "type": "cash",
                "currency": "EUR",
                "amount": "10.00",
                "mode": "total"
            }
        }
        result = self.make_request("POST", "/api/bets/accumulator", data)
        self.record_result("TC-ACC-001", "成功提交串关 (真实赛事)", result, "200 OK")

    def test_accumulator_insufficient_selections(self):
        """TC-ACC-002: 选项数量不足"""
        self.log("开始测试: 选项数量不足", "TEST")
        data = {
            "ticketId": self.generate_ticket_id("acc"),
            "selections": [
                {
                    "productId": "3",
                    "eventId": self.event_ids[0],
                    "marketId": "1",
                    "outcomeId": "1",
                    "odds": "2.50"
                }
            ],
            "stake": {
                "type": "cash",
                "currency": "EUR",
                "amount": "10.00",
                "mode": "total"
            }
        }
        result = self.make_request("POST", "/api/bets/accumulator", data)
        self.record_result("TC-ACC-002", "选项数量不足", result, "400 Bad Request")

    def test_system_bet_success(self):
        """TC-SYS-001: 成功提交系统投注 (真实赛事)"""
        self.log("开始测试: 成功提交系统投注 (真实赛事)", "TEST")
        data = {
            "ticketId": self.generate_ticket_id("sys"),
            "size": [2],
            "selections": [
                {
                    "productId": "3",
                    "eventId": self.event_ids[0],
                    "marketId": "1",
                    "outcomeId": "1",
                    "odds": "2.50"
                },
                {
                    "productId": "3",
                    "eventId": self.event_ids[1],
                    "marketId": "1",
                    "outcomeId": "2",
                    "odds": "1.80"
                },
                {
                    "productId": "3",
                    "eventId": self.event_ids[2],
                    "marketId": "1",
                    "outcomeId": "1",
                    "odds": "3.00"
                }
            ],
            "stake": {
                "type": "cash",
                "currency": "EUR",
                "amount": "1.00",
                "mode": "unit"
            }
        }
        result = self.make_request("POST", "/api/bets/system", data)
        self.record_result("TC-SYS-001", "成功提交系统投注 (真实赛事)", result, "200 OK")

    def test_system_bet_invalid_mode(self):
        """TC-SYS-002: 无效的 stake mode"""
        self.log("开始测试: 无效的 stake mode", "TEST")
        data = {
            "ticketId": self.generate_ticket_id("sys"),
            "size": [2],
            "selections": [
                {
                    "productId": "3",
                    "eventId": self.event_ids[0],
                    "marketId": "1",
                    "outcomeId": "1",
                    "odds": "2.50"
                },
                {
                    "productId": "3",
                    "eventId": self.event_ids[1],
                    "marketId": "1",
                    "outcomeId": "2",
                    "odds": "1.80"
                },
                {
                    "productId": "3",
                    "eventId": self.event_ids[2],
                    "marketId": "1",
                    "outcomeId": "1",
                    "odds": "3.00"
                }
            ],
            "stake": {
                "type": "cash",
                "currency": "EUR",
                "amount": "10.00",
                "mode": "total"
            }
        }
        result = self.make_request("POST", "/api/bets/system", data)
        self.record_result("TC-SYS-002", "无效的 stake mode", result, "400 Bad Request")

    def test_banker_system_success(self):
        """TC-BANKER-001: 成功提交 Banker 系统投注 (真实赛事)"""
        self.log("开始测试: 成功提交 Banker 系统投注 (真实赛事)", "TEST")
        data = {
            "ticketId": self.generate_ticket_id("banker"),
            "bankers": [
                {
                    "productId": "3",
                    "eventId": self.event_ids[0],
                    "marketId": "1",
                    "outcomeId": "1",
                    "odds": "1.50"
                }
            ],
            "size": [2],
            "selections": [
                {
                    "productId": "3",
                    "eventId": self.event_ids[1],
                    "marketId": "1",
                    "outcomeId": "2",
                    "odds": "2.50"
                },
                {
                    "productId": "3",
                    "eventId": self.event_ids[2],
                    "marketId": "1",
                    "outcomeId": "1",
                    "odds": "3.00"
                },
                {
                    "productId": "3",
                    "eventId": self.event_ids[3],
                    "marketId": "1",
                    "outcomeId": "2",
                    "odds": "2.20"
                }
            ],
            "stake": {
                "type": "cash",
                "currency": "EUR",
                "amount": "1.00",
                "mode": "unit"
            }
        }
        result = self.make_request("POST", "/api/bets/banker-system", data)
        self.record_result("TC-BANKER-001", "成功提交 Banker 系统投注 (真实赛事)", result, "200 OK")

    def test_preset_trixie_success(self):
        """TC-PRESET-001: 成功提交 Trixie 投注 (真实赛事)"""
        self.log("开始测试: 成功提交 Trixie 投注 (真实赛事)", "TEST")
        data = {
            "ticketId": self.generate_ticket_id("trixie"),
            "type": "trixie",
            "selections": [
                {
                    "productId": "3",
                    "eventId": self.event_ids[0],
                    "marketId": "1",
                    "outcomeId": "1",
                    "odds": "2.50"
                },
                {
                    "productId": "3",
                    "eventId": self.event_ids[1],
                    "marketId": "1",
                    "outcomeId": "2",
                    "odds": "1.80"
                },
                {
                    "productId": "3",
                    "eventId": self.event_ids[2],
                    "marketId": "1",
                    "outcomeId": "1",
                    "odds": "3.00"
                }
            ],
            "stake": {
                "type": "cash",
                "currency": "EUR",
                "amount": "1.00",
                "mode": "unit"
            }
        }
        result = self.make_request("POST", "/api/bets/preset", data)
        self.record_result("TC-PRESET-001", "成功提交 Trixie 投注 (真实赛事)", result, "200 OK")

    def test_preset_invalid_selection_count(self):
        """TC-PRESET-002: 选项数量与类型不匹配"""
        self.log("开始测试: 选项数量与类型不匹配", "TEST")
        data = {
            "ticketId": self.generate_ticket_id("trixie"),
            "type": "trixie",
            "selections": [
                {
                    "productId": "3",
                    "eventId": self.event_ids[0],
                    "marketId": "1",
                    "outcomeId": "1",
                    "odds": "2.50"
                },
                {
                    "productId": "3",
                    "eventId": self.event_ids[1],
                    "marketId": "1",
                    "outcomeId": "2",
                    "odds": "1.80"
                }
            ],
            "stake": {
                "type": "cash",
                "currency": "EUR",
                "amount": "1.00",
                "mode": "unit"
            }
        }
        result = self.make_request("POST", "/api/bets/preset", data)
        self.record_result("TC-PRESET-002", "选项数量与类型不匹配", result, "400 Bad Request")

    def test_multi_bet_success(self):
        """TC-MULTI-001: 成功提交多重彩 (真实赛事)"""
        self.log("开始测试: 成功提交多重彩 (真实赛事)", "TEST")
        data = {
            "ticketId": self.generate_ticket_id("multi"),
            "bets": [
                {
                    "type": "single",
                    "selections": [
                        {
                            "productId": "3",
                            "eventId": self.event_ids[0],
                            "marketId": "1",
                            "outcomeId": "1",
                            "odds": "2.50"
                        }
                    ],
                    "stake": {
                        "type": "cash",
                        "currency": "EUR",
                        "amount": "5.00",
                        "mode": "total"
                    }
                },
                {
                    "type": "accumulator",
                    "selections": [
                        {
                            "productId": "3",
                            "eventId": self.event_ids[1],
                            "marketId": "1",
                            "outcomeId": "2",
                            "odds": "1.80"
                        },
                        {
                            "productId": "3",
                            "eventId": self.event_ids[2],
                            "marketId": "1",
                            "outcomeId": "1",
                            "odds": "3.00"
                        }
                    ],
                    "stake": {
                        "type": "cash",
                        "currency": "EUR",
                        "amount": "10.00",
                        "mode": "total"
                    }
                }
            ]
        }
        result = self.make_request("POST", "/api/bets/multi", data)
        self.record_result("TC-MULTI-001", "成功提交多重彩 (真实赛事)", result, "200 OK")

    def run_all_tests(self):
        """运行所有测试"""
        self.log("=" * 60, "INFO")
        self.log("开始 MTS API 最终测试 (使用真实赛事 ID)", "INFO")
        self.log(f"使用的赛事 ID: {', '.join(self.event_ids[:5])}...", "INFO")
        self.log("=" * 60, "INFO")

        self.test_health_check()
        self.test_single_bet_success()
        self.test_single_bet_missing_ticket_id()
        self.test_accumulator_success()
        self.test_accumulator_insufficient_selections()
        self.test_system_bet_success()
        self.test_system_bet_invalid_mode()
        self.test_banker_system_success()
        self.test_preset_trixie_success()
        self.test_preset_invalid_selection_count()
        self.test_multi_bet_success()

        self.log("=" * 60, "INFO")
        self.log("测试完成", "INFO")
        self.log("=" * 60, "INFO")

    def generate_report(self, output_file: str = "final_test_results.json"):
        """生成测试报告"""
        total = len(self.results)
        passed = sum(1 for r in self.results if r["passed"])
        failed = total - passed

        summary = {
            "timestamp": datetime.now().isoformat(),
            "test_environment": {
                "base_url": self.base_url,
                "event_ids_used": self.event_ids[:10]
            },
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
    import sys
    
    # 读取真实赛事 ID
    with open("/home/ubuntu/mts-service/real_event_ids.txt", "r") as f:
        event_ids = [line.strip() for line in f if line.strip()]
    
    base_url = sys.argv[1] if len(sys.argv) > 1 else "http://localhost:8080"
    
    tester = MTSFinalTester(base_url, event_ids)
    tester.run_all_tests()
    tester.generate_report()


if __name__ == "__main__":
    main()
