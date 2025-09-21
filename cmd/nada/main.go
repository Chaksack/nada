package main

import (
	"fmt"
	"os"

	"github.com/chaksack/nada/internal/cli"
)

// Version information (set by build flags)
var (
	Version   = "dev"
	BuildTime = "unknown"
	Commit    = "unknown"
)

func main() {
	// Set version information in CLI
	cli.SetVersionInfo(Version, BuildTime, Commit)

	if err := cli.Execute(); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
}
