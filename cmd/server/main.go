package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fangelod/webrtc-test/internal/server"
)

var (
	// Version is set at build time
	Version = "Development Version"
)

var (
	ver = flag.Bool("version", false, "Print the version and exit.")
)

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}

	if *ver {
		fmt.Println(Version)
		os.Exit(0)
	}

	server.Run(Version)
}
