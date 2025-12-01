#!/usr/bin/env python3
import requests, json, time, uuid
from datetime import datetime

class AsyncResultVerifier:
    def __init__(self, base_url: str = "http://localhost:8080"):
        self.base_url = base_url
        self.session = requests.Session()
        self.session.headers.update({"Content-Type": "application/json"})
        self.submitted_tickets = []

    def log(self, message: str, level: str = "INFO"):
        timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        print(f"[{timestamp}] [{level}] {message}")

    def generate_ticket_id(self, prefix: str = "async-test") -> str:
        return f"{prefix}-{uuid.uuid4().hex[:12]}-{int(time.time())}"

    def submit_ticket(self, bet_type: str, data: dict):
        endpoint_map = {
            "single": "/api/bets/single",
            "accumulator": "/api/bets/accumulator",
            "preset": "/api/bets/preset"
        }
        endpoint = endpoint_map.get(bet_type)
        if not endpoint:
            return None
        
        url = f"{self.base_url}{endpoint}"
        try:
            self.log(f"提交 {bet_type} ticket: {data.get('ticketId')}", "INFO")
            response = self.session.post(url, json=data)
            
            if response.status_code == 200:
                result = response.json()
                ticket_info = {
                    "ticketId": data.get("ticketId"),
                    "bet_type": bet_type,
                    "submit_time": datetime.now().isoformat(),
                    "response": result.get("data", {}),
                    "status": result.get("data", {}).get("content", {}).get("status"),
                    "signature": result.get("data", {}).get("content", {}).get("signature"),
                    "code": result.get("data", {}).get("content", {}).get("code"),
                    "message": result.get("data", {}).get("content", {}).get("message")
                }
                self.submitted_tickets.append(ticket_info)
                self.log(f"✓ Ticket 提交成功: {data.get('ticketId')} - Status: {ticket_info['status']}", "SUCCESS")
                return ticket_info
            else:
                self.log(f"✗ HTTP 错误: {response.status_code}", "ERROR")
                return None
        except Exception as e:
            self.log(f"✗ 异常: {str(e)}", "ERROR")
            return None

    def run_async_verification_test(self):
        self.log("=" * 60, "INFO")
        self.log("开始异步验证测试", "INFO")
        self.log("=" * 60, "INFO")

        # 测试 1: 单注
        self.submit_ticket("single", {
            "ticketId": self.generate_ticket_id("single"),
            "selection": {
                "productId": "3",
                "eventId": "sr:match:12345",
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
        })

        # 测试 2: 串关
        self.submit_ticket("accumulator", {
            "ticketId": self.generate_ticket_id("acc"),
            "selections": [
                {"productId": "3", "eventId": "sr:match:12345", "marketId": "1", "outcomeId": "1", "odds": "2.50"},
                {"productId": "3", "eventId": "sr:match:12346", "marketId": "1", "outcomeId": "2", "odds": "1.80"}
            ],
            "stake": {
                "type": "cash",
                "currency": "EUR",
                "amount": "10.00",
                "mode": "total"
            }
        })

        # 测试 3: Trixie
        self.submit_ticket("preset", {
            "ticketId": self.generate_ticket_id("trixie"),
            "type": "trixie",
            "selections": [
                {"productId": "3", "eventId": "sr:match:12345", "marketId": "1", "outcomeId": "1", "odds": "2.50"},
                {"productId": "3", "eventId": "sr:match:12346", "marketId": "1", "outcomeId": "2", "odds": "1.80"},
                {"productId": "3", "eventId": "sr:match:12347", "marketId": "1", "outcomeId": "1", "odds": "3.00"}
            ],
            "stake": {
                "type": "cash",
                "currency": "EUR",
                "amount": "1.00",
                "mode": "unit"
            }
        })

        self.log("=" * 60, "INFO")
        self.log("异步验证测试完成", "INFO")
        self.log("=" * 60, "INFO")

        return {
            "submitted_tickets": self.submitted_tickets,
            "verification_results": [
                {
                    "ticketId": t["ticketId"],
                    "bet_type": t["bet_type"],
                    "status": t["status"],
                    "verified": t["status"] in ["accepted", "rejected"]
                } for t in self.submitted_tickets
            ]
        }

    def save_results(self, results: dict, output_file: str = "async_verification_results.json"):
        with open(output_file, "w", encoding="utf-8") as f:
            json.dump(results, f, indent=2, ensure_ascii=False)
        
        self.log(f"验证结果已保存到: {output_file}", "INFO")
        total = len(results["verification_results"])
        verified = sum(1 for r in results["verification_results"] if r["verified"])
        self.log(f"总计: {total} | 验证通过: {verified} | 验证失败: {total - verified}", "SUMMARY")

import sys
base_url = sys.argv[1] if len(sys.argv) > 1 else "http://localhost:8080"
verifier = AsyncResultVerifier(base_url)
results = verifier.run_async_verification_test()
verifier.save_results(results)
