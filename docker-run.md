# Docker Setup Guide for DeFi Portfolio

## ✅ Prerequisites Completed
- [x] Environment variables configured (.env file)
- [x] Docker images built successfully
- [x] Nginx configuration ready
- [x] Frontend Docker setup with memory optimization

## 🚀 Step-by-Step Startup

### Option 1: Full Stack (Recommended)
```bash
# Start all services
docker-compose up

# Or in detached mode
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs
```

### Option 2: Manual Step-by-Step (For Troubleshooting)
```bash
# 1. Start infrastructure services
docker-compose up postgres redis -d

# 2. Wait for database to be ready, then run migrations
docker-compose up migrate

# 3. Start backend services
docker-compose up api worker -d

# 4. Start frontend
docker-compose up web -d
```

### Option 3: Individual Service Testing
```bash
# Database only
docker-compose up postgres -d

# Frontend only (for development)
docker-compose up web -d

# Backend API only
docker-compose up postgres redis migrate api -d
```

## 📊 Service Ports
- **Frontend (Web)**: http://localhost:8080
- **Backend API**: http://localhost:3000  
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379
- **pgAdmin** (dev): http://localhost:5050 (run with `--profile dev`)

## 🔧 Environment Variables
Already configured in `.env`:
- Database: postgres service (defi_dashboard)
- Redis: redis service
- API URL: Container-to-container communication configured

## 🐛 Troubleshooting Commands
```bash
# Check container status
docker-compose ps

# View logs for specific service
docker-compose logs [service-name]

# Restart a specific service
docker-compose restart [service-name]

# Rebuild and restart
docker-compose up --build [service-name]

# Clean restart
docker-compose down && docker-compose up
```

## 🎯 Quick Start
```bash
# Ensure you're in the project directory
cd c:\project\defip\defip

# Start everything
docker-compose up -d

# Check if all services are running
docker-compose ps

# View logs to verify startup
docker-compose logs
```

Once running, access the application at: **http://localhost:8080**