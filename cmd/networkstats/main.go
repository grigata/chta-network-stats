package main

import (
	"context"
	"fmt"
	"log"

	"github.com/grigata/chta-network-stats/internal/config"
	"github.com/grigata/chta-network-stats/internal/rpc"
	"github.com/grigata/chta-network-stats/internal/scanner"
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

	blockScanner := scanner.New(client)

	blocks, err := blockScanner.ReadLastBlocks(ctx, 10)
	if err != nil {
		log.Fatalf("Scan error: %v", err)
	}

	fmt.Println("Connected.")
	fmt.Println()

	fmt.Println("Last 10 network blocks")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf(
		"%-10s %-18s %-14s %-10s %-8s %-5s\n",
		"Height",
		"Pool",
		"Difficulty",
		"Type",
		"Gap",
		"Tx",
	)
	fmt.Println("--------------------------------------------------------------------------------")

	for _, block := range blocks {
		fmt.Printf(
			"%-10d %-18s %-14s %-10s %+7ds %-5d\n",
			block.Height,
			block.Pool,
			formatDifficulty(block.Difficulty),
			block.Type,
			block.Gap,
			block.TxCount,
		)
	}
}

func formatDifficulty(difficulty float64) string {
	switch {
	case difficulty >= 1_000_000_000:
		return fmt.Sprintf("%.2fG", difficulty/1_000_000_000)

	case difficulty >= 1_000_000:
		return fmt.Sprintf("%.2fM", difficulty/1_000_000)

	case difficulty >= 1_000:
		return fmt.Sprintf("%.2fK", difficulty/1_000)

	default:
		return fmt.Sprintf("%.4f", difficulty)
	}
}
