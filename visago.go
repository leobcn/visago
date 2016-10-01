package main

import (
	"fmt"
	"os"

	_ "github.com/zquestz/visago/plugins/clarifai"
	_ "github.com/zquestz/visago/plugins/google_vision"

	"github.com/zquestz/visago/cmd"
)

func main() {
	setupSignalHandlers()

	if err := cmd.FilesCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
