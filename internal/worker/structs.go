package worker

// Catalog holds info about all image repositories
type Catalog struct {
	Repositories []Repository `json:"repositories"`
}

// CatalogResponseError holds info about a possible error
type CatalogResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// CatalogResponse is the json api response
type CatalogResponse struct {
	Repositories []string               `json:"repositories"`
	Errors       []CatalogResponseError `json:"errors"`
}

// Repository is our struct to hold parsed prefix, name, tag of image
type Repository struct {
	RepositoryPrefix string
	ImageName        string
	Tag              string
}

// TagsResponse is the json api response
type TagsResponse struct {
	Child    []interface{}     `json:"child"`
	Manifest map[string]Digest `json:"manifest"`
	Name     string            `json:"name"`
	Tags     []string          `json:"tags"`
}

// Digest holds info about image digest
type Digest struct {
	ImageSizeBytes string   `json:"imageSizeBytes"`
	LayerID        string   `json:"layerId"`
	MediaType      string   `json:"mediaType"`
	Tag            []string `json:"tag"`
	TimeCreatedMs  string   `json:"timeCreatedMs"`
	TimeUploadedMs string   `json:"timeUploadedMs"`
	Name           string
}

// ListResponse is the json api response
type ListResponse struct {
	TagsResponses []TagsResponse `json:"tagsResponses"`
}

// Config is our application config for filtering, listing, cleaning
type Config struct {
	CredsFile       string   // path of credentials json file
	RepoFilter      []string // list of regions we want to check
	KeepTags        int
	RetentionDays   int
	KubeconfigPath  string
	RegistryURL     string
	SortBy          string
	ProtectTagRegex string
	DryRun          bool
}

// FilteredList holds list of tags that were already filtered out
type FilteredList struct {
	TagsResponses []TagsResponse `json:"tagsResponses"`
}
