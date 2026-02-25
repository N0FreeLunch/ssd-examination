.PHONY: seed-admin up down logs shell
WITH_SECRETS := ./tools/with-secrets.sh
ENV ?= local

# Data Seeding
seed-admin:
	@echo "Seeding admin user..."
	@docker exec -it examination-app-local go run cmd/seeder/main.go

seed-exam-preview:
	@echo "Seeding exam preview data..."
	@docker exec -it examination-app-local go run cmd/seeder/main.go -name=exam_preview -clean

build-seeder-linux:
	@echo "🔨 Building the seeder binary for Linux (amd64)..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o seeder ./cmd/seeder/main.go

seed-dev: build-seeder-linux
	@echo "📤 Uploading seeder to dev server and executing..."
	@HOST=$$(aws ssm get-parameter --name "/sdd-exam/dev/DEPLOY_HOST" --with-decryption --query "Parameter.Value" --output text) && \
	scp -i sdd-exam-key.pem seeder "$$HOST:/app/seeder" && \
	ssh -i sdd-exam-key.pem "$$HOST" "sudo docker cp /app/seeder examination-app:/app/seeder && sudo docker exec examination-app /app/seeder -name=exam_preview -clean"
	@echo "✅ Database seeded successfully on dev server."

# Docker Compose Helpers
up:
	@$(WITH_SECRETS) $(ENV) "docker-compose up -d"

down:
	@$(WITH_SECRETS) $(ENV) "docker-compose down"

logs:
	docker-compose logs -f app

shell:
	docker exec -it examination-app-local /bin/sh

build-linux:
	@echo "🔨 Building the server binary for Linux (amd64)..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server ./cmd/server

deploy: build-linux
	@$(WITH_SECRETS) dev "./deploy.sh"

rebuild:
	@$(WITH_SECRETS) $(ENV) "docker-compose up -d --build"
