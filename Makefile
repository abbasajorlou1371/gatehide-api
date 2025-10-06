.PHONY: help run build test clean install lint fmt dev hot migrate-status migrate-up migrate-down migrate-create migrate-build migrate-reset migrate-fresh

# Variables
BINARY_NAME=gatehide-api
BUILD_DIR=bin
MAIN_PATH=cmd/app/main.go

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy

run: ## Run the application
	@echo "🚀 Starting application..."
	@go run $(MAIN_PATH)

build: ## Build the application
	@echo "🔨 Building application..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test -v ./...

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean

lint: ## Run linter
	@echo "🔍 Running linter..."
	@golangci-lint run

fmt: ## Format code
	@echo "✨ Formatting code..."
	@go fmt ./...

dev: ## Run in development mode with auto-reload (requires air)
	@echo "🔄 Starting development mode with hot reload..."
	@export PATH=$$PATH:$$(go env GOPATH)/bin && air

hot: ## Run with hot reloading (alias for dev)
	@echo "🔥 Starting with hot reload..."
	@export PATH=$$PATH:$$(go env GOPATH)/bin && air

# Migration commands
migrate-build: ## Build migration CLI tool
	@echo "🔨 Building migration CLI..."
	@go build -o $(BUILD_DIR)/migrate cmd/migrate/main.go
	@echo "✅ Migration CLI built: $(BUILD_DIR)/migrate"

migrate-status: ## Show migration status
	@echo "📊 Checking migration status..."
	@DB_AUTO_CREATE=true go run cmd/migrate/main.go -command=status

migrate-up: ## Run pending migrations (optionally specify steps with STEPS=n)
	@echo "⬆️  Running pending migrations..."
	@DB_AUTO_CREATE=true go run cmd/migrate/main.go -command=up -steps=$${STEPS:-1}

migrate-down: ## Rollback migrations (optionally specify steps with STEPS=n)
	@echo "⬇️  Rolling back migrations..."
	@DB_AUTO_CREATE=true go run cmd/migrate/main.go -command=down -steps=$${STEPS:-1}

migrate-create: ## Create a new migration file (usage: make migrate-create NAME="create_users_table")
	@if [ -z "$(NAME)" ]; then \
		echo "❌ Please specify migration name: make migrate-create NAME=\"create_users_table\""; \
		exit 1; \
	fi
	@echo "📝 Creating migration: $(NAME)..."
	@go run cmd/migrate/main.go -command=create -name="$(NAME)"

migrate-reset: ## Reset database (rollback all migrations)
	@echo "🔄 Resetting database..."
	@DB_AUTO_CREATE=true go run cmd/migrate/main.go -command=down -steps=999

migrate-fresh: ## Fresh migration (reset and run all migrations)
	@echo "🆕 Fresh migration..."
	@$(MAKE) migrate-reset
	@$(MAKE) migrate-up STEPS=999

