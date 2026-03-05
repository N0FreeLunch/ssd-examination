#!/bin/bash
set -ne

echo "🛑 [ApplicationStop] Stopping existing application gracefully..."
cd /app || exit 0

if [ -f "docker-compose.dev.yml" ] && command -v docker &> /dev/null; then
    docker compose -f docker-compose.dev.yml down || echo "No active containers."
elif [ -f "docker-compose.prod.yml" ] && command -v docker &> /dev/null; then
    docker compose -f docker-compose.prod.yml down || echo "No active containers."
fi
