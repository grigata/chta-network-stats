package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type RPCConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Config struct {
	RPC RPCConfig `json:"rpc"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if cfg.RPC.Host == "" {
		return nil, fmt.Errorf("rpc host is required")
	}
	if cfg.RPC.Port <= 0 {
		return nil, fmt.Errorf("rpc port must be greater than zero")
	}

	return &cfg, nil
}
