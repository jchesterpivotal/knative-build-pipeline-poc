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

// As this is not original code, I have included the Concourse license header.
// From: https://github.com/concourse/testflight/blob/b883e2e40a4172452a5b26123c10e6dd34660b9f/helpers/concourse_client.go
package pipeline

import (
	"github.com/concourse/go-concourse/concourse"
	"net/http"
	"golang.org/x/oauth2"
	"crypto/tls"
	"context"
	"log"
)

func ConcourseClient(atcURL string, username, password string) (concourse.Client, error) {
	token, err := fetchToken(atcURL, username, password)
	if err != nil {
		log.Printf("legacyAuth failed: %s\n", err)
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

	return concourse.NewClient(atcURL, httpClient, false), nil
}

func fetchToken(atcURL string, username, password string) (*oauth2.Token, error) {

	oauth2Config := oauth2.Config{
		ClientID:     "fly",
		ClientSecret: "Zmx5",
		Endpoint:     oauth2.Endpoint{TokenURL: atcURL + "/sky/token"},
		Scopes:       []string{"openid", "federated:id"},
	}

	return oauth2Config.PasswordCredentialsToken(context.Background(), username, password)
}
