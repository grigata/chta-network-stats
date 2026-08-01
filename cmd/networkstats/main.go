package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/grigata/chta-network-stats/internal/config"
	"github.com/grigata/chta-network-stats/internal/datasource"
	"github.com/grigata/chta-network-stats/internal/scanner"
	"github.com/grigata/chta-network-stats/internal/statistics"
	"github.com/grigata/chta-network-stats/internal/version"
)

func main() {
	blockCount := 100

	if len(os.Args) > 1 {
		value, err := strconv.Atoi(os.Args[1])
		if err != nil {
			exitWithError(
				"Invalid block count %q: expected a whole number",
				os.Args[1],
			)
		}

		if value <= 0 {
			exitWithError("Block count must be greater than zero")
		}

		blockCount = value
	}
	fmt.Printf("%s v%s\n", version.AppName, version.AppVersion)

	cfg, err := config.Load("config.json")
	if err != nil {
		exitWithError("Configuration error: %v", err)
	}

	ctx := context.Background()

	fmt.Printf("Data source mode: %s\n", configMode(cfg.Mode))
	fmt.Println("Connecting...")

	source, err := datasource.Connect(ctx, cfg)
	if err != nil {
		exitWithError("Connection error: %v", err)
	}

	fmt.Printf("Connected using: %s\n", source.Name)
	fmt.Printf("Current height : %d\n", source.Height)
	fmt.Printf("Scanning last %d network blocks...\n", blockCount)

	blockScanner := scanner.New(source.Client)

	scanStarted := time.Now()

	blocks, err := blockScanner.ReadLastBlocks(
		ctx,
		blockCount,
		func(current, total int) {
			fmt.Printf("\rScanning: %3d/%d", current, total)
		},
	)

	fmt.Println()

	if err != nil {
		exitWithError("Scan error: %v", err)
	}

	scanDuration := time.Since(scanStarted)

	fmt.Printf(
		"Scan completed in %.2f seconds\n",
		scanDuration.Seconds(),
	)

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
	fmt.Printf("Total Blocks    : %d\n", stats.TotalBlocks)
	fmt.Printf("Normal Blocks   : %d\n", stats.NormalBlocks)
	fmt.Printf("LOW-DIFF Blocks : %d\n", stats.LowDiffBlocks)
	fmt.Printf("Average Gap     : %.1f sec\n", stats.AverageGap)
	fmt.Printf("Minimum Gap     : %d sec\n", stats.MinGap)
	fmt.Printf("Maximum Gap     : %d sec\n", stats.MaxGap)

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

func configMode(mode string) string {
	if mode == "" {
		return "auto"
	}

	return mode
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
	os.Exit(1)
}

func waitForEnter() {
	fmt.Println()
	fmt.Print("Press Enter to exit...")
	fmt.Scanln()
}
