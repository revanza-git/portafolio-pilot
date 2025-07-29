# Local Development Setup Guide

This guide explains how to run the DeFi Portfolio application with **infrastructure and API in Docker** while **running the frontend locally** for optimal development experience with hot reloading.

## 🏗️ Architecture Overview

- **Frontend**: React/TypeScript with Vite (runs locally on `http://localhost:5173`)
- **Backend API**: Go with Fiber (runs in Docker, exposed on `http://localhost:3000`)
- **Worker**: Background service (runs in Docker)
- **Database**: PostgreSQL (runs in Docker, exposed on `localhost:5432`)
- **Cache**: Redis (runs in Docker, exposed on `localhost:6379`)
- **Admin**: pgAdmin (optional, runs in Docker on `http://localhost:5050`)

## 🚀 Quick Start

### Prerequisites

1. **Docker & Docker Compose** installed
2. **Node.js** (v18+) and **npm** installed
3. **Git** installed

### Step 1: Environment Setup

```bash
# Copy the local development environment file
cp .env.local .env

# Edit .env and add your API keys:
# - ALCHEMY_API_KEY=your-actual-key
# - INFURA_API_KEY=your-actual-key  
# - ETHERSCAN_API_KEY=your-actual-key
```

### Step 2: Start Backend Services (Docker)

```bash
# Start all backend services (PostgreSQL, Redis, API, Worker)
docker-compose -f docker-compose.local.yml up -d

# Or with logs visible:
docker-compose -f docker-compose.local.yml up

# Optional: Include pgAdmin for database management
docker-compose -f docker-compose.local.yml --profile tools up -d
```

### Step 3: Install Frontend Dependencies

```bash
# Install npm dependencies
npm install
```

### Step 4: Start Frontend (Local)

```bash
# Start Vite dev server with hot reloading
npm run dev
```

### Step 5: Access the Application

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:3000
- **pgAdmin** (if enabled): http://localhost:5050 (admin@defiportfolio.com / admin)

## 🔧 Development Workflow

### Backend Development
- Backend code changes require rebuilding Docker images:
  ```bash
  docker-compose -f docker-compose.local.yml up --build api worker
  ```

### Frontend Development
- Frontend changes are automatically hot-reloaded
- No restart needed for React/TypeScript changes

### Database Management
- **pgAdmin**: http://localhost:5050 (if started with `--profile tools`)
- **Direct connection**: `localhost:5432` with credentials from `.env`

## 📊 Service Status

Check if all services are running:

```bash
# View running containers
docker-compose -f docker-compose.local.yml ps

# View logs
docker-compose -f docker-compose.local.yml logs api
docker-compose -f docker-compose.local.yml logs worker
```

## 🛠️ Common Commands

### Docker Management
```bash
# Stop all services
docker-compose -f docker-compose.local.yml down

# Stop and remove volumes (clean database)
docker-compose -f docker-compose.local.yml down -v

# Rebuild and restart specific service
docker-compose -f docker-compose.local.yml up --build api

# View logs for specific service
docker-compose -f docker-compose.local.yml logs -f api
```

### Frontend Development
```bash
# Start development server
npm run dev

# Build for production
npm run build

# Run linting
npm run lint

# Run tests (if available)
npm test
```

### Database Operations
```bash
# Run migrations manually (if needed)
docker-compose -f docker-compose.local.yml run --rm migrate

# Connect to database directly
docker exec -it defi-dashboard-db psql -U defi -d defi_dashboard
```

## 🔍 Troubleshooting

### API Connection Issues
1. Ensure Docker services are running: `docker-compose -f docker-compose.local.yml ps`
2. Check API logs: `docker-compose -f docker-compose.local.yml logs api`
3. Verify `.env` has correct `VITE_API_BASE_URL=http://localhost:3000`

### Database Connection Issues
1. Check PostgreSQL is healthy: `docker-compose -f docker-compose.local.yml ps postgres`
2. Verify database credentials in `.env`
3. Check if migrations ran: `docker-compose -f docker-compose.local.yml logs migrate`

### Frontend Build Issues
1. Clear node_modules: `rm -rf node_modules && npm install`
2. Clear Vite cache: `rm -rf node_modules/.vite`
3. Check Node.js version: `node --version` (should be 18+)

### Port Conflicts
If ports are already in use, modify the ports in `.env`:
- `API_PORT=3001` (change from 3000)
- `WEB_PORT=5174` (change from 5173)

## 📋 Environment Variables

Key variables in `.env`:

| Variable | Purpose | Default |
|----------|---------|---------|
| `VITE_API_BASE_URL` | Frontend API endpoint | `http://localhost:3000` |
| `API_PORT` | Backend API port | `3000` |
| `WEB_PORT` | Frontend dev server port | `5173` |
| `DB_USER` | Database username | `defi` |
| `DB_PASSWORD` | Database password | `defi123` |
| `ALCHEMY_API_KEY` | Ethereum RPC provider | Required |
| `ETHERSCAN_API_KEY` | Transaction data | Required |

## 🚦 Service Health Checks

The Docker setup includes health checks for critical services:

- **PostgreSQL**: `pg_isready` check
- **Redis**: `redis-cli ping` check
- **API**: Depends on healthy database
- **Worker**: Depends on healthy database and Redis

## 🎯 Next Steps

1. Add your real API keys to `.env`
2. Start developing features in the frontend
3. Use pgAdmin to inspect database structure
4. Check backend logs for any API errors
5. Test wallet connections and DeFi integrations

## 📝 Notes

- The frontend automatically connects to the Dockerized backend
- Database and Redis are only accessible from Docker network + exposed ports
- Hot reloading works for frontend changes
- Backend changes require Docker rebuild
- Use `--profile tools` to include pgAdmin for database management