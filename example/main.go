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
// Here to illustrate the simplest way to use the visagoapi
// in your projects.
//
// Example: go run main.go <imageUrls...>

func main() {
	pluginConfig := &visagoapi.PluginConfig{
		URLs: os.Args[1:],
		// To only enable select features, set them below.
		// By default all features are enabled.
		// Features: []string{visagoapi.TagsFeature, visagoapi.ColorsFeature, visagoapi.FacesFeature},
	}

	output, err := visagoapi.RunPlugins(pluginConfig, true)
	if err != nil {
		fmt.Printf("[Error] %s\n", err.Error())
		return
	}

	fmt.Printf(output)
}
