package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/zquestz/go-ucl"
)

// Config stores all the application configuration.
type Config struct {
	Blacklist      []string `json:"blacklist"`
	DisplayVersion bool     `json:"-"`
	ListPlugins    bool     `json:"-"`
	Verbose        bool     `json:"verbose,string"`
	Whitelist      []string `json:"whitelist"`
	JSONOutput     bool     `json:"json_output,string"`
	TagScore       float64  `json:"tag_score,string"`
	Tags           bool     `json:"tags,string"`
	Faces          bool     `json:"faces,string"`
	Colors         bool     `json:"colors,string"`
}

// Load reads the configuration from ~/.visago/config and loads it into the Config struct.
// The config is in UCL format.
func (c *Config) Load() error {
	conf, err := c.loadConfig()
	if err != nil {
		return err
	}

	// There are cases when we don't have a configuration.
	if conf != nil {
		err = c.applyConf(conf)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) loadConfig() ([]byte, error) {
	h, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filepath.Join(h, ".visago", "config"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}
	defer f.Close()

	ucl.Ucldebug = false
	data, err := ucl.NewParser(f).Ucl()
	if err != nil {
		return nil, err
	}

	conf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *Config) applyConf(conf []byte) error {
	err := json.Unmarshal(conf, c)
	if err != nil {
		return err
	}

	return nil
}
