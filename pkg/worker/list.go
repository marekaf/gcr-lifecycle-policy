package worker

import (

	// kubernetes

	// we need this for our oauth2 token
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// HandleListCatalog function
func HandleListCatalog(c Config) Catalog {

	token := getToken(c.CredsFile)

	catalog := fetchCatalog(c, token)

	filteredCatalog := filterCatalog(catalog, c.RepoFilter)

	return filteredCatalog
}

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
