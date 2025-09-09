# Default target
.DEFAULT_GOAL := help

# help
.PHONY: help
help: ## Show help messages
	@echo "Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

# format
.PHONY: format
format: ## Format the code
	go fmt ./...

# lint
.PHONY: lint
lint: ## Lint the code
	golangci-lint run

# check
.PHONY: check
check: ## Check the code
	go fmt ./... && golangci-lint run && go test ./...

# mock_gen
.PHONY: mock_gen
mock_gen: ## Generate mock files
	go generate ./...

# sqlc
.PHONY: sqlc
sqlc: ## Generate the sqlc files
	go tool sqlc generate

# run
.PHONY: run
run: ## Run the application
	docker compose up -d
	@echo "Waiting for Databases to be ready..."
	@until docker compose ps app-db | grep -q "healthy" && docker compose ps vector-db | grep -q "healthy"; do \
		echo "Databases are starting up..."; \
		sleep 2; \
	done
	@echo "PostgreSQL databases are ready!"
	@echo "Applying migrations..."
	make migrate-app-up
	make migrate-vector-up

# stop
.PHONY: stop
stop: ## Stop the application
	docker compose down --volumes --remove-orphans

# migrate-app-up
.PHONY: migrate-app-up
migrate-app-up: ## Run the migrations app up
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path ../infrastructure/schemas/migrations/app -database postgres://postgres:password@localhost:5433/app_db?sslmode=disable up

# migrate-vector-up
.PHONY: migrate-vector-up
migrate-vector-up: ## Run the migrations vector up
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path ../infrastructure/schemas/migrations/vector -database postgres://postgres:password@localhost:5432/vector_db?sslmode=disable up

# migrate-down
.PHONY: migrate-app-down
migrate-app-down: ## Run the migrations app down
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path ../infrastructure/schemas/migrations/app -database postgres://postgres:password@localhost:5433/app_db?sslmode=disable down -all

# migrate-vector-down
.PHONY: migrate-vector-down
migrate-vector-down: ## Run the migrations vector down
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path ../infrastructure/schemas/migrations/vector -database postgres://postgres:password@localhost:5432/vector_db?sslmode=disable down -all

# oapi-codegen
.PHONY: oapi-codegen
oapi-codegen: ## generate oapi-codegen
	go tool oapi-codegen --config=./oapi-codegen.yaml ../infrastructure/schemas/openapi/admin/openapi.yaml

# wire_gen
.PHONY: wire_gen
wire_gen: ## generate wire
	go run github.com/google/wire/cmd/wire ./di	