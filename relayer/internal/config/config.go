package config

import (
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Chains  []ChainConfig `yaml:"chains"`
	Relayer RelayerConfig `yaml:"relayer"`
}

type ChainConfig struct {
	Name           string `yaml:"name"`
	ChainID        int64  `yaml:"chain_id"`
	RpcURL         string `yaml:"rpc_url"`
	SourceContract string `yaml:"source_contract"`
	DestContract   string `yaml:"dest_contract"`
	StartBlock     uint64 `yaml:"start_block"`
	Confirmations  uint64 `yaml:"confirmations"`
}

type RelayerConfig struct {
	PrivateKey   string `yaml:"private_key"`
	PollInterval string `yaml:"poll_interval"`
	MaxRetries   int    `yaml:"max_retries"`
	GasLimit     uint64 `yaml:"gas_limit"`
	DBPath       string `yaml:"db_path"`
}

func LoadConfig(path string) (*Config, error) {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// Expand environment variables like ${VAR} or $VAR
	expanded := os.ExpandEnv(string(data))

	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

func (c *ChainConfig) GetChainID() *big.Int {
	return big.NewInt(c.ChainID)
}

func (c *ChainConfig) GetSourceContract() common.Address {
	return common.HexToAddress(c.SourceContract)
}

func (c *ChainConfig) GetDestContract() common.Address {
	return common.HexToAddress(c.DestContract)
}
