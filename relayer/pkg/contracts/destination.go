package contracts

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// DestinationMessengerABI is the ABI of the DestinationMessenger contract
const DestinationMessengerABI = `[{"inputs":[{"internalType":"address","name":"_relayer","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"AlreadyProcessed","type":"error"},{"inputs":[],"name":"InvalidSourceChain","type":"error"},{"inputs":[],"name":"OnlyRelayer","type":"error"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"messageHash","type":"bytes32"},{"indexed":true,"internalType":"uint256","name":"sourceChainId","type":"uint256"},{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"bytes","name":"payload","type":"bytes"},{"indexed":false,"internalType":"uint256","name":"nonce","type":"uint256"}],"name":"MessageReceived","type":"event"},{"inputs":[{"internalType":"bytes32","name":"_messageHash","type":"bytes32"}],"name":"isProcessed","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"processedMessages","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"_nonce","type":"uint256"},{"internalType":"uint256","name":"_sourceChainId","type":"uint256"},{"internalType":"address","name":"_sender","type":"address"},{"internalType":"bytes","name":"_payload","type":"bytes"},{"internalType":"uint256","name":"_timestamp","type":"uint256"}],"name":"receiveMessage","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"receivedCount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"relayer","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_newRelayer","type":"address"}],"name":"updateRelayer","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

// DestinationMessenger is Go binding for the DestinationMessenger contract
type DestinationMessenger struct {
	DestinationMessengerCaller
	DestinationMessengerTransactor
	DestinationMessengerFilterer
}

type DestinationMessengerCaller struct {
	contract *bind.BoundContract
}

type DestinationMessengerTransactor struct {
	contract *bind.BoundContract
}

type DestinationMessengerFilterer struct {
	contract *bind.BoundContract
}

// NewDestinationMessenger creates a new instance of DestinationMessenger bound to a contract
func NewDestinationMessenger(address common.Address, backend bind.ContractBackend) (*DestinationMessenger, error) {
	parsed, err := abi.JSON(strings.NewReader(DestinationMessengerABI))
	if err != nil {
		return nil, err
	}
	contract := bind.NewBoundContract(address, parsed, backend, backend, backend)
	return &DestinationMessenger{
		DestinationMessengerCaller:     DestinationMessengerCaller{contract: contract},
		DestinationMessengerTransactor: DestinationMessengerTransactor{contract: contract},
		DestinationMessengerFilterer:   DestinationMessengerFilterer{contract: contract},
	}, nil
}

// ReceiveMessage relays a message to the destination chain
func (t *DestinationMessengerTransactor) ReceiveMessage(
	opts *bind.TransactOpts,
	nonce *big.Int,
	sourceChainId *big.Int,
	sender common.Address,
	payload []byte,
	timestamp *big.Int,
) (*types.Transaction, error) {
	return t.contract.Transact(opts, "receiveMessage", nonce, sourceChainId, sender, payload, timestamp)
}

// IsProcessed checks if a message has been processed
func (c *DestinationMessengerCaller) IsProcessed(opts *bind.CallOpts, messageHash [32]byte) (bool, error) {
	var out []interface{}
	err := c.contract.Call(opts, &out, "isProcessed", messageHash)
	if err != nil {
		return false, err
	}
	return *abi.ConvertType(out[0], new(bool)).(*bool), nil
}

// Relayer returns the relayer address
func (c *DestinationMessengerCaller) Relayer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := c.contract.Call(opts, &out, "relayer")
	if err != nil {
		return common.Address{}, err
	}
	return *abi.ConvertType(out[0], new(common.Address)).(*common.Address), nil
}

// ReceivedCount returns the count of messages received from an address
func (c *DestinationMessengerCaller) ReceivedCount(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := c.contract.Call(opts, &out, "receivedCount", addr)
	if err != nil {
		return nil, err
	}
	return *abi.ConvertType(out[0], new(*big.Int)).(**big.Int), nil
}

// Convenience methods on the main struct
func (d *DestinationMessenger) ReceiveMessage(
	opts *bind.TransactOpts,
	nonce *big.Int,
	sourceChainId *big.Int,
	sender common.Address,
	payload []byte,
	timestamp *big.Int,
) (*types.Transaction, error) {
	return d.DestinationMessengerTransactor.ReceiveMessage(opts, nonce, sourceChainId, sender, payload, timestamp)
}

func (d *DestinationMessenger) IsProcessed(opts *bind.CallOpts, messageHash [32]byte) (bool, error) {
	return d.DestinationMessengerCaller.IsProcessed(opts, messageHash)
}

func (d *DestinationMessenger) Relayer(opts *bind.CallOpts) (common.Address, error) {
	return d.DestinationMessengerCaller.Relayer(opts)
}
