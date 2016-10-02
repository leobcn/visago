package visagoapi

import "github.com/zquestz/visago/util"

// Result is the struct passed back to the user.
type Result struct {
	Assets []*Asset `json:"assets,omitempty"`
	Errors []string `json:"errors,omitempty"`
}

// Asset represents each item fetched.
type Asset struct {
	Name string   `json:"name,omitempty"`
	Tags []string `json:"tags,omitempty"`
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

		for _, a := range v {
			mergedAsset.Tags = append(mergedAsset.Tags, a.Tags...)
		}

		mergedAsset.Tags = util.RemoveDuplicatesUnordered(mergedAsset.Tags)

		mergedAssets = append(mergedAssets, &mergedAsset)
	}

	return mergedAssets
}
