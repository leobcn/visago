package visagoapi

// Result is the struct passed back to the user.
type Result struct {
	Assets []*Asset `json:"assets,omitempty"`
	Errors []string `json:"errors,omitempty"`
}

// Asset represents each item fetched.
type Asset struct {
	Name   string                          `json:"name,omitempty"`
	Tags   map[string][]*PluginTagResult   `json:"tags,omitempty"`
	Colors map[string][]*PluginColorResult `json:"colors,omitempty"`
	Faces  []*PluginFaceResult             `json:"faces,omitempty"`
	Source string                          `json:"-"`
}

func mergeAssets(assets []*Asset) []*Asset {
	mergedAssets := []*Asset{}

	assetMap := make(map[string][]*Asset)
	for _, asset := range assets {
		if _, ok := assetMap[asset.Name]; !ok {
			assetMap[asset.Name] = []*Asset{}
		}

		assetMap[asset.Name] = append(assetMap[asset.Name], asset)
	}

	for k, v := range assetMap {
		mergedAsset := Asset{
			Name: k,
		}

		mergedAsset.Tags = make(map[string][]*PluginTagResult)
		mergedAsset.Faces = []*PluginFaceResult{}
		mergedAsset.Colors = make(map[string][]*PluginColorResult)

		for _, a := range v {
			for tk := range a.Tags {
				for _, t := range a.Tags[tk] {
					nt := &PluginTagResult{
						Name:   t.Name,
						Score:  t.Score,
						Source: a.Source,
					}

					mergedAsset.Tags[tk] = append(mergedAsset.Tags[tk], nt)
				}
			}

			for ck := range a.Colors {
				for _, c := range a.Colors[ck] {
					nc := &PluginColorResult{
						Source:        a.Source,
						Score:         c.Score,
						Alpha:         c.Alpha,
						Hex:           c.Hex,
						Blue:          c.Blue,
						Green:         c.Green,
						Red:           c.Red,
						PixelFraction: c.PixelFraction,
					}

					mergedAsset.Colors[ck] = append(mergedAsset.Colors[ck], nc)
				}
			}

			for _, f := range a.Faces {
				nf := &PluginFaceResult{
					Source:                 a.Source,
					BoundingPoly:           f.BoundingPoly,
					DetectionScore:         f.DetectionScore,
					JoyLikelihood:          f.JoyLikelihood,
					SorrowLikelihood:       f.SorrowLikelihood,
					AngerLikelihood:        f.AngerLikelihood,
					SurpriseLikelihood:     f.SurpriseLikelihood,
					UnderExposedLikelihood: f.UnderExposedLikelihood,
					BlurredLikelihood:      f.BlurredLikelihood,
					HeadwearLikelihood:     f.HeadwearLikelihood,
				}

				mergedAsset.Faces = append(mergedAsset.Faces, nf)
			}
		}

		mergedAssets = append(mergedAssets, &mergedAsset)
	}

	return mergedAssets
}
