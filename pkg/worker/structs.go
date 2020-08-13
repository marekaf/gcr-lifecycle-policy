package worker

type Catalog struct {
	Repositories []Repository `json:"repositories"`
}

type CatalogResponse struct {
	Repositories []string `json:"repositories"`
}

type Repository struct {
	RepositoryPrefix string
	ImageName        string
	Tag              string
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
	CredsFile      string   // path of credentials json file
	RepoFilter     []string // list of regions we want to check
	KeepTags       int
	RetentionDays  int
	KubeconfigPath string
	RegistryURL    string
}

type FilteredList struct {
	TagsResponses []TagsResponse `json:"tagsResponses"`
}
