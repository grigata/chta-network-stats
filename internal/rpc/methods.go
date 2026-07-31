package rpc

import (
	"context"
	"fmt"
)

func (c *Client) GetBlockCount(ctx context.Context) (int64, error) {
	var height int64

	if err := c.Call(ctx, "getblockcount", []any{}, &height); err != nil {
		return 0, fmt.Errorf("get block count: %w", err)
	}

	return height, nil
}

func (c *Client) GetBlockHash(ctx context.Context, height int64) (string, error) {
	var hash string

	if err := c.Call(ctx, "getblockhash", []any{height}, &hash); err != nil {
		return "", fmt.Errorf("get block hash at height %d: %w", height, err)
	}

	return hash, nil
}

func (c *Client) GetBlock(ctx context.Context, hash string) (*Block, error) {
	var block Block

	if err := c.Call(ctx, "getblock", []any{hash, true}, &block); err != nil {
		return nil, fmt.Errorf("get block %s: %w", hash, err)
	}

	return &block, nil
}
func (c *Client) GetRawTransaction(ctx context.Context, txid string) (*RawTransaction, error) {
	var tx RawTransaction

	if err := c.Call(ctx, "getrawtransaction", []any{txid, 1}, &tx); err != nil {
		return nil, fmt.Errorf("get raw transaction %s: %w", txid, err)
	}

	return &tx, nil
}
