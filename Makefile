.PHONY: build run test test.coverage fmt lint clean help

BINARY_NAME := uncomment-cli
MAIN_PATH   := .
BIN_DIR     := bin
BIN_PATH    := $(BIN_DIR)/$(BINARY_NAME)

build: ## Build the CLI binary
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_PATH) $(MAIN_PATH)

run: ## Run the CLI using go run
	go run $(MAIN_PATH)

run.bin: build ## Run the built binary
	./$(BIN_PATH)

test: ## Run all tests
	go test -v ./...

test.coverage: ## Run tests with coverage report
	go test -v ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

test.html: test.coverage ## Generate HTML coverage report
	go tool cover -html=coverage.out -o coverage.html
	@echo "Отчёт: coverage.html"

test.show: test.html ## Show HTML coverage report in browser
	go tool cover -html=coverage.out

fmt: ## Format code with go fmt
	go fmt ./...

lint: ## Run all linters (go vet, golangci-lint, gosec)
	go vet ./...
	golangci-lint run ./...
	gosec ./...

deps: ## Tidy module dependencies
	go mod tidy

clean: ## Remove build artifacts, coverage files and test cache
	rm -rf $(BIN_DIR)
	rm -f coverage.out coverage.html
	go clean -testcache

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'