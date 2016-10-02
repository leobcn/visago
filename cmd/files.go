package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

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
