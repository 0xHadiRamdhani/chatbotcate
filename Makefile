.PHONY: build run test clean docker-build docker-run migrate seed lint format help

# Variables
APP_NAME=whatsapp-bot
DOCKER_IMAGE=whatsapp-bot:latest
GO=go
AIR=air
DOCKER_COMPOSE=docker-compose

# Default target
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application in development mode"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose"
	@echo "  migrate      - Run database migrations"
	@echo "  seed         - Seed database with test data"
	@echo "  lint         - Run linter"
	@echo "  format       - Format code"
	@echo "  deps         - Download dependencies"
	@echo "  mod-tidy     - Tidy Go modules"
	@echo "  coverage     - Generate test coverage report"
	@echo "  docs         - Generate documentation"
	@echo "  security     - Run security checks"

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	$(GO) build -o bin/$(APP_NAME) main.go
	@echo "Build complete!"

# Run in development mode
run:
	@echo "Running $(APP_NAME) in development mode..."
	$(GO) run main.go

# Run with hot reload (requires air)
dev:
	@echo "Running with hot reload..."
	$(AIR)

# Run tests
test:
	@echo "Running tests..."
	$(GO) test -v ./test

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -cover ./test

# Generate coverage report
coverage:
	@echo "Generating coverage report..."
	$(GO) test -coverprofile=coverage.out ./test
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

# Run with Docker Compose
docker-run:
	@echo "Starting services with Docker Compose..."
	$(DOCKER_COMPOSE) up -d

# Stop Docker services
docker-stop:
	@echo "Stopping Docker services..."
	$(DOCKER_COMPOSE) down

# View Docker logs
docker-logs:
	@echo "Viewing Docker logs..."
	$(DOCKER_COMPOSE) logs -f

# Database migration
migrate:
	@echo "Running database migrations..."
	$(GO) run main.go migrate

# Seed database
seed:
	@echo "Seeding database..."
	$(GO) run main.go seed

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
format:
	@echo "Formatting code..."
	$(GO) fmt ./...
	@echo "Code formatted!"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download

# Tidy Go modules
mod-tidy:
	@echo "Tidying Go modules..."
	$(GO) mod tidy

# Generate documentation
docs:
	@echo "Generating documentation..."
	$(GO) doc -all > docs/godoc.txt
	@echo "Documentation generated!"

# Run security checks
security:
	@echo "Running security checks..."
	$(GO) list -json -m all | nancy sleuth
	$(GO) run github.com/securecodewarrior/gosec/v2/cmd/gosec@latest ./...

# Install development tools
install-tools:
	@echo "Installing development tools..."
	$(GO) install github.com/cosmtrek/air@latest
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "Development tools installed!"

# Database operations
db-create:
	@echo "Creating database..."
	psql -U postgres -c "CREATE DATABASE whatsapp_bot;"

db-drop:
	@echo "Dropping database..."
	psql -U postgres -c "DROP DATABASE IF EXISTS whatsapp_bot;"

db-reset:
	@echo "Resetting database..."
	$(MAKE) db-drop
	$(MAKE) db-create
	$(MAKE) migrate
	$(MAKE) seed

# Redis operations
redis-flush:
	@echo "Flushing Redis..."
	redis-cli FLUSHALL

# Production build
build-prod:
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags="-w -s" -o bin/$(APP_NAME) main.go

# Run production build
run-prod:
	@echo "Running production build..."
	./bin/$(APP_NAME)

# Setup development environment
setup-dev:
	@echo "Setting up development environment..."
	$(MAKE) install-tools
	$(MAKE) deps
	$(MAKE) mod-tidy
	@echo "Development environment setup complete!"

# Run all tests
test-all:
	@echo "Running all tests..."
	$(MAKE) test
	$(MAKE) test-coverage

# CI/CD pipeline
ci:
	@echo "Running CI pipeline..."
	$(MAKE) lint
	$(MAKE) test
	$(MAKE) security
	$(MAKE) build

# Deploy to staging
deploy-staging:
	@echo "Deploying to staging..."
	$(MAKE) build-prod
	# Add your staging deployment commands here

# Deploy to production
deploy-prod:
	@echo "Deploying to production..."
	$(MAKE) build-prod
	# Add your production deployment commands here

# Health check
health:
	@echo "Checking application health..."
	curl -f http://localhost:8080/health || exit 1

# Load test (requires hey or similar tool)
load-test:
	@echo "Running load tests..."
	hey -n 1000 -c 10 http://localhost:8080/health

# Benchmark
benchmark:
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./test

# Generate API documentation
api-docs:
	@echo "Generating API documentation..."
	# Add your API documentation generation commands here
	@echo "API documentation generated!"

# Run database backup
backup:
	@echo "Creating database backup..."
	pg_dump -U postgres whatsapp_bot > backups/whatsapp_bot_$(shell date +%Y%m%d_%H%M%S).sql

# Restore database
restore:
	@echo "Restoring database..."
	# Usage: make restore FILE=backups/your-backup.sql
	psql -U postgres whatsapp_bot < $(FILE)

# Watch for changes and rebuild
watch:
	@echo "Watching for changes..."
	# Requires entr or similar tool
	find . -name "*.go" | entr -r $(MAKE) build

# Quick development cycle
dev-cycle: format lint test
	@echo "Development cycle complete!"

# Full clean and rebuild
rebuild: clean build
	@echo "Rebuild complete!"

# Show Go version and environment
version:
	@echo "Go version: $(shell $(GO) version)"
	@echo "Go environment:"
	@$(GO) env

# List all make targets
targets:
	@echo "Available targets:"
	@$(MAKE) -qp | grep -v '\.PHONY:' | grep -v '^[[:space:]]*#' | grep -v '^$$' | grep '^[a-zA-Z0-9_-]*:' | sed 's/:.*$$//' | sort