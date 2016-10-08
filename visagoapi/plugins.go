package visagoapi

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
// different Visual AI backends. Plugins should
// initialize themselves with a PluginConfig.
type Plugin interface {
	Perform(*PluginConfig) (string, PluginResult, error)
	Setup() error
	Reset()
	RequestIDs() ([]string, error)
}

// PluginResult is an interface on returned objects
// from the API. All result methods will require
// a string containing the requestID returned from
// Perform() and a minimum confidence score. Passing
// in 0 gets you all tags.
type PluginResult interface {
	Tags(string, float64) (map[string]map[string]*PluginTagResult, error)
	Faces(string) (map[string][]*PluginFaceResult, error)
	Colors(string) (map[string]map[string]*PluginColorResult, error)
}

// PluginTagResult are the attributes on a tag. The score
// is a value from 0 and 1.
type PluginTagResult struct {
	Name  string  `json:"name,omitempty"`
	Score float64 `json:"score,omitempty"`
}

// PluginColorResult are the attributes for a color.
type PluginColorResult struct {
	Hex           string  `json:"hex,omitempty"`
	Score         float64 `json:"score,omitempty"`
	PixelFraction float64 `json:"pixel_fraction,omitempty"`
	Red           float64 `json:"red,omitempty"`
	Green         float64 `json:"green,omitempty"`
	Blue          float64 `json:"blue,omitempty"`
	Alpha         float64 `json:"alpha,omitempty"`
}

// PluginFaceResult are the attributes on a face match.
type PluginFaceResult struct {
	BoundingPoly           *BoundingPoly `json:"bounding_poly,omitempty"`
	DetectionScore         float64       `json:"detection_score,omitempty"`
	JoyLikelihood          string        `json:"joy_likelihood,omitempty"`
	SorrowLikelihood       string        `json:"sorrow_likelihood,omitempty"`
	AngerLikelihood        string        `json:"anger_likelihood,omitempty"`
	SurpriseLikelihood     string        `json:"surprise_likelihood,omitempty"`
	UnderExposedLikelihood string        `json:"under_exposed_likelihood,omitempty"`
	BlurredLikelihood      string        `json:"blurred_likelihood,omitempty"`
	HeadwearLikelihood     string        `json:"headwear_likelihood,omitempty"`
}

// PluginConfig is used to pass configuration
// data to plugins when they load.
type PluginConfig struct {
	URLs     []string `json:"urls"`
	Files    []string `json:"files"`
	Verbose  bool     `json:"verbose"`
	TagScore float64  `json:"tag_score"`
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
