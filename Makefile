APP_NAME := go_blog
BIN_PATH := ./tmp/main
SRC_PATH := ./cmd/api
CONFIG := config/local.yml

.PHONY: build run clean test

build:
    go build -o $(BIN_PATH) $(SRC_PATH)

run: build
    docker run --name postgres -e POSTGRES_PASSWORD=123 -dp 5432:5432 postgres
    docker exec -it postgres psql -U postgres -c "CREATE DATABASE $(APP_NAME);"
    APP_ENV=dev $(BIN_PATH) -config $(CONFIG)

clean:
    rm -f $(BIN_PATH)

test:
    go test ./...