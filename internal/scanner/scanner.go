package scanner

import (
	"context"
	"fmt"
	"time"

	"github.com/grigata/chta-network-stats/internal/models"
	"github.com/grigata/chta-network-stats/internal/parser"
	"github.com/grigata/chta-network-stats/internal/rpc"
)

type Scanner struct {
	client *rpc.Client
}

func New(client *rpc.Client) *Scanner {
	return &Scanner{
		client: client,
	}
}

func (s *Scanner) ReadLastBlocks(
	ctx context.Context,
	count int,
) ([]models.NetworkBlock, error) {
	if count <= 0 {
		return nil, fmt.Errorf("block count must be greater than zero")
	}

	latestHeight, err := s.client.GetBlockCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("get latest block height: %w", err)
	}

	if int64(count) > latestHeight {
		count = int(latestHeight)
	}

	blocks := make([]models.NetworkBlock, 0, count)

	for height := latestHeight; height > latestHeight-int64(count); height-- {

		hash, err := s.client.GetBlockHash(ctx, height)
		if err != nil {
			return nil, fmt.Errorf(
				"get block hash at height %d: %w",
				height,
				err,
			)
		}

		block, err := s.client.GetBlock(ctx, hash)
		if err != nil {
			return nil, fmt.Errorf(
				"get block at height %d: %w",
				height,
				err,
			)
		}

		var gap int64

		if height > 0 {
			previousHash, err := s.client.GetBlockHash(ctx, height-1)
			if err != nil {
				return nil, fmt.Errorf(
					"get previous block hash at height %d: %w",
					height-1,
					err,
				)
			}

			previousBlock, err := s.client.GetBlock(ctx, previousHash)
			if err != nil {
				return nil, fmt.Errorf(
					"get previous block at height %d: %w",
					height-1,
					err,
				)
			}

			gap = block.Time - previousBlock.Time
		}

		coinbaseTxID := ""
		coinbaseHex := ""
		pool := "Unknown"

		if len(block.Tx) > 0 {
			coinbaseTxID = block.Tx[0]

			tx, err := s.client.GetRawTransaction(ctx, coinbaseTxID)
			if err != nil {
				return nil, fmt.Errorf(
					"get coinbase transaction at height %d: %w",
					height,
					err,
				)
			}

			if len(tx.Vin) > 0 {
				coinbaseHex = tx.Vin[0].Coinbase

				coinbaseText := parser.DecodeCoinbaseText(coinbaseHex)
				pool = parser.DetectPool(coinbaseText)
			}
		}

		blocks = append(blocks, models.NetworkBlock{
			Height:       block.Height,
			Hash:         block.Hash,
			Difficulty:   block.Difficulty,
			Bits:         block.Bits,
			Time:         time.Unix(block.Time, 0).Local(),
			Gap:          gap,
			TxCount:      len(block.Tx),
			Type:         blockType(block.Difficulty),
			CoinbaseTxID: coinbaseTxID,
			CoinbaseHex:  coinbaseHex,
			Pool:         pool,
		})
	}

	return blocks, nil
}

func blockType(difficulty float64) string {
	if difficulty < 1 {
		return "CHEETAH"
	}

	return "NORMAL"
}
