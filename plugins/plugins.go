package plugins

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

var (
	enableBlacklist = false
	enableWhitelist = false
	blacklist       = make(map[string]interface{})
	whitelist       = make(map[string]interface{})
)

// Plugin interface provides a way to query
// different Visual AI backends. Plugins should
// initialize themselves with a PluginConfig.
type Plugin interface {
	Perform(PluginConfig) (string, PluginResult, error)
	Setup() error
	Reset()
	RequestIDs() ([]string, error)
}

// PluginResult is an interface on returned objects
// from the API. All result methods will require
// a string containing the requestID returned from
// Perform().
type PluginResult interface {
	Tags(string) (map[string][]string, error)
}

// PluginConfig is used to pass configuration
// data to plugins when they load.
type PluginConfig struct {
	URLs  []string   `json:"-"`
	Files []*os.File `json:"-"`
}

// Plugins tracks loaded plugins.
var Plugins map[string]Plugin

func init() {
	Plugins = make(map[string]Plugin)
}

// AddPlugin should be called within your plugin's init() func.
// This will register the plugin so it can be used.
func AddPlugin(name string, plugin Plugin) {
	Plugins[name] = plugin
}

// SetBlacklist filters out unneeded plugins.
func SetBlacklist(b []string) {
	for _, v := range b {
		blacklist[v] = nil
	}

	if len(blacklist) > 0 {
		enableBlacklist = true
	}
}

// SetWhitelist sets an exact list of supported plugins.
func SetWhitelist(w []string) {
	for _, v := range w {
		whitelist[v] = nil
	}

	if len(whitelist) > 0 {
		enableWhitelist = true
	}
}

// DisplayPlugins displays all the loaded plugins.
func DisplayPlugins() string {
	names := PluginNames()

	return fmt.Sprintf("%s\n", strings.Join(names, "\n"))
}

// PluginNames returns a sorted slice of plugin names
// applying both the whitelist and then the blacklist.
func PluginNames() []string {
	names := []string{}

	for key := range Plugins {
		if enableWhitelist {
			_, ok := whitelist[key]
			if !ok {
				continue
			}
		}

		if enableBlacklist {
			_, ok := blacklist[key]
			if ok {
				continue
			}
		}

		names = append(names, key)
	}

	sort.Strings(names)
	return names
}
