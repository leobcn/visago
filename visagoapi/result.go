package visagoapi

// Result is the struct passed back to the user.
type Result struct {
	Assets []*Asset `json:"assets,omitempty"`
	Errors []string `json:"errors,omitempty"`
}

// Asset represents each item fetched.
type Asset struct {
	Name string                 `json:"name,omitempty"`
	Tags map[string][]*AssetTag `json:"tags,omitempty"`
}

// AssetTag reprsents the details for that tag.
// The score is a value from 0 to 1.
type AssetTag struct {
	Name  string  `json:"name,omitempty"`
	Score float64 `json:"score,omitempty"`
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

		mergedAsset.Tags = make(map[string][]*AssetTag)

		for _, a := range v {
			for tk := range a.Tags {
				mergedAsset.Tags[tk] = append(mergedAsset.Tags[tk], a.Tags[tk]...)
			}
		}

		mergedAssets = append(mergedAssets, &mergedAsset)
	}

	return mergedAssets
}
