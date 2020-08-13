package worker

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

func fetchCatalog(c Config, auth *oauth2.Token) Catalog {

	// TODO: make this configurable
	url := "https://" + c.RegistryURL + "/v2/_catalog"

	spaceClient := http.Client{
		Timeout: time.Second * 10, // Timeout after 10 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
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

func fetchTags(c Config, auth *oauth2.Token, catalog Catalog) ListResponse {

	list := ListResponse{}

	for _, repo := range catalog.Repositories {

		url := "https://" + c.RegistryURL + "/v2/" + repo.RepositoryPrefix + repo.ImageName + "/tags/list"

		spaceClient := http.Client{
			Timeout: time.Second * 10, // Timeout after 10 seconds
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+auth.AccessToken)

		res, getErr := spaceClient.Do(req)
		if getErr != nil {
			log.Fatal(getErr)
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

		tags := TagsResponse{}
		jsonErr := json.Unmarshal(body, &tags)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		//log.Println(string(body))

		list.TagsResponses = append(list.TagsResponses, tags)
	}

	return list
}
