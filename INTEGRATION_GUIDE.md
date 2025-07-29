# DeFi Portfolio Full-Stack Integration Guide

This guide documents the complete integration between the Go Fiber backend and Next.js frontend, with all mock data replaced by real API endpoints.

## 🎯 Integration Overview

### What's Been Integrated

✅ **Authentication (SIWE)**
- Frontend uses API client for nonce generation and signature verification
- JWT tokens stored in localStorage and attached to all requests
- Auth context manages user state across the app

✅ **Portfolio Management**
- Real-time balance fetching from backend
- Transaction history with pagination
- Token approval management with revoke functionality

✅ **Yield Farming**
- Pool discovery and APY tracking
- Position management per wallet
- Rewards claiming (stubbed for chain interaction)

✅ **Alerts System**
- CRUD operations for price/APY/allowance alerts
- Alert history tracking
- Email/webhook notifications (backend ready)

✅ **Bridge & Swap**
- Quote aggregation from multiple providers
- Route optimization
- Transaction building (execution stubbed)

✅ **Analytics**
- P&L calculation and export
- Portfolio history charting
- CSV export functionality

✅ **Admin Features**
- Feature flags management
- System banner controls
- User management

## 📦 API Client Package

Generated TypeScript client located at `packages/api-client/`:

```typescript
import { createDefiAPIClient } from '@defip/api-client';

const apiClient = createDefiAPIClient({
  baseURL: process.env.VITE_API_BASE_URL,
  timeout: 30000
});
```

### Key Client Methods

```typescript
// Authentication
await apiClient.getNonce(address);
await apiClient.verifySiwe(message, signature);
await apiClient.getCurrentUser();

// Portfolio
await apiClient.getBalances(address, { chainId, hideSmall });
await apiClient.getPortfolioHistory(address, { period, interval });

// Transactions
await apiClient.getTransactions(address, { page, limit });
await apiClient.getApprovals(address);
await apiClient.revokeApproval(address, tokenAddress);

// Yield
await apiClient.getYieldPools({ chainId, protocol });
await apiClient.getYieldPositions(address);

// Alerts
await apiClient.getAlerts();
await apiClient.createAlert(alertData);
await apiClient.updateAlert(alertId, updates);
```

## 🔧 Development Workflow

### 1. Start Full Stack
```bash
# Using Make
make dev

# OR using Docker Compose directly
docker-compose --profile dev up --build
```

### 2. Access Services
- Frontend: http://localhost:8080
- Backend API: http://localhost:3000
- API Documentation: http://localhost:3000/swagger
- pgAdmin: http://localhost:5050 (dev profile only)

### 3. Run Tests
```bash
# All tests
make test

# Specific test suites
make test-backend    # Go tests
make test-frontend   # React tests
make test-e2e        # Playwright tests
```

## 🔄 Data Flow Architecture

```mermaid
graph TB
    A[React Frontend] -->|API Client| B[Go Fiber Backend]
    B --> C[PostgreSQL]
    B --> D[Redis Cache]
    B --> E[External APIs]
    E --> F[CoinGecko]
    E --> G[DefiLlama]
    E --> H[0x/1inch]
    E --> I[LiFi/Socket]
    J[Background Worker] --> C
    J --> E
```

## 🚀 Production Deployment

### Environment Variables

**Frontend (.env)**:
```env
VITE_API_BASE_URL=https://api.your-domain.com
VITE_WALLETCONNECT_PROJECT_ID=your_project_id
```

**Backend (.env)**:
```env
DATABASE_URL=postgresql://user:pass@host:5432/db
JWT_SECRET=your-production-secret
ALLOW_ORIGINS=https://your-domain.com
# Add all API keys from .env.example
```

### Docker Production Build

```bash
# Build production images
docker-compose -f docker-compose.prod.yml build

# Deploy
docker-compose -f docker-compose.prod.yml up -d
```

## 🧪 Testing Strategy

### Unit Tests
- Backend: `go test ./...` with mocked dependencies
- Frontend: Jest/Vitest for components and hooks

### Integration Tests
- API endpoints tested with real database
- Frontend API client tested against running backend

### E2E Tests
- Playwright tests cover critical user flows
- Run against dockerized full stack

## 📝 Adding New Features

### 1. Update OpenAPI Spec
```yaml
# spec/openapi.yaml
paths:
  /api/v1/new-feature:
    get:
      tags:
        - feature
      summary: New feature endpoint
      responses:
        '200':
          description: Success
```

### 2. Implement Backend
```go
// handlers/feature_handler.go
func (h *FeatureHandler) GetFeature(c *fiber.Ctx) error {
    // Implementation
}

// router/router.go
feature.Get("/new-feature", featureHandler.GetFeature)
```

### 3. Generate Client Types
```bash
cd packages/api-client
npm run generate
npm run build
```

### 4. Use in Frontend
```typescript
const data = await apiClient.getNewFeature();
```

## 🐛 Troubleshooting

### Common Issues

1. **CORS Errors**
   - Check `ALLOW_ORIGINS` in backend .env
   - Ensure frontend URL is whitelisted

2. **Auth Token Issues**
   - Clear localStorage and re-authenticate
   - Check JWT_SECRET matches between restarts

3. **Database Connection**
   - Ensure PostgreSQL is healthy: `docker-compose ps`
   - Check migration status: `make migrate`

4. **API Client Type Errors**
   - Regenerate types: `make build-api-client`
   - Ensure OpenAPI spec is valid

## 📚 Additional Resources

- [Backend README](backend/README.md) - Go API documentation
- [OpenAPI Spec](spec/openapi.yaml) - Complete API specification
- [Database Schema](backend/db/migrations/) - SQL migrations
- [E2E Tests](tests/e2e/) - Test examples