package artifactory

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"net/http"
)

type RepoDetails struct {
	PackageType string `json:"packageType"`
	Key         string `json:"key"`
}

func ArtifactoryHttpGet(serverDetails *config.ServerDetails, path string) (*http.Response, []byte, error) {
	authConfig, _ := serverDetails.CreateArtAuthConfig()
	httpClientDetails := authConfig.CreateHttpClientDetails()
	jfrogHttpClient, _ := jfroghttpclient.JfrogClientBuilder().Build()
	url := fmt.Sprintf("%s%s", serverDetails.ArtifactoryUrl, path)
	get, body, _, err := jfrogHttpClient.SendGet(url, false, &httpClientDetails)
	return get, body, err
}

func GetRepoDetails(serverDetails *config.ServerDetails, repoKey string) (*RepoDetails, error) {
	get, body, err := ArtifactoryHttpGet(serverDetails, fmt.Sprintf("api/repositories/%s", repoKey))
	if err != nil {
		return nil, err
	}
	if get.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("could not find repository %s", repoKey))
	}
	repoDetails := &RepoDetails{}
	err = json.Unmarshal(body, &repoDetails)
	if err != nil {
		return nil, err
	}
	return repoDetails, nil
}

