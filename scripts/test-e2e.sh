#!/bin/bash

set -e

echo "Running E2E tests..."

# Start local chains
echo "Starting Anvil instances..."
anvil --port 8545 --chain-id 31337 &
ANVIL_PID1=$!
anvil --port 8546 --chain-id 31338 &
ANVIL_PID2=$!

sleep 2

# Deploy contracts
echo "Deploying contracts..."
cd contracts
forge script script/Deploy.s.sol --rpc-url http://localhost:8545 --broadcast
forge script script/Deploy.s.sol --rpc-url http://localhost:8546 --broadcast

# Start relayer
echo "Starting relayer..."
cd ../relayer
go run cmd/relayerd/main.go &
RELAYER_PID=$!

sleep 5

# Send test message
echo "Sending test message..."
cd ../cli
go run cmd/messenger-cli/main.go send \
  --rpc http://localhost:8545 \
  --contract 0x5FbDB2315678afecb367f032d93F642f64180aa3 \
  --key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
  --dest-chain 31338 \
  --message "E2E test"

# Wait for relay
sleep 10

# Check status
go run cmd/messenger-cli/main.go status \
  --rpc http://localhost:8546 \
  --contract 0x5FbDB2315678afecb367f032d93F642f64180aa3 \
  --hash 0x...

# Cleanup
kill $ANVIL_PID1 $ANVIL_PID2 $RELAYER_PID

echo "E2E tests passed!"