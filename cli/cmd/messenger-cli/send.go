package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"

	"cli/pkg/contracts"
)

var (
	destChainID int64
	message     string
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a cross-chain message",
	Long:  "Send a message from one blockchain to another",
	Example: `  messenger-cli send \
    --rpc https://sepolia // ky
    --contract 0x1234... \
    --key YOUR_PRIVATE_KEY \
    --dest-chain 80001 \
    --message "Hello!"`,
	RunE: runSend,
}

func init() {
	sendCmd.Flags().StringVar(&rpcURL, "rpc", "", "RPC URL of source chain (required)")
	sendCmd.Flags().StringVar(&contractAddr, "contract", "", "Source messenger contract address (required)")
	sendCmd.Flags().StringVar(&privateKey, "key", "", "Private key (required)")
	sendCmd.Flags().Int64Var(&destChainID, "dest-chain", 0, "Destination chain ID (required)")
	sendCmd.Flags().StringVar(&message, "message", "", "Message to send (required)")

	sendCmd.MarkFlagRequired("rpc")
	sendCmd.MarkFlagRequired("contract")
	sendCmd.MarkFlagRequired("key")
	sendCmd.MarkFlagRequired("dest-chain")
	sendCmd.MarkFlagRequired("message")
}

func runSend(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Connect to chain
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC: %w", err)
	}
	defer client.Close()

	// Get chain ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	fmt.Printf("Connected to chain ID: %s\n", chainID.String())

	// Parse private key
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return fmt.Errorf("invalid private key: %w", err)
	}

	// Create transactor
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}

	// Load contract
	contract, err := contracts.NewSourceMessenger(
		common.HexToAddress(contractAddr),
		client,
	)
	if err != nil {
		return fmt.Errorf("failed to load contract: %w", err)
	}

	fmt.Printf("\n Sending message to chain %d...\n", destChainID)
	fmt.Printf("Message: %s\n\n", message)

	// Send message
	tx, err := contract.SendMessage(
		auth,
		big.NewInt(destChainID),
		[]byte(message),
	)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	fmt.Printf(" Transaction sent!\n")
	fmt.Printf("Tx Hash: %s\n", tx.Hash().Hex())

	// Wait for receipt
	fmt.Printf("\n Waiting for confirmation...\n")
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	if receipt.Status == 1 {
		fmt.Printf(" Transaction confirmed in block %d\n", receipt.BlockNumber.Uint64())
		fmt.Printf("\n Message sent successfully!\n")
		fmt.Printf("The relayer will now pick it up and deliver it to chain %d\n", destChainID)
	} else {
		return fmt.Errorf("transaction reverted")
	}

	return nil
}
