.PHONY: help install test build deploy clean

help:
	@echo "Cross-Chain Messenger - Make targets"
	@echo ""
	@echo "  install       Install all dependencies"
	@echo "  test          Run all tests"
	@echo "  build         Build all components"
	@echo "  deploy        Deploy contracts and services"
	@echo "  clean         Clean build artifacts"

install:
	@echo "Installing dependencies..."
	cd contracts && forge install
	cd relayer && go mod download
	cd cli && go mod download

test:
	@echo "Running tests..."
	cd contracts && forge test
	cd relayer && go test ./...
	cd cli && go test ./...

build:
	@echo "Building..."
	cd contracts && forge build
	cd relayer && go build -o relayerd ./cmd/relayerd
	cd cli && go build -o messenger-cli ./cmd/messenger-cli

deploy-contracts-sepolia:
	./scripts/deploy-contracts.sh sepolia

deploy-contracts-amoy:
	./scripts/deploy-contracts.sh amoy

docker-build:
	docker build -t cross-chain-relayer:latest -f docker/relayer.Dockerfile ./relayer
	docker build -t messenger-cli:latest -f docker/cli.Dockerfile ./cli

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

k8s-deploy:
	kubectl apply -f k8s/

clean:
	cd contracts && forge clean
	cd relayer && rm -f relayerd
	cd cli && rm -f messenger-cli