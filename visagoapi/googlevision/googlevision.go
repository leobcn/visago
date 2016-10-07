package googlevision

import (
	"fmt"
	"os"

	"google.golang.org/api/vision/v1"

	"github.com/kaneshin/pigeon"
	"github.com/kaneshin/pigeon/credentials"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/nats-io/nuid"
	"github.com/zquestz/visago/visagoapi"
)

func init() {
	visagoapi.AddPlugin("googlevision", &Plugin{})
}

// Plugin implements the Plugin interface and stores
// configuration data needed by the pigeon.
type Plugin struct {
	configured bool
	creds      string
	responses  map[string]*vision.BatchAnnotateImagesResponse
	items      map[string][]string
}

// Perform gathers metadata from the Google Vision API.
func (p *Plugin) Perform(c *visagoapi.PluginConfig) (string, visagoapi.PluginResult, error) {
	if p.configured == false {
		return "", nil, fmt.Errorf("not configured")
	}

	if len(c.URLs) == 0 && len(c.Files) == 0 {
		return "", nil, fmt.Errorf("must supply files/URLs")
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
		pigeon.NewFeature(pigeon.FaceDetection),
		pigeon.NewFeature(pigeon.ImageProperties),
	}

	items := []string{}
	items = append(items, c.URLs...)
	items = append(items, c.Files...)

	batch, err := client.NewBatchAnnotateImageRequest(items, features...)
	if err != nil {
		return "", nil, err
	}

	requestID := nuid.Next()

	p.items[requestID] = items
	p.responses[requestID], err = client.ImagesService().Annotate(batch).Do()
	if err != nil {
		return "", nil, err
	}

	return requestID, p, nil
}

// Tags returns the tags on an entry
func (p *Plugin) Tags(requestID string, score float64) (tags map[string]map[string]*visagoapi.PluginTagResult, err error) {
	tags = make(map[string]map[string]*visagoapi.PluginTagResult)

	if p.responses[requestID] == nil {
		return tags, fmt.Errorf("request has not been made to google")
	}

	for i, response := range p.responses[requestID].Responses {
		tags[p.items[requestID][i]] = make(map[string]*visagoapi.PluginTagResult)

		for _, annotation := range response.LabelAnnotations {
			if annotation.Score > score {
				tag := &visagoapi.PluginTagResult{
					Name:  annotation.Description,
					Score: annotation.Score,
				}

				tags[p.items[requestID][i]][annotation.Description] = tag
			}
		}
	}

	return
}

// Faces returns the faces on an entry
func (p *Plugin) Faces(requestID string) (faces map[string][]*visagoapi.PluginFaceResult, err error) {
	faces = make(map[string][]*visagoapi.PluginFaceResult)

	if p.responses[requestID] == nil {
		return faces, fmt.Errorf("request has not been made to google")
	}

	for i, response := range p.responses[requestID].Responses {
		for _, faceA := range response.FaceAnnotations {
			poly := &visagoapi.BoundingPoly{}
			for _, v := range faceA.BoundingPoly.Vertices {
				vertex := visagoapi.Vertex{
					X: v.X,
					Y: v.Y,
				}
				poly.Vertices = append(poly.Vertices, &vertex)
			}

			face := &visagoapi.PluginFaceResult{
				BoundingPoly:           poly,
				DetectionScore:         faceA.DetectionConfidence,
				JoyLikelihood:          faceA.JoyLikelihood,
				SorrowLikelihood:       faceA.SorrowLikelihood,
				AngerLikelihood:        faceA.AngerLikelihood,
				SurpriseLikelihood:     faceA.SurpriseLikelihood,
				UnderExposedLikelihood: faceA.UnderExposedLikelihood,
				BlurredLikelihood:      faceA.BlurredLikelihood,
				HeadwearLikelihood:     faceA.HeadwearLikelihood,
			}

			faces[p.items[requestID][i]] = append(faces[p.items[requestID][i]], face)
		}
	}

	return
}

// Colors returns the colors on an entry
func (p *Plugin) Colors(requestID string) (colors map[string][]*visagoapi.PluginColorResult, err error) {
	colors = make(map[string][]*visagoapi.PluginColorResult)

	if p.responses[requestID] == nil {
		return colors, fmt.Errorf("request has not been made to google")
	}

	for i, response := range p.responses[requestID].Responses {
		for _, c := range response.ImagePropertiesAnnotation.DominantColors.Colors {
			cf := colorful.Color{c.Color.Red / 255, c.Color.Green / 255, c.Color.Blue / 255}

			color := &visagoapi.PluginColorResult{
				Hex:           cf.Hex(),
				Red:           c.Color.Red,
				Green:         c.Color.Green,
				Blue:          c.Color.Blue,
				Alpha:         c.Color.Alpha,
				Score:         c.Score,
				PixelFraction: c.PixelFraction,
			}

			colors[p.items[requestID][i]] = append(colors[p.items[requestID][i]], color)
		}
	}

	return
}

// Reset clears the cache of existing responses.
func (p *Plugin) Reset() {
	p.responses = make(map[string]*vision.BatchAnnotateImagesResponse)
	p.items = make(map[string][]string)
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
	p.items = make(map[string][]string)

	p.creds = creds
	p.configured = true

	return nil
}
