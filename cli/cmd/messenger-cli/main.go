package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rpcURL       string
	contractAddr string
	privateKey   string
	configPath   string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "messenger-cli",
		Short: "Cross-chain messenger CLI",
		Long:  "A CLI tool for sending and tracking cross-chain messages",
	}

	rootCmd.PersistentFlags().StringVar(&configPath, "config", "config.yaml", "Path to config file")

	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(statusCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
