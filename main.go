package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bedirmirac/glipboard/cmd"
	"github.com/bedirmirac/glipboard/helper"
	"github.com/bedirmirac/glipboard/storage"
	"github.com/bedirmirac/glipboard/tui"
)

func main() {
	tuiMode := flag.Bool("tui", false, "Open with the TUI")
	limitFlag := flag.Int("o", 0, "Number of items that you want to store (overrides config file)")

	flag.Parse()

	s, err := storage.NewStorage()
	if err != nil {
		fmt.Printf("there's been an error during connecting database: %v", err)
		os.Exit(1)
	}

	if *limitFlag > 0 {
		recordCount, err := s.Count()
		if err == nil && recordCount > *limitFlag {
			x := recordCount - *limitFlag
			err := s.DeleteFromX(x)
			if err != nil {
				fmt.Printf("there's been error deleting old records to apply new limit, if they aren't deleted you can delete manually\n")
			}
		}
		helper.WriteLimit(*limitFlag)
	} else {
		_, err := helper.ReadLimit()
		if err != nil {
			recordCount, err := s.Count()
			if err == nil && recordCount > *limitFlag {
				x := recordCount - *limitFlag
				err := s.DeleteFromX(x)
				if err != nil {
					fmt.Printf("there's been error deleting old records to apply new limit, if they aren't deleted you can delete manually\n")
				}
			}
			helper.WriteLimit(50)
		}
	}
	if *tuiMode {
		go cmd.StartDaemon(s)
		tui.StartTUI(s)
		os.Exit(0)
	}

	cmd.StartDaemon(s)
}
