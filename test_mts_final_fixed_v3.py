#!/usr/bin/env python3
"""
MTS Service API 最终测试脚本 (修复 Banker 投注 v3)
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
        
        if not passed and result.get("response"):
            resp = result["response"]
            if resp.get("data", {}).get("content", {}).get("message"):
                message = resp.get("data", {}).get("content", {}).get("message", "No message")
                self.log(f"  错误信息: {message}", "DEBUG")
            elif resp.get("error"):
                error_message = resp.get("error", "No error message")
                self.log(f"  错误信息: {error_message}", "DEBUG")

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

    def test_banker_system_success_final(self):
        """TC-BANKER-001-FINAL: 成功提交 Banker 系统投注 (最终修复版)"""
        self.log("开始测试: 成功提交 Banker 系统投注 (最终修复版)", "TEST")
        
        data = {
            "ticketId": self.generate_ticket_id("banker-final"),
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
        self.record_result("TC-BANKER-001-FINAL", "成功提交 Banker 系统投注 (最终修复版)", result, "200 OK")

    def run_all_tests(self):
        """运行所有测试"""
        self.log("=" * 60, "INFO")
        self.log("开始 MTS API 最终测试 (修复 Banker v3)", "INFO")
        event_id_preview = ", ".join(self.event_ids[:5])
        self.log(f"使用的赛事 ID: {event_id_preview}...", "INFO")
        self.log("=" * 60, "INFO")

        self.test_banker_system_success_final()

        self.log("=" * 60, "INFO")
        self.log("测试完成", "INFO")
        self.log("=" * 60, "INFO")

    def generate_report(self, output_file: str = "final_test_results_fixed_v3.json"):
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
        pass_rate = summary["pass_rate"]
        self.log(f"总计: {total} | 通过: {passed} | 失败: {failed} | 通过率: {pass_rate}", "SUMMARY")

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
