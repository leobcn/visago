package visagoapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
)

const (
	errorKey = "errors"
)

type runner struct {
	Name   string
	TagMap map[string][]string
	Errors []error
}

// RunPlugins runs all the plugins with the provided pluginConfig.
// Output is directed at stdout. Not intended for API use.
func RunPlugins(pluginConfig *PluginConfig, jsonOutput bool) (string, error) {
	wg := &sync.WaitGroup{}
	dwg := &sync.WaitGroup{}

	outputChan := make(chan string)
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
		go r.run(name, pluginConfig, wg, runChan)
	}

	// Wait for plugins to finish.
	wg.Wait()

	finishedChan <- true

	output := <-outputChan

	// Wait for display goroutine to end.
	dwg.Wait()

	return output, nil
}

func processOutput(wg *sync.WaitGroup, outputChan chan<- string, runChan <-chan *runner, finishedChan <-chan bool, jsonOutput bool) {
	defer wg.Done()

	runners := []*runner{}

Loop:
	for {
		select {
		case runner := <-runChan:
			runners = append(runners, runner)
		case <-finishedChan:
			break Loop
		}
	}

	outputChan <- displayOutput(buildOutput(runners), jsonOutput)

	return
}

func buildOutput(runners []*runner) map[string]*Result {
	output := make(map[string]*Result)

	for _, r := range runners {
		if _, ok := output[r.Name]; !ok {
			output[r.Name] = &Result{}
		}

		if len(r.Errors) > 0 {
			output[r.Name].Errors = []string{}
			for _, e := range r.Errors {
				output[r.Name].Errors = append(output[r.Name].Errors, e.Error())
			}
		}

		for k, v := range r.TagMap {
			asset := Asset{
				Name: k,
				Tags: v,
			}

			output[r.Name].Assets = append(output[r.Name].Assets, &asset)
		}
	}

	return output
}

func displayOutput(output map[string]*Result, jsonOutput bool) string {
	outputBuf := bytes.NewBuffer([]byte{})

	if len(output) == 0 {
		return ""
	}

	if jsonOutput {
		b, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			outputBuf.WriteString(fmt.Sprintf("%s", err))
		}
		outputBuf.WriteString(fmt.Sprintf("%s\n", b))
	} else {
		for k, v := range output {
			outputBuf.WriteString(fmt.Sprintf("%s\n", k))

			for _, asset := range v.Assets {
				outputBuf.WriteString(fmt.Sprintf("- %s\n", asset.Name))
				outputBuf.WriteString(fmt.Sprintf("%v\n", asset.Tags))
			}

			for _, err := range output[k].Errors {
				outputBuf.WriteString(fmt.Sprintf("- %v\n", err))
			}
		}
	}

	return outputBuf.String()
}

func (r *runner) run(name string, pluginConfig *PluginConfig, wg *sync.WaitGroup, runChan chan<- *runner) {
	defer wg.Done()

	defer func() { runChan <- r }()

	err := Plugins[name].Setup()
	if err != nil {
		r.Errors = append(r.Errors, err)
		return
	}

	requestID, pluginResponse, err := Plugins[name].Perform(pluginConfig)
	if err != nil {
		r.Errors = append(r.Errors, err)
		return
	}

	tagMap, err := pluginResponse.Tags(requestID)
	if err != nil {
		r.Errors = append(r.Errors, err)
		return
	}

	r.TagMap = tagMap

	return
}