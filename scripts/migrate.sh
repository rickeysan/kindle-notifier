#!/bin/bash

# Load environment variables
source .env

# Set migrate command path
MIGRATE_CMD="$HOME/go/bin/migrate"

# Build the database URL
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# Function to display usage
usage() {
    echo "Usage: $0 [up|down|create]"
    echo "  up      - Apply all migrations"
    echo "  down    - Rollback all migrations"
    echo "  create  - Create a new migration (requires name argument)"
    exit 1
}

# Check if migrate command exists
if [ ! -f "$MIGRATE_CMD" ]; then
    echo "Error: migrate command not found at $MIGRATE_CMD"
    echo "Please run: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    exit 1
fi

# Check if command is provided
if [ $# -lt 1 ]; then
    usage
fi

# Handle commands
case "$1" in
    "up")
        $MIGRATE_CMD -database "${DB_URL}" -path migrations up
        ;;
    "down")
        $MIGRATE_CMD -database "${DB_URL}" -path migrations down
        ;;
    "create")
        if [ $# -lt 2 ]; then
            echo "Error: Migration name required"
            usage
        fi
        $MIGRATE_CMD create -ext sql -dir migrations -seq "$2"
        ;;
    *)
        usage
        ;;
esac 