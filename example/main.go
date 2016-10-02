package main

import (
	"fmt"
	"os"

	"github.com/zquestz/visago/visagoapi"

	_ "github.com/zquestz/visago/visagoapi/clarifai"
	_ "github.com/zquestz/visago/visagoapi/googlevision"
	_ "github.com/zquestz/visago/visagoapi/imagga"
)

// Simple app that processes URLs from the command line.
// Here to illustrate the simplest way to use visago
// in your projects.
//
// Exmaple: go run main.go <imageUrls...>

func main() {
	pluginConfig := &visagoapi.PluginConfig{
		URLs: os.Args[1:],
	}

	output, err := visagoapi.RunPlugins(pluginConfig, true)
	if err != nil {
		fmt.Printf("[Error] %s\n", err.Error())
		return
	}

	fmt.Printf(output)
}
