package types

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type CrossChainMessage struct {
	Nonce            *big.Int
	SourceChainID    *big.Int
	DestChainID      *big.Int
	Sender           common.Address
	Payload          []byte
	Timestamp        *big.Int
	MessageHash      common.Hash
	SourceTxHash     common.Hash
	DestTxHash       common.Hash
	Status           MessageStatus
	CreatedAt        time.Time
	ProcessedAt      *time.Time
	RetryCount       int
	LastRetryAt      *time.Time
}

type MessageStatus string

const (
	StatusPending   MessageStatus = "pending"
	StatusRelaying  MessageStatus = "relaying"
	StatusCompleted MessageStatus = "completed"
	StatusFailed    MessageStatus = "failed"
)

type ChainConfig struct {
	Name            string
	ChainID         *big.Int
	RpcURL          string
	SourceContract  common.Address
	DestContract    common.Address
	StartBlock      uint64
	Confirmations   uint64
}