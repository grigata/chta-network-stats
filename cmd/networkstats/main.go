package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/grigata/chta-network-stats/internal/config"
	"github.com/grigata/chta-network-stats/internal/rpc"
	"github.com/grigata/chta-network-stats/internal/version"
)

func main() {
	fmt.Printf("%s v%s\n", version.AppName, version.AppVersion)

	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	client := rpc.NewClient(cfg.RPC)
	ctx := context.Background()

	fmt.Println("Connecting to CHTA Core RPC...")

	height, err := client.GetBlockCount(ctx)
	if err != nil {
		log.Fatalf("RPC error: %v", err)
	}

	hash, err := client.GetBlockHash(ctx, height)
	if err != nil {
		log.Fatalf("RPC error: %v", err)
	}

	block, err := client.GetBlock(ctx, hash)
	if err != nil {
		log.Fatalf("RPC error: %v", err)
	}
	previousHash, err := client.GetBlockHash(ctx, height-1)
	if err != nil {
		log.Fatalf("RPC error: %v", err)
	}

	previousBlock, err := client.GetBlock(ctx, previousHash)
	if err != nil {
		log.Fatalf("RPC error: %v", err)
	}

	gap := block.Time - previousBlock.Time

	blockTime := time.Unix(block.Time, 0).Local()

	fmt.Println("Connected.")
	fmt.Println()
	fmt.Println("Latest network block")
	fmt.Println("--------------------")
	fmt.Printf("Height      : %d\n", block.Height)
	fmt.Printf("Hash        : %s\n", block.Hash)
	fmt.Printf("Difficulty  : %.8f\n", block.Difficulty)
	fmt.Printf("Bits        : %s\n", block.Bits)
	fmt.Printf("Time        : %s\n", blockTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Transactions: %d\n", len(block.Tx))
	fmt.Printf("Gap         : %ds\n", gap)
}
