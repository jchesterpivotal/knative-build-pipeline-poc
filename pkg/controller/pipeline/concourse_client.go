/*
Copyright 2014-2017 Alex Suraci, Chris Brown, and Pivotal Software, Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the
License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
*/

// Modified from original: https://github.com/concourse/testflight/blob/master/helpers/concourse_client.go
// As this is not original code, I have included the license header from the concourse/testflight repo.
package pipeline

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/concourse/go-concourse/concourse"
	"golang.org/x/oauth2"
	"log"
)

func ConcourseClient(atcURL string, username, password string) (concourse.Client, error) {
	token, err := fetchToken(atcURL, username, password)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2.StaticTokenSource(token),
			Base: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}

	return concourse.NewClient(atcURL, httpClient, true), nil
}

func fetchToken(atcURL string, username, password string) (*oauth2.Token, error) {

	oauth2Config := oauth2.Config{
		ClientID:     "fly",
		ClientSecret: "Zmx5",
		Endpoint: oauth2.Endpoint{
			TokenURL: atcURL + "/oauth/token",
		},
		Scopes: []string{"openid", "federated:id"},
	}

	token, err := oauth2Config.PasswordCredentialsToken(context.Background(), username, password)

	log.Printf("Error fetching Concourse oauth token: %s\nConfig was: %+v\nToken was: %+v\n", err, oauth2Config, token)

	return token, err
}
