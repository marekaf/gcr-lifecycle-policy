package worker

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"golang.org/x/oauth2"
)

func fetchCatalog(c Config, auth *oauth2.Token) Catalog {

	url := "https://" + c.RegistryURL + "/v2/_catalog"

	log.Debugf("fetching catalog from url %s", url)

	spaceClient := http.Client{
		Timeout: time.Second * 10, // Timeout after 10 seconds
	}

	body, err := getWithAuth(spaceClient, url, auth.AccessToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("got reponse, will unmarshal now %s", string(body))

	catalogResp := CatalogResponse{}
	jsonErr := json.Unmarshal(body, &catalogResp)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	if catalogResp.Errors != nil {
		log.Fatalf("response contains an error: %s", catalogResp.Errors)
	}

	catalog := Catalog{}
	for _, item := range catalogResp.Repositories {
		catalog.Repositories = append(catalog.Repositories, extractRepositoryFromImage(item))
	}

	return catalog
}

type request struct {
	url   string
	token string
}

func concFetchTags(req request, wg *sync.WaitGroup, c chan<- TagsResponse) {

	spaceClient := http.Client{
		Timeout: time.Second * 10, // Timeout after 10 seconds
	}

	body, err := getWithAuth(spaceClient, req.url, req.token)
	if err != nil {
		log.Fatal(err)
	}

	tags := TagsResponse{}
	jsonErr := json.Unmarshal(body, &tags)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	if body != nil {
		c <- tags
	} else {
		log.Fatalf("error fetching tags")
	}

	wg.Done()
}

func fetchTags(c Config, auth *oauth2.Token, catalog Catalog) ListResponse {

	list := ListResponse{}

	var wg sync.WaitGroup // create waitgroup (empty struct)

	queue := make(chan TagsResponse, len(catalog.Repositories))

	for _, repo := range catalog.Repositories {

		wg.Add(1)

		url := "https://" + c.RegistryURL + "/v2/" + repo.RepositoryPrefix + repo.ImageName + "/tags/list"

		req := request{
			url:   url,
			token: auth.AccessToken,
		}

		go concFetchTags(req, &wg, queue)

	}

	wg.Wait() // blocks here

	for range catalog.Repositories {
		item := <-queue

		list.TagsResponses = append(list.TagsResponses, item)
	}

	return list
}
