package clarifai

import (
	"fmt"
	"os"

	"github.com/Clarifai/clarifai-go"
	"github.com/nats-io/nuid"
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
	responses  map[string]*clarifai.TagResp
}

// Perform gathers metadata from Clarifai, for the first pass
// it only supports urls. When PR's in the clarifai repo add
// file support, it will be added here.
func (p *Plugin) Perform(c plugins.PluginConfig) (string, plugins.PluginResult, error) {
	if p.configured == false {
		return "", nil, fmt.Errorf("not configured")
	}

	if len(c.URLs) == 0 {
		return "", nil, fmt.Errorf("must supply urls")
	}

	client := clarifai.NewClient(p.clientID, p.secret)
	_, err := client.Info()
	if err != nil {
		return "", nil, err
	}

	nuid := nuid.Next()

	p.responses[nuid], err = client.Tag(clarifai.TagRequest{
		URLs: c.URLs,
	})

	if err != nil {
		return "", nil, err
	}

	return nuid, p, nil
}

// Tags returns the tags on an entry
func (p *Plugin) Tags(requestID string) (tags map[string][]string, err error) {
	tags = make(map[string][]string)

	if p.responses[requestID] == nil {
		return tags, fmt.Errorf("request has not been made to clarifai")
	}

	for _, result := range p.responses[requestID].Results {
		tags[result.URL] = result.Result.Tag.Classes
	}

	return
}

// Reset clears the cache of existing responses.
func (p *Plugin) Reset() {
	p.responses = make(map[string]*clarifai.TagResp)
}

// RequestIDs returns a list of all cached response
// requestIDs.
func (p *Plugin) RequestIDs() ([]string, error) {
	if p.configured == false {
		return nil, fmt.Errorf("not configured")
	}

	keys := []string{}
	for k := range p.responses {
		keys = append(keys, k)
	}

	return keys, nil
}

// Setup sets up the plugin for use. This should only
// be called once per plugin.
func (p *Plugin) Setup() error {
	id := os.Getenv("CLARIFAI_CLIENT_ID")
	secret := os.Getenv("CLARIFAI_CLIENT_SECRET")

	if id == "" || secret == "" {
		p.configured = false
		return fmt.Errorf("credentials not found")
	}

	p.responses = make(map[string]*clarifai.TagResp)

	p.clientID = id
	p.secret = secret
	p.configured = true

	return nil
}
