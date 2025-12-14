#!/bin/bash

set -e

echo "Starting deployment..."

# Configuration
PROJECT_DIR="$HOME/Projects/angagrar-backend"
ENV_FILE="$PROJECT_DIR/.env"

cd "$PROJECT_DIR" || exit 1

# Check .env file
if [ ! -f "$ENV_FILE" ]; then
    echo "Error: .env file not found"
    echo "Please create .env from .env.example"
    exit 1
fi

# Stop existing containers
echo "Stopping containers..."
docker compose down || true

# Pull latest code
echo "Pulling latest code..."
git pull origin main || echo "Already up to date"

# Tidy dependencies
echo "Tidying Go dependencies..."
go mod tidy

# Build and start
echo "Building images..."
docker compose build --no-cache

echo "Starting containers..."
docker compose up -d

# Wait for startup
sleep 10

# Show status
echo "Container status:"
docker compose ps

echo "Recent logs:"
docker compose logs --tail=20

# Verify deployment
if docker compose ps | grep -q "angagrar-backend.*Up"; then
    echo "Deployment successful"
    echo "Backend running on port 8080"

    # Optional health check
    sleep 3
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        echo "Health check passed"
    fi
else
    echo "Deployment failed"
    docker compose logs --tail=50
    exit 1
fi
