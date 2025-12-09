package listener

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"relayer/internal/config"
	customTypes "relayer/internal/types"
	"relayer/pkg/contracts"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Listener struct {
	client         *ethclient.Client
	chainConfig    *config.ChainConfig
	sourceContract *contracts.SourceMessenger
	messageChan    chan *customTypes.CrossChainMessage
}

func NewListener(
	client *ethclient.Client,
	chainConfig *config.ChainConfig,
	messageChan chan *customTypes.CrossChainMessage,
) (*Listener, error) {
	sourceContract, err := contracts.NewSourceMessenger(
		chainConfig.GetSourceContract(),
		client,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate contract: %w", err)
	}

	return &Listener{
		client:         client,
		chainConfig:    chainConfig,
		sourceContract: sourceContract,
		messageChan:    messageChan,
	}, nil
}

func (l *Listener) Start(ctx context.Context) error {
	log.Printf("Starting listener for %s (Chain ID: %d)", l.chainConfig.Name, l.chainConfig.ChainID)

	// Subscribe to new blocks
	headers := make(chan *types.Header)
	sub, err := l.client.SubscribeNewHead(ctx, headers)
	if err != nil {
		return fmt.Errorf("failed to subscribe to new heads: %w", err)
	}
	defer sub.Unsubscribe()

	fromBlock := l.chainConfig.StartBlock

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-sub.Err():
			return fmt.Errorf("subscription error: %w", err)
		case header := <-headers:
			// Wait for confirmations
			confirmedBlock := header.Number.Uint64()
			if confirmedBlock < l.chainConfig.Confirmations {
				continue
			}
			confirmedBlock -= l.chainConfig.Confirmations

			if fromBlock > confirmedBlock {
				continue
			}

			// Query logs
			if err := l.processBlocks(ctx, fromBlock, confirmedBlock); err != nil {
				log.Printf("Error processing blocks: %v", err)
				continue
			}

			fromBlock = confirmedBlock + 1
		}
	}
}

func (l *Listener) processBlocks(ctx context.Context, from, to uint64) error {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(from)),
		ToBlock:   big.NewInt(int64(to)),
		Addresses: []common.Address{l.chainConfig.GetSourceContract()},
	}

	logs, err := l.client.FilterLogs(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to filter logs: %w", err)
	}

	for _, vLog := range logs {
		if err := l.handleLog(vLog); err != nil {
			log.Printf("Error handling log: %v", err)
		}
	}

	return nil
}

func (l *Listener) handleLog(vLog types.Log) error {
	// Parse MessageSent event
	event, err := l.sourceContract.ParseMessageSent(vLog)
	if err != nil {
		return fmt.Errorf("failed to parse event: %w", err)
	}

	message := &customTypes.CrossChainMessage{
		Nonce:         event.Nonce,
		SourceChainID: l.chainConfig.GetChainID(),
		DestChainID:   event.DestinationChainId,
		Sender:        event.Sender,
		Payload:       event.Payload,
		Timestamp:     event.Timestamp,
		SourceTxHash:  vLog.TxHash,
		Status:        customTypes.StatusPending,
		CreatedAt:     time.Now(),
		RetryCount:    0,
	}

	// Compute message hash
	messageHash := computeMessageHash(
		event.Nonce,
		l.chainConfig.GetChainID(),
		event.DestinationChainId,
		event.Sender,
		event.Payload,
		event.Timestamp,
	)
	message.MessageHash = messageHash

	log.Printf(" New message detected: Nonce=%s, From=%s, To Chain=%s",
		event.Nonce.String(), event.Sender.Hex(), event.DestinationChainId.String())

	// Send to executor
	l.messageChan <- message

	return nil
}

func computeMessageHash(
	nonce, sourceChainID, destChainID *big.Int,
	sender common.Address,
	payload []byte,
	timestamp *big.Int,
) common.Hash {
	return crypto.Keccak256Hash(
		nonce.Bytes(),
		sourceChainID.Bytes(),
		destChainID.Bytes(),
		sender.Bytes(),
		payload,
		timestamp.Bytes(),
	)
}
