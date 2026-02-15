APP_NAME := laluer
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP_NAME)

.PHONY: all build run dev test fmt vet tidy clean

all: build

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

build: $(BIN_DIR)
	go build -o $(BIN) .

run:
	go run .

dev:
	air

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -rf $(BIN) tmp
