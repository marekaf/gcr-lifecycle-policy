package worker

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"golang.org/x/oauth2"
)

func cleanup(list FilteredList, c Config, auth *oauth2.Token) {

	spaceClient := http.Client{
		Timeout: time.Second * 10, // Timeout after 10 seconds
	}

	deleted := 0

	for _, repo := range list.TagsResponses {

		for _, digest := range repo.Manifest {

			// first delete all tags to prevent GOOGLE_MANIFEST_DANGLING_TAG error
			for _, tag := range digest.Tag {

				log.Infof("cleaning up tag: %s", tag)

				url := "https://" + c.RegistryURL + "/v2/" + repo.Name + "/manifests/" + tag

				if c.DryRun {
					log.Debugf("not calling HTTP DELETE %s because dryRun is enabled", url)
				} else {
					body, err := deleteWithAuth(spaceClient, url, auth.AccessToken)
					if err != nil {
						log.Fatal(err)
					}

					log.Debugf(string(body))
					// TODO: we should probably read the response body
				}
			}

			// second delete the image by referencing the manifest itself
			url := "https://" + c.RegistryURL + "/v2/" + repo.Name + "/manifests/" + digest.Name

			log.Infof("cleaning up image: %s", digest.Name)
			deleted++

			if c.DryRun {
				log.Debugf("not calling HTTP DELETE %s because dryRun is enabled", url)
			} else {
				body, err := deleteWithAuth(spaceClient, url, auth.AccessToken)
				if err != nil {
					log.Fatal(err)
				}

				log.Debugf(string(body))
				// TODO: we should probably read the response body
			}
		}

	}

	if c.DryRun {
		log.Infof("Total number of images deleted: '%d' (dryRun is enabled, so nothing was deleted!)", deleted)
	} else {
		log.Infof("Total number of images deleted: '%d'", deleted)
	}
}

// HandleCleanup function
func HandleCleanup(c Config) {

	token := getToken(c.CredsFile)

	catalog := fetchCatalog(c, token)
	filteredCatalog := filterCatalog(catalog, c.RepoFilter)

	registryList := fetchTags(c, token, filteredCatalog)
	clusterList := fetchImagesFromCluster(c)

	cleanupList := filter(c, registryList, clusterList)

	printBeforeCleanup(cleanupList)

	cleanup(cleanupList, c, token)
}
