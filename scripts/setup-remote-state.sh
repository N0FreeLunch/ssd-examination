#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Starting Terraform Remote State Setup Automation...${NC}"

# 1. Bootstrap Backend Resources (S3 & DynamoDB)
echo -e "${BLUE}📦 Step 1: Bootstrapping remote backend resources...${NC}"
cd infra/bootstrap
terraform init
terraform apply -auto-approve
cd ../..

# 2. Migrate Main State to Remote
echo -e "${BLUE}🔄 Step 2: Migrating local state to S3 backend...${NC}"
cd infra/terraform
terraform init -migrate-state -force-copy
cd ../..

# 3. Apply Main Infrastructure (Deploy Bucket & Metadata Registry)
echo -e "${BLUE}🌍 Step 3: Applying main infrastructure and creating metadata...${NC}"
cd infra/terraform
terraform apply -auto-approve
cd ../..

echo -e "${GREEN}✅ Remote state setup and migration completed successfully!${NC}"
