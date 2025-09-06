# Default target
.DEFAULT_GOAL := help

# help
.PHONY: help
help: ## Show help messages
	@echo "Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

# format
.PHONY: format
format: ## Format the code
	go fmt ./...

# lint
.PHONY: lint
lint: ## Lint the code
	golangci-lint run

# check
.PHONY: check
check: ## Check the code
	go fmt ./... && golangci-lint run && go test ./...

# mock_gen
.PHONY: mock_gen
mock_gen: ## Generate mock files
	go generate ./...