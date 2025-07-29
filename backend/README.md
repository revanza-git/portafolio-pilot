# DeFi Dashboard Backend

Production-ready Go Fiber API implementing the DeFi Dashboard OpenAPI specification.

## Tech Stack

- **Go 1.22+** - Programming language
- **Fiber v2** - Web framework
- **PostgreSQL** - Database
- **pgx + sqlc** - Database access
- **golang-migrate** - Database migrations
- **JWT** - Authentication
- **Docker & Docker Compose** - Containerization

## Architecture

```
backend/
├── cmd/
│   ├── api/           # API server entry point
│   └── worker/        # Background worker entry point
├── internal/          # Private application code
│   ├── config/        # Configuration management
│   ├── handlers/      # HTTP request handlers
│   ├── middleware/    # HTTP middleware
│   ├── models/        # Domain models
│   ├── repos/         # Data repositories
│   ├── router/        # Route definitions
│   ├── services/      # Business logic
│   └── jobs/          # Background job implementations
├── pkg/               # Public packages
│   ├── errors/        # Error handling
│   ├── logger/        # Logging utilities
│   ├── utils/         # Helper functions
│   └── external/      # External API clients
├── db/                # Database files
│   ├── migrations/    # SQL migrations
│   └── queries/       # SQL queries for sqlc
└── scripts/           # Utility scripts
```

## Worker Service

The backend includes a worker service for running scheduled background jobs:

### Jobs

1. **Price Refresh Job** (every 10 minutes)
   - Fetches latest token prices from CoinGecko
   - Updates yield pool APR/TVL data from DefiLlama
   - Stores historical price data

2. **Alert Evaluator Job** (every 5 minutes)
   - Evaluates active user alerts
   - Checks price thresholds, APR changes, large transfers
   - Triggers notifications when conditions are met

### Running the Worker

```bash
# With Docker Compose (includes worker service)
docker-compose up --build

# Manually
make run-worker

# Development with hot reload
make dev-worker
```

## Quick Start

### Using Docker Compose (Recommended)

1. Clone the repository and navigate to the backend directory:
```bash
cd backend
```

2. Copy the environment variables:
```bash
cp .env.example .env
```

3. Start all services:
```bash
docker-compose up --build
```

This will:
- Start PostgreSQL database
- Run database migrations
- Build and start the API server

The API will be available at `http://localhost:3000`

### Manual Setup

1. Install dependencies:
```bash
go mod download
```

2. Install development tools:
```bash
make install-tools
```

3. Start PostgreSQL:
```bash
docker run -d \
  --name defi-postgres \
  -e POSTGRES_USER=defi \
  -e POSTGRES_PASSWORD=defi123 \
  -e POSTGRES_DB=defi_dashboard \
  -p 5432:5432 \
  postgres:16-alpine
```

4. Run migrations:
```bash
make migrate-up
```

5. Generate code:
```bash
make generate
```

6. Run the server:
```bash
make run
```

## Development

### Available Make Commands

```bash
make help          # Show all available commands
make run           # Run the API server
make run-worker    # Run the worker service
make dev           # Run API with hot reload (requires air)
make dev-worker    # Run worker with hot reload
make test          # Run tests
make lint          # Run linter
make migrate-up    # Run database migrations
make migrate-down  # Rollback migrations
make seed          # Seed the database
make generate      # Generate code (sqlc, OpenAPI)
make docker-up     # Start all services with Docker
make docker-down   # Stop all services
make build         # Build both API and worker binaries
```

### Environment Variables

See `.env.example` for all available configuration options. Key variables:

- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret key for JWT tokens
- `PORT` - Server port (default: 3000)
- `LOG_LEVEL` - Logging level (debug, info, warn, error)

### API Documentation

The API implements the OpenAPI specification located at `../spec/openapi.yaml`.

Key endpoints:

- `GET /health` - Health check
- `POST /api/v1/auth/verify` - SIWE authentication
- `GET /api/v1/portfolio/{address}/balances` - Get token balances
- `GET /api/v1/transactions/{address}` - Get transaction history
- `GET /api/v1/yield/pools` - Get available yield pools

### Authentication

The API uses Sign-In with Ethereum (SIWE) for authentication:

1. Get a nonce: `GET /api/v1/auth/nonce?address=0x...`
2. Sign the message with your wallet
3. Verify signature: `POST /api/v1/auth/verify`
4. Use the returned JWT token in the Authorization header: `Bearer <token>`

### Database

The project uses PostgreSQL with migrations managed by golang-migrate.

To create a new migration:
```bash
make migrate-create name=add_new_table
```

### Testing

Run unit tests:
```bash
make test
```

Run tests with coverage:
```bash
make test-coverage
```

Example test included in `tests/` directory.

## Production Deployment

1. Build the Docker image:
```bash
docker build -t defi-dashboard-api .
```

2. Run with production configuration:
```bash
docker run -d \
  --name defi-api \
  -e DATABASE_URL=$DATABASE_URL \
  -e JWT_SECRET=$JWT_SECRET \
  -e LOG_LEVEL=info \
  -p 3000:3000 \
  defi-dashboard-api
```

### Security Considerations

- Always use HTTPS in production
- Set strong JWT_SECRET
- Enable rate limiting
- Use environment-specific database credentials
- Regularly update dependencies
- Enable CORS only for trusted origins

## Development Status

Currently implemented with mock data:
- ✅ Authentication flow (SIWE)
- ✅ Portfolio balances
- ✅ Transaction history
- ✅ Token approvals
- ✅ Yield pools
- ✅ Bridge routes
- ✅ Analytics/PnL
- ✅ Alerts CRUD
- ✅ Watchlists CRUD

TODO for production:
- [ ] Integrate with actual blockchain nodes (Alchemy/Infura)
- [ ] Implement real SIWE verification
- [ ] Connect to price feeds
- [ ] Add caching layer (Redis)
- [ ] Implement WebSocket for real-time updates
- [ ] Add comprehensive test coverage
- [ ] Set up CI/CD pipeline

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linter
5. Submit a pull request

## License

MIT