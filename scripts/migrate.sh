#!/bin/bash

# migrate.sh - Run database migrations

echo "Running database migrations..."

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | xargs)
else
    echo "Error: .env file not found"
    exit 1
fi

# Check if PostgreSQL is installed
if ! command -v psql &> /dev/null
then
    echo "Warning: psql is not installed. Run migrations in application instead."
    exit 0
fi

# Create database if it doesn't exist
echo "Creating database if not exists..."
psql -U $DB_USER -h $DB_HOST -p $DB_PORT -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1 || \
    psql -U $DB_USER -h $DB_HOST -p $DB_PORT -c "CREATE DATABASE $DB_NAME;"

echo "Database ready. Run: go run cmd/server/main.go"
echo "Migrations will run automatically on server startup."
