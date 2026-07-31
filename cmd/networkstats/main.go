package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/grigata/chta-network-stats/internal/config"
	"github.com/grigata/chta-network-stats/internal/rpc"
	"github.com/grigata/chta-network-stats/internal/version"
)

func main() {
	fmt.Printf("%s v%s\n", version.AppName, version.AppVersion)

	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	client := rpc.NewClient(cfg.RPC)

	fmt.Println("Connecting to CHTA Core RPC...")

	height, err := client.GetBlockCount(context.Background())
	if err != nil {
		log.Fatalf("RPC error: %v", err)
	}

	fmt.Printf("Connected.\nCurrent block height: %d\n", height)

	os.Exit(0)
}
