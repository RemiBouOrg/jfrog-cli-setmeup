package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/artifactory"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/commons"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
)

var testMavenRepoKey = ""
var testNpmRepoKey = ""
var serverDetails *config.ServerDetails

func TestMain(m *testing.M) {
	// this gets the first maven repository from the default artifactory instance
	var err error
	serverDetails, err = config.GetDefaultServerConf()
	if err != nil {
		panic(err)
	}
	testMavenRepoKey = getRepoListFromDefaultServer("maven")[0].Key
	testNpmRepoKey = getRepoListFromDefaultServer("npm")[0].Key

	code := m.Run()
	os.Exit(code)
}

func getRepoListFromDefaultServer(repoType string) []artifactory.RepoDetails {
	authConfig, err := serverDetails.CreateArtAuthConfig()
	if err != nil {
		panic(err)
	}
	httpClientDetails := authConfig.CreateHttpClientDetails()
	jfrogHttpClient, err := jfroghttpclient.JfrogClientBuilder().Build()
	if err != nil {
		panic(fmt.Errorf("error occured when building http client: %w", err))
	}
	get, body, _, err := jfrogHttpClient.SendGet(fmt.Sprintf("%sapi/repositories?packageType=%s", authConfig.GetUrl(), repoType), false, &httpClientDetails)
	if err != nil {
		panic(fmt.Errorf("error occured when getting repository : %w", err))
	}
	if get.StatusCode != http.StatusOK {
		panic(fmt.Errorf("unexpected http status when getting repositories : %d", get.StatusCode))
	}
	repos := &[]artifactory.RepoDetails{}
	err = json.Unmarshal(body, repos)
	return *repos
}

func TestFailIfServerIdDoesntExists(t *testing.T) {
	badConfig := config.ServerDetails{}
	err := setMeUpCommand(context.Background(), []string{"test"}, &badConfig)
	require.Error(t, err)
}

func TestFailsIfRepoDoesntExists(t *testing.T) {
	err := setMeUpCommand(context.Background(), []string{testMavenRepoKey + "$$$"}, serverDetails)
	require.Error(t, err)
}

func TestOkIfRepoExists(t *testing.T) {
	_ = createTempDotM2(t)
	err := setMeUpCommand(context.Background(), []string{testMavenRepoKey}, serverDetails)
	require.NoError(t, err)
}

func TestGetSetMeUpCommandFlags(t *testing.T) {
	got := GetSetMeUpCommand(nil)
	var flagNames []string
	for _, flag := range got.Flags {
		flagNames = append(flagNames, flag.GetName())
	}
	assert.ElementsMatch(t, []string{commons.ServerIdFlag}, flagNames)
}

func TestGetSetMeUpCommandArgs(t *testing.T) {
	got := GetSetMeUpCommand(nil)
	var argNames []string
	for _, arg := range got.Arguments {
		argNames = append(argNames, arg.Name)
	}
	assert.ElementsMatch(t, []string{"repoKey"}, argNames)
}

type dummyFindRepo struct {
}

func (r dummyFindRepo) FindRepo(serverDetails *config.ServerDetails) (*artifactory.RepoDetails, error) {
	return &artifactory.RepoDetails{Key: "default-npm-local", PackageType: "npm"}, nil
}

func TestGetEnvInitCommandFuncNoArgs(t *testing.T) {
	findRepoService := &dummyFindRepo{}
	got := GetSetMeUpCommand(findRepoService)

	err := got.Action(&components.Context{})
	require.NoError(t, err)
}
