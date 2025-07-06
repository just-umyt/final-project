APP_NAME := cart
BIN_DIR := bin
PORT := 8080

.PHONY: build run test clean

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/main.go

run: build
	@echo "Running $(APP_NAME) on port $(PORT)..."
	@./$(BIN_DIR)/$(APP_NAME) -port $(PORT)

test:
	@echo "Testing $(APP_NAME)..."
	@go test -v ./...

clean:
	@echo "Cleaning $(APP_NAME)..."
	@rm -rf $(BIN_DIR)
	@go clean
	
test-env-up:
	@echo "Test-environments up..."
	@docker compose -f docker-compose.test.yaml up -d

test-env-down:
	@echo "Test-environments down..."
	@docker compose -f docker-compose.test.yaml down -v

test-int:
	@echo "Testing integration test..."
	@INTEGRATION_TEST=1 go test -v ./integration

test-unit:
	@echo "Testing unit tests..."
	@go test -v -cover  ./internal/router/http/controller/
	@go test -v -cover  ./internal/usecase/

