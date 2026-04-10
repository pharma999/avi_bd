.PHONY: help build run test clean deps lint fmt vet

# Variables
APP_NAME=aviator_backend
MAIN_PATH=cmd/server/main.go
BINARY_NAME=aviator_backend

help:
	@echo "Aviator Backend - Available targets:"
	@echo "  make build      - Build the application"
	@echo "  make run        - Run the application"
	@echo "  make test       - Run tests"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make deps       - Download dependencies"
	@echo "  make fmt        - Format code"
	@echo "  make vet        - Run go vet"
	@echo "  make lint       - Run linter (golangci-lint)"
	@echo "  make db-setup   - Setup PostgreSQL database"
	@echo "  make db-reset   - Reset database"

build:
	@echo "Building $(APP_NAME)..."
	@go build -o bin/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: bin/$(BINARY_NAME)"

run:
	@echo "Running $(APP_NAME)..."
	@go run $(MAIN_PATH)

test:
	@echo "Running tests..."
	@go test -v ./...

test-cover:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out
	@go clean

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

fmt:
	@echo "Formatting code..."
	@go fmt ./...

vet:
	@echo "Running go vet..."
	@go vet ./...

lint:
	@echo "Running linter..."
	@golangci-lint run ./...

db-setup:
	@echo "Setting up database..."
	@psql -U postgres -c "CREATE DATABASE aviator_db;" 2>/dev/null || true
	@echo "Database ready"

db-reset:
	@echo "Resetting database..."
	@psql -U postgres -c "DROP DATABASE IF EXISTS aviator_db;"
	@psql -U postgres -c "CREATE DATABASE aviator_db;"
	@echo "Database reset complete"

dev:
	@echo "Running with air (live reload)..."
	@air

all: clean deps fmt vet test build
	@echo "All done!"
