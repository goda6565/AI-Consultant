.PHONY: backend-check
backend-check:
	cd backend && go fmt ./... && golangci-lint run

.PHONY: backend-format
backend-format:
	cd backend && go fmt ./...

.PHONY: backend-lint
backend-lint:
	cd backend && golangci-lint run