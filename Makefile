APP_NAME := kyc-sim
ENV_FILE ?= .env.local

# ==========================
# Help
# ==========================
.PHONY: help
help:
	@echo "Targets:"
	@echo "  make run-sqlite        Run API using SQLite (quick local dev)"
	@echo "  make run-postgres      Run API using Postgres (docker-compose)"
	@echo "  make run-env           Run API loading .env.local"
	@echo "  make up-postgres       Start Postgres via docker-compose"
	@echo "  make down-postgres     Stop Postgres (keep volume)"
	@echo "  make down-postgres-v   Stop Postgres and remove volume"
	@echo "  make test              Run tests"
	@echo "  make lint              Run golangci-lint"

# ==========================
# Utils
# ==========================
.PHONY: ensure-data-dir
ensure-data-dir:
	@mkdir -p ./data

# ==========================
# Docker / Postgres
# ==========================
.PHONY: up-postgres
up-postgres:
	docker compose up -d postgres

.PHONY: down-postgres
down-postgres:
	docker compose down

.PHONY: down-postgres-v
down-postgres-v:
	docker compose down -v

# ==========================
# Run targets
# ==========================
.PHONY: run-sqlite
run-sqlite: ensure-data-dir
	DB_DRIVER=sqlite \
	DB_PATH=./data/kyc.db \
	PORT=8080 \
	DB_LOG_LEVEL=info \
	go run ./cmd/api

.PHONY: run-postgres
run-postgres: up-postgres
	DB_DRIVER=postgres \
	DB_HOST=localhost \
	DB_PORT=5432 \
	DB_USER=postgres \
	DB_PASSWORD=postgres \
	DB_NAME=kyc \
	DB_SSLMODE=disable \
	DB_TIMEZONE=UTC \
	PORT=8080 \
	DB_LOG_LEVEL=info \
	go run ./cmd/api

.PHONY: run-env
run-env:
	@set -a; . $(ENV_FILE); set +a; \
	go run ./cmd/api

# ==========================
# Quality
# ==========================
.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run
