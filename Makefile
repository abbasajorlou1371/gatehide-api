.PHONY: help run build test clean install lint fmt dev hot

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

