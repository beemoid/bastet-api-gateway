# Makefile for API Gateway
# Provides convenient commands for common development tasks

.PHONY: help install run build test clean swagger docker-build docker-run

# Default target - show help
help:
	@echo "Available commands:"
	@echo "  make install       - Install dependencies"
	@echo "  make swagger       - Generate Swagger documentation"
	@echo "  make run          - Run the application"
	@echo "  make build        - Build the application binary"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run Docker container"

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod download
	go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g main.go --output docs

# Run the application
run: swagger
	@echo "Starting API Gateway..."
	go run main.go

# Build the application
build: swagger
	@echo "Building application..."
	go build -o bin/api-gateway main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf docs/

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t api-gateway:latest .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env api-gateway:latest
