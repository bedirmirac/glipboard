package main

import (
	"flag"
	"os"

	"github.com/bedirmirac/glipboard/cmd"
	"github.com/bedirmirac/glipboard/tui"
)

func main() {
	tuiMode := flag.Bool("tui", false, "Open with the TUI")

	flag.Parse()

	if *tuiMode {
		go cmd.StartDaemon()
		tui.StartTUI()
		os.Exit(0)
	}

	cmd.StartDaemon()
}
