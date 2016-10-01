package clarifai

import (
	"fmt"
	"os"

	"github.com/Clarifai/clarifai-go"
	"github.com/zquestz/visago/plugins"
)

func init() {
	plugins.AddPlugin("clarifai", &Plugin{})
}

// Plugin implements the Plugin interface and stores
// configuration data needed by the clarifai library.
type Plugin struct {
	configured     bool
	clientID       string
	secret         string
	tagResp        *clarifai.TagResp
	processedURLs  []string
	processedFiles []*os.File
}

// Perform gathers metadata from Clarifai, for the first pass
// it only supports urls.
func (p *Plugin) Perform(c plugins.PluginConfig) (plugins.PluginResult, error) {
	if p.configured == false {
		return nil, fmt.Errorf("not configured")
	}

	if len(c.URLs) == 0 {
		return nil, fmt.Errorf("must supply urls")
	}

	client := clarifai.NewClient(p.clientID, p.secret)
	_, err := client.Info()
	if err != nil {
		return nil, err
	}

	p.tagResp, err = client.Tag(clarifai.TagRequest{
		URLs: c.URLs,
	})
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Tags returns the tags on an entry
func (p *Plugin) Tags() (tags map[string][]string, err error) {
	tags = make(map[string][]string)

	if p.tagResp == nil {
		return tags, fmt.Errorf("request has not been made to clarifai")
	}

	for _, result := range p.tagResp.Results {
		tags[result.URL] = result.Result.Tag.Classes
	}

	return
}

// Setup sets up the plugin for use.
func (p *Plugin) Setup() error {
	id := os.Getenv("CLARIFAI_CLIENT_ID")
	secret := os.Getenv("CLARIFAI_CLIENT_SECRET")

	if id == "" || secret == "" {
		p.configured = false
		return fmt.Errorf("credentials not found")
	}

	p.clientID = id
	p.secret = secret
	p.configured = true

	return nil
}
