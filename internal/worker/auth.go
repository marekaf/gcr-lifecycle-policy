package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

func getToken(credsFile string) *oauth2.Token {
	token, err := serviceAccount(credsFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return token
}

func serviceAccount(credentialFile string) (*oauth2.Token, error) {

	b, err := ioutil.ReadFile(credentialFile)
	if err != nil {
		return nil, err
	}
	var c = struct {
		Email      string `json:"client_email"`
		PrivateKey string `json:"private_key"`
	}{}

	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}
	config := &jwt.Config{
		Email:      c.Email,
		PrivateKey: []byte(c.PrivateKey),
		Scopes: []string{
			"https://www.googleapis.com/auth/cloud-platform.read-only",
			"https://www.googleapis.com/auth/devstorage.read_write",
		},
		TokenURL: google.JWTTokenURL,
	}
	token, err := config.TokenSource(context.TODO()).Token()
	if err != nil {
		return nil, err
	}
	return token, nil
}
