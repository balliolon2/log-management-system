#!/bin/bash
# run.sh - Helper script to start the system
# Usage: ./run.sh [dev|prod|down]

MODE=${1:-dev}

if [ "$MODE" == "down" ]; then
    echo "Stopping system..."
    cd deploy && docker compose down
    exit 0
fi

echo "Starting Log Management System in [$MODE] mode..."

cd deploy

# Check for .env
if [ ! -f .env ]; then
    echo " .env not found, creating from example..."
    cp .env.example .env
fi

# Check for Certs
if [ ! -f certs/cert.pem ]; then
    echo "Certificates not found, generating self-signed..."
    mkdir -p certs
    MSYS_NO_PATHCONV=1 docker run --rm -v "${PWD}/certs:/certs" alpine/openssl req -x509 -newkey rsa:4096 -keyout /certs/key.pem -out /certs/cert.pem -days 365 -nodes -subj "/CN=localhost"
fi

# Run
echo "Docker Compose Up..."
docker compose up --build -d

echo "System started!"
echo "Frontend: https://localhost"
