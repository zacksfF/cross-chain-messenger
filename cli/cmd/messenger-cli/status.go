package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"

	"cli/pkg/contracts"
)

var (
	messageHash  string
	destRPC      string
	destContract string
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check message status",
	Long:  "Check if a message has been processed on the destination chain",
	Example: `  messenger-cli status \
    --rpc https://Amoy.polygonscan.com/... \
    --contract 0x5678... \
    --hash 0xabcd...`,
	RunE: runStatus,
}

func init() {
	statusCmd.Flags().StringVar(&destRPC, "rpc", "", "Destination chain RPC URL (required)")
	statusCmd.Flags().StringVar(&destContract, "contract", "", "Destination contract address (required)")
	statusCmd.Flags().StringVar(&messageHash, "hash", "", "Message hash to check (required)")

	statusCmd.MarkFlagRequired("rpc")
	statusCmd.MarkFlagRequired("contract")
	statusCmd.MarkFlagRequired("hash")
}

func runStatus(cmd *cobra.Command, args []string) error {
	_ = context.Background() // Reserved for future use

	// Connect to destination chain
	client, err := ethclient.Dial(destRPC)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC: %w", err)
	}
	defer client.Close()

	// Load contract
	contract, err := contracts.NewDestinationMessenger(
		common.HexToAddress(destContract),
		client,
	)
	if err != nil {
		return fmt.Errorf("failed to load contract: %w", err)
	}

	// Check if processed
	hash := common.HexToHash(messageHash)
	processed, err := contract.IsProcessed(nil, hash)
	if err != nil {
		return fmt.Errorf("failed to check status: %w", err)
	}

	fmt.Printf("\n Message Status\n")
	fmt.Printf("─────────────────\n")
	fmt.Printf("Hash: %s\n", messageHash)

	if processed {
		fmt.Printf("Status: Processed\n")
	} else {
		fmt.Printf("Status: Pending\n")
	}

	return nil
}
