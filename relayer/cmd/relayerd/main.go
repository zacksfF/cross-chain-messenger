package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"relayer/internal/config"
	"relayer/internal/executor"
	"relayer/internal/listener"
	"relayer/internal/signer"
	"syscall"

	"github.com/ethereum/go-ethereum/ethclient"

	customTypes "relayer/internal/types"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load config
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize signer
	sign, err := signer.NewSigner(cfg.Relayer.PrivateKey)
	if err != nil {
		log.Fatalf("Failed to create signer: %v", err)
	}

	log.Printf(" Relayer address: %s", sign.GetAddress().Hex())

	// Initialize clients and chains
	clients := make(map[int64]*ethclient.Client)
	chains := make(map[int64]*config.ChainConfig)

	for _, chain := range cfg.Chains {
		client, err := ethclient.Dial(chain.RpcURL)
		if err != nil {
			log.Fatalf("Failed to connect to %s: %v", chain.Name, err)
		}
		clients[chain.ChainID] = client
		chains[chain.ChainID] = &chain
		log.Printf(" Connected to %s (Chain ID: %d)", chain.Name, chain.ChainID)
	}

	// Message channel
	messageChan := make(chan *customTypes.CrossChainMessage, 100)

	// Start listeners
	for _, chain := range cfg.Chains {
		chainListener, err := listener.NewListener(
			clients[chain.ChainID],
			&chain,
			messageChan,
		)
		if err != nil {
			log.Fatalf("Failed to create listener for %s: %v", chain.Name, err)
		}

		go func(l *listener.Listener) {
			if err := l.Start(ctx); err != nil {
				log.Printf("Listener error: %v", err)
			}
		}(chainListener)
	}

	// Start executor
	exec := executor.NewExecutor(
		clients,
		chains,
		sign,
		cfg.Relayer.MaxRetries,
		cfg.Relayer.GasLimit,
		messageChan,
	)

	go func() {
		if err := exec.Start(ctx); err != nil {
			log.Printf("Executor error: %v", err)
		}
	}()

	log.Println(" Relayer started successfully!")

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	cancel()
}
