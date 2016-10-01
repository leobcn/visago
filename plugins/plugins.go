package plugins

import (
	"fmt"
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
// different Visual AI backends
type Plugin interface {
	Perform(string) string
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
func DisplayPlugins(verbose bool) string {
	names := PluginNames(verbose)

	return fmt.Sprintf("%s\n", strings.Join(names, "\n"))
}

// PluginNames returns a sorted slice of plugin names
// applying both the whitelist and then the blacklist.
func PluginNames(verbose bool) []string {
	names := []string{}

	for key, _ := range Plugins {
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
