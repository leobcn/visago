package imagga

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
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
	contentIDs map[string]map[string]string
}

// Perform gathers metadata from imagga.
func (p *Plugin) Perform(c *visagoapi.PluginConfig) (string, visagoapi.PluginResult, error) {
	if p.configured == false {
		return "", nil, fmt.Errorf("not configured")
	}

	if len(c.URLs) == 0 && len(c.Files) == 0 {
		return "", nil, fmt.Errorf("must supply files/URLs")
	}

	client := &http.Client{}
	urlParams := []string{}

	responseID := nuid.Next()

	b, contentType, err := prepareUploadRequest(c)
	if err != nil {
		return "", nil, err
	}

	fileReq, _ := http.NewRequest("POST", "https://api.imagga.com/v1/content", b)
	fileReq.SetBasicAuth(p.apiKey, p.apiSecret)
	fileReq.Header.Set("Content-Type", contentType)

	fileResp, err := client.Do(fileReq)
	if err != nil {
		return "", nil, err
	}
	defer fileResp.Body.Close()

	fileBody, err := ioutil.ReadAll(fileResp.Body)
	if err != nil {
		return "", nil, err
	}

	fResp := resultFile{}
	err = json.Unmarshal(fileBody, &fResp)
	if err != nil {
		return "", nil, err
	}

	p.contentIDs[responseID] = make(map[string]string)

	for _, uploaded := range fResp.Uploaded {
		urlParams = append(urlParams, "content="+url.QueryEscape(uploaded.ID))
		p.contentIDs[responseID][uploaded.ID] = uploaded.Filename
	}

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
		k := result.Image

		if p.contentIDs[requestID][result.Image] != "" {
			k = p.contentIDs[requestID][result.Image]
		}

		tags[k] = make(map[string]*visagoapi.PluginTagResult)

		for _, t := range result.Tags {
			tag := &visagoapi.PluginTagResult{
				Name:       t.Tag,
				Confidence: t.Confidence,
			}

			tags[k][t.Tag] = tag
		}
	}

	return
}

// Reset clears the cache of existing responses.
func (p *Plugin) Reset() {
	p.responses = make(map[string]*response)
	p.contentIDs = make(map[string]map[string]string)
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
	p.contentIDs = make(map[string]map[string]string)

	p.apiKey = id
	p.apiSecret = secret
	p.configured = true

	return nil
}

func prepareUploadRequest(c *visagoapi.PluginConfig) (*bytes.Buffer, string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	defer w.Close()

	for idx, uFile := range c.Files {
		fw, err := w.CreateFormFile(fmt.Sprintf("%d", idx), uFile)
		if err != nil {
			return nil, "", err
		}

		f, err := os.Open(uFile)
		if err != nil {
			return nil, "", err
		}
		defer f.Close()

		_, err = io.Copy(fw, f)
		if err != nil {
			return nil, "", err
		}
	}

	contentType := w.FormDataContentType()

	return &b, contentType, nil
}
