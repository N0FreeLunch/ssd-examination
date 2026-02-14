.PHONY: seed-admin up down logs shell


# Data Seeding
seed-admin:
	@echo "Seeding admin user..."
	@docker exec -it examination-app-local go run cmd/seeder/main.go

seed-exam-preview:
	@echo "Seeding exam preview data..."
	@docker exec -it examination-app-local go run cmd/seeder/main.go -name=exam_preview -clean

# Docker Compose Helpers
up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f app

shell:
	docker exec -it examination-app-local /bin/sh
