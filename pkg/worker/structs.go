package worker

type Catalog struct {
	Repositories []string `json:"repositories"`
}

type TagsResponse struct {
	Child    []interface{}     `json:"child"`
	Manifest map[string]Digest `json:"manifest"`
	Name     string            `json:"name"`
	Tags     []string          `json:"tags"`
}
type Digest struct {
	ImageSizeBytes string   `json:"imageSizeBytes"`
	LayerID        string   `json:"layerId"`
	MediaType      string   `json:"mediaType"`
	Tag            []string `json:"tag"`
	TimeCreatedMs  string   `json:"timeCreatedMs"`
	TimeUploadedMs string   `json:"timeUploadedMs"`
	Name           string
}

type ListResponse struct {
	TagsResponses []TagsResponse `json:"tagsResponses"`
}

type Config struct {
	CredsFile     string   // path of credentials json file
	RepoFilter    []string // list of regions we want to check
	KeepTags      int
	RetentionDays int
	ClusterID     string
	RegistryURL   string
}

type FilteredList struct {
	TagsResponses []TagsResponse `json:"tagsResponses"`
}
