package worker

import (
	log "github.com/sirupsen/logrus"
	// kubernetes

	// we need this for our oauth2 token
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

	log.Debugf("got token %s", token.AccessToken)

	catalog := fetchCatalog(c, token)

	log.Debugf("got catalog %s", catalog)

	return catalog
}

// HandleListCluster function
func HandleListCluster(c Config) Catalog {

	return fetchImagesFromCluster(c)
}
