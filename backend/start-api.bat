@echo off
echo Starting Backend API locally...
set DATABASE_URL=postgresql://defi:defi123@localhost:5432/defi_dashboard?sslmode=disable
set JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
set PORT=3000
set LOG_LEVEL=debug
set ALLOW_ORIGINS=http://localhost:5173,http://localhost:8080
go run cmd/api/main.go