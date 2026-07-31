package rpc

import "context"

func (c *Client) GetBlockCount(ctx context.Context) (int64, error) {
	var height int64

	if err := c.Call(ctx, "getblockcount", []any{}, &height); err != nil {
		return 0, err
	}

	return height, nil
}
