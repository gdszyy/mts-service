# MTS Service - RESTful API Documentation

**Version**: 1.0.0
**Author**: Manus AI

## 1. Introduction

This document provides a comprehensive guide to using the RESTful API of the MTS Service. The service acts as a bridge to Sportradar's Managed Trading Services (MTS), allowing you to place bets and manage tickets via a simple HTTP interface.

## 2. Base URL

The base URL for all API endpoints will be provided by your Railway deployment, for example:

```
https://mts-service-production.up.railway.app
```

## 3. Authentication

All API endpoints are currently public, but it is recommended to add an API key or other authentication mechanism in a production environment.

## 4. API Endpoints

### 4.1. Health Check

This endpoint allows you to check the health of the service and its connection to MTS.

- **Endpoint**: `GET /health`
- **Method**: `GET`
- **Description**: Returns the current status of the service.

#### Response

**Success (Connected)**:
```json
{
  "status": "healthy",
  "timestamp": 1698412800,
  "service": "mts-service"
}
```

**Failure (Disconnected)**:
```json
{
  "status": "disconnected",
  "timestamp": 1698412800,
  "service": "mts-service"
}
```

### 4.2. Place Ticket

This is the primary endpoint for placing bets. It accepts a ticket with one or more bets and forwards it to MTS for processing.

- **Endpoint**: `POST /api/tickets`
- **Method**: `POST`
- **Description**: Places a new ticket with one or more bets.

#### Request Body

`Content-Type: application/json`

```json
{
  "ticketId": "TICKET_12345",
  "customerId": "customer_456",
  "currency": "EUR",
  "totalStake": 1000,
  "testSource": true,
  "oddsChange": "any",
  "bets": [
    {
      "id": "bet_1",
      "stake": 1000,
      "customBet": false,
      "selections": [
        {
          "id": "selection_1",
          "eventId": "sr:match:12345678",
          "odds": 20000,
          "banker": false
        }
      ]
    }
  ]
}
```

#### Field Descriptions

| Field | Type | Required | Description |
|:---|:---|:---:|:---|
| `ticketId` | string | Yes | Unique ID for the ticket |
| `customerId` | string | Yes | Your customer's ID |
| `currency` | string | Yes | 3-letter ISO currency code (e.g., "EUR") |
| `totalStake` | integer | Yes | Total stake for the ticket in cents (e.g., 1000 = 10.00 EUR) |
| `testSource` | boolean | Yes | `true` for test tickets, `false` for real money |
| `oddsChange` | string | Yes | Odds change strategy: `none`, `higher`, `any` |
| `bets` | array | Yes | Array of one or more bet objects |
| `bets[].id` | string | Yes | Unique ID for the bet |
| `bets[].stake` | integer | Yes | Stake for the bet in cents |
| `bets[].customBet` | boolean | Yes | `true` for same-game parlays (Custom Bet) |
| `bets[].selections` | array | Yes | Array of one or more selection objects |
| `bets[].selections[].id` | string | Yes | Unique ID for the selection |
| `bets[].selections[].eventId` | string | Yes | Sportradar event ID (e.g., "sr:match:12345678") |
| `bets[].selections[].odds` | integer | Yes | Odds in 1/10000 format (e.g., 20000 = 2.0000) |
| `bets[].selections[].banker` | boolean | No | `true` if this is a banker selection in a system bet |

#### Response

**Success (Accepted)**:
```json
{
  "operation": "ticket-placement-response",
  "content": {
    "type": "ticket-response",
    "ticketId": "TICKET_12345",
    "status": "accepted",
    "signature": "...",
    "betDetails": [...]
  },
  "correlationId": "...",
  "timestampUtc": 1698412800000,
  "version": "2.4"
}
```

**Success (Rejected)**:
```json
{
  "operation": "ticket-placement-response",
  "content": {
    "type": "ticket-response",
    "ticketId": "TICKET_12345",
    "status": "rejected",
    "reason": {
      "code": 1001,
      "message": "Stake exceeds limit"
    },
    "betDetails": [...]
  },
  "correlationId": "...",
  "timestampUtc": 1698412800000,
  "version": "2.4"
}
```

**Failure (API Error)**:
```json
{
  "error": "Validation failed",
  "details": "ticketId is required"
}
```

## 5. Bet Types

### 5.1. Single Bet

One bet with one selection.

```json
{
  "bets": [
    {
      "id": "bet_single",
      "stake": 1000,
      "customBet": false,
      "selections": [
        {
          "id": "sel_1",
          "eventId": "sr:match:1",
          "odds": 15000
        }
      ]
    }
  ]
}
```

### 5.2. Accumulator (Parlay/Combo)

One bet with multiple selections.

```json
{
  "bets": [
    {
      "id": "bet_combo",
      "stake": 500,
      "customBet": false,
      "selections": [
        {
          "id": "sel_1",
          "eventId": "sr:match:1",
          "odds": 15000
        },
        {
          "id": "sel_2",
          "eventId": "sr:match:2",
          "odds": 20000
        }
      ]
    }
  ]
}
```

### 5.3. System Bet

Multiple bets with different combinations of selections.

```json
{
  "bets": [
    {
      "id": "bet_1",
      "stake": 100,
      "selections": [{"id": "sel_1", ...}, {"id": "sel_2", ...}]
    },
    {
      "id": "bet_2",
      "stake": 100,
      "selections": [{"id": "sel_1", ...}, {"id": "sel_3", ...}]
    },
    {
      "id": "bet_3",
      "stake": 100,
      "selections": [{"id": "sel_2", ...}, {"id": "sel_3", ...}]
    }
  ]
}
```

### 5.4. Custom Bet (Same-Game Parlay)

Set `customBet: true` for the bet.

```json
{
  "bets": [
    {
      "id": "bet_sgp",
      "stake": 200,
      "customBet": true,
      "selections": [
        {
          "id": "sel_1",
          "eventId": "sr:match:1",
          "odds": 18000
        },
        {
          "id": "sel_2",
          "eventId": "sr:match:1",
          "odds": 25000
        }
      ]
    }
  ]
}
```

## 6. Error Handling

The API returns standard HTTP status codes:

| Code | Meaning | Description |
|:---:|:---|:---|
| 200 | OK | Request successful |
| 400 | Bad Request | Invalid request body or validation failed |
| 500 | Internal Server Error | Failed to send ticket to MTS or other server error |

## 7. cURL Examples

### Health Check

```bash
curl https://your-app.up.railway.app/health
```

### Place Single Bet

```bash
curl -X POST https://your-app.up.railway.app/api/tickets \
  -H "Content-Type: application/json" \
  -d '{
    "ticketId": "TICKET_SINGLE_001",
    "customerId": "customer_123",
    "currency": "USD",
    "totalStake": 1000,
    "testSource": true,
    "oddsChange": "any",
    "bets": [
      {
        "id": "bet_1",
        "stake": 1000,
        "customBet": false,
        "selections": [
          {
            "id": "sel_1",
            "eventId": "sr:match:12345678",
            "odds": 18000
          }
        ]
      }
    ]
  }'
```

## 8. References

- [Sportradar MTS Documentation](https://docs.sportradar.com/transaction30api/)
- [MTS Service GitHub Repository](https://github.com/gdsZyy/mts-service)

