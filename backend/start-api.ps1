Write-Host "Starting Backend API locally..." -ForegroundColor Green
$env:DATABASE_URL="postgresql://defi:defi123@127.0.0.1:5432/defi_dashboard?sslmode=disable"
$env:JWT_SECRET="your-super-secret-jwt-key-change-this-in-production"
$env:PORT="3000"
$env:LOG_LEVEL="debug"
$env:ALLOW_ORIGINS="http://localhost:5173,http://localhost:8080"
go run cmd/api/main.go