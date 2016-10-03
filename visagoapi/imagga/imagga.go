package imagga

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/nats-io/nuid"
	"github.com/zquestz/visago/visagoapi"
)

func init() {
	visagoapi.AddPlugin("imagga", &Plugin{})
}

// Plugin implements the Plugin interface and stores
// configuration data needed by the imagga library.
type Plugin struct {
	configured bool
	apiKey     string
	apiSecret  string
	responses  map[string]*response
}

// Perform gathers metadata from imagga, for the first pass
// it only supports urls.
// TODO: Implement file handling.
func (p *Plugin) Perform(c *visagoapi.PluginConfig) (string, visagoapi.PluginResult, error) {
	if p.configured == false {
		return "", nil, fmt.Errorf("not configured")
	}

	if len(c.URLs) == 0 {
		return "", nil, fmt.Errorf("must supply urls")
	}

	client := &http.Client{}

	urlParams := []string{}
	for _, uri := range c.URLs {
		urlParams = append(urlParams, "url="+url.QueryEscape(uri))
	}

	req, _ := http.NewRequest("GET", "https://api.imagga.com/v1/tagging?"+strings.Join(urlParams, "&"), nil)
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

	responseID := nuid.Next()

	p.responses[responseID] = &uResp

	return responseID, p, nil
}

// Tags returns the tagging information for the request
func (p *Plugin) Tags(requestID string) (tags map[string]map[string]*visagoapi.PluginTagResult, err error) {
	if p.responses[requestID] == nil {
		return tags, fmt.Errorf("request has not been made to imagga")
	}

	tags = make(map[string]map[string]*visagoapi.PluginTagResult)

	for _, result := range p.responses[requestID].Results {
		tags[result.Image] = make(map[string]*visagoapi.PluginTagResult)

		for _, t := range result.Tags {
			tag := &visagoapi.PluginTagResult{
				Name:       t.Tag,
				Confidence: t.Confidence,
			}

			tags[result.Image][t.Tag] = tag
		}
	}

	return
}

// Reset clears the cache of existing responses.
func (p *Plugin) Reset() {
	p.responses = make(map[string]*response)
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

	p.responses = make(map[string]*response)

	p.apiKey = id
	p.apiSecret = secret
	p.configured = true

	return nil
}
