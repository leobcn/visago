package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/asaskevich/govalidator"
	"github.com/spf13/cobra"
	"github.com/zquestz/visago/plugins"
)

const (
	appName = "visago"
	version = "0.0.2"
)

// Stores configuration data.
var config Config

// FilesCmd is the main command for Cobra.
var FilesCmd = &cobra.Command{
	Use:   "visago <files/urls>",
	Short: "Visual AI Aggregator",
	Long:  `Visual AI Aggregator`,
	Run: func(cmd *cobra.Command, args []string) {
		err := filesCommand(cmd, args)
		if err != nil {
			bail(err)
		}
	},
}

func init() {
	err := config.Load()
	if err != nil {
		bail(fmt.Errorf("Failed to load configuration: %s", err))
	}

	prepareFlags()
}

func bail(err error) {
	fmt.Fprintf(os.Stderr, "[Error] %s\n", err)
	os.Exit(1)
}

func prepareFlags() {
	FilesCmd.PersistentFlags().BoolVarP(
		&config.DisplayVersion, "version", "", false, "display version")
	FilesCmd.PersistentFlags().BoolVarP(
		&config.Verbose, "verbose", "v", config.Verbose, "verbose mode")
	FilesCmd.PersistentFlags().BoolVarP(
		&config.ListPlugins, "list-plugins", "l", false, "list supported plugins")
	FilesCmd.PersistentFlags().BoolVarP(
		&config.JSONOutput, "json", "j", false, "provide JSON output")
}

// Where all the work happens.
func filesCommand(cmd *cobra.Command, args []string) error {
	if config.DisplayVersion {
		fmt.Printf("%s %s\n", appName, version)
		return nil
	}

	plugins.SetBlacklist(config.Blacklist)
	plugins.SetWhitelist(config.Whitelist)

	if config.ListPlugins {
		fmt.Printf(plugins.DisplayPlugins())
		return nil
	}

	items := strings.Join(args, " ")

	st, err := os.Stdin.Stat()
	if err != nil {
		// os.Stdin.Stat() can be unavailable on Windows.
		if runtime.GOOS != "windows" {
			return fmt.Errorf("Failed to stat Stdin: %s", err)
		}
	} else {
		if st.Mode()&os.ModeNamedPipe != 0 {
			bytes, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("Failed to read from Stdin: %s", err)
			}

			items = strings.TrimSpace(fmt.Sprintf("%s %s", items, bytes))
		}
	}

	if items != "" {
		itemList := strings.Split(items, " ")

		urls, files, errs := sortItems(itemList)
		for _, err := range errs {
			fmt.Printf("[Warn] Failed to sort input: %v\n", err)
		}

		if len(urls) == 0 && len(files) == 0 {
			return fmt.Errorf("failed to find any valid files or URLs")
		}

		pluginConfig := &plugins.PluginConfig{
			URLs:  urls,
			Files: files,
		}

		err := runPlugins(pluginConfig)
		if err != nil {
			return err
		}

		return nil
	}

	// Don't return an error, help screen is more appropriate.
	help := cmd.HelpFunc()
	help(cmd, args)

	return nil
}

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

func sortItems(items []string) (urls []string, files []*os.File, errs []error) {
	for _, item := range items {
		_, err := os.Stat(item)
		if err == nil {
			file, err := os.Open(item)
			if err == nil {
				files = append(files, file)
			} else {
				fmt.Println(err)
			}

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
