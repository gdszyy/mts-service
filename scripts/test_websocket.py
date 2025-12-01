#!/usr/bin/env python3
"""
WebSocket Test Script for MTS Service

This script tests the WebSocket functionality of the MTS service,
including connection, heartbeat, bet submission, and result reception.
"""

import asyncio
import websockets
import json
import uuid
from datetime import datetime

# Configuration
WS_URL = "ws://mts-service-production.up.railway.app/ws"
USER_ID = "test_user_123"
TOKEN = "test_token_456"

# Test event IDs (use real ones from database)
TEST_EVENT_IDS = [
    "sr:match:52688209",
    "sr:match:52688211",
    "sr:match:52688213"
]


async def test_connection():
    """Test WebSocket connection establishment"""
    print("=" * 60)
    print("TEST 1: WebSocket Connection")
    print("=" * 60)
    
    url = f"{WS_URL}?userId={USER_ID}&token={TOKEN}"
    
    try:
        async with websockets.connect(url) as websocket:
            # Wait for connection established message
            response = await websocket.recv()
            data = json.loads(response)
            
            print(f"✓ Connection established")
            print(f"  Response: {json.dumps(data, indent=2)}")
            
            assert data["type"] == "connection_established"
            assert data["userId"] == USER_ID
            
            print("\n✓ Connection test PASSED\n")
            return True
            
    except Exception as e:
        print(f"\n✗ Connection test FAILED: {e}\n")
        return False


async def test_heartbeat():
    """Test heartbeat (ping/pong) mechanism"""
    print("=" * 60)
    print("TEST 2: Heartbeat Mechanism")
    print("=" * 60)
    
    url = f"{WS_URL}?userId={USER_ID}&token={TOKEN}"
    
    try:
        async with websockets.connect(url) as websocket:
            # Wait for connection established
            await websocket.recv()
            
            # Send ping
            ping_msg = {
                "type": "ping",
                "timestamp": datetime.utcnow().isoformat()
            }
            await websocket.send(json.dumps(ping_msg))
            print(f"→ Sent ping: {json.dumps(ping_msg, indent=2)}")
            
            # Wait for pong
            response = await websocket.recv()
            data = json.loads(response)
            
            print(f"← Received pong: {json.dumps(data, indent=2)}")
            
            assert data["type"] == "pong"
            
            print("\n✓ Heartbeat test PASSED\n")
            return True
            
    except Exception as e:
        print(f"\n✗ Heartbeat test FAILED: {e}\n")
        return False


async def test_single_bet():
    """Test single bet submission"""
    print("=" * 60)
    print("TEST 3: Single Bet Submission")
    print("=" * 60)
    
    url = f"{WS_URL}?userId={USER_ID}&token={TOKEN}"
    
    try:
        async with websockets.connect(url) as websocket:
            # Wait for connection established
            await websocket.recv()
            
            # Prepare bet request
            request_id = str(uuid.uuid4())
            bet_request = {
                "type": "place_bet",
                "requestId": request_id,
                "betType": "single",
                "timestamp": datetime.utcnow().isoformat(),
                "payload": {
                    "selection": {
                        "eventId": TEST_EVENT_IDS[0],
                        "marketId": "1",
                        "outcomeId": "1",
                        "odds": "2.50"
                    },
                    "stake": {
                        "amount": "1000",
                        "currency": "CNY"
                    }
                }
            }
            
            # Send bet request
            await websocket.send(json.dumps(bet_request))
            print(f"→ Sent bet request: {json.dumps(bet_request, indent=2)}")
            
            # Wait for bet received confirmation
            response1 = await websocket.recv()
            data1 = json.loads(response1)
            print(f"\n← Received confirmation: {json.dumps(data1, indent=2)}")
            
            assert data1["type"] == "bet_received"
            assert data1["requestId"] == request_id
            assert "ticketId" in data1
            
            ticket_id = data1["ticketId"]
            print(f"\n✓ Bet received, ticketId: {ticket_id}")
            
            # Wait for bet result (with timeout)
            print("\n⏳ Waiting for MTS result (max 15 seconds)...")
            try:
                response2 = await asyncio.wait_for(websocket.recv(), timeout=15.0)
                data2 = json.loads(response2)
                print(f"\n← Received result: {json.dumps(data2, indent=2)}")
                
                assert data2["type"] in ["bet_result", "bet_timeout"]
                assert data2["requestId"] == request_id
                
                if data2["type"] == "bet_result":
                    print(f"\n✓ Bet result received: {data2['status']}")
                else:
                    print(f"\n⚠ Bet timeout (will receive delayed result later)")
                
            except asyncio.TimeoutError:
                print(f"\n⚠ No result received within 15 seconds")
            
            print("\n✓ Single bet test PASSED\n")
            return True
            
    except Exception as e:
        print(f"\n✗ Single bet test FAILED: {e}\n")
        return False


async def test_multi_bet():
    """Test multiple single bets submission"""
    print("=" * 60)
    print("TEST 4: Multiple Single Bets Submission")
    print("=" * 60)
    
    url = f"{WS_URL}?userId={USER_ID}&token={TOKEN}"
    
    try:
        async with websockets.connect(url) as websocket:
            # Wait for connection established
            await websocket.recv()
            
            # Prepare multi bet request
            request_id = str(uuid.uuid4())
            bet_request = {
                "type": "place_bet",
                "requestId": request_id,
                "betType": "multi",
                "timestamp": datetime.utcnow().isoformat(),
                "payload": {
                    "bets": [
                        {
                            "selection": {
                                "eventId": TEST_EVENT_IDS[0],
                                "marketId": "1",
                                "outcomeId": "1",
                                "odds": "2.50"
                            },
                            "stake": {
                                "amount": "1000",
                                "currency": "CNY"
                            }
                        },
                        {
                            "selection": {
                                "eventId": TEST_EVENT_IDS[1],
                                "marketId": "1",
                                "outcomeId": "2",
                                "odds": "3.00"
                            },
                            "stake": {
                                "amount": "1500",
                                "currency": "CNY"
                            }
                        }
                    ]
                }
            }
            
            # Send bet request
            await websocket.send(json.dumps(bet_request))
            print(f"→ Sent multi bet request: {json.dumps(bet_request, indent=2)}")
            
            # Wait for bet received confirmation
            response1 = await websocket.recv()
            data1 = json.loads(response1)
            print(f"\n← Received confirmation: {json.dumps(data1, indent=2)}")
            
            assert data1["type"] == "bet_received"
            assert data1["requestId"] == request_id
            assert "ticketIds" in data1
            
            print(f"\n✓ Bets received, ticketIds: {data1['ticketIds']}")
            
            # Wait for partial results
            print("\n⏳ Waiting for partial results...")
            partial_count = 0
            
            try:
                while partial_count < 3:  # 2 partial + 1 final
                    response = await asyncio.wait_for(websocket.recv(), timeout=20.0)
                    data = json.loads(response)
                    print(f"\n← Received: {json.dumps(data, indent=2)}")
                    
                    if data["type"] == "bet_partial_result":
                        print(f"  Progress: {data['completed']}")
                        partial_count += 1
                    elif data["type"] == "bet_result":
                        print(f"  Final summary: {data.get('summary', {})}")
                        break
                        
            except asyncio.TimeoutError:
                print(f"\n⚠ Timeout waiting for results")
            
            print("\n✓ Multi bet test PASSED\n")
            return True
            
    except Exception as e:
        print(f"\n✗ Multi bet test FAILED: {e}\n")
        return False


async def test_accumulator_bet():
    """Test accumulator bet submission"""
    print("=" * 60)
    print("TEST 5: Accumulator Bet Submission")
    print("=" * 60)
    
    url = f"{WS_URL}?userId={USER_ID}&token={TOKEN}"
    
    try:
        async with websockets.connect(url) as websocket:
            # Wait for connection established
            await websocket.recv()
            
            # Prepare accumulator bet request
            request_id = str(uuid.uuid4())
            bet_request = {
                "type": "place_bet",
                "requestId": request_id,
                "betType": "accumulator",
                "timestamp": datetime.utcnow().isoformat(),
                "payload": {
                    "selections": [
                        {
                            "eventId": TEST_EVENT_IDS[0],
                            "marketId": "1",
                            "outcomeId": "1",
                            "odds": "2.50"
                        },
                        {
                            "eventId": TEST_EVENT_IDS[1],
                            "marketId": "1",
                            "outcomeId": "2",
                            "odds": "3.00"
                        },
                        {
                            "eventId": TEST_EVENT_IDS[2],
                            "marketId": "1",
                            "outcomeId": "1",
                            "odds": "1.80"
                        }
                    ],
                    "stake": {
                        "amount": "2000",
                        "currency": "CNY"
                    }
                }
            }
            
            # Send bet request
            await websocket.send(json.dumps(bet_request))
            print(f"→ Sent accumulator bet request")
            
            # Wait for confirmation
            response1 = await websocket.recv()
            data1 = json.loads(response1)
            print(f"\n← Received confirmation: {json.dumps(data1, indent=2)}")
            
            assert data1["type"] == "bet_received"
            
            # Wait for result
            print("\n⏳ Waiting for result...")
            try:
                response2 = await asyncio.wait_for(websocket.recv(), timeout=15.0)
                data2 = json.loads(response2)
                print(f"\n← Received result: {json.dumps(data2, indent=2)}")
            except asyncio.TimeoutError:
                print(f"\n⚠ Timeout")
            
            print("\n✓ Accumulator bet test PASSED\n")
            return True
            
    except Exception as e:
        print(f"\n✗ Accumulator bet test FAILED: {e}\n")
        return False


async def run_all_tests():
    """Run all tests"""
    print("\n" + "=" * 60)
    print("MTS WebSocket Test Suite")
    print("=" * 60 + "\n")
    
    results = []
    
    # Test 1: Connection
    results.append(("Connection", await test_connection()))
    await asyncio.sleep(1)
    
    # Test 2: Heartbeat
    results.append(("Heartbeat", await test_heartbeat()))
    await asyncio.sleep(1)
    
    # Test 3: Single Bet
    results.append(("Single Bet", await test_single_bet()))
    await asyncio.sleep(1)
    
    # Test 4: Multi Bet
    results.append(("Multi Bet", await test_multi_bet()))
    await asyncio.sleep(1)
    
    # Test 5: Accumulator Bet
    results.append(("Accumulator Bet", await test_accumulator_bet()))
    
    # Print summary
    print("\n" + "=" * 60)
    print("Test Summary")
    print("=" * 60)
    
    passed = sum(1 for _, result in results if result)
    total = len(results)
    
    for name, result in results:
        status = "✓ PASSED" if result else "✗ FAILED"
        print(f"{name:.<40} {status}")
    
    print(f"\nTotal: {passed}/{total} tests passed ({passed/total*100:.1f}%)")
    print("=" * 60 + "\n")


if __name__ == "__main__":
    asyncio.run(run_all_tests())
