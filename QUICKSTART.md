# Quick Start Guide

## Prerequisites
- Go 1.21 or higher: https://golang.org/dl/
- PostgreSQL 12 or higher: https://www.postgresql.org/download/
- Git

## Installation Steps

### 1. Set up the database

**On Windows (PowerShell):**
```powershell
# Start PostgreSQL if not running
# Then connect to psql and create database
psql -U postgres

# In psql prompt:
CREATE DATABASE aviator_db;
\q
```

**On Mac/Linux:**
```bash
createdb aviator_db
```

### 2. Configure environment

```bash
cd d:\go project\avi_bd

# Copy example config
cp .env.example .env

# Edit .env with your database credentials
# Example:
# DB_HOST=localhost
# DB_PORT=5432
# DB_USER=postgres
# DB_PASSWORD=postgres
# DB_NAME=aviator_db
```

### 3. Install dependencies

```bash
go mod download
go mod tidy
```

### 4. Run the server

```bash
# Option 1: Run directly
go run cmd/server/main.go

# Option 2: Build then run
go build -o aviator_backend cmd/server/main.go
./aviator_backend

# Option 3: Use script (Linux/Mac)
bash scripts/run.sh
```

The server will start on `http://localhost:8080`

### 5. Test the API

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Register user
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

## Next Steps

1. **Connect WebSocket in your Flutter app:**
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log(data.type, data.payload);
};
```

2. **Use JWT tokens for authenticated requests:**
- Save the token from login response
- Add `Authorization: Bearer <token>` header to protected endpoints

3. **Place bets during countdown:**
- Listen for `countdown` WebSocket events
- POST to `/api/v1/bet` with bet details during waiting phase

4. **Cash out before crash:**
- Listen for `multiplier_update` events from WebSocket
- POST to `/api/v1/cashout` with current multiplier before crash event

## Troubleshooting

### "psql: command not found"
- PostgreSQL is not in your PATH
- On Windows: Add PostgreSQL bin folder to PATH
- On Mac: `brew install postgresql`
- On Linux: `apt-get install postgresql`

### "database connection refused"
```bash
# Check PostgreSQL is running
# Windows: Services > PostgreSQL
# Mac/Linux: psql -U postgres
```

### "port 8080 already in use"
```bash
# Change port in .env
APP_PORT=8081

# Or kill existing process
# Windows: netstat -ano | findstr :8080
# Mac/Linux: lsof -i :8080 | kill -9 <PID>
```

### "module not found" errors
```bash
go clean -modcache
go mod download
go mod tidy
go run cmd/server/main.go
```

## Development Tips

### Enable Live Reload with Air
```bash
go install github.com/cosmtrek/air@latest
air
```

### View Database
```bash
psql -U postgres -d aviator_db

\dt  # List tables
SELECT * FROM users;
SELECT * FROM wallets;
SELECT * FROM game_rounds;
SELECT * FROM bets;
```

### Check Server Logs
The server outputs logs to console and includes timestamps and error details.

## Performance Testing

### Load test with Apache Bench
```bash
# Test 1000 requests with 10 concurrent
ab -n 1000 -c 10 http://localhost:8080/api/v1/health
```

### Monitor WebSocket connections
The WebSocket hub connection count is logged.

## Security Checklist

Before production deployment:
- [ ] Change JWT_SECRET to strong random value
- [ ] Set up HTTPS/TLS
- [ ] Configure CORS with specific origins
- [ ] Enable database backups
- [ ] Set up rate limiting
- [ ] Add request logging
- [ ] Set up error monitoring
- [ ] Configure firewall rules
- [ ] Use environment-specific configs

## Docker Setup (Optional)

Create `Dockerfile`:
```dockerfile
FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o aviator_backend cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/aviator_backend .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./aviator_backend"]
```

Run with Docker:
```bash
docker build -t aviator_backend .
docker run -p 8080:8080 --env-file .env aviator_backend
```

---

For full documentation, see:
- [README.md](README.md) - Project overview and setup
- [docs/api.md](docs/api.md) - Complete API reference
