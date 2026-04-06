# Multi-Chain Staking POC — Makefile
# Run `make help` to list available targets.

DATABASE_URL ?= postgres://staking:staking@localhost:5432/staking?sslmode=disable

.PHONY: help up down db-migrate db-reset sqlc proto build test test-cover \
        test-contracts test-contracts-fuzz deploy-local test-solana test-all \
        lint clean demo

##@ Infrastructure

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}' | \
		sort

up: ## Start all services, wait, then run migrations
	docker compose up -d
	@echo "Waiting for services to be healthy..."
	@sleep 3
	@$(MAKE) db-migrate

down: ## Stop all services
	docker compose down

##@ Database

db-migrate: ## Apply all pending migrations
	cd backend && goose -dir db/migrations postgres "$(DATABASE_URL)" up

db-reset: ## Roll back all migrations then re-apply
	cd backend && goose -dir db/migrations postgres "$(DATABASE_URL)" reset
	@$(MAKE) db-migrate

##@ Code Generation

sqlc: ## Generate type-safe Go code from SQL queries
	cd backend && sqlc generate

proto: ## Lint and generate code from protobuf definitions
	cd backend/proto && buf lint && buf generate

##@ Build

build: ## Compile the gRPC server binary to backend/bin/server
	cd backend && go build -o bin/server ./cmd/server

##@ Testing

test: ## Run all Go unit tests with verbose output
	cd backend && go test ./... -v -count=1

test-cover: ## Run Go tests with coverage report (HTML)
	cd backend && go test ./... -v -count=1 \
		-coverprofile=coverage.out -covermode=atomic
	cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report written to backend/coverage.html"

test-contracts: ## Run EVM contract tests with Forge (verbose)
	cd contracts/evm && forge test -vvv

test-contracts-fuzz: ## Run EVM contract fuzz tests (1000 runs)
	cd contracts/evm && forge test --fuzz-runs 1000 -vvv

deploy-local: ## Deploy EVM contracts to local Anvil node
	cd contracts/evm && forge script Deploy.s.sol \
		--rpc-url http://localhost:8545 \
		--broadcast

deploy-pegasus: ## Deploy EVM contracts to LightLink Pegasus testnet
	cd contracts/evm && PRIVATE_KEY=$(PRIVATE_KEY) forge script script/DeployTestnet.s.sol \
		--rpc-url https://replicator.pegasus.lightlink.io/rpc/v1 \
		--broadcast --legacy

deploy-sepolia: ## Deploy EVM contracts to Sepolia testnet
	cd contracts/evm && PRIVATE_KEY=$(PRIVATE_KEY) forge script script/DeployTestnet.s.sol \
		--rpc-url $(SEPOLIA_RPC_URL) \
		--broadcast --verify

test-solana: ## Run Solana program tests with Anchor
	cd contracts/solana && anchor test

test-all: test test-contracts test-solana ## Run every test suite (Go + EVM + Solana)

##@ Quality

lint: ## Run golangci-lint on the backend
	golangci-lint run ./backend/...

##@ Misc

clean: ## Remove build artifacts
	rm -rf backend/bin
	rm -f backend/coverage.out backend/coverage.html

demo: ## Spin up the stack and open the gRPC UI
	@$(MAKE) up
	@echo "Opening gRPC UI at localhost:8080..."
	grpcui -plaintext localhost:8080
