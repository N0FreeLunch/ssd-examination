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

# Docker Compose Helpers
up:
	@$(WITH_SECRETS) $(ENV) "docker-compose up -d"

down:
	@$(WITH_SECRETS) $(ENV) "docker-compose down"

logs:
	docker-compose logs -f app

shell:
	docker exec -it examination-app-local /bin/sh

deploy:
	@$(WITH_SECRETS) dev "./deploy.sh"
