APP_NAME := stocks
BIN_DIR := bin
PORT := 8081

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

protoc:
	@echo "Generating proto files"
	@protoc -I ../proto/stock \
  --go_out=pkg/api/stock --go_opt=paths=source_relative \
  --go-grpc_out=pkg/api/stock --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=pkg/api/stock --grpc-gateway_opt=paths=source_relative \
  stock.proto
  