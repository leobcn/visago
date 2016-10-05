package clarifai

import (
	"fmt"
	"os"

	"github.com/nats-io/nuid"
	"github.com/zquestz/clarifai-go"
	"github.com/zquestz/visago/visagoapi"
)

func init() {
	visagoapi.AddPlugin("clarifai", &Plugin{})
}

// Plugin implements the Plugin interface and stores
// configuration data needed by the clarifai library.
type Plugin struct {
	configured bool
	clientID   string
	secret     string
	responses  map[string][]*clarifai.TagResp
	files      map[string][]string
}

// Perform gathers metadata from Clarifai.
func (p *Plugin) Perform(c *visagoapi.PluginConfig) (string, visagoapi.PluginResult, error) {
	if p.configured == false {
		return "", nil, fmt.Errorf("not configured")
	}

	if len(c.URLs) == 0 && len(c.Files) == 0 {
		return "", nil, fmt.Errorf("must supply files/URLs")
	}

	client := clarifai.NewClient(p.clientID, p.secret)
	_, err := client.Info()
	if err != nil {
		return "", nil, err
	}

	requestID := nuid.Next()

	// Unfortunately the Clarifai API doesn't support URLs and Files
	// in a single request.
	if len(c.URLs) > 0 {
		urlResp, err := client.Tag(clarifai.TagRequest{
			URLs: c.URLs,
		})
		if err != nil {
			return "", nil, err
		}

		p.responses[requestID] = append(p.responses[requestID], urlResp)
	}

	if len(c.Files) > 0 {
		p.files[requestID] = c.Files

		filesResp, err := client.Tag(clarifai.TagRequest{
			Files: c.Files,
		})
		if err != nil {
			return "", nil, err
		}

		p.responses[requestID] = append(p.responses[requestID], filesResp)
	}

	return requestID, p, nil
}

// Tags returns the tags on an entry
func (p *Plugin) Tags(requestID string, score float64) (tags map[string]map[string]*visagoapi.PluginTagResult, err error) {
	tags = make(map[string]map[string]*visagoapi.PluginTagResult)

	if p.responses[requestID] == nil {
		return tags, fmt.Errorf("request has not been made to clarifai")
	}

	for _, req := range p.responses[requestID] {
		for i, result := range req.Results {
			var k string

			if result.URL != "" {
				k = result.URL
			} else {
				k = p.files[requestID][i]
			}

			tags[k] = make(map[string]*visagoapi.PluginTagResult)

			for i, t := range result.Result.Tag.Classes {
				confidence := float64(result.Result.Tag.Probs[i])

				if confidence > score {
					tag := &visagoapi.PluginTagResult{
						Name:  t,
						Score: confidence,
					}

					tags[k][t] = tag
				}
			}
		}
	}

	return
}

// Faces returns the faces on an entry
func (p *Plugin) Faces(requestID string) (faces map[string][]*visagoapi.PluginFaceResult, err error) {
	faces = make(map[string][]*visagoapi.PluginFaceResult)

	return
}

// Reset clears the cache of existing responses.
func (p *Plugin) Reset() {
	p.responses = make(map[string][]*clarifai.TagResp)
	p.files = make(map[string][]string)
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

	p.responses = make(map[string][]*clarifai.TagResp)
	p.files = make(map[string][]string)

	p.clientID = id
	p.secret = secret
	p.configured = true

	return nil
}
