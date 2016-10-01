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
	configured bool
	clientID   string
	secret     string
}

// Perform gathers metadata from Clarifai, for the first pass
// it only supports urls.
func (p *Plugin) Perform(c plugins.PluginConfig) (string, error) {
	if p.configured == false {
		return "", fmt.Errorf("not configured")
	}

	if len(c.URLs) == 0 {
		return "", fmt.Errorf("must supply urls")
	}

	client := clarifai.NewClient(p.clientID, p.secret)
	info, err := client.Info()
	if err != nil {
		return "", err
	} else {
		fmt.Printf("%+v\n", info)
	}

	tagData, err := client.Tag(clarifai.TagRequest{
		URLs: c.URLs,
	})

	if err != nil {
		return "", err
	} else {
		fmt.Printf("%+v\n", tagData) // See what we got!
	}

	return "Clarifai performed", nil
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
