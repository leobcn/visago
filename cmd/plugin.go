package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/asaskevich/govalidator"
	"github.com/zquestz/visago/plugins"
)

func runPlugins(pluginConfig *plugins.PluginConfig) error {
	wg := &sync.WaitGroup{}

	outputChan := make(chan string)
	defer close(outputChan)

	finishedChan := make(chan bool)
	defer close(finishedChan)

	wg.Add(1)
	go processPluginOutput(wg, outputChan, finishedChan)

	for _, name := range plugins.PluginNames() {
		wg.Add(1)
		go runPlugin(name, pluginConfig, wg, outputChan)
	}

	wg.Wait()

	finishedChan <- true

	return nil
}

func processPluginOutput(wg *sync.WaitGroup, outputChan <-chan string, finishedChan <-chan bool) {
	wg.Done()

Loop:
	for {
		select {
		case output := <-outputChan:
			fmt.Printf(output)
		case <-finishedChan:
			break Loop
		}
	}
	return
}

func runPlugin(name string, pluginConfig *plugins.PluginConfig, wg *sync.WaitGroup, outputChan chan<- string) {
	defer wg.Done()

	err := plugins.Plugins[name].Setup()
	if err != nil {
		outputChan <- fmt.Sprintf("[Warn] Failed to setup %q plugin: %s\n", name, err)
		return
	}

	requestID, pluginResponse, err := plugins.Plugins[name].Perform(pluginConfig)
	if err != nil {
		outputChan <- fmt.Sprintf("[Warn] Failed running perform on %q plugin: %s\n", name, err)
		return
	}

	tagMap, err := pluginResponse.Tags(requestID)
	if err != nil {
		outputChan <- fmt.Sprintf("[Warn] Failed to fetch tags from plugin %q: %s\n", name, err)
		return
	}

	outputChan <- displayTags(name, tagMap)

	return
}

func displayTags(name string, tagMap map[string][]string) string {
	var buf []byte
	output := bytes.NewBuffer(buf)

	if config.JSONOutput {
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

func sortItems(items []string) (urls []string, files []string, errs []error) {
	for _, item := range items {
		_, err := os.Stat(item)
		if err == nil {
			files = append(files, item)

			continue
		}

		valid := govalidator.IsURL(item)
		if valid {
			urls = append(urls, item)
			continue
		}

		errs = append(errs, fmt.Errorf("%q is not a file or url", item))
	}

	return
}
