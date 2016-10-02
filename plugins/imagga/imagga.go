package imagga

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/nats-io/nuid"
	"github.com/zquestz/visago/plugins"
)

func init() {
	plugins.AddPlugin("imagga", &Plugin{})
}

// Plugin implements the Plugin interface and stores
// configuration data needed by the imagga library.
type Plugin struct {
	configured bool
	apiKey     string
	apiSecret  string
	responses  map[string][]*response
}

// Perform gathers metadata from imagga, for the first pass
// it only supports urls.
func (p *Plugin) Perform(c *plugins.PluginConfig) (string, plugins.PluginResult, error) {
	if p.configured == false {
		return "", nil, fmt.Errorf("not configured")
	}

	if len(c.URLs) == 0 {
		return "", nil, fmt.Errorf("must supply urls")
	}

	responses := []*response{}

	client := &http.Client{}

	for _, uri := range c.URLs {
		req, _ := http.NewRequest("GET", "https://api.imagga.com/v1/tagging?url="+url.QueryEscape(uri), nil)
		req.SetBasicAuth(p.apiKey, p.apiSecret)

		resp, err := client.Do(req)
		if err != nil {
			return "", nil, err
		}
		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", nil, err
		}

		uResp := response{}
		err = json.Unmarshal(respBody, &uResp)
		if err != nil {
			return "", nil, err
		}

		responses = append(responses, &uResp)
	}

	responseID := nuid.Next()

	p.responses[responseID] = responses

	return responseID, p, nil
}

// Tags returns the tagging information for the request
func (p *Plugin) Tags(requestID string) (tags map[string][]string, err error) {
	tags = make(map[string][]string)

	if p.responses[requestID] == nil {
		return tags, fmt.Errorf("request has not been made to imagga")
	}

	for _, response := range p.responses[requestID] {
		for _, result := range response.Results {
			collectTags := []string{}

			for _, t := range result.Tags {
				collectTags = append(collectTags, t.Tag)
			}

			tags[result.Image] = collectTags
		}
	}

	return
}

// Reset clears the cache of existing responses.
func (p *Plugin) Reset() {
	p.responses = make(map[string][]*response)
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
	id := os.Getenv("IMAGGA_API_KEY")
	secret := os.Getenv("IMAGGA_API_SECRET")

	if id == "" || secret == "" {
		p.configured = false
		return fmt.Errorf("credentials not found")
	}

	p.responses = make(map[string][]*response)

	p.apiKey = id
	p.apiSecret = secret
	p.configured = true

	return nil
}
