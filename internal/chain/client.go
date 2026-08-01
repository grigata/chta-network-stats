package chain

import (
	"context"

	"github.com/grigata/chta-network-stats/internal/rpc"
)

type Client interface {
	GetBlockCount(ctx context.Context) (int64, error)
	GetBlockHash(ctx context.Context, height int64) (string, error)
	GetBlock(ctx context.Context, hash string) (*rpc.Block, error)
	GetRawTransaction(ctx context.Context, txid string) (*rpc.RawTransaction, error)
}
