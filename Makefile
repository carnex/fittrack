# Makefile
# Run `make help` to see all commands

.PHONY: help dev stop backend test lint

help:  ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

dev:  ## Start all services (db + backend + frontend)
	docker compose up

dev-db:  ## Start only the database
	docker compose up db

stop:  ## Stop all services
	docker compose down

backend:  ## Run the Go backend locally (without Docker)
	cd backend && go run main.go

test:  ## Run all Go tests
	cd backend && go test ./... -v

test-unit:  ## Run only unit tests
	cd backend && go test ./tests/unit/... -v

lint:  ## Run Go linter
	cd backend && go vet ./...

tidy:  ## Tidy Go modules
	cd backend && go mod tidy