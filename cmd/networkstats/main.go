package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/grigata/chta-network-stats/internal/config"
	"github.com/grigata/chta-network-stats/internal/datasource"
	"github.com/grigata/chta-network-stats/internal/models"
	"github.com/grigata/chta-network-stats/internal/scanner"
	"github.com/grigata/chta-network-stats/internal/statistics"
	"github.com/grigata/chta-network-stats/internal/version"
)

const defaultBlockCount = 100

func main() {
	blockCount := readBlockCount()

	cfg, err := config.Load("config.json")
	if err != nil {
		exitWithError("Configuration error: %v", err)
	}

	ctx := context.Background()

	fmt.Println("Connecting to Cheetahcoin network...")

	source, err := datasource.Connect(ctx, cfg)
	if err != nil {
		exitWithError("Connection error: %v", err)
	}

	printStartup(
		cfg.Mode,
		source.Name,
		source.Height,
		blockCount,
	)

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
	poolStats := statistics.AnalyzePools(blocks)

	printBlockTable(blocks)
	printStatistics(
		stats,
		source.Name,
		scanDuration,
	)
	printPoolAnalysis(poolStats)

	waitForEnter()
}

func readBlockCount() int {
	blockCount := defaultBlockCount

	if len(os.Args) <= 1 {
		return blockCount
	}

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

	return value
}

func printStartup(
	mode string,
	sourceName string,
	height int64,
	blockCount int,
) {
	fmt.Println()
	fmt.Println("============================================================")
	fmt.Printf("%s v%s\n", version.AppName, version.AppVersion)
	fmt.Println("============================================================")
	fmt.Println()

	fmt.Printf("Mode        : %s\n", strings.ToUpper(configMode(mode)))
	fmt.Printf("Source      : %s\n", sourceName)
	fmt.Printf("Height      : %d\n", height)
	fmt.Printf("Blocks      : %d\n", blockCount)
	fmt.Println()

	if sourceName == "Cheetahcoin Public API" && blockCount >= 1000 {
		fmt.Println("Warning:")
		fmt.Println("Large scans through the Public API may take several minutes.")
		fmt.Println()
	}

	fmt.Printf("Scanning last %d network blocks...\n", blockCount)
}

func printBlockTable(blocks []models.NetworkBlock) {
	fmt.Println()
	fmt.Printf("Last %d network blocks\n", len(blocks))
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

func printStatistics(
	stats statistics.Statistics,
	sourceName string,
	scanDuration time.Duration,
) {
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

	fmt.Printf("Data Source     : %s\n", sourceName)
	fmt.Printf("Scan Time       : %.2f sec\n", scanDuration.Seconds())

	fmt.Println()

}

func configMode(mode string) string {
	mode = strings.TrimSpace(mode)

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
func printPoolAnalysis(pools []models.PoolStats) {

	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("Pool Analysis")
	fmt.Println("============================================================")

	fmt.Printf(
		"%-18s %6s %7s %10s %7s\n",
		"Pool",
		"Total",
		"Normal",
		"LOW-DIFF",
		"%",
	)

	fmt.Println("------------------------------------------------------------")

	for _, p := range pools {

		fmt.Printf(
			"%-18s %6d %7d %10d %6.1f%%\n",
			p.Name,
			p.TotalBlocks,
			p.NormalBlocks,
			p.LowDiffBlocks,
			p.Percent,
		)
	}
}
