-include .env
export

COMPOSE = docker compose -f docker-compose.yml
APP = mindcards
BIN = bin/$(APP)
MAIN_PATH = cmd/app/main.go
DB_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: migrate-up migrate-down migrate-create migrate-version
migrate-up:
	migrate -path ./internal/repository/db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path ./internal/repository/db/migrations -database "$(DB_URL)" down 1

migrate-create:
ifeq ($(strip $(NAME)),)
	@echo Usage: make migrate-create NAME=some_name
	@exit 1
endif
	migrate create -ext sql -dir internal/repository/db/migrations -seq $(NAME)

migrate-version:
	migrate -path ./internal/repository/db/migrations -database "$(DB_URL)" version

build: deps
	@mkdir -p bin
	go build -o $(BIN) $(MAIN_PATH)

run:
	go run $(MAIN_PATH)

deps:
	go mod tidy
	go mod download

test:
	go test -v ./internal/...

clean:
	rm -rf bin/

compose-build:
	$(COMPOSE) build


up:
	$(COMPOSE) up -d --build

up-infra:
	$(COMPOSE) up -d test-db mind-redis

down:
	$(COMPOSE) down

down-full:
	$(COMPOSE) down -v

logs:
	$(COMPOSE) logs -f

db:
	$(COMPOSE) exec test-db psql -U $(DB_USER) $(DB_NAME)

tables:
	$(COMPOSE) exec test-db psql -U $(DB_USER) $(DB_NAME) -c "\dt"

tidy:
	go mod tidy

dev: up-infra build
	./$(BIN)

reset: down-full compose-build up

.PHONY: tidy build run deps test clean \
        compose-build up up-infra down down-full logs \
        db tables dev reset