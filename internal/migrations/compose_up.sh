#!/bin/bash

# Script to run all migrations up
# Usage: ./compose_up.sh

set -e

echo "Running database migrations up..."

# Check if goose is installed
if ! command -v goose &> /dev/null; then
    echo "Error: goose is not installed. Please install it first:"
    echo "go install github.com/pressly/goose/v3/cmd/goose@latest"
    exit 1
fi

# Database connection parameters
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-3306}
DB_USER=${DB_USER:-root}
DB_PASSWORD=${DB_PASSWORD:-password}
DB_NAME=${DB_NAME:-magic_video}

# Connection string
DB_CONNECTION="$DB_USER:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME?parseTime=true"

echo "Connecting to database: $DB_HOST:$DB_PORT/$DB_NAME"

# Run migrations
goose -dir . mysql "$DB_CONNECTION" up

echo "All migrations completed successfully!"
