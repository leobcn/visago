package main

import (
	"fmt"
	"os"

	"github.com/zquestz/visago/cmd"
)

func main() {
	setupSignalHandlers()

	if err := cmd.FilesCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
