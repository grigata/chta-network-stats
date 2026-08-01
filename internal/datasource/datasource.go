package datasource

import (
	"context"
	"fmt"
	"strings"

	"github.com/grigata/chta-network-stats/internal/api"
	"github.com/grigata/chta-network-stats/internal/chain"
	"github.com/grigata/chta-network-stats/internal/config"
	"github.com/grigata/chta-network-stats/internal/rpc"
)

const defaultAPIURL = "http://chtaexplorer.mooo.com:3002"

type Result struct {
	Client chain.Client
	Name   string
	Height int64
}

func Connect(ctx context.Context, cfg *config.Config) (*Result, error) {
	mode := strings.ToLower(strings.TrimSpace(cfg.Mode))
	if mode == "" {
		mode = "auto"
	}

	switch mode {
	case "rpc":
		return connectRPC(ctx, cfg)

	case "api":
		return connectAPI(ctx, cfg)

	case "auto":
		rpcResult, rpcErr := connectRPC(ctx, cfg)
		if rpcErr == nil {
			return rpcResult, nil
		}

		apiResult, apiErr := connectAPI(ctx, cfg)
		if apiErr == nil {
			return apiResult, nil
		}

		return nil, fmt.Errorf(
			"no data source available; local RPC: %v; public API: %v",
			rpcErr,
			apiErr,
		)

	default:
		return nil, fmt.Errorf(
			"invalid mode %q: expected auto, rpc, or api",
			cfg.Mode,
		)
	}
}

func connectRPC(ctx context.Context, cfg *config.Config) (*Result, error) {
	client := rpc.NewClient(cfg.RPC)

	height, err := client.GetBlockCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("local RPC unavailable: %w", err)
	}

	return &Result{
		Client: client,
		Name:   "Local CHTA Core RPC",
		Height: height,
	}, nil
}

func connectAPI(ctx context.Context, cfg *config.Config) (*Result, error) {
	baseURL := strings.TrimSpace(cfg.API.BaseURL)
	if baseURL == "" {
		baseURL = defaultAPIURL
	}

	client := api.NewClient(baseURL)

	height, err := client.GetBlockCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("public API unavailable: %w", err)
	}

	return &Result{
		Client: client,
		Name:   "Cheetahcoin Public API",
		Height: height,
	}, nil
}
