.PHONY: build run test clean docker-build docker-run docker-down migrate-up migrate-down deps

# Build the application
build:
	go build -o bin/main cmd/main.go

# Build for production
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o bin/main cmd/main.go

# Run the application
run:
	go run cmd/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -cover ./...

# Run tests with detailed coverage report
test-coverage-detail:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run only unit tests (exclude integration tests)
test-unit:
	go test -v -short ./internal/...

# Run specific test package
test-repositories:
	go test -v ./internal/repositories/...

# Run specific test package
test-services:
	go test -v ./internal/services/...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Generate Swagger documentation
swagger:
	swag init -g cmd/main.go -o docs

# Docker commands
docker-build:
	docker-compose build

docker-run:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f api

# Database migration commands (handled by application)
# These commands require external migrate tool - migrations are now handled by the application
# migrate-up:
# 	migrate -path migrations -database "postgres://postgres:password@localhost:5432/flowsight_db?sslmode=disable" up
# 
# migrate-down:
# 	migrate -path migrations -database "postgres://postgres:password@localhost:5432/flowsight_db?sslmode=disable" down

# Development setup
dev-setup: deps docker-run
	@echo "Waiting for database to be ready..."
	@sleep 10
	@echo "Migrations will be handled by the application on startup"
	@echo "Development environment is ready!"

# Format code
fmt:
	go fmt ./...

# Lint code (assuming golangci-lint is installed)
lint:
	golangci-lint run

# Security check (assuming gosec is installed)
security:
	gosec ./...

# Run development server with hot reload (assuming air is installed)
dev:
	air
