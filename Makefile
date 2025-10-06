.PHONY: help run build test clean install lint fmt dev hot migrate-status migrate-up migrate-down migrate-create migrate-build migrate-reset migrate-fresh migrate-up-seed migrate-fresh-seed seed-admin seed-build

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
	@DB_AUTO_CREATE=true go run cmd/migrate/main.go -command=up -steps=$${STEPS:-999}

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

migrate-up-seed: ## Run pending migrations with seeding (optionally specify seeder with SEEDER=name)
	@echo "⬆️  Running pending migrations with seeding..."
	@DB_AUTO_CREATE=true go run cmd/migrate/main.go -command=up -steps=$${STEPS:-999} -seed=$${SEEDER:-all}

migrate-fresh-seed: ## Fresh migration with seeding (optionally specify seeder with SEEDER=name)
	@echo "🆕 Fresh migration with seeding..."
	@$(MAKE) migrate-reset
	@DB_AUTO_CREATE=true go run cmd/migrate/main.go -command=up -steps=999 -seed=$${SEEDER:-all}

# Seeder commands
seed-build: ## Build seeder CLI tool
	@echo "🔨 Building seeder CLI..."
	@go build -o $(BUILD_DIR)/seed cmd/seed/main.go
	@echo "✅ Seeder CLI built: $(BUILD_DIR)/seed"

seed-admin: ## Seed admin user
	@echo "👤 Seeding admin user..."
	@go run cmd/seed/main.go -command=admin

# Test commands
test: ## Run all tests
	@echo "🧪 Running all tests..."
	@go test -v ./tests/...

test-unit: ## Run unit tests only
	@echo "🔬 Running unit tests..."
	@go test -v -short ./tests/unit/...

test-integration: ## Run integration tests only
	@echo "🔗 Running integration tests..."
	@go test -v -run "Integration" ./tests/integration/...

test-auth: ## Run authentication tests only
	@echo "🔐 Running authentication tests..."
	@go test -v ./tests/unit -run "TestJWT\|TestUserRepository\|TestAdminRepository\|TestAuthService\|TestAuthHandler\|TestAuth"
	@go test -v ./tests/integration -run "TestAuthentication"

test-coverage: ## Run tests with coverage
	@echo "📊 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./tests/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

test-coverage-auth: ## Run authentication tests with coverage
	@echo "📊 Running authentication tests with coverage..."
	@go test -v -coverprofile=auth_coverage.out ./tests/unit/... ./tests/integration/...
	@go tool cover -html=auth_coverage.out -o auth_coverage.html
	@echo "✅ Authentication coverage report generated: auth_coverage.html"

test-watch: ## Run tests in watch mode (requires entr)
	@echo "👀 Running tests in watch mode..."
	@find tests/ -name "*.go" | entr -c go test -v ./tests/...

test-benchmark: ## Run benchmark tests
	@echo "⚡ Running benchmark tests..."
	@go test -bench=. -benchmem ./tests/...

test-race: ## Run tests with race detection
	@echo "🏁 Running tests with race detection..."
	@go test -v -race ./tests/...

test-db: ## Setup test database
	@echo "🗄️  Setting up test database..."
	@mysql -u root -e "CREATE DATABASE IF NOT EXISTS gatehide_test;" || echo "⚠️  Could not create test database. Make sure MySQL is running."

test-db-drop: ## Drop test database
	@echo "🗑️  Dropping test database..."
	@mysql -u root -e "DROP DATABASE IF NOT EXISTS gatehide_test;" || echo "⚠️  Could not drop test database."

test-clean: ## Clean test artifacts
	@echo "🧹 Cleaning test artifacts..."
	@rm -f coverage.out coverage.html auth_coverage.out auth_coverage.html
	@echo "✅ Test artifacts cleaned"

