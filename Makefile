# =============================================================================
# Makefile — API Gateway
# =============================================================================

APP_NAME   := api-gateway
BIN_DIR    := bin
DIST_DIR   := dist
DOCKER_IMG := api-gateway

.PHONY: help install build run test clean dist docker-build

# ── Default ──────────────────────────────────────────────────────────────────
help:
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@echo "Development"
	@echo "  install        Download Go module dependencies"
	@echo "  build          Compile binary to bin/$(APP_NAME)"
	@echo "  run            Build then run locally (reads .env)"
	@echo "  test           Run all tests"
	@echo "  clean          Remove build and dist artifacts"
	@echo ""
	@echo "Deployment (bare-metal / systemd)"
	@echo "  dist           Build release package → dist/$(APP_NAME)/"
	@echo "                 Includes: binary, templates/, docs/, .env.example, service.sh"
	@echo ""
	@echo "Docker"
	@echo "  docker-build   Build Docker image ($(DOCKER_IMG):latest)"
	@echo "  docker-dist    Build Docker release package → dist/docker/"
	@echo "                 Includes: Dockerfile, docker-compose.yml, .env.example"
	@echo ""

# ── Dependencies ──────────────────────────────────────────────────────────────
install:
	@echo "→ Downloading dependencies..."
	go mod download
	go mod tidy

# ── Build ─────────────────────────────────────────────────────────────────────
build:
	@echo "→ Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	go build -ldflags="-s -w" -o $(BIN_DIR)/$(APP_NAME) main.go
	@echo "✓ Binary: $(BIN_DIR)/$(APP_NAME)"

# ── Run locally ───────────────────────────────────────────────────────────────
run: build
	@echo "→ Starting $(APP_NAME)..."
	./$(BIN_DIR)/$(APP_NAME)

# ── Tests ─────────────────────────────────────────────────────────────────────
test:
	@echo "→ Running tests..."
	go test -v ./...

# ── Clean ─────────────────────────────────────────────────────────────────────
clean:
	@echo "→ Cleaning..."
	rm -rf $(BIN_DIR) $(DIST_DIR)

# ── Dist (bare-metal / systemd package) ──────────────────────────────────────
dist: build
	@echo "→ Creating release package..."
	@rm -rf $(DIST_DIR)/$(APP_NAME)
	@mkdir -p $(DIST_DIR)/$(APP_NAME)

	# Binary
	cp $(BIN_DIR)/$(APP_NAME)    $(DIST_DIR)/$(APP_NAME)/$(APP_NAME)

	# Runtime assets
	cp -r templates              $(DIST_DIR)/$(APP_NAME)/templates
	cp -r docs                   $(DIST_DIR)/$(APP_NAME)/docs

	# Config & scripts
	cp .env.example              $(DIST_DIR)/$(APP_NAME)/.env.example
	cp service.sh                $(DIST_DIR)/$(APP_NAME)/service.sh
	chmod +x                     $(DIST_DIR)/$(APP_NAME)/service.sh

	@echo ""
	@echo "✓ Release package ready: $(DIST_DIR)/$(APP_NAME)/"
	@echo ""
	@echo "  Deploy steps:"
	@echo "    1. Copy $(DIST_DIR)/$(APP_NAME)/ to the server"
	@echo "    2. cp .env.example .env  →  fill in your values"
	@echo "    3. sudo ./service.sh install"
	@echo "    4. sudo ./service.sh start"

# ── Docker build ──────────────────────────────────────────────────────────────
docker-build:
	@echo "→ Building Docker image $(DOCKER_IMG):latest..."
	docker build -t $(DOCKER_IMG):latest .
	@echo "✓ Image built: $(DOCKER_IMG):latest"

# ── Docker dist (deployment package) ─────────────────────────────────────────
docker-dist:
	@echo "→ Creating Docker release package..."
	@rm -rf $(DIST_DIR)/docker
	@mkdir -p $(DIST_DIR)/docker

	cp Dockerfile                $(DIST_DIR)/docker/Dockerfile
	cp docker-compose.yml        $(DIST_DIR)/docker/docker-compose.yml
	cp .env.example              $(DIST_DIR)/docker/.env.example

	@echo ""
	@echo "✓ Docker package ready: $(DIST_DIR)/docker/"
	@echo ""
	@echo "  Deploy steps:"
	@echo "    1. Copy $(DIST_DIR)/docker/ to the server"
	@echo "    2. cp .env.example .env  →  fill in your values"
	@echo "    3. docker compose up -d"
