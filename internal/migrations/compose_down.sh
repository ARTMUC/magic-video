#!/bin/bash

# Script to rollback all migrations
# Usage: ./compose_down.sh [number_of_migrations]

set -e

echo "Rolling back database migrations..."

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

# Number of migrations to rollback (default: all)
ROLLBACK_COUNT=${1:-0}

if [ "$ROLLBACK_COUNT" -eq 0 ]; then
    echo "Rolling back ALL migrations..."
    goose -dir . mysql "$DB_CONNECTION" reset
else
    echo "Rolling back $ROLLBACK_COUNT migration(s)..."
    for ((i=1; i<=ROLLBACK_COUNT; i++)); do
        goose -dir . mysql "$DB_CONNECTION" down
    done
fi

echo "Migration rollback completed successfully!"
