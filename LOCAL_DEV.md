# 🚀 Local Development Setup Guide

This guide explains how to run the DeFi Portfolio application locally for faster development and debugging.

## 📋 Prerequisites

- **Node.js** 18+ and npm
- **Go** 1.23+
- **Docker** (for PostgreSQL and Redis only)
- **Git**

## 🛠️ Quick Start

### 1. Set Up Environment Variables

```bash
# Copy the local development environment file
cp .env.local .env

# Edit .env and add your API keys:
# - ALCHEMY_API_KEY
# - INFURA_API_KEY
# - ETHERSCAN_API_KEY
# - VITE_WALLETCONNECT_PROJECT_ID
```

### 2. Start Infrastructure Services

```bash
# Start PostgreSQL and Redis in Docker
docker-compose -f docker-compose.dev.yml up -d postgres redis

# Run database migrations
docker-compose -f docker-compose.dev.yml run --rm migrate
```

### 3. Install Dependencies

```bash
# Install frontend dependencies
npm install

# Install backend dependencies
cd backend && go mod download && cd ..
```

### 4. Start Development Servers

Open **3 separate terminal windows**:

#### Terminal 1: Backend API
```bash
cd backend
go run cmd/api/main.go
```
The API will start at http://localhost:3000

#### Terminal 2: Frontend
```bash
npm run dev
```
The frontend will start at http://localhost:5173

#### Terminal 3: Worker (Optional)
```bash
cd backend
go run cmd/worker/main.go
```

## 🎯 Development URLs

- **Frontend**: http://localhost:5173 (Vite dev server with hot reload)
- **Backend API**: http://localhost:3000
- **PostgreSQL**: localhost:5432 (user: defi, password: defi123)
- **Redis**: localhost:6379
- **pgAdmin** (optional): http://localhost:5050

## 🔧 Common Commands

### Infrastructure Management
```bash
# Start infrastructure
docker-compose -f docker-compose.dev.yml up -d postgres redis

# Stop infrastructure
docker-compose -f docker-compose.dev.yml down

# View logs
docker-compose -f docker-compose.dev.yml logs -f

# Reset database
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up -d postgres redis
docker-compose -f docker-compose.dev.yml run --rm migrate
```

### Development Commands
```bash
# Frontend with hot reload
npm run dev

# Backend API with file watching (install air first)
cd backend && air

# Or manually restart on changes
cd backend && go run cmd/api/main.go

# Run tests
cd backend && go test ./...
npm test
```

## 🐛 Debugging Tips

### Frontend Debugging
1. Open Chrome DevTools (F12)
2. Check Console for errors
3. Check Network tab for API calls
4. React Developer Tools extension recommended

### Backend Debugging
1. API logs appear in terminal
2. Set LOG_LEVEL=debug in .env
3. Use VS Code debugger with Go extension
4. Check http://localhost:3000/health

### Common Issues

**Port already in use**
```bash
# Kill process on port 3000 (API)
lsof -ti:3000 | xargs kill -9

# Kill process on port 5173 (Frontend)
lsof -ti:5173 | xargs kill -9
```

**Database connection issues**
```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Check logs
docker logs defi-dashboard-db
```

**Authentication issues**
1. Clear browser localStorage
2. Check JWT_SECRET in .env matches between restarts
3. Verify CORS settings include your frontend URL

## 🔄 Hot Reload Setup

### Backend Hot Reload with Air
```bash
# Install air
go install github.com/air-verse/air@latest

# Create .air.toml in backend directory
cd backend
air init

# Run with hot reload
air
```

### Frontend Hot Reload
Vite provides hot reload out of the box:
```bash
npm run dev
```

## 📝 Environment Variables

Key variables for local development:

```env
# API runs on
API_PORT=3000

# Frontend runs on
VITE_API_BASE_URL=http://localhost:3000

# Database (local Docker)
DATABASE_URL=postgresql://defi:defi123@localhost:5432/defi_dashboard?sslmode=disable

# Redis (local Docker)
REDIS_URL=redis://localhost:6379

# CORS for local dev
ALLOW_ORIGINS=http://localhost:5173,http://localhost:8080
```

## 🎉 Benefits of Local Development

- ⚡ **Instant hot reload** - See changes immediately
- 🐛 **Better debugging** - Use your IDE's debugger
- 📊 **Direct logs** - See all console output
- 🚀 **Faster iteration** - No Docker rebuild delays
- 💻 **IDE integration** - Full IntelliSense and Go to Definition

Happy coding! 🚀