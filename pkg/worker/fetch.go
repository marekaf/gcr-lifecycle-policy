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

	catalog := Catalog{}
	jsonErr := json.Unmarshal(body, &catalog)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return catalog
}

func filterCatalog(c Catalog, filter []string) Catalog {

	// by default not filter anything
	if len(filter) == 0 {
		return c
	}

	filtered := Catalog{}

	for _, repo := range c.Repositories {

		for _, filterRepo := range filter {
			if repo == filterRepo {
				filtered.Repositories = append(filtered.Repositories, filterRepo)
				break
			}
		}
	}

	return filtered
}

func fetchTags(c Config, auth *oauth2.Token, catalog Catalog) ListResponse {

	list := ListResponse{}

	for _, image := range catalog.Repositories {

		url := "https://" + c.RegistryURL + "/v2/" + image + "/tags/list"

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

		list.TagsResponses = append(list.TagsResponses, tags)
	}

	return list
}
