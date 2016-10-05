package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/zquestz/visago/util"
	"github.com/zquestz/visago/visagoapi"

	"github.com/asaskevich/govalidator"
	"github.com/spf13/cobra"
)

const (
	appName = "visago"
	version = "0.1.0"
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
	FilesCmd.PersistentFlags().Float64VarP(
		&config.TagScore, "tag-score", "t", 0, "minimum tag score")
}

// Where all the work happens.
func filesCommand(cmd *cobra.Command, args []string) error {
	if config.DisplayVersion {
		fmt.Printf("%s %s\n", appName, version)
		return nil
	}

	visagoapi.SetBlacklist(config.Blacklist)
	visagoapi.SetWhitelist(config.Whitelist)

	if config.ListPlugins {
		fmt.Printf(visagoapi.DisplayPlugins())
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
			util.SmartPrint("warn", fmt.Sprintf("Failed to sort input: %v\n", err), config.JSONOutput)
		}

		if len(urls) == 0 && len(files) == 0 {
			util.SmartPrint("error", fmt.Sprintf("failed to find any valid files or URLs"), config.JSONOutput)
			return nil
		}

		pluginConfig := &visagoapi.PluginConfig{
			URLs:     urls,
			Files:    files,
			Verbose:  config.Verbose,
			TagScore: config.TagScore,
		}

		output, err := visagoapi.RunPlugins(pluginConfig, config.JSONOutput)
		if err != nil {
			return err
		}

		fmt.Printf(output)

		return nil
	}

	// Don't return an error, help screen is more appropriate.
	help := cmd.HelpFunc()
	help(cmd, args)

	return nil
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
