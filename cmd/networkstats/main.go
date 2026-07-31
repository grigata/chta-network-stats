package main

import (
	"fmt"

	"github.com/grigata/chta-network-stats/internal/version"
)

func main() {
	fmt.Println("========================================")
	fmt.Printf("   %s v%s\n", version.AppName, version.AppVersion)
	fmt.Println("========================================")
}
