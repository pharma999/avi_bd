# Quick Start Development Script
# Save as: dev-setup.ps1 (PowerShell)
# Run: .\dev-setup.ps1

$ProjectRoot = "D:\go project\avi_bd"
cd $ProjectRoot

Write-Host "=== Aviator Backend Development Setup ===" -ForegroundColor Cyan

# Menu
Write-Host "`nSelect an option:" -ForegroundColor Yellow
Write-Host "1. Start Backend Server"
Write-Host "2. Run Tests"
Write-Host "3. Format & Clean Code"
Write-Host "4. Check Code Quality"
Write-Host "5. Build Production Binary"
Write-Host "6. Reset Database"
Write-Host "7. View Database"
Write-Host "8. Test API"
Write-Host "9. Setup New Project"
Write-Host "0. Exit"

$choice = Read-Host "`nEnter your choice (0-9)"

switch($choice) {
    "1" {
        Write-Host "`n▶ Starting Backend Server..." -ForegroundColor Green
        go run cmd/server/main.go
    }
    
    "2" {
        Write-Host "`n▶ Running Tests..." -ForegroundColor Green
        go test -v ./...
    }
    
    "3" {
        Write-Host "`n▶ Formatting Code..." -ForegroundColor Green
        go fmt ./...
        Write-Host "✓ Format complete" -ForegroundColor Green
        
        Write-Host "`n▶ Cleaning Build Artifacts..." -ForegroundColor Green
        Remove-Item -Path "bin\" -Recurse -Force -ErrorAction SilentlyContinue
        Write-Host "✓ Clean complete" -ForegroundColor Green
    }
    
    "4" {
        Write-Host "`n▶ Running go vet..." -ForegroundColor Green
        go vet ./...
        Write-Host "✓ Vet complete" -ForegroundColor Green
    }
    
    "5" {
        Write-Host "`n▶ Building Production Binary..." -ForegroundColor Green
        go build -o bin/aviator_backend.exe cmd/server/main.go
        Write-Host "✓ Build complete: bin/aviator_backend.exe" -ForegroundColor Green
        Get-Item bin/aviator_backend.exe | Select-Object FullName, @{Name="Size(MB)";Expression={"{0:N2}" -f ($_.Length/1MB)}}
    }
    
    "6" {
        Write-Host "`n▶ Resetting Database..." -ForegroundColor Yellow
        $confirm = Read-Host "WARNING: This will delete all data. Continue? (yes/no)"
        
        if($confirm -eq "yes") {
            # Create new database
            psql -U postgres -c "DROP DATABASE IF EXISTS aviator_db;"
            psql -U postgres -c "CREATE DATABASE aviator_db;"
            Write-Host "✓ Database reset" -ForegroundColor Green
        } else {
            Write-Host "Cancelled" -ForegroundColor Red
        }
    }
    
    "7" {
        Write-Host "`n▶ Database Contents:" -ForegroundColor Green
        Write-Host "`nUsers:" -ForegroundColor Cyan
        psql -U postgres -d aviator_db -c "SELECT id, name, email, created_at FROM users LIMIT 5;"
        
        Write-Host "`nWallets:" -ForegroundColor Cyan
        psql -U postgres -d aviator_db -c "SELECT id, user_id, balance FROM wallets LIMIT 5;"
        
        Write-Host "`nBets:" -ForegroundColor Cyan
        psql -U postgres -d aviator_db -c "SELECT id, user_id, amount, result FROM bets LIMIT 5;"
    }
    
    "8" {
        Write-Host "`n=== Testing API Endpoints ===" -ForegroundColor Cyan
        
        $baseUrl = "http://localhost:8080"
        
        Write-Host "`n1. Testing Server Health..." -ForegroundColor Green
        $health = curl -s "$baseUrl/api/v1/health" | ConvertFrom-Json
        Write-Host "Status: $($health.data.status)" -ForegroundColor Green
        
        Write-Host "`n2. Register New User..." -ForegroundColor Green
        $registerResponse = curl -s -X POST "$baseUrl/api/v1/register" `
            -H "Content-Type: application/json" `
            -d "{`"name`":`"Test User`",`"email`":`"test@$(Get-Random).com`",`"password`":`"TestPass123!`"}" | ConvertFrom-Json
        
        if($registerResponse.data.token) {
            $token = $registerResponse.data.token
            Write-Host "✓ User Created: $($registerResponse.data.user.email)" -ForegroundColor Green
            Write-Host "Token: $($token.Substring(0, 20))..." -ForegroundColor Gray
            
            Write-Host "`n3. Get User Profile..." -ForegroundColor Green
            $profile = curl -s -H "Authorization: Bearer $token" "$baseUrl/api/v1/profile" | ConvertFrom-Json
            Write-Host "Name: $($profile.data.name)" -ForegroundColor Green
            Write-Host "Balance: $($profile.data.balance)" -ForegroundColor Green
        } else {
            Write-Host "Error: $($registerResponse.error.message)" -ForegroundColor Red
        }
    }
    
    "9" {
        Write-Host "`n▶ New Project Setup..." -ForegroundColor Cyan
        
        Write-Host "`nStep 1: Download dependencies"
        go mod download
        go mod tidy
        Write-Host "✓ Dependencies downloaded" -ForegroundColor Green
        
        Write-Host "`nStep 2: Create database"
        psql -U postgres -c "DROP DATABASE IF EXISTS aviator_db;"
        psql -U postgres -c "CREATE DATABASE aviator_db;"
        Write-Host "✓ Database created" -ForegroundColor Green
        
        Write-Host "`nStep 3: Setup environment"
        if(!(Test-Path .env)) {
            $envContent = @"
APP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=aviator_db
JWT_SECRET=dev-secret-change-in-production
GAME_TICK_MS=100
COUNTDOWN_SECONDS=5
"@
            Set-Content -Path .env -Value $envContent
            Write-Host "✓ .env created" -ForegroundColor Green
        }
        
        Write-Host "`nSetup complete! Run 'make run' to start the server." -ForegroundColor Green
    }
    
    "0" {
        Write-Host "Goodbye!" -ForegroundColor Cyan
        exit
    }
    
    default {
        Write-Host "Invalid option" -ForegroundColor Red
    }
}

Write-Host "`n" -ForegroundColor Gray
