# DeFi Portfolio Management System

A professional full-stack DeFi portfolio management dashboard built with React, TypeScript, Go, and modern web technologies.

## 🚀 Features

- **Portfolio Tracking**: Monitor token balances and portfolio value across multiple networks
- **Transaction History**: View and filter your DeFi activity with detailed transaction data
- **Security Management**: Review and revoke token approvals to protect your assets
- **Token Swaps**: Exchange tokens using integrated DEX aggregators (UI ready, integration coming soon)
- **Responsive Design**: Beautiful dark/light theme with professional DeFi styling

## 🛠 Tech Stack

### Frontend
- **React 18** with TypeScript
- **Vite** for fast development and building
- **TailwindCSS** + **shadcn/ui** for styling
- **Zustand** for state management
- **TanStack Query** for data fetching
- **React Router** for navigation
- **Recharts** for data visualization

### Backend
- **Go 1.22+** with Fiber framework
- **PostgreSQL** database with migrations
- **Docker** containerization with Docker Compose
- **REST API** with comprehensive endpoints
- **Background worker** for data fetching
- **Redis** for caching and sessions

## 🏃‍♂️ Quick Start

### Prerequisites
- **Docker & Docker Compose** (recommended)
- **OR** Node.js 18+ and Go 1.22+ for local development
- Git

## 🐳 Docker Setup (Recommended) - Infrastructure + API in Docker, Frontend Local

### Hybrid Development Setup (Current Implementation)

This setup runs all backend services and infrastructure in Docker while allowing the frontend to run locally for optimal development experience with hot reloading.

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd defip
   ```

2. **Environment setup**
   ```bash
   # The .env file is already configured for local development
   # Check .env.example for reference values
   ```
   Key environment variables:
   - `VITE_API_BASE_URL=http://localhost:3000` (Frontend API endpoint)
   - `API_PORT=3000` (Backend API port)
   - `ALCHEMY_API_KEY`, `INFURA_API_KEY`, `ETHERSCAN_API_KEY` (Add your real API keys)

3. **Start backend services (Docker)**
   ```bash
   # Step 1: Start infrastructure services first
   docker-compose -f docker-compose.local.yml up -d postgres redis
   
   # Step 2: Run database migrations
   docker-compose -f docker-compose.local.yml up migrate
   
   # Step 3: Start API and worker services
   docker-compose -f docker-compose.local.yml up -d api worker
   
   # Optional: Start pgAdmin for database management
   docker-compose -f docker-compose.local.yml --profile tools up -d pgadmin
   ```

4. **Start frontend (Local)**
   ```bash
   # Install dependencies
   npm install
   
   # Start development server with hot reloading
   npm run dev
   ```

5. **Access the applications**
   - **Frontend**: http://localhost:8080 (Vite dev server)
   - **Backend API**: http://localhost:3000
   - **Database**: localhost:5432 (defi/defi123)
   - **Redis**: localhost:6379
   - **pgAdmin** (optional): http://localhost:5050 (admin@defiportfolio.com / admin)

### Service Management

```bash
# View logs for all services
docker-compose -f docker-compose.local.yml logs -f

# View logs for specific service
docker-compose -f docker-compose.local.yml logs -f api

# Stop all services
docker-compose -f docker-compose.local.yml down

# Stop and remove volumes (clean database)
docker-compose -f docker-compose.local.yml down -v

# Restart specific service
docker-compose -f docker-compose.local.yml restart api

# View running services
docker-compose -f docker-compose.local.yml ps

# Rebuild and restart after code changes
docker-compose -f docker-compose.local.yml up --build -d api worker
```

### Available Endpoints

- `GET /health` - API health check
- `GET /api/v1/auth/nonce` - Get authentication nonce
- `POST /api/v1/auth/verify` - Verify wallet signature
- `GET /api/v1/portfolio/:address/balances` - Get portfolio balances
- `GET /api/v1/transactions/:address` - Get transaction history

## 💻 Local Development

### Frontend Only

1. **Install dependencies**
   ```bash
   npm install
   ```

2. **Start development server**
   ```bash
   npm run dev
   ```

3. **Open in browser**
   ```
   http://localhost:5173
   ```

### Full Stack Local

1. **Start backend services**
   ```bash
   # Start database and infrastructure
   docker-compose up postgres redis migrate -d
   ```

2. **Run backend locally**
   ```bash
   cd backend
   cp .env.example .env
   # Edit .env with your configuration
   go run cmd/api/main.go
   ```

3. **Run frontend locally**
   ```bash
   npm run dev
   ```

## 📁 Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── ui/             # shadcn/ui components
│   ├── dashboard/      # Dashboard-specific components
│   ├── wallet/         # Wallet connection components
│   ├── transactions/   # Transaction management
│   ├── approvals/      # Token approval management
│   └── swap/           # Swap interface components
├── hooks/              # Custom React hooks
├── lib/                # Utilities and configurations
├── pages/              # Main application pages
├── stores/             # Zustand state management
├── types/              # TypeScript type definitions
└── assets/             # Static assets
```

## 🎨 Design System

The application uses a comprehensive design system with:
- Professional DeFi-focused color palette
- Custom gradients and shadows
- Semantic color tokens for profit/loss indicators
- Responsive typography and spacing
- Dark/light theme support

## 🔗 Wallet Integration

**Fully Implemented:**
- ✅ Real wagmi/viem integration with production RPC endpoints
- ✅ MetaMask, WalletConnect, and injected wallet support
- ✅ Multi-chain support (Ethereum, Polygon, Arbitrum, Optimism)
- ✅ SIWE (Sign-In With Ethereum) authentication
- ✅ Real ERC20 contract interactions and transaction signing
- ✅ Wallet state persistence and chain detection

## 📊 Data & APIs

### ✅ Implemented Integrations
- **Blockchain RPCs**: Alchemy, Infura APIs with production endpoints
- **Price Data**: CoinGecko, DefiLlama APIs (backend configured)
- **Transaction Data**: Etherscan API integration
- **Swap Quotes**: 0x Protocol, 1inch Network APIs (backend ready)
- **Bridge Services**: LiFi, Socket APIs for cross-chain functionality

### 🔄 Integration Status
- **Backend APIs**: Fully configured and implemented
- **Frontend**: Uses mock data for development, real API integration available
- **Database**: Complete schema for wallets, transactions, balances, positions

## 🔐 Security Features

- Token approval monitoring and revocation
- Security warnings for unlimited approvals
- Safe transaction confirmation dialogs
- Best practices for DeFi security

## 📱 Pages

1. **Landing Page** (`/`) - Hero section with wallet connection
2. **Dashboard** (`/dashboard`) - Portfolio overview with charts and balances
3. **Transactions** (`/transactions`) - Complete transaction history with filters
4. **Approvals** (`/approvals`) - Token approval management and revocation
5. **Swap** (`/swap`) - Token swap interface (UI complete, execution coming soon)

## 🚀 Production Deployment

### Prerequisites for Production
- **Docker & Docker Compose** (recommended for containerized deployment)
- **PostgreSQL 15+** (for database)
- **Redis 7+** (for caching and sessions)
- **SSL Certificate** (for HTTPS)
- **Domain name** with proper DNS configuration
- **CDN** (optional but recommended for static assets)

### Environment Configuration

#### 1. Production Environment Variables
Copy `.env.example` to `.env` and configure production values:

```bash
# Essential Production Variables
NODE_ENV=production
JWT_SECRET=your-ultra-secure-jwt-secret-min-32-chars
DATABASE_URL=postgresql://username:password@db-host:5432/database_name
REDIS_URL=redis://redis-host:6379

# API Configuration
API_PORT=3000
WEB_PORT=8080
ALLOW_ORIGINS=https://yourdomain.com

# Required API Keys
ALCHEMY_API_KEY=your-production-alchemy-key
INFURA_API_KEY=your-production-infura-key
ETHERSCAN_API_KEY=your-production-etherscan-key
VITE_WALLETCONNECT_PROJECT_ID=your-walletconnect-project-id

# Security Settings
ENABLE_RATE_LIMIT=true
LOG_LEVEL=warn
```

#### 2. Security Checklist
- [ ] Generate a strong JWT secret (minimum 32 characters)
- [ ] Use environment variables for all secrets
- [ ] Enable HTTPS/TLS encryption
- [ ] Configure proper CORS origins
- [ ] Enable rate limiting
- [ ] Set up database connection pooling
- [ ] Configure secure headers
- [ ] Enable audit logging

### Docker Production Deployment

#### 1. Production Docker Compose
Create `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile.web
    ports:
      - "8080:8080"
    environment:
      - NODE_ENV=production
    env_file:
      - .env
    depends_on:
      - api
      - postgres
      - redis
    restart: unless-stopped

  api:
    build:
      context: ./backend
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    env_file:
      - .env
    depends_on:
      - postgres
      - redis
    restart: unless-stopped
    volumes:
      - api-logs:/app/logs

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=${DB_NAME}
      - POSTGRES_USER=${DB_USER}
      - POSTGRES_PASSWORD=${DB_PASSWORD}
    volumes:
      - postgres-data:/var/lib/postgresql/data
      - ./backend/db/migrations:/docker-entrypoint-initdb.d
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes
    volumes:
      - redis-data:/data
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.prod.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - app
    restart: unless-stopped

volumes:
  postgres-data:
  redis-data:
  api-logs:
```

#### 2. Deploy to Production
```bash
# Build and deploy
docker-compose -f docker-compose.prod.yml up -d --build

# Monitor deployment
docker-compose -f docker-compose.prod.yml logs -f

# Check service health
docker-compose -f docker-compose.prod.yml ps
```

### Database Management

#### 1. Database Migrations
```bash
# Run migrations in production
docker-compose exec api migrate -path=/app/db/migrations -database="$DATABASE_URL" up

# Backup database before migrations
docker-compose exec postgres pg_dump -U $DB_USER $DB_NAME > backup_$(date +%Y%m%d_%H%M%S).sql
```

#### 2. Database Backup Strategy
```bash
# Daily backup script
#!/bin/bash
BACKUP_DIR="/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
docker-compose exec postgres pg_dump -U $DB_USER $DB_NAME | gzip > $BACKUP_DIR/db_backup_$TIMESTAMP.sql.gz

# Keep only last 30 days of backups
find $BACKUP_DIR -name "db_backup_*.sql.gz" -mtime +30 -delete
```

### Monitoring & Logging

#### 1. Application Monitoring
```bash
# Health check endpoints
curl https://yourdomain.com/health
curl https://yourdomain.com/api/v1/health

# Monitor logs
docker-compose logs -f api
docker-compose logs -f app
```

#### 2. Log Management
Configure structured logging in production:
```yaml
# In docker-compose.prod.yml
logging:
  driver: "json-file"
  options:
    max-size: "10m"
    max-file: "3"
```

### CI/CD Pipeline

#### 1. GitHub Actions Example
Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy to Production

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '18'
          cache: 'npm'
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      
      - name: Run tests
        run: |
          npm ci
          npm run test
          cd backend && go test ./...
      
      - name: Build and deploy
        run: |
          docker-compose -f docker-compose.prod.yml build
          docker-compose -f docker-compose.prod.yml up -d
```

### Performance Optimization

#### 1. Frontend Optimization
- Enable build optimizations in `vite.config.ts`
- Configure CDN for static assets
- Enable gzip compression
- Implement proper caching headers

#### 2. Backend Optimization
- Configure database connection pooling
- Enable Redis caching
- Implement rate limiting
- Set up database indexes

### Security Hardening

#### 1. Network Security
- Configure firewall rules
- Use VPC/private networks
- Enable DDoS protection
- Implement rate limiting

#### 2. Application Security
- Regular security audits
- Dependency vulnerability scanning
- Secure API endpoints
- Input validation and sanitization

## 🔧 Recent Updates & Fixes (January 2025)

### 🎯 Authentication System Overhaul
- **✅ SIWE (Sign-In With Ethereum) Implementation**: Complete authentication flow with proper message formatting
- **✅ Wagmi v2 Integration**: Fixed connector compatibility issues and message signing
- **✅ Provider Context Optimization**: Resolved React context ordering for wagmi/auth integration
- **✅ Error Handling Enhancement**: Comprehensive error handling for wallet connection and signing flows

### 🐛 Technical Fixes
- **Fixed**: `TypeError: connection.connector.getChainId is not a function` error
- **Fixed**: "Invalid SIWE message format" backend validation error
- **Fixed**: React context provider ordering issues
- **Fixed**: Wagmi signMessageAsync parameter handling
- **Improved**: Wallet connection state management and persistence

### 🏗️ Development Experience
- **Enhanced**: Step-by-step Docker service startup to prevent timeouts
- **Added**: Comprehensive error logging and debugging
- **Improved**: Development workflow with proper service dependencies
- **Updated**: Documentation with actual working setup instructions

### 🔒 Security Improvements
- **Implemented**: Proper SIWE message domain validation
- **Enhanced**: JWT token handling and storage
- **Added**: Rate limiting protection for authentication endpoints
- **Improved**: Error messages for better UX without exposing sensitive details

## 🚧 Development Status

### ✅ Completed
- [x] Professional UI/UX design system
- [x] Responsive layout and navigation
- [x] Portfolio dashboard with charts
- [x] Transaction history and filtering
- [x] Token approval management UI
- [x] Swap interface mockup
- [x] State management setup
- [x] Theme switching (dark/light)
- [x] **Backend API development (Go + Fiber)**
- [x] **Database schema and migrations**
- [x] **Docker containerization with Docker Compose**
- [x] **Full-stack integration**
- [x] **Authentication system (JWT + Web3)**
- [x] **Production deployment configuration**
- [x] **Security hardening and monitoring setup**
- [x] **Real wallet integration (wagmi/viem with SIWE authentication)**
- [x] **Blockchain API integration (Alchemy, Infura, Etherscan)**
- [x] **DeFi protocol APIs (0x, 1inch, LiFi, Socket)**
- [x] **Complete SIWE authentication flow with proper message formatting**
- [x] **Wagmi v2 connector compatibility and error handling**
- [x] **Production-ready Docker development setup**

### 🔄 In Progress
- [ ] Frontend-Backend API integration (replacing mock data with real blockchain data)
- [ ] Real-time WebSocket price feeds implementation
- [ ] Transaction broadcasting and execution flows
- [ ] Cross-chain balance aggregation

### 📋 Planned
- [ ] Real-time WebSocket price feeds
- [ ] Advanced portfolio analytics and PnL tracking
- [ ] Mobile app (React Native)
- [ ] Push notification system
- [ ] Load balancing and auto-scaling
- [ ] Advanced DeFi yield farming strategies
- [ ] Cross-chain portfolio aggregation

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 👨‍💻 Author

**Revanza Raytama**
- 📧 Email: [revanza.raytama@gmail.com](mailto:revanza.raytama@gmail.com)
- 🌐 Website: [https://revanza.vercel.app](https://revanza.vercel.app)
- 💼 Full-Stack Developer & DeFi Enthusiast

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

**Copyright © 2025 Revanza Raytama**

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

## 🙏 Acknowledgments

- [shadcn/ui](https://ui.shadcn.com/) for the excellent component library
- [Lucide React](https://lucide.dev/) for beautiful icons
- [TailwindCSS](https://tailwindcss.com/) for the utility-first CSS framework
- DeFi community for inspiration and best practices

---

**Note**: This is a full-stack DeFi portfolio management application with a complete backend API and database integration. The frontend includes professional UI/UX with comprehensive DeFi features. Real wallet integration and blockchain data are the next development priorities.