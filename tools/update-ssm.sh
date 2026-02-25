#!/bin/bash
set -e

echo "⏳ Fetching Terraform outputs..."
cd infra/terraform
DEPLOY_HOST=$(terraform output -raw static_ip)
LITESTREAM_ACCESS_KEY_ID=$(terraform output -raw litestream_access_key_id)
LITESTREAM_SECRET_ACCESS_KEY=$(terraform output -raw litestream_secret_access_key)
LITESTREAM_BUCKET=$(terraform output -raw backup_bucket_name)

cd ../../

echo "✅ Terraform outputs fetched. Updating SSM Parameter Store..."

aws ssm put-parameter --region ap-northeast-2 --name "/sdd-exam/dev/DEPLOY_HOST" --value "ubuntu@$DEPLOY_HOST" --type "SecureString" --overwrite > /dev/null
aws ssm put-parameter --region ap-northeast-2 --name "/sdd-exam/dev/LITESTREAM_ACCESS_KEY_ID" --value "$LITESTREAM_ACCESS_KEY_ID" --type "SecureString" --overwrite > /dev/null
aws ssm put-parameter --region ap-northeast-2 --name "/sdd-exam/dev/LITESTREAM_SECRET_ACCESS_KEY" --value "$LITESTREAM_SECRET_ACCESS_KEY" --type "SecureString" --overwrite > /dev/null
aws ssm put-parameter --region ap-northeast-2 --name "/sdd-exam/dev/LITESTREAM_BUCKET" --value "$LITESTREAM_BUCKET" --type "SecureString" --overwrite > /dev/null

aws ssm put-parameter --region ap-northeast-2 --name "/sdd-exam/prod/DEPLOY_HOST" --value "ubuntu@$DEPLOY_HOST" --type "SecureString" --overwrite > /dev/null
aws ssm put-parameter --region ap-northeast-2 --name "/sdd-exam/prod/LITESTREAM_ACCESS_KEY_ID" --value "$LITESTREAM_ACCESS_KEY_ID" --type "SecureString" --overwrite > /dev/null
aws ssm put-parameter --region ap-northeast-2 --name "/sdd-exam/prod/LITESTREAM_SECRET_ACCESS_KEY" --value "$LITESTREAM_SECRET_ACCESS_KEY" --type "SecureString" --overwrite > /dev/null
aws ssm put-parameter --region ap-northeast-2 --name "/sdd-exam/prod/LITESTREAM_BUCKET" --value "$LITESTREAM_BUCKET" --type "SecureString" --overwrite > /dev/null

echo "🎉 All SSM Parameters updated successfully!"
