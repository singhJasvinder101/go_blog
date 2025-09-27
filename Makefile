# this file will help you setup the project quickly by using 
# make commands in your terminal
# e.g make build to build the go binary

APP_NAME := go_blog
BIN_PATH := ./tmp/main
SRC_PATH := ./cmd/api
CONFIG := config/local.yml


# test makefile formatting of tabs
check:
	cat -e -t -v Makefile

# build go binary (^I denotes tab and $ denotes line ending)
build:
	go build -o $(BIN_PATH) $(SRC_PATH)

# run the dev server
run-dev:
	docker stop postgres || true
	docker stop redis || true
	docker compose -f docker/docker-compose.yml up -d
	air

run-test:
	docker stop postgres || true
	docker stop redis || true
	docker compose -f docker/docker-compose.yml up -d
	go run $(SRC_PATH) -config $(CONFIG)

# clean up binary
clean:
	rm -f $(BIN_PATH)

# run tests
test:
	go test ./...


.PHONY: build run clean test check