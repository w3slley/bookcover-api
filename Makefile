.PHONY: build test run clean docker-build docker-up docker-down

# Build the application
build:
	go build -o bin/server ./cmd/api

# Run tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Generate detailed coverage report
test-coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Run the application
run:
	go run cmd/api/main.go

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html

# Build Docker image
docker-build:
	docker build -t bookcover-api .

# Start with Docker Compose
docker-up:
	docker compose up --build

# Stop Docker Compose
docker-down:
	docker compose down
