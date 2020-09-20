package worker

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func extractRepositoryFromImage(input string) Repository {

	prefix := ""
	image := ""
	tag := ""

	// +1 to keep the trailing slash
	whereToSplit := strings.LastIndex(input, "/") + 1
	prefix = input[:whereToSplit]
	imagetag := input[whereToSplit:]

	// no +1 because we don't want the ':'
	wts := strings.LastIndex(imagetag, ":")
	if wts == -1 {
		image = imagetag
	} else {
		image = imagetag[:wts]
		tag = imagetag[wts+1:]
	}

	repo := Repository{
		RepositoryPrefix: prefix,
		ImageName:        image,
		Tag:              tag,
	}

	return repo
}

func daysToTime(days int) time.Time {
	return time.Now().AddDate(0, -days, 0)
}

func reqWithAuth(req *http.Request, c http.Client, url string, token string) ([]byte, error) {

	req.Header.Set("Authorization", "Bearer "+token)

	req.Header.Set("Content-Type", "application/json")

	res, getErr := c.Do(req)
	if getErr != nil {
		return nil, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	return body, nil
}

func getWithAuth(c http.Client, url string, token string) ([]byte, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := reqWithAuth(req, c, url, token)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func deleteWithAuth(c http.Client, url string, token string) ([]byte, error) {

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := reqWithAuth(req, c, url, token)
	if err != nil {
		return nil, err
	}

	return body, nil
}
