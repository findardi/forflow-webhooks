-include .env
export

GOOSE_DRIVER ?= $(GOOSE_DRIVER)
GOOSE_DBSTRING ?= $(DATABASE_URL)
GOOSE_MIGRATION_DIR ?= $(GOOSE_MIGRATION_DIR)
DB_DEV_COMPOSE = docker compose -f compose.dev.yaml

export GOOSE_DRIVER GOOSE_DBSTRING GOOSE_MIGRATION_DIR

.PHONY: migrate-create migrate-up migrate-down migrate-status generate format lint db-up db-down

migrate-create:
	@test -n "$(name)" || (echo "usage: make migrate-create name=<name>" && exit 1)
	mkdir -p $(GOOSE_MIGRATION_DIR)
	goose -s create $(name) sql

migrate-up:
	goose up

migrate-down:
	goose down

migrate-status:
	goose status

generate:
	cd server && sqlc generate

format:
	cd server && gofmt -w .
	cd web && bun run format

lint:
	cd server && go mod tidy -diff && golangci-lint run
	cd web && bun run lint && bun run check

db-up:
	$(DB_DEV_COMPOSE) up -d --wait

db-down:
	$(DB_DEV_COMPOSE) down