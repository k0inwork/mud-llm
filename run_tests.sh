#!/bin/bash

set -e

SERVER_PID=""
DB_FILE="./mud.db"

cleanup() {
    echo "Cleaning up..."
    # Force kill any processes using ports 4000 or 8080
    lsof -ti:4000 | xargs -r kill -9 || true
    lsof -ti:8080 | xargs -r kill -9 || true
    if [ -n "$SERVER_PID" ]; then
        echo "Killing server process $SERVER_PID"
        kill -SIGTERM "$SERVER_PID" || true
        wait "$SERVER_PID" 2>/dev/null || true
    fi
    if [ -f "$DB_FILE" ]; then
        echo "Removing database file: $DB_FILE"
        rm "$DB_FILE"
    fi
    echo "Cleanup complete."
}

trap cleanup EXIT

# Ensure a clean environment before starting
lsof -ti:4000 | xargs -r kill -9 || true
lsof -ti:8080 | xargs -r kill -9 || true

# Ensure a clean database for each run
if [ -f "$DB_FILE" ]; then
    echo "Removing existing database file: $DB_FILE"
    rm "$DB_FILE"
fi

LOG_FILE="test_run_$(date +%Y%m%d_%H%M%S).log"
echo "Starting MUD server in background... Log file: $LOG_FILE"
go run main.go > "$LOG_FILE" 2>&1 & 
SERVER_PID=$!

echo "Server started with PID: $SERVER_PID. Giving it some time to initialize..."
sleep 5 # Give the server 5 seconds to start up

echo "Running Go tests..."
go test ./... >> "$LOG_FILE" 2>&1

echo "Tests finished. Server will be killed by cleanup trap."