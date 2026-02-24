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

# 4. Upload source code and tools
echo "💻 Uploading source code and tools..."
# Use rsync to exclude unnecessary files
rsync -avz -e "ssh -i $KEY" \
  --exclude '.git' \
  --exclude '.idea' \
  --exclude 'bin' \
  --exclude 'data' \
  --exclude 'node_modules' \
  --exclude '*.test' \
  --exclude 'sdd-exam-key.pem' \
  cmd internal tools \
  $HOST:$REMOTE_DIR/

# 5. Execute deployment on server
echo "🚀 Running deployment on remote server..."
echo "   (Using Access Key: ${LITESTREAM_ACCESS_KEY_ID:?LITESTREAM_ACCESS_KEY_ID is missing locally})"
# Pass Litestream credentials (which have SSM read permissions) to the server
# so it can fetch the remaining secrets via with-secrets.sh
ssh -i $KEY $HOST "cd $REMOTE_DIR && \
  chmod +x tools/with-secrets.sh && \
  AWS_ACCESS_KEY_ID=$LITESTREAM_ACCESS_KEY_ID \
  AWS_SECRET_ACCESS_KEY=$LITESTREAM_SECRET_ACCESS_KEY \
  AWS_DEFAULT_REGION=${AWS_REGION:-ap-northeast-2} \
  ./tools/with-secrets.sh prod 'docker compose -f docker-compose.prod.yml up -d --build'"

echo "✅ Deployment completed successfully!"
