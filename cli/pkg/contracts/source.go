package contracts

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// SourceMessengerABI is the ABI of the SourceMessenger contract
const SourceMessengerABI = `[{"inputs":[{"internalType":"uint256","name":"_destChainId","type":"uint256"},{"internalType":"bytes","name":"_payload","type":"bytes"}],"name":"sendMessage","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"_nonce","type":"uint256"},{"internalType":"uint256","name":"_sourceChainId","type":"uint256"},{"internalType":"uint256","name":"_destChainId","type":"uint256"},{"internalType":"address","name":"_sender","type":"address"},{"internalType":"bytes","name":"_payload","type":"bytes"},{"internalType":"uint256","name":"_timestamp","type":"uint256"}],"name":"getMessageHash","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"bytes32","name":"_messageHash","type":"bytes32"}],"name":"verifyMessage","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"nonce","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"messageExists","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"nonce","type":"uint256"},{"indexed":true,"internalType":"uint256","name":"destinationChainId","type":"uint256"},{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"bytes","name":"payload","type":"bytes"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"MessageSent","type":"event"}]`

// SourceMessenger is Go binding for the SourceMessenger contract
type SourceMessenger struct {
	SourceMessengerCaller
	SourceMessengerTransactor
	SourceMessengerFilterer
}

type SourceMessengerCaller struct {
	contract *bind.BoundContract
}

type SourceMessengerTransactor struct {
	contract *bind.BoundContract
}

type SourceMessengerFilterer struct {
	contract *bind.BoundContract
}

// SourceMessengerMessageSent represents a MessageSent event
type SourceMessengerMessageSent struct {
	Nonce              *big.Int
	DestinationChainId *big.Int
	Sender             common.Address
	Payload            []byte
	Timestamp          *big.Int
	Raw                types.Log
}

// SourceMessengerMessageSentIterator is returned from FilterMessageSent
type SourceMessengerMessageSentIterator struct {
	Event    *SourceMessengerMessageSent
	contract *bind.BoundContract
	event    string
	logs     chan types.Log
	sub      event.Subscription
	done     bool
	fail     error
}

func (it *SourceMessengerMessageSentIterator) Next() bool {
	if it.fail != nil {
		return false
	}
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SourceMessengerMessageSent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true
		default:
			return false
		}
	}
	select {
	case log := <-it.logs:
		it.Event = new(SourceMessengerMessageSent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true
	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *SourceMessengerMessageSentIterator) Error() error {
	return it.fail
}

func (it *SourceMessengerMessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NewSourceMessenger creates a new instance of SourceMessenger bound to a contract
func NewSourceMessenger(address common.Address, backend bind.ContractBackend) (*SourceMessenger, error) {
	parsed, err := abi.JSON(strings.NewReader(SourceMessengerABI))
	if err != nil {
		return nil, err
	}
	contract := bind.NewBoundContract(address, parsed, backend, backend, backend)
	return &SourceMessenger{
		SourceMessengerCaller:     SourceMessengerCaller{contract: contract},
		SourceMessengerTransactor: SourceMessengerTransactor{contract: contract},
		SourceMessengerFilterer:   SourceMessengerFilterer{contract: contract},
	}, nil
}

// SendMessage sends a cross-chain message
func (t *SourceMessengerTransactor) SendMessage(opts *bind.TransactOpts, destChainId *big.Int, payload []byte) (*types.Transaction, error) {
	return t.contract.Transact(opts, "sendMessage", destChainId, payload)
}

// Nonce returns the current nonce
func (c *SourceMessengerCaller) Nonce(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := c.contract.Call(opts, &out, "nonce")
	if err != nil {
		return nil, err
	}
	return *abi.ConvertType(out[0], new(*big.Int)).(**big.Int), nil
}

// ParseMessageSent parses a MessageSent event from a log
func (f *SourceMessengerFilterer) ParseMessageSent(log types.Log) (*SourceMessengerMessageSent, error) {
	event := new(SourceMessengerMessageSent)
	if err := f.contract.UnpackLog(event, "MessageSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Convenience method on the main struct
func (s *SourceMessenger) SendMessage(opts *bind.TransactOpts, destChainId *big.Int, payload []byte) (*types.Transaction, error) {
	return s.SourceMessengerTransactor.SendMessage(opts, destChainId, payload)
}

func (s *SourceMessenger) ParseMessageSent(log types.Log) (*SourceMessengerMessageSent, error) {
	return s.SourceMessengerFilterer.ParseMessageSent(log)
}
