#!/usr/bin/env python3
"""
MTS 异步结果验证脚本

该脚本用于验证 MTS ticket 的异步处理流程：
1. 提交一个 ticket
2. 保存 ticketId 和 signature
3. 模拟异步查询（通过日志或其他方式）
4. 验证结果的一致性

注意：由于 mts-service 没有提供查询接口，本脚本主要演示异步验证的概念。
在实际生产环境中，可以通过 MTS 的回调机制或查询 API 来获取最终结果。
"""

import requests
import json
import time
import uuid
from datetime import datetime
from typing import Dict, Optional


class AsyncResultVerifier:
    """异步结果验证器"""

    def __init__(self, base_url: str = "http://localhost:8080"):
        self.base_url = base_url
        self.session = requests.Session()
        self.session.headers.update({"Content-Type": "application/json"})
        self.submitted_tickets = []

    def log(self, message: str, level: str = "INFO"):
        """打印日志"""
        timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        print(f"[{timestamp}] [{level}] {message}")

    def generate_ticket_id(self, prefix: str = "async-test") -> str:
        """生成唯一的 ticket ID"""
        return f"{prefix}-{uuid.uuid4().hex[:12]}-{int(time.time())}"

    def submit_ticket(self, bet_type: str, data: Dict) -> Optional[Dict]:
        """提交 ticket"""
        endpoint_map = {
            "single": "/api/bets/single",
            "accumulator": "/api/bets/accumulator",
            "system": "/api/bets/system",
            "banker": "/api/bets/banker-system",
            "preset": "/api/bets/preset",
            "multi": "/api/bets/multi"
        }

        endpoint = endpoint_map.get(bet_type)
        if not endpoint:
            self.log(f"不支持的投注类型: {bet_type}", "ERROR")
            return None

        url = f"{self.base_url}{endpoint}"
        
        try:
            self.log(f"提交 {bet_type} ticket: {data.get('ticketId')}", "INFO")
            response = self.session.post(url, json=data)
            
            if response.status_code == 200:
                result = response.json()
                if result.get("success"):
                    ticket_info = {
                        "ticketId": data.get("ticketId"),
                        "bet_type": bet_type,
                        "submit_time": datetime.now().isoformat(),
                        "response": result.get("data", {}),
                        "status": result.get("data", {}).get("content", {}).get("status"),
                        "signature": result.get("data", {}).get("content", {}).get("signature")
                    }
                    self.submitted_tickets.append(ticket_info)
                    self.log(f"✓ Ticket 提交成功: {data.get('ticketId')} - Status: {ticket_info['status']}", "SUCCESS")
                    return ticket_info
                else:
                    self.log(f"✗ Ticket 提交失败: {result.get('error')}", "ERROR")
                    return None
            else:
                self.log(f"✗ HTTP 错误: {response.status_code}", "ERROR")
                return None
                
        except Exception as e:
            self.log(f"✗ 异常: {str(e)}", "ERROR")
            return None

    def verify_ticket_status(self, ticket_info: Dict) -> bool:
        """
        验证 ticket 状态
        
        注意：由于 mts-service 没有提供查询接口，这里只是演示概念。
        在实际环境中，可以：
        1. 调用 MTS 的查询 API
        2. 监听 MTS 的回调
        3. 查询数据库中的 ticket 状态
        """
        self.log(f"验证 ticket: {ticket_info['ticketId']}", "INFO")
        
        # 模拟异步处理延迟
        time.sleep(1)
        
        # 在实际环境中，这里应该调用查询 API
        # 例如: GET /api/tickets/{ticketId}/status
        # 但由于 mts-service 没有提供这个接口，我们只能基于提交时的响应进行验证
        
        expected_status = ticket_info.get("status")
        self.log(f"  - 提交时状态: {expected_status}", "INFO")
        self.log(f"  - Signature: {ticket_info.get('signature', 'N/A')[:50]}...", "INFO")
        
        # 在实际环境中，这里应该比较查询到的状态与提交时的状态
        # 目前我们只能假设状态保持一致
        verified = expected_status in ["accepted", "rejected"]
        
        if verified:
            self.log(f"✓ Ticket 验证通过: {ticket_info['ticketId']}", "SUCCESS")
        else:
            self.log(f"✗ Ticket 验证失败: {ticket_info['ticketId']}", "ERROR")
        
        return verified

    def run_async_verification_test(self):
        """运行异步验证测试"""
        self.log("=" * 60, "INFO")
        self.log("开始异步验证测试", "INFO")
        self.log("=" * 60, "INFO")

        # 测试 1: 提交单注
        single_ticket = self.submit_ticket("single", {
            "ticketId": self.generate_ticket_id("single"),
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
        })

        # 测试 2: 提交串关
        accumulator_ticket = self.submit_ticket("accumulator", {
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
        })

        # 测试 3: 提交 Trixie
        trixie_ticket = self.submit_ticket("preset", {
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
        })

        self.log("=" * 60, "INFO")
        self.log("开始验证提交的 tickets", "INFO")
        self.log("=" * 60, "INFO")

        # 验证所有提交的 tickets
        verification_results = []
        for ticket in self.submitted_tickets:
            verified = self.verify_ticket_status(ticket)
            verification_results.append({
                "ticketId": ticket["ticketId"],
                "bet_type": ticket["bet_type"],
                "verified": verified
            })

        self.log("=" * 60, "INFO")
        self.log("异步验证测试完成", "INFO")
        self.log("=" * 60, "INFO")

        return {
            "submitted_tickets": self.submitted_tickets,
            "verification_results": verification_results
        }

    def save_results(self, results: Dict, output_file: str = "async_verification_results.json"):
        """保存验证结果"""
        with open(output_file, "w", encoding="utf-8") as f:
            json.dump(results, f, indent=2, ensure_ascii=False)
        
        self.log(f"验证结果已保存到: {output_file}", "INFO")
        
        total = len(results["verification_results"])
        verified = sum(1 for r in results["verification_results"] if r["verified"])
        self.log(f"总计: {total} | 验证通过: {verified} | 验证失败: {total - verified}", "SUMMARY")


def main():
    """主函数"""
    import sys

    base_url = sys.argv[1] if len(sys.argv) > 1 else "http://localhost:8080"
    
    verifier = AsyncResultVerifier(base_url)
    results = verifier.run_async_verification_test()
    verifier.save_results(results)


if __name__ == "__main__":
    main()
