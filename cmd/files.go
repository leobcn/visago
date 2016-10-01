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
	version = "0.0.1"
)

// Stores configuration data.
var config Config

// FilesCmd is the main command for Cobra.
var FilesCmd = &cobra.Command{
	Use:   "visago <files>",
	Short: "Visual AI Aggregator",
	Long:  `Visual AI Aggregator`,
	Run: func(cmd *cobra.Command, args []string) {
		err := FilesCommand(cmd, args)
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
}

// Where all the work happens.
func FilesCommand(cmd *cobra.Command, args []string) error {
	if config.DisplayVersion {
		fmt.Printf("%s %s\n", appName, version)
		return nil
	}

	plugins.SetBlacklist(config.Blacklist)
	plugins.SetWhitelist(config.Whitelist)

	if config.ListPlugins {
		fmt.Printf(plugins.DisplayPlugins(config.Verbose))
		return nil
	}

	files := strings.Join(args, " ")

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

			files = strings.TrimSpace(fmt.Sprintf("%s %s", files, bytes))
		}
	}

	if files != "" {
	} else {
		// Don't return an error, help screen is more appropriate.
		help := cmd.HelpFunc()
		help(cmd, args)
	}

	return nil
}
