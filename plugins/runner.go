package plugins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/zquestz/visago/util"
)

// RunPlugins runs all the plugins with the provided pluginConfig.
// Output is directed at stdout. Not intended for API use.
func RunPlugins(pluginConfig *PluginConfig, jsonOutput bool) error {
	wg := &sync.WaitGroup{}

	outputChan := make(chan []string)
	defer close(outputChan)

	finishedChan := make(chan bool)
	defer close(finishedChan)

	wg.Add(1)
	go processPluginOutput(wg, outputChan, finishedChan, jsonOutput)

	for _, name := range PluginNames() {
		wg.Add(1)
		go runPlugin(name, pluginConfig, wg, outputChan, jsonOutput)
	}

	wg.Wait()

	finishedChan <- true

	return nil
}

func processPluginOutput(wg *sync.WaitGroup, outputChan <-chan []string, finishedChan <-chan bool, jsonOutput bool) {
	wg.Done()

Loop:
	for {
		select {
		case output := <-outputChan:
			util.SmartPrint(output[0], output[1], jsonOutput)
		case <-finishedChan:
			break Loop
		}
	}
	return
}

func runPlugin(name string, pluginConfig *PluginConfig, wg *sync.WaitGroup, outputChan chan<- []string, jsonOutput bool) {
	defer wg.Done()

	err := Plugins[name].Setup()
	if err != nil {
		outputChan <- []string{"warn", fmt.Sprintf("Failed to setup %q plugin: %s\n", name, err)}
		return
	}

	requestID, pluginResponse, err := Plugins[name].Perform(pluginConfig)
	if err != nil {
		outputChan <- []string{"warn", fmt.Sprintf("Failed running perform on %q plugin: %s\n", name, err)}
		return
	}

	tagMap, err := pluginResponse.Tags(requestID)
	if err != nil {
		outputChan <- []string{"warn", fmt.Sprintf("Failed to fetch tags from plugin %q: %s\n", name, err)}
		return
	}

	outputChan <- []string{"", displayTags(name, tagMap, jsonOutput)}

	return
}

func displayTags(name string, tagMap map[string][]string, jsonOutput bool) string {
	var buf []byte
	output := bytes.NewBuffer(buf)

	if jsonOutput {
		b, _ := json.Marshal(tagMap)
		output.WriteString(fmt.Sprintf("%s\n", b))
	} else {
		if len(tagMap) > 0 {
			for asset, tags := range tagMap {
				output.WriteString(fmt.Sprintf("%s - %s\n", name, asset))
				output.WriteString(fmt.Sprintf("%v\n", tags))
			}
		}
	}

	return output.String()
}
