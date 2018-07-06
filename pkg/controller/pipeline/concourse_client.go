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
// Frankenstein'd from originals:
//		https://github.com/concourse/testflight/blob/065d3854302c3d3c60f119362e7311cdc3949104/helpers/concourse_client.go
//		https://github.com/concourse/fly/blob/b470a7d130522711157fdcaad2629ed0fc9fe484/commands/login.go
// Using earlier versions which are compatible with Concourse 3.x.
// The code on master at time of writing was aimed at the upcoming 4.x authentication flow.
package pipeline

import (
	"crypto/tls"
	"encoding/json"
	"github.com/concourse/go-concourse/concourse"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

func ConcourseClient(atcURL string) (concourse.Client, error) {
	tokenType, tokenVal, err := legacyAuth(atcURL)
	if err != nil {
		log.Printf("legacyAuth failed: %s\n", err)
		return nil, err
	}

	httpClient := oauthClient(tokenType, tokenVal)

	return concourse.NewClient(atcURL, httpClient, true), nil
}

func legacyAuth(atcUrl string) (string, string, error) {
	request, err := http.NewRequest("GET", atcUrl+"/api/v1/teams/main/auth/token", nil)
	if err != nil {
		return "", "", err
	}
	request.SetBasicAuth("concourse", "concourse") // from the Helm chart defaults

	tokenResponse, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", "", err
	}

	type authToken struct {
		Type  string `json:"token_type"`
		Value string `json:"token_value"`
	}

	defer tokenResponse.Body.Close()

	var token authToken
	json.NewDecoder(tokenResponse.Body).Decode(&token)

	return token.Type, token.Value, nil
}

func oauthClient(tokenType, tokenVal string) *http.Client {
	return &http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2.StaticTokenSource(&oauth2.Token{
				TokenType:   tokenType,
				AccessToken: tokenVal,
			}),
			Base: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}
