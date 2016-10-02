package google_vision

import (
	"fmt"
	"os"

	"google.golang.org/api/vision/v1"

	"github.com/kaneshin/pigeon"
	"github.com/kaneshin/pigeon/credentials"
	"github.com/nats-io/nuid"
	"github.com/zquestz/visago/plugins"
)

func init() {
	plugins.AddPlugin("google_vision", &Plugin{})
}

// Plugin implements the Plugin interface and stores
// configuration data needed by the pigeon.
type Plugin struct {
	configured bool
	creds      string
	responses  map[string]*vision.BatchAnnotateImagesResponse
	urls       map[string][]string
}

// Perform gathers metadata from the Google Vision API. Currently
// we only support URLs but files will be added.
func (p *Plugin) Perform(c *plugins.PluginConfig) (string, plugins.PluginResult, error) {
	if p.configured == false {
		return "", nil, fmt.Errorf("not configured")
	}

	if len(c.URLs) == 0 {
		return "", nil, fmt.Errorf("must supply URLs")
	}

	// Reads from ENV var GOOGLE_APPLICATION_CREDENTIALS when blank.
	creds := credentials.NewApplicationCredentials("")

	config := pigeon.NewConfig().WithCredentials(creds)

	client, err := pigeon.New(config)
	if err != nil {
		return "", nil, err
	}

	features := []*vision.Feature{
		pigeon.NewFeature(pigeon.LabelDetection),
	}

	batch, err := client.NewBatchAnnotateImageRequest(c.URLs, features...)
	if err != nil {
		return "", nil, err
	}

	requestID := nuid.Next()

	p.urls[requestID] = c.URLs
	p.responses[requestID], err = client.ImagesService().Annotate(batch).Do()
	if err != nil {
		return "", nil, err
	}

	return requestID, p, nil
}

// Tags returns the tags on an entry
func (p *Plugin) Tags(requestID string) (tags map[string][]string, err error) {
	tags = make(map[string][]string)

	if p.responses[requestID] == nil {
		return tags, fmt.Errorf("request has not been made to google")
	}

	for i, response := range p.responses[requestID].Responses {
		retTags := []string{}

		for _, annotation := range response.LabelAnnotations {
			retTags = append(retTags, annotation.Description)
		}

		tags[p.urls[requestID][i]] = retTags
	}

	return
}

// Reset clears the cache of existing responses.
func (p *Plugin) Reset() {
	p.responses = make(map[string]*vision.BatchAnnotateImagesResponse)
	p.urls = make(map[string][]string)
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
	creds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	if creds == "" {
		p.configured = false
		return fmt.Errorf("credentials not found")
	}

	p.responses = make(map[string]*vision.BatchAnnotateImagesResponse)
	p.urls = make(map[string][]string)

	p.creds = creds
	p.configured = true

	return nil
}
