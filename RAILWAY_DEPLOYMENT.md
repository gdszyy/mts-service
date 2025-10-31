# Railway Deployment Guide

## Quick Deploy to Railway

### Step 1: Connect GitHub Repository

1. Visit [Railway](https://railway.app/)
2. Click "New Project"
3. Select "Deploy from GitHub repo"
4. Choose `gdszyy/mts-service`

### Step 2: Configure Environment Variables

In Railway dashboard, add the following environment variables:

```
MTS_CLIENT_ID=your_client_id_here
MTS_CLIENT_SECRET=your_client_secret_here
MTS_BOOKMAKER_ID=your_bookmaker_id_here # Optional if UOF_ACCESS_TOKEN is provided
UOF_ACCESS_TOKEN=your_uof_token_here # Optional: Auto-fetch Bookmaker ID
MTS_PRODUCTION=false
```

**Note**: You can either provide `MTS_BOOKMAKER_ID` directly or `UOF_ACCESS_TOKEN` to auto-fetch it.

**Note**: Railway automatically sets the `PORT` variable, no need to configure it.

### Step 3: Deploy

Railway will automatically:
- Detect the Dockerfile
- Build the Docker image
- Deploy the service
- Assign a public URL (e.g., `https://mts-service-production.up.railway.app`)

### Step 4: Verify Deployment

Once deployed, test the health endpoint:

```bash
curl https://your-app.up.railway.app/health
```

Expected response:
```json
{
  "status": "healthy",
  "timestamp": 1698412800,
  "service": "mts-service"
}
```

## Testing the API

### Test Ticket Placement

```bash
curl -X POST https://your-app.up.railway.app/api/tickets \
  -H "Content-Type: application/json" \
  -d '{
    "ticketId": "TICKET_TEST_001",
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

## Monitoring

### View Logs

In Railway dashboard:
1. Click on your service
2. Go to "Deployments" tab
3. Click on the latest deployment
4. View real-time logs

### Check Service Status

Monitor the `/health` endpoint to ensure the service is running and connected to MTS.

## Troubleshooting

### Service Not Starting

Check logs for:
- Missing environment variables
- Invalid MTS credentials
- Connection issues

### Connection Issues

Ensure:
- MTS credentials are correct
- MTS_PRODUCTION is set correctly (false for integration, true for production)
- Railway has outbound internet access (it does by default)

### API Errors

Common issues:
- Invalid request format (check JSON structure)
- Missing required fields
- Invalid event IDs (must be from UOF)

## Updating the Service

To deploy updates:

1. Push changes to GitHub:
   ```bash
   git add .
   git commit -m "Your update message"
   git push
   ```

2. Railway will automatically detect the push and redeploy

## Production Checklist

Before going to production:

- [ ] Set `MTS_PRODUCTION=true`
- [ ] Use production MTS credentials
- [ ] Set up monitoring and alerts
- [ ] Configure custom domain (optional)
- [ ] Enable Railway's built-in metrics
- [ ] Set up log aggregation
- [ ] Test with real event IDs from UOF

## Support

- Railway Documentation: https://docs.railway.app/
- Railway Discord: https://discord.gg/railway
- MTS Support: support@betradar.com

