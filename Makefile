# Project variables
BINARY_NAME=mindcards
MAIN_PATH=cmd/app/main.go
DOCKER_COMPOSE=docker-compose.yml

.PHONY: all
all: help

.PHONY: build
build: tidy
	@echo "Building binary..."
	@mkdir -p bin
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

.PHONY: run
run: tidy
	@echo "Starting application..."
	go run $(MAIN_PATH)

.PHONY: up
up: docker-up
	@echo "Starting application..."
	go run $(MAIN_PATH)

.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./internal/...

.PHONY: tidy
tidy:
	@echo "Cleaning dependencies..."
	go mod tidy

.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	go clean -cache -modcache

.PHONY: fix
fix:
	@echo "Fixing Go environment..."
	go clean -cache -modcache -i -r
	go mod tidy

.PHONY: docker-up
docker-up:
	@echo "Starting Docker containers..."
	docker-compose -f $(DOCKER_COMPOSE) up -d

.PHONY: docker-down
docker-down:
	@echo "Stopping Docker containers..."
	docker-compose -f $(DOCKER_COMPOSE) down

.PHONY: help
help:
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'