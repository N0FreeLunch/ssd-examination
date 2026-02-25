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
scp -i $KEY docker-compose.dev.yml $HOST:$REMOTE_DIR/
scp -i $KEY infra/config/litestream.dev.yml $HOST:$REMOTE_DIR/infra/config/
scp -i $KEY infra/docker/restore-db.sh $HOST:$REMOTE_DIR/infra/docker/
# Make restore script executable
ssh -i $KEY $HOST "chmod +x $REMOTE_DIR/infra/docker/restore-db.sh"

# 3. Upload required files (binary, Dockerfile, tools)
echo "💻 Uploading pre-built binary and tools..."
# Use rsync for efficient upload
rsync -avz -e "ssh -i $KEY" \
  server \
  Dockerfile.deploy \
  tools \
  $HOST:$REMOTE_DIR/

# 4. Execute deployment on server
echo "🚀 Running deployment on remote server..."
echo "   (Using Access Key: ${LITESTREAM_ACCESS_KEY_ID:?LITESTREAM_ACCESS_KEY_ID is missing locally})"
# Pass Litestream credentials (which have SSM read permissions) to the server
# so it can fetch the remaining secrets via with-secrets.sh
ssh -i $KEY $HOST "cd $REMOTE_DIR && \
  chmod +x tools/with-secrets.sh && \
  AWS_ACCESS_KEY_ID=$LITESTREAM_ACCESS_KEY_ID \
  AWS_SECRET_ACCESS_KEY=$LITESTREAM_SECRET_ACCESS_KEY \
  AWS_DEFAULT_REGION=${AWS_REGION:-ap-northeast-2} \
  ./tools/with-secrets.sh dev 'docker compose -f docker-compose.dev.yml up -d --build'"

echo "✅ Deployment completed successfully!"
