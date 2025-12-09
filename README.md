# Cross-Chain Messenger

A production-ready cross-chain messaging protocol enabling secure message passing between EVM-compatible blockchains. The system consists of Solidity smart contracts, a Go-based relayer service, and a command-line interface for sending and tracking messages.

## Architecture

The project comprises three main components:

**Smart Contracts**: Deployed on source and destination chains to emit and receive cross-chain messages.

**Relayer**: A Go service that monitors source chains for message events and relays them to destination chains.

**CLI**: A command-line tool for sending messages and checking their processing status.

### How It Works

1. Users call `sendMessage()` on the SourceMessenger contract
2. The contract emits a `MessageSent` event with message details
3. The relayer listens for these events on configured source chains
4. Upon detection, the relayer submits the message to the DestinationMessenger contract
5. The destination contract verifies and processes the message, emitting a `MessageReceived` event

## Prerequisites

- Go 1.22 or higher
- Foundry (for smart contract development)
- Docker and Docker Compose (optional, for containerized deployment)
- Kubernetes cluster (optional, for production deployment)
- RPC access to supported chains (Ethereum Sepolia and Polygon Amoy testnets)

## Installation

### Clone the Repository

```bash
git clone https://github.com/yourusername/cross-chain-messenger.git
cd cross-chain-messenger
```

### Install Dependencies

```bash
make install
```

This command installs:
- Foundry dependencies for smart contracts
- Go modules for the relayer
- Go modules for the CLI

## Smart Contracts

### Compile Contracts

```bash
cd contracts
forge build
```

### Run Tests

```bash
forge test
```

### Deploy Contracts

Deploy to Ethereum Sepolia:

```bash
make deploy-contracts-sepolia
```

Deploy to Polygon Amoy:

```bash
make deploy-contracts-amoy
```

Note: Update deployment scripts with your private key and RPC URLs before deploying.

## Configuration

### Environment Variables

Create a `.env` file in the `relayer` directory:

```bash
cp relayer/.env.example relayer/.env
```

Edit the file with your values:

```
RELAYER_PRIVATE_KEY=your_private_key_without_0x_prefix
SEPOLIA_RPC_URL=https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
AMOY_RPC_URL=https://polygon-amoy.g.alchemy.com/v2/YOUR_API_KEY
```

### Configuration File

Update `relayer/config.yaml` with deployed contract addresses:

```yaml
chains:
  - name: "sepolia"
    chain_id: 11155111
    rpc_url: "${SEPOLIA_RPC_URL}"
    source_contract: "0xYourSepoliaSourceContractAddress"
    dest_contract: "0xYourSepoliaDestContractAddress"
    start_block: 5000000
    confirmations: 3

  - name: "amoy"
    chain_id: 80002
    rpc_url: "${AMOY_RPC_URL}"
    source_contract: "0xYourAmoySourceContractAddress"
    dest_contract: "0xYourAmoyDestContractAddress"
    start_block: 1000000
    confirmations: 5

relayer:
  private_key: "${RELAYER_PRIVATE_KEY}"
  poll_interval: "5s"
  max_retries: 3
  gas_limit: 300000
  db_path: "./data/messages.db"
```

## Running the Relayer

### Build the Relayer

```bash
cd relayer
go build -o relayerd ./cmd/relayerd
```

### Run the Relayer

```bash
./relayerd
```

The relayer will:
- Load configuration from `config.yaml`
- Connect to all configured RPC endpoints
- Begin monitoring for cross-chain message events
- Automatically relay messages to destination chains

### Expected Output

```
2025/12/09 16:00:00 Relayer address: 0x14dC79964da2C08b23698B3D3cc7Ca32193d9955
2025/12/09 16:00:00 Connected to sepolia (Chain ID: 11155111)
2025/12/09 16:00:00 Connected to amoy (Chain ID: 80002)
2025/12/09 16:00:00 Relayer started successfully!
```

## Using the CLI

### Build the CLI

```bash
cd cli
go build -o messenger-cli ./cmd/messenger-cli
```

### Send a Cross-Chain Message

```bash
./messenger-cli send \
  --rpc https://eth-sepolia.g.alchemy.com/v2/YOUR_KEY \
  --contract 0xYourSourceContractAddress \
  --key your_private_key_without_0x \
  --dest-chain 80002 \
  --message "Hello from Sepolia to Amoy!"
```

### Check Message Status

```bash
./messenger-cli status \
  --rpc https://polygon-amoy.g.alchemy.com/v2/YOUR_KEY \
  --contract 0xYourDestContractAddress \
  --hash 0xYourMessageHash
```

## Docker Deployment

### Build Docker Images

```bash
make docker-build
```

### Run with Docker Compose

```bash
# Copy and configure environment variables
cp .env.template .env
# Edit .env with your values

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f relayer

# Stop services
docker-compose down
```

The Docker Compose setup includes:
- Relayer service
- Prometheus for metrics collection
- Grafana for monitoring dashboards

## Kubernetes Deployment

### Prerequisites

- Kubernetes cluster access
- kubectl configured
- Container registry access

### Prepare Secrets

Edit `k8s/secrets.yaml` with your actual values, then apply:

```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/secrets.yaml
```

### Deploy the Relayer

```bash
kubectl apply -f k8s/
```

### Verify Deployment

```bash
kubectl get pods -n cross-chain
kubectl logs -f deployment/relayer -n cross-chain
```

## Development

### Project Structure

```
cross-chain-messenger/
├── contracts/          # Solidity smart contracts
│   ├── src/
│   ├── script/
│   └── test/
├── relayer/           # Go relayer service
│   ├── cmd/relayerd/
│   ├── internal/
│   └── pkg/contracts/
├── cli/               # Command-line interface
│   ├── cmd/messenger-cli/
│   └── pkg/contracts/
├── docker/            # Dockerfiles
├── k8s/              # Kubernetes manifests
└── scripts/          # Deployment scripts
```

### Run Tests

```bash
# All tests
make test

# Contract tests only
cd contracts && forge test

# Relayer tests only
cd relayer && go test ./...

# CLI tests only
cd cli && go test ./...
```

### Build All Components

```bash
make build
```

## CI/CD

The project includes GitHub Actions workflows for:

- **contracts-ci.yml**: Compile and test smart contracts
- **relayer-ci.yml**: Build, test, and lint the relayer
- **cli-ci.yml**: Build and test the CLI
- **deploy-contracts.yml**: Automated contract deployment
- **deploy-relayer.yml**: Build Docker images and deploy to staging/production

### Required GitHub Secrets

- `DOCKER_USERNAME`: Docker Hub username
- `DOCKER_PASSWORD`: Docker Hub access token
- `KUBE_CONFIG`: Kubernetes configuration (base64 encoded)
- `AWS_ACCESS_KEY_ID`: AWS credentials (if using ECS)
- `AWS_SECRET_ACCESS_KEY`: AWS credentials (if using ECS)

## Monitoring

### Metrics

The relayer exposes Prometheus metrics on port 9090:

- Message processing rates
- Success/failure counts
- RPC connection status
- Gas costs

### Grafana Dashboards

When using Docker Compose, access Grafana at http://localhost:3000 with credentials:

- Username: admin
- Password: (set in `.env` as `GRAFANA_PASSWORD`)

## Troubleshooting

### Relayer Not Starting

Check that:
- Environment variables are correctly set in `.env`
- RPC URLs are accessible
- Private key is valid (without 0x prefix)
- Contract addresses are correct

### Messages Not Being Relayed

Verify:
- Relayer has sufficient gas on destination chain
- Start block is set correctly in config
- Contract addresses match deployed contracts
- RPC endpoints support WebSocket connections (or switch to polling)

### CLI Commands Failing

Ensure:
- RPC URL is correct and accessible
- Contract address is valid
- Private key has sufficient balance for gas
- Chain ID matches the network

## Security Considerations

- Never commit private keys or sensitive data to version control
- Use environment variables for all secrets
- The relayer wallet should have enough funds for gas but minimal excess
- Implement rate limiting in production deployments
- Use WebSocket RPC endpoints with authentication
- Regularly rotate relayer private keys
- Monitor for unusual activity or failed transactions

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome. Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Submit a pull request with a clear description

## Support

For issues, questions, or feature requests, please open an issue on GitHub.
