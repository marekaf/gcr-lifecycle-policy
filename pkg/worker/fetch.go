package worker

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

func fetchCatalog(c Config, auth *oauth2.Token) Catalog {

	// TODO: make this configurable
	url := "https://" + c.RegistryURL + "/v2/_catalog"

	spaceClient := http.Client{
		Timeout: time.Second * 10, // Timeout after 10 seconds
	}

	body, err := getWithAuth(spaceClient, url, auth.AccessToken)
	if err != nil {
		log.Fatal(err)
	}

	catalogResp := CatalogResponse{}
	jsonErr := json.Unmarshal(body, &catalogResp)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	//log.Println(string(body))

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

		//log.Println(string(body))
		//log.Println("")

	}

	wg.Wait() // blocks here

	for range catalog.Repositories {
		item := <-queue

		list.TagsResponses = append(list.TagsResponses, item)
	}

	return list
}
