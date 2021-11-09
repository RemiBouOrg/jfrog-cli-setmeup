package commands

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

var testMavenRepoKey = ""
var serverDetails *config.ServerDetails

func TestMain(m *testing.M) {
	// this gets the first maven repository from the default artifactory instance
	var err error
	serverDetails, err = config.GetDefaultServerConf()
	if err != nil {
		panic(err)
	}
	authConfig, err := serverDetails.CreateArtAuthConfig()
	if err != nil {
		panic(err)
	}
	httpClientDetails := authConfig.CreateHttpClientDetails()
	jfrogHttpClient, err := jfroghttpclient.JfrogClientBuilder().Build()
	if err != nil {
		panic(fmt.Errorf("error occured when building http client: %w", err))
	}
	get, body, _, err := jfrogHttpClient.SendGet(fmt.Sprintf("%sapi/repositories?packageType=maven", authConfig.GetUrl()), false, &httpClientDetails)
	if err != nil {
		panic(fmt.Errorf("error occured when getting repository : %w", err))
	}
	if get.StatusCode != http.StatusOK {
		panic(fmt.Errorf("unexpected http status when getting repositories : %d", get.StatusCode))
	}
	repos := &[]RepoDetails{}
	err = json.Unmarshal(body, repos)
	testMavenRepoKey = (*repos)[0].Key
	m.Run()
}

func TestFailIfServerIdDoesntExists(t *testing.T) {
	err := setMeUpCommand("test", "donotexists")
	require.Error(t, err)
}

func TestFailsIfRepoDoesntExists(t *testing.T) {
	err := setMeUpCommand(testMavenRepoKey+"$$$", "")
	require.Error(t, err)
}

func TestOkIfRepoExists(t *testing.T) {
	err := setMeUpCommand(testMavenRepoKey, "")
	require.NoError(t, err)
}
