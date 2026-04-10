#!/bin/bash

# run.sh - Start the Aviator backend server

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | xargs)
else
    echo "Error: .env file not found. Please create it from .env.example"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null
then
    echo "Error: Go is not installed"
    exit 1
fi

# Download dependencies
echo "Downloading dependencies..."
go mod download
go mod tidy

# Run the server
echo "Starting Aviator backend server on port $APP_PORT..."
go run cmd/server/main.go
