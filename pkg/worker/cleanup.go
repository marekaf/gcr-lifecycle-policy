package worker

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/oauth2"
)

func cleanup(list FilteredList, c Config, auth *oauth2.Token) {

	spaceClient := http.Client{
		Timeout: time.Second * 10, // Timeout after 10 seconds
	}

	for _, repo := range list.TagsResponses {

		for _, digest := range repo.Manifest {

			// first delete all tags to prevent GOOGLE_MANIFEST_DANGLING_TAG error
			for _, tag := range repo.Tags {

				url := "https://" + c.RegistryURL + "/v2/" + repo.Name + "/manifests/" + tag

				body, err := deleteWithAuth(spaceClient, url, auth.AccessToken)
				if err != nil {
					log.Fatal(err)
				}

				// tags := TagsResponse{}
				// jsonErr := json.Unmarshal(body, &tags)
				// if jsonErr != nil {
				// 	log.Fatal(jsonErr)
				// }

				log.Println(string(body))

			}

			// second delete the image by referencing the manifest itself
			url := "https://" + c.RegistryURL + "/v2/" + repo.Name + "/manifests/" + digest.Name

			body, err := deleteWithAuth(spaceClient, url, auth.AccessToken)
			if err != nil {
				log.Fatal(err)
			}

			// tags := TagsResponse{}
			// jsonErr := json.Unmarshal(body, &tags)
			// if jsonErr != nil {
			// 	log.Fatal(jsonErr)
			// }

			log.Println(string(body))

		}

	}
}

func olderThanRetention(d Digest, retention time.Time) bool {

	timecreated, err := strconv.ParseInt(d.TimeCreatedMs, 10, 64)

	if err != nil {
		log.Fatal(err)
	}

	// retention is in days, convert it to ms
	return time.Unix(timecreated/1000, 0).Before(retention)
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
