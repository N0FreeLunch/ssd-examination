#!/bin/bash
set -e
export AWS_DEFAULT_REGION=ap-northeast-2

echo "🚀 [ApplicationStart] Starting new application container..."
cd /app

chmod +x server

# SSM 구성을 통과시켜 Docker Compose 기동 (환경에 맞는 yml 사용)
# CodeDeploy는 내부적으로 env를 주입받으므로 with-secrets.sh를 통해 의존성 주입 확인 필수
TARGET_ENV=${DEPLOYMENT_GROUP_NAME:-dev}

if [[ "$TARGET_ENV" == *"prod"* ]]; then
    ./tools/with-secrets.sh prod "docker compose -f docker-compose.prod.yml up -d --build"
else
    ./tools/with-secrets.sh dev "docker compose -f docker-compose.dev.yml up -d --build"
fi

echo "✅ Deployment completed smoothly."
