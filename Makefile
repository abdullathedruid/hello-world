.PHONY: test test-verbose test-coverage build run clean lint fmt vet docker-build docker-up docker-down docker-restart

# Default target
all: test build

# Run tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run benchmark tests
bench:
	go test -bench=. -benchmem ./...

# Run integration tests
test-integration:
	go test -v -tags=integration ./...

# Build the application
build:
	go build -o main .

# Run the application
run:
	go run .

# Clean build artifacts
clean:
	rm -f main coverage.out coverage.html

# Format code
fmt:
	go fmt ./...

# Vet code for potential issues
vet:
	go vet ./...

# Run linter (requires golangci-lint)
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, running go vet instead"; \
		go vet ./...; \
	fi

# Run all checks
check: fmt vet test bench

# Install dependencies
deps:
	go mod tidy
	go mod download

# Docker targets
COMMIT_HASH := $(shell git rev-parse HEAD)

# Build Docker images with commit hash
docker-build:
	@echo "Building with commit hash: $(COMMIT_HASH)"
	COMMIT_HASH=$(COMMIT_HASH) docker-compose build

# Start services with automatic commit hash
docker-up:
	@echo "Starting with commit hash: $(COMMIT_HASH)"
	COMMIT_HASH=$(COMMIT_HASH) docker-compose up

# Start services in background with automatic commit hash
docker-up-detached:
	@echo "Starting in background with commit hash: $(COMMIT_HASH)"
	COMMIT_HASH=$(COMMIT_HASH) docker-compose up -d

# Stop services
docker-down:
	docker-compose down

# Restart services with latest commit hash
docker-restart: docker-down docker-up

# Rebuild and start services
docker-rebuild: docker-build docker-up

# View logs
docker-logs:
	docker-compose logs -f