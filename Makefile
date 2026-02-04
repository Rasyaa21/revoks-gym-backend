# Makefile for Fiber-GORM application

.PHONY: help dev-up dev-down dev-logs dev-rebuild prod-up prod-down prod-logs prod-rebuild clean status

# Default target
help:
	@echo "Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  make dev-up        - Start development environment with hot reload"
	@echo "  make dev-down      - Stop development environment"
	@echo "  make dev-logs      - View development logs"
	@echo "  make dev-rebuild   - Rebuild development environment"
	@echo ""
	@echo "Production:"
	@echo "  make prod-up       - Start production environment"
	@echo "  make prod-down     - Stop production environment"
	@echo "  make prod-logs     - View production logs"
	@echo "  make prod-rebuild  - Rebuild production environment"
	@echo ""
	@echo "Utilities:"
	@echo "  make status        - Show container status"
	@echo "  make clean         - Remove all Docker resources"
	@echo "  make go-mod        - Download Go modules"
	@echo "  make go-tidy       - Tidy Go modules"

# Development commands
dev-up:
	docker compose -f docker-compose.dev.yml up --build -d
	@echo "Development environment started!"
	@echo "App: http://localhost:3000"
	@echo "pgAdmin: http://localhost:5050"

dev-down:
	docker compose -f docker-compose.dev.yml down

dev-logs:
	docker compose -f docker-compose.dev.yml logs -f

dev-logs-app:
	docker compose -f docker-compose.dev.yml logs -f app

dev-restart:
	docker compose -f docker-compose.dev.yml restart

dev-rebuild:
	docker compose -f docker-compose.dev.yml down
	docker compose -f docker-compose.dev.yml up --build -d

# Production commands
prod-up:
	docker compose -f docker-compose.prod.yml up --build -d
	@echo "Production environment started!"
	@echo "App: http://localhost:3000"
	@echo "pgAdmin: http://localhost:5050"

prod-down:
	docker compose -f docker-compose.prod.yml down

prod-logs:
	docker compose -f docker-compose.prod.yml logs -f

prod-logs-app:
	docker compose -f docker-compose.prod.yml logs -f app

prod-restart:
	docker compose -f docker-compose.prod.yml restart

prod-rebuild:
	docker compose -f docker-compose.prod.yml down
	docker compose -f docker-compose.prod.yml up --build -d

# Database commands
db-shell-dev:
	docker exec -it fiber-postgres-dev psql -U postgres -d fiber_gorm_db

db-shell-prod:
	docker exec -it fiber-postgres-prod psql -U postgres -d fiber_gorm_db

db-backup-dev:
	docker exec fiber-postgres-dev pg_dump -U postgres fiber_gorm_db > backup_dev_$$(date +%Y%m%d_%H%M%S).sql

db-backup-prod:
	docker exec fiber-postgres-prod pg_dump -U postgres fiber_gorm_db > backup_prod_$$(date +%Y%m%d_%H%M%S).sql

# Utility commands
status:
	docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

clean:
	docker compose -f docker-compose.dev.yml down -v --rmi all --remove-orphans 2>/dev/null || true
	docker compose -f docker-compose.prod.yml down -v --rmi all --remove-orphans 2>/dev/null || true

clean-volumes:
	docker compose -f docker-compose.dev.yml down -v 2>/dev/null || true
	docker compose -f docker-compose.prod.yml down -v 2>/dev/null || true

# Go commands (for local development without Docker)
go-mod:
	go mod download

go-tidy:
	go mod tidy

go-run:
	go run cmd/main.go

go-build:
	go build -o bin/main cmd/main.go

go-test:
	go test -v ./...

# Air hot reload (local development)
air:
	air -c .air.toml
