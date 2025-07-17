#!/bin/bash

set -e

SERVER_PID=""
DB_FILE="./mud.db"

cleanup() {
    echo "Cleaning up..."
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

# Ensure a clean database for each run
if [ -f "$DB_FILE" ]; then
    echo "Removing existing database file: $DB_FILE"
    rm "$DB_FILE"
fi

echo "Starting MUD server in background..."
go run main.go & 
SERVER_PID=$!

echo "Server started with PID: $SERVER_PID. Giving it some time to initialize..."
sleep 2 # Give the server 2 seconds to start up

echo "Running Go tests..."
go test ./...

echo "Tests finished. Server will be killed by cleanup trap."
