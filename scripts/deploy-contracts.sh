#!/bin/bash

set -e

NETWORK=$1
VERIFY=${2:-true}

if [ -z "$NETWORK" ]; then
  echo "Usage: ./deploy-contracts.sh <network> [verify]"
  exit 1
fi

echo "Deploying contracts to $NETWORK..."

cd contracts

# Load environment variables
source ../.env.$NETWORK

# Deploy
forge script script/Deploy.s.sol:DeployScript \
  --rpc-url $RPC_URL \
  --broadcast \
  $([ "$VERIFY" = "true" ] && echo "--verify") \
  -vvvv

# Save deployment addresses
DEPLOYMENT_FILE="../deployments/$NETWORK.json"
mkdir -p ../deployments
cp broadcast/Deploy.s.sol/**/run-latest.json $DEPLOYMENT_FILE

echo "Deployment complete!"
echo "Deployment info saved to: $DEPLOYMENT_FILE"