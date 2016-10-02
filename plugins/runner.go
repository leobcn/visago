package plugins

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/zquestz/visago/util"
)

// RunPlugins runs all the plugins with the provided pluginConfig.
// Output is directed at stdout. Not intended for API use.
func RunPlugins(pluginConfig *PluginConfig, jsonOutput bool) error {
	wg := &sync.WaitGroup{}
	dwg := &sync.WaitGroup{}

	outputChan := make(chan []string)
	defer close(outputChan)

	runChan := make(chan *runner)
	defer close(runChan)

	finishedChan := make(chan bool)
	defer close(finishedChan)

	dwg.Add(1)
	go processOutput(dwg, outputChan, runChan, finishedChan, jsonOutput)

	for _, name := range PluginNames() {
		wg.Add(1)
		r := runner{
			Name: name,
		}
		go r.run(name, pluginConfig, wg, outputChan, runChan)
	}

	// Wait for plugins to finish.
	wg.Wait()

	finishedChan <- true

	// Wait for display goroutine to end.
	dwg.Wait()

	return nil
}

type runner struct {
	Name   string
	TagMap map[string][]string
}

func processOutput(wg *sync.WaitGroup, outputChan <-chan []string, runChan <-chan *runner, finishedChan <-chan bool, jsonOutput bool) {
	defer wg.Done()

	runners := []*runner{}

Loop:
	for {
		select {
		case output := <-outputChan:
			util.SmartPrint(output[0], output[1], jsonOutput)
		case runner := <-runChan:
			if len(runner.TagMap) > 0 {
				runners = append(runners, runner)
			}
		case <-finishedChan:
			break Loop
		}
	}

	displayOutput(runners, jsonOutput)

	return
}

func displayOutput(runners []*runner, jsonOutput bool) {
	output := make(map[string]map[string][]string)

	for _, r := range runners {
		if _, ok := output[r.Name]; !ok {
			output[r.Name] = make(map[string][]string)
		}

		for k, v := range r.TagMap {
			output[r.Name][k] = v
		}
	}

	if jsonOutput {
		b, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			fmt.Printf("%s", err)
		}
		fmt.Printf("%s\n", b)
	} else {
		for k, v := range output {
			fmt.Printf("%s\n", k)

			for asset, tags := range v {
				fmt.Printf("- %s\n", asset)
				fmt.Printf("  %v\n", tags)
			}
		}
	}
}

func (r *runner) run(name string, pluginConfig *PluginConfig, wg *sync.WaitGroup, outputChan chan<- []string, runChan chan<- *runner) {
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

	r.TagMap = tagMap

	runChan <- r

	return
}
