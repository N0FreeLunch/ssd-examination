#!/bin/bash
set -e

# Configuration
# Configuration (Injected via tools/with-secrets.sh)
HOST="${DEPLOY_HOST:?DEPLOY_HOST is required}"
KEY="${DEPLOY_KEY:?DEPLOY_KEY is required}"
REMOTE_DIR="/app"

echo "🚀 Starting deployment to $HOST..."

# 1. Create directory structure
echo "📁 Creating remote directories..."
ssh -i $KEY $HOST "mkdir -p $REMOTE_DIR/infra/config $REMOTE_DIR/infra/docker"

# 2. Upload configuration files
echo "📤 Uploading configuration files..."
scp -i $KEY docker-compose.prod.yml $HOST:$REMOTE_DIR/
scp -i $KEY infra/config/litestream.prod.yml $HOST:$REMOTE_DIR/infra/config/
scp -i $KEY infra/docker/restore-db.sh $HOST:$REMOTE_DIR/infra/docker/
# Make restore script executable
ssh -i $KEY $HOST "chmod +x $REMOTE_DIR/infra/docker/restore-db.sh"

# 3. Upload build dependencies
echo "📦 Uploading build dependencies..."
scp -i $KEY Dockerfile $HOST:$REMOTE_DIR/
scp -i $KEY go.mod go.sum $HOST:$REMOTE_DIR/

# 4. Upload source code (using rsync for efficiency)
echo "💻 Uploading source code..."
# Use rsync to exclude unnecessary files
rsync -avz -e "ssh -i $KEY" \
  --exclude '.git' \
  --exclude '.idea' \
  --exclude 'bin' \
  --exclude 'data' \
  --exclude 'node_modules' \
  --exclude '*.test' \
  --exclude 'sdd-exam-key.pem' \
  cmd internal \
  $HOST:$REMOTE_DIR/

echo "✅ All files uploaded successfully!"
echo "👉 Next steps:"
echo "   1. SSH into the server: ssh -i $KEY $HOST"
echo "   2. Go to app directory: cd /app"
echo "   3. Create .env file with secrets (if not exists)"
echo "   4. Run: docker compose -f docker-compose.prod.yml up -d --build"
