# MTS Service

A microservice for integrating with Sportradar's Managed Trading Services (MTS) for ticket placement.

## Features

- ✅ REST API for ticket placement
- ✅ WebSocket connection to MTS
- ✅ OAuth 2.0 authentication
- ✅ Automatic reconnection
- ✅ Health check endpoint
- ✅ Ready for Railway deployment

## Environment Variables

Required environment variables:

```bash
MTS_CLIENT_ID=your_client_id
MTS_CLIENT_SECRET=your_client_secret
MTS_BOOKMAKER_ID=your_bookmaker_id # Optional if UOF_ACCESS_TOKEN is provided
UOF_ACCESS_TOKEN=your_uof_token # Optional: Auto-fetch Bookmaker ID from whoami.xml
MTS_PRODUCTION=false  # Set to true for production
PORT=8080  # Optional, defaults to 8080
```

**Note**: You can either:
- Provide `MTS_BOOKMAKER_ID` directly, OR
- Provide `UOF_ACCESS_TOKEN` to auto-fetch the Bookmaker ID from `whoami.xml`

## API Endpoints

### Health Check

```
GET /health
```

Response:
```json
{
  "status": "healthy",
  "timestamp": 1698412800,
  "service": "mts-service"
}
```

### Place Ticket

```
POST /api/tickets
Content-Type: application/json
```

Request body:
```json
{
  "ticketId": "TICKET_123",
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

Response (accepted):
```json
{
  "operation": "ticket-placement-response",
  "content": {
    "type": "ticket-response",
    "ticketId": "TICKET_123",
    "status": "accepted",
    "signature": "...",
    "betDetails": [...]
  },
  "correlationId": "...",
  "timestampUtc": 1698412800000,
  "version": "2.4"
}
```

Response (rejected):
```json
{
  "operation": "ticket-placement-response",
  "content": {
    "type": "ticket-response",
    "ticketId": "TICKET_123",
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

## Field Descriptions

### Stake and Odds Format

- **Stake**: Amount in cents (e.g., `1000` = 10.00 EUR)
- **Odds**: Odds in 1/10000 format (e.g., `20000` = 2.0000)

### Odds Change Strategies

- `none`: Reject any odds changes
- `higher`: Accept only higher odds
- `any`: Accept any odds changes

### Bet Types

- **Single Bet**: One selection in a bet
- **Accumulator**: Multiple selections in a bet (all must win)
- **Custom Bet**: Same-game parlay (set `customBet: true`)

## Local Development

### Prerequisites

- Go 1.24+
- Sportradar MTS credentials

### Run Locally

```bash
# Set environment variables
export MTS_CLIENT_ID=your_client_id
export MTS_CLIENT_SECRET=your_client_secret
export MTS_BOOKMAKER_ID=your_bookmaker_id # Get from whoami.xml or Sportradar support
export MTS_PRODUCTION=false

# Run the service
go run cmd/server/main.go
```

### Build

```bash
go build -o mts-service ./cmd/server
./mts-service
```

## Railway Deployment

### Deploy to Railway

1. Push this repository to GitHub
2. Connect Railway to your GitHub repository
3. Set environment variables in Railway dashboard:
   - `MTS_CLIENT_ID`
   - `MTS_CLIENT_SECRET`
   - `MTS_BOOKMAKER_ID`
   - `MTS_PRODUCTION`
4. Railway will automatically detect the Dockerfile and deploy

### Railway Configuration

Railway will automatically:
- Detect the Dockerfile
- Build the Docker image
- Deploy the service
- Assign a public URL
- Set PORT environment variable

## Testing

### Example cURL Request

```bash
curl -X POST http://localhost:8080/api/tickets \
  -H "Content-Type: application/json" \
  -d '{
    "ticketId": "TICKET_'$(date +%s)'",
    "customerId": "customer_123",
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
  }'
```

## Architecture

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │ HTTP
       ↓
┌─────────────┐
│  REST API   │
│  (Gin/HTTP) │
└──────┬──────┘
       │
       ↓
┌─────────────┐
│ MTS Service │
│ (WebSocket) │
└──────┬──────┘
       │ WSS
       ↓
┌─────────────┐
│ Sportradar  │
│     MTS     │
└─────────────┘
```

## Project Structure

```
mts-service/
├── cmd/
│   └── server/
│       └── main.go          # Entry point
├── internal/
│   ├── api/
│   │   └── handlers.go      # HTTP handlers
│   ├── config/
│   │   └── config.go        # Configuration
│   ├── models/
│   │   └── ticket.go        # Data models
│   └── service/
│       └── mts.go           # MTS WebSocket client
├── Dockerfile               # Docker configuration
├── go.mod                   # Go dependencies
├── go.sum
└── README.md
```

## Error Handling

The service handles various error scenarios:

- **Connection errors**: Automatic reconnection with exponential backoff
- **Authentication errors**: Token refresh before expiry
- **Validation errors**: Clear error messages in API responses
- **Timeout errors**: 10-second timeout for MTS responses

## Monitoring

Monitor the following metrics:

- Health check status (`/health`)
- Connection status (in health check response)
- API response times
- Error rates

## Support

For issues or questions:
- Sportradar Support: support@betradar.com
- Sportradar Sales: sales@betradar.com
- Documentation: https://docs.sportradar.com/transaction30api/

## License

Proprietary - Sportradar MTS Integration

