package imagga

type response struct {
	Results []*resultEntry `json:"results"`
}

type resultEntry struct {
	Image string            `json:"image"`
	Tags  []*resultEntryTag `json:"tags"`
}

type resultEntryTag struct {
	Confidence float64 `json:"confidence"`
	Tag        string  `json:"tag"`
}

type resultFile struct {
	Status   string          `json:"status"`
	Uploaded []*resultUpload `json:"uploaded"`
}

type resultUpload struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
}
