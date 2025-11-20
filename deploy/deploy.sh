#!/bin/bash

# deploy.sh - Simple deployment script for EchoMind
# Usage: ./deploy.sh <repo_owner> <db_password>

REPO_OWNER=$1
DB_PASSWORD=$2

if [ -z "$REPO_OWNER" ] || [ -z "$DB_PASSWORD" ]; then
  echo "Usage: ./deploy.sh <repo_owner> <db_password>"
  exit 1
fi

echo "Deploying EchoMind..."

# Export variables for docker-compose
export REPO_OWNER=$REPO_OWNER
export DB_PASSWORD=$DB_PASSWORD

# Pull latest images
docker-compose -f docker-compose.prod.yml pull

# Restart services
docker-compose -f docker-compose.prod.yml up -d

echo "Deployment complete!"
