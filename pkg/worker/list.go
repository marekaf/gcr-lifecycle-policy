package worker

import (

	// kubernetes

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// HandleList function
func HandleList(c Config) ListResponse {

	token := getToken(c.CredsFile)

	catalog := fetchCatalog(c, token)

	filteredCatalog := filterCatalog(catalog, c.RepoFilter)

	list := fetchTags(c, token, filteredCatalog)

	return list
}

// HandleListRepos function
func HandleListRepos(c Config) Catalog {

	token := getToken(c.CredsFile)

	catalog := fetchCatalog(c, token)

	return catalog
}

// HandleListCluster function
func HandleListCluster(c Config) Catalog {

	return fetchImagesFromCluster(c)
}
