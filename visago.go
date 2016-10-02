package main

import (
	"fmt"
	"os"

	_ "github.com/zquestz/visago/visagoapi/clarifai"
	_ "github.com/zquestz/visago/visagoapi/googlevision"
	_ "github.com/zquestz/visago/visagoapi/imagga"

	"github.com/zquestz/visago/cmd"
)

func main() {
	setupSignalHandlers()

	if err := cmd.FilesCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
