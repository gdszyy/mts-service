# MTS Service

A comprehensive microservice for integrating with Sportradar's Managed Trading Services (MTS) for sports betting ticket placement.

## 🎯 Features

### Core Features
- ✅ **Complete REST API** for all bet types (single, accumulator, system, banker, presets)
- ✅ **WebSocket connection** to MTS with automatic reconnection
- ✅ **OAuth 2.0 authentication** with token refresh
- ✅ **Cashout support** (full and partial)
- ✅ **Health check** endpoint
- ✅ **Production-ready** with comprehensive error handling

### Supported Bet Types (14 Total)
- ✅ Single Bet
- ✅ Accumulator (Parlay)
- ✅ System Bet (custom combinations)
- ✅ Banker System Bet
- ✅ **10 Preset System Bets**:
  - Trixie, Patent, Yankee, Lucky 15, Lucky 31, Lucky 63
  - Super Yankee, Heinz, Super Heinz, Goliath

### API Features
- 🔄 **7 dedicated endpoints** for different bet types
- 📝 **Comprehensive validation** with detailed error messages
- 📊 **Structured logging** for debugging and monitoring
- 🧪 **Test scripts** for automated testing
- 📚 **Extensive documentation** with examples

## 📦 Quick Start

### Environment Variables

Required environment variables:

```bash
# MTS Credentials
MTS_CLIENT_ID=your_client_id
MTS_CLIENT_SECRET=your_client_secret

# Bookmaker Configuration (Option 1: Direct)
MTS_BOOKMAKER_ID=your_bookmaker_id
MTS_VIRTUAL_HOST=mts-api-ci.betradar.com

# Bookmaker Configuration (Option 2: Auto-fetch)
UOF_ACCESS_TOKEN=your_uof_token
UOF_API_BASE_URL=https://global.api.betradar.com

# Environment
MTS_PRODUCTION=false  # Set to true for production
PORT=8080  # Optional, defaults to 8080
```

**Note**: You can either:
- Provide `MTS_BOOKMAKER_ID` and `MTS_VIRTUAL_HOST` directly, OR
- Provide `UOF_ACCESS_TOKEN` to auto-fetch them from `whoami.xml`

### Installation

```bash
# Clone the repository
git clone https://github.com/gdszyy/mts-service.git
cd mts-service

# Install dependencies
go mod download

# Run the service
go run cmd/server/mts_main.go
```

### Docker (Optional)

```bash
# Build Docker image
docker build -t mts-service .

# Run container
docker run -p 8080:8080 \
  -e MTS_CLIENT_ID=your_client_id \
  -e MTS_CLIENT_SECRET=your_client_secret \
  -e UOF_ACCESS_TOKEN=your_token \
  mts-service
```

## 🚀 API Endpoints

### Overview

| Endpoint | Method | Description |
|:---|:---:|:---|
| `/health` | GET | Health check |
| `/api/bets/single` | POST | Place single bet |
| `/api/bets/accumulator` | POST | Place accumulator bet |
| `/api/bets/system` | POST | Place system bet |
| `/api/bets/banker-system` | POST | Place banker system bet |
| `/api/bets/preset` | POST | Place preset system bet |
| `/api/bets/multi` | POST | Place multi-bet ticket |
| `/api/cashout` | POST | Request cashout |

### Quick Examples

#### Single Bet

```bash
curl -X POST http://localhost:8080/api/bets/single \
  -H "Content-Type: application/json" \
  -d '{
    "ticketId": "single-001",
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
```

#### Trixie Bet

```bash
curl -X POST http://localhost:8080/api/bets/preset \
  -H "Content-Type: application/json" \
  -d '{
    "ticketId": "trixie-001",
    "type": "trixie",
    "selections": [
      {"productId": "3", "eventId": "sr:match:12345", "marketId": "1", "outcomeId": "1", "odds": 2.50},
      {"productId": "3", "eventId": "sr:match:12346", "marketId": "1", "outcomeId": "2", "odds": 1.80},
      {"productId": "3", "eventId": "sr:match:12347", "marketId": "1", "outcomeId": "1", "odds": 3.00}
    ],
    "stake": {
      "type": "cash",
      "currency": "EUR",
      "amount": 1.00,
      "mode": "unit"
    }
  }'
```

For more examples, see [EXAMPLES.md](./docs/technical/EXAMPLES.md).

## 📚 Documentation

### Product Documentation
- **[Frontend Interaction Design](./docs/product/Frontend_Interaction_Design.md)** - Complete frontend interaction logic including single/multi-bet modes, beginner/expert modes, and exception handling
- **[WebSocket Protocol](./docs/product/WebSocket_Protocol.md)** - Real-time betting communication protocol
- **[Product Flowcharts](./docs/product/flowcharts/)** - All product design flowcharts in Mermaid format

### Technical Documentation
- **[API Documentation](./docs/technical/API_DOCUMENTATION.md)** - Complete API reference with all endpoints
- **[Examples](./docs/technical/EXAMPLES.md)** - Practical examples in cURL, Python, and JavaScript
- **[MTS Specification Analysis](./docs/technical/MTS_Specification_Analysis.md)** - Analysis of MTS 3.0 API specifications

### Testing Documentation
- **[Test Cases](./docs/testing/test_cases.md)** - Comprehensive test case definitions
- **[Test Report](./docs/testing/final_test_report.md)** - Latest test execution report
- **[Testing Summary](./docs/testing/TESTING_SUMMARY.md)** - Overview of testing approach and results

### Deployment Documentation
- **[Railway Deployment](./docs/deployment/RAILWAY_DEPLOYMENT.md)** - Guide for deploying to Railway platform

## 🧪 Testing

### Automated Test Script

```bash
# Run all API tests
./scripts/test_api.sh

# Test against custom URL
BASE_URL=http://your-server:8080 ./scripts/test_api.sh
```

### Manual Testing

```bash
# Health check
curl http://localhost:8080/health

# Get API documentation
curl http://localhost:8080/
```

## 🏗️ Architecture

### Project Structure

```
mts-service/
├── cmd/
│   └── server/
│       └── mts_main.go          # Main entry point
├── internal/
│   ├── api/
│   │   ├── bet_handlers.go      # Bet endpoint handlers
│   │   ├── cashout_handlers.go  # Cashout handler
│   │   ├── helpers.go           # Validation & conversion
│   │   ├── logging.go           # Logging utilities
│   │   └── request_models.go    # API request/response models
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── models/
│   │   ├── ticket.go            # MTS ticket models
│   │   ├── cashout.go           # Cashout models
│   │   └── ticket_builder.go    # Ticket builder
│   └── service/
│       └── mts.go               # MTS WebSocket service
├── scripts/
│   └── test_api.sh              # Test script
├── API_DOCUMENTATION.md         # API docs
├── EXAMPLES.md                  # Usage examples
└── README.md                    # This file
```

### Data Flow

```
Client Request
    ↓
API Handler (validation)
    ↓
TicketBuilder (construct MTS request)
    ↓
MTSService (WebSocket)
    ↓
MTS Platform
    ↓
Response (ticket-reply)
    ↓
Acknowledgement (ticket-ack)
    ↓
Client Response
```

## 🔧 Configuration

### Limit ID

The service supports configuring a `LIMIT_ID` for risk management:

```bash
MTS_LIMIT_ID=4268
```

This will be included in all ticket requests sent to MTS.

### Production Mode

For production deployment:

```bash
MTS_PRODUCTION=true
MTS_VIRTUAL_HOST=mts-api.betradar.com  # Production endpoint
```

## 📊 Response Format

All API endpoints return a standardized response:

### Success Response

```json
{
  "success": true,
  "data": {
    "content": {
      "type": "ticket-reply",
      "ticketId": "single-001",
      "status": "accepted",
      "signature": "...",
      "betDetails": [...]
    },
    "correlationId": "corr-123456789",
    "timestampUtc": 1732612345000,
    "operation": "ticket-placement",
    "version": "3.0"
  }
}
```

### Error Response

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": 400,
    "message": "Validation failed",
    "details": "ticketId is required"
  }
}
```

## 🛠️ Development

### Prerequisites

- Go 1.21 or higher
- MTS credentials from Sportradar
- Access to MTS test environment

### Building

```bash
# Build binary
go build -o mts-service cmd/server/mts_main.go

# Run
./mts-service
```

### Adding New Bet Types

1. Add request model in `internal/api/request_models.go`
2. Add validation in `internal/api/helpers.go`
3. Add handler in `internal/api/bet_handlers.go`
4. Add route in `cmd/server/mts_main.go`
5. Update documentation

## 🐛 Troubleshooting

### Common Issues

**Connection Failed**
```
Error: not connected to MTS
```
- Check MTS credentials
- Verify network connectivity
- Check MTS_VIRTUAL_HOST setting

**Validation Error**
```
Error: ticketId is required
```
- Ensure all required fields are provided
- Check field formats (odds, amounts, etc.)

**MTS Rejection**
```
Status: rejected, Code: -401
```
- Event not found in MTS
- Check event ID format
- Verify event is available for betting

### Logs

The service provides detailed logging:

```
[2025-11-26 10:30:45] [SingleBet] Sending to MTS: TicketID=single-001
[2025-11-26 10:30:46] [SingleBet] ✓ Ticket ACCEPTED: TicketID=single-001, Status=accepted
```

## 📝 License

This project is licensed under the MIT License.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📞 Support

For issues and questions:
- GitHub Issues: https://github.com/gdszyy/mts-service/issues
- Sportradar MTS Documentation: https://docs.betradar.com/display/BD/MTS+-+Transaction+3.0

## 🔗 Related Resources

- [Sportradar MTS Documentation](https://docs.betradar.com/display/BD/MTS+-+Transaction+3.0)
- [Unified Odds Feed (UOF)](https://docs.betradar.com/display/BD/UOF+-+Unified+Odds+Feed)
- [Betradar API](https://docs.betradar.com/)

---

**Version**: 2.0.0  
**Last Updated**: 2025-11-26  
**Author**: gdszyy
