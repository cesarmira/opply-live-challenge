BINARY := bin/server
PKG    := ./cmd/server
PORT   ?= 8080

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the server binary into bin/
	go build -o $(BINARY) $(PKG)

.PHONY: run
run: ## Run the server (PORT=8080 by default)
	PORT=$(PORT) go run $(PKG)

.PHONY: test
test: ## Run unit tests
	go test ./...

.PHONY: smoke
smoke: ## Build, boot the server, and verify a known request
	@PORT=$(PORT) ./scripts/smoke.sh

.PHONY: fmt
fmt: ## Format all Go code
	go fmt ./...

.PHONY: tidy
tidy: ## Tidy go.mod / go.sum
	go mod tidy

.PHONY: build-lambda
build-lambda: ## Build Linux/amd64 binary and zip for Lambda deployment
	./scripts/build-lambda.sh

.PHONY: deploy
deploy: build-lambda ## Build and deploy to AWS Lambda via Terraform
	cd terraform && terraform apply

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf bin dist
