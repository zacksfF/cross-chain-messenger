package executor

import (
	"context"
	"fmt"
	"log"
	"relayer/internal/config"
	"relayer/internal/signer"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	customTypes "relayer/internal/types"
	"relayer/pkg/contracts"
)

type Executor struct {
	clients     map[int64]*ethclient.Client
	chains      map[int64]*config.ChainConfig
	signer      *signer.Signer
	maxRetries  int
	gasLimit    uint64
	messageChan chan *customTypes.CrossChainMessage
}

func NewExecutor(
	clients map[int64]*ethclient.Client,
	chains map[int64]*config.ChainConfig,
	signer *signer.Signer,
	maxRetries int,
	gasLimit uint64,
	messageChan chan *customTypes.CrossChainMessage,
) *Executor {
	return &Executor{
		clients:     clients,
		chains:      chains,
		signer:      signer,
		maxRetries:  maxRetries,
		gasLimit:    gasLimit,
		messageChan: messageChan,
	}
}

func (e *Executor) Start(ctx context.Context) error {
	log.Println("Starting executor...")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-e.messageChan:
			if err := e.processMessage(ctx, msg); err != nil {
				log.Printf(" Failed to process message: %v", err)
				// Implement retry logic here
			}
		}
	}
}

func (e *Executor) processMessage(ctx context.Context, msg *customTypes.CrossChainMessage) error {
	destChainID := msg.DestChainID.Int64()

	client, ok := e.clients[destChainID]
	if !ok {
		return fmt.Errorf("no client for chain %d", destChainID)
	}

	chainConfig, ok := e.chains[destChainID]
	if !ok {
		return fmt.Errorf("no config for chain %d", destChainID)
	}

	destContract, err := contracts.NewDestinationMessenger(
		chainConfig.GetDestContract(),
		client,
	)
	if err != nil {
		return fmt.Errorf("failed to instantiate destination contract: %w", err)
	}

	// Check if already processed
	processed, err := destContract.IsProcessed(nil, msg.MessageHash)
	if err != nil {
		return fmt.Errorf("failed to check if processed: %w", err)
	}

	if processed {
		log.Printf(" Message already processed: %s", msg.MessageHash.Hex())
		return nil
	}

	// Get transactor
	auth, err := e.signer.GetTransactor(msg.DestChainID)
	if err != nil {
		return fmt.Errorf("failed to get transactor: %w", err)
	}

	auth.GasLimit = e.gasLimit

	log.Printf(" Relaying message to chain %d...", destChainID)

	// Send transaction
	tx, err := destContract.ReceiveMessage(
		auth,
		msg.Nonce,
		msg.SourceChainID,
		msg.Sender,
		msg.Payload,
		msg.Timestamp,
	)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	log.Printf(" Message relayed! Tx: %s", tx.Hash().Hex())

	// Wait for confirmation
	receipt, err := waitForConfirmation(ctx, client, tx.Hash())
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	if receipt.Status == 1 {
		log.Printf("✨ Message confirmed on chain %d", destChainID)
		msg.Status = customTypes.StatusCompleted
		msg.DestTxHash = tx.Hash()
		now := time.Now()
		msg.ProcessedAt = &now
	} else {
		return fmt.Errorf("transaction reverted")
	}

	return nil
}

func waitForConfirmation(ctx context.Context, client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
}
