package main

import (
	"context"
	"fmt"

	"github.com/grigata/chta-network-stats/internal/api"
	"github.com/grigata/chta-network-stats/internal/scanner"
	"github.com/grigata/chta-network-stats/internal/statistics"
	"github.com/grigata/chta-network-stats/internal/version"
)

const (
	blockCount = 100
	apiURL     = "http://chtaexplorer.mooo.com:3002"
)

func main() {
	fmt.Printf("%s v%s\n", version.AppName, version.AppVersion)

	ctx := context.Background()

	fmt.Println("Connecting to Cheetahcoin Public API...")

	apiClient := api.NewClient(apiURL)

	if _, err := apiClient.GetBlockCount(ctx); err != nil {
		exitWithError("Public API connection error: %v", err)
	}

	fmt.Println("Connected.")
	fmt.Printf("Scanning last %d network blocks...\n", blockCount)

	blockScanner := scanner.New(apiClient)

	blocks, err := blockScanner.ReadLastBlocks(ctx, blockCount)
	if err != nil {
		exitWithError("Scan error: %v", err)
	}

	stats := statistics.Calculate(blocks)

	fmt.Println()
	fmt.Printf("Last %d network blocks\n", blockCount)
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

	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("Network Statistics")
	fmt.Println("============================================================")
	fmt.Printf("Total Blocks   : %d\n", stats.TotalBlocks)
	fmt.Printf("Normal Blocks  : %d\n", stats.NormalBlocks)
	fmt.Printf("LOW-DIFF Blocks: %d\n", stats.LowDiffBlocks)
	fmt.Printf("Average Gap    : %.1f sec\n", stats.AverageGap)
	fmt.Printf("Minimum Gap    : %d sec\n", stats.MinGap)
	fmt.Printf("Maximum Gap    : %d sec\n", stats.MaxGap)

	fmt.Println()
	fmt.Println("Pool Distribution")
	fmt.Println("----------------------------")

	for _, pool := range statistics.SortedPools(stats) {
		fmt.Printf(
			"%-18s %3d (%5.1f%%)\n",
			pool.Name,
			pool.Blocks,
			pool.Percent,
		)
	}

	waitForEnter()
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

func exitWithError(format string, args ...any) {
	fmt.Printf("\nError: "+format+"\n", args...)
	waitForEnter()
}

func waitForEnter() {
	fmt.Println()
	fmt.Print("Press Enter to exit...")
	fmt.Scanln()
}
