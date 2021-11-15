package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/repository"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"strconv"
)

func GetSetMeUpCommand() components.Command {
	return components.Command{
		Name:        "repository",
		Description: "Set up your environment to use Artifactory repository",
		Aliases:     []string{"r"},
		Arguments:   getRepositorySetMeUpArguments(),
		Flags:       getRepositorySetMeUpFlags(),
		Action: func(c *components.Context) error {
			if len(c.Arguments) != 1 {
				return errors.New("Wrong number of arguments. Expected: 1, " + "Received: " + strconv.Itoa(len(c.Arguments)))
			}
			return setMeUpCommand(c.Arguments[0], c.GetStringFlagValue("server-id"))
		},
	}
}

func getRepositorySetMeUpArguments() []components.Argument {
	return []components.Argument{
		{
			Name:        "repoKey",
			Description: "The repositoryKey you want to use",
		},
	}
}

func getRepositorySetMeUpFlags() []components.Flag {
	return []components.Flag{
		components.StringFlag{
			Name:         "server-id",
			Description:  "The Jfrog Platform you want to use, if not set then the default one is used",
			DefaultValue: "",
		},
	}
}

type SetMeUpConfiguration struct {
	repositoryKey string
	serverDetails *config.ServerDetails
	repoDetails   RepoDetails
}

func (c SetMeUpConfiguration) artifactoryHttpGet(path string) (*http.Response, []byte, error) {
	authConfig, _ := c.serverDetails.CreateArtAuthConfig()
	httpClientDetails := authConfig.CreateHttpClientDetails()
	jfrogHttpClient, _ := jfroghttpclient.JfrogClientBuilder().Build()
	url := fmt.Sprintf("%s%s", c.serverDetails.ArtifactoryUrl, path)
	get, body, _, err := jfrogHttpClient.SendGet(url, false, &httpClientDetails)
	return get, body, err
}

type RepoDetails struct {
	PackageType string `json:"packageType"`
	Key         string `json:"key"`
}

var handlers = map[string]func(SetMeUpConfiguration) error{
	repository.Maven:  handleMaven,
	repository.Nuget:  handleNuget,
	repository.Docker: handleDocker,
}

func setMeUpCommand(repoKey string, serverId string) error {
	var conf = SetMeUpConfiguration{}
	conf.repositoryKey = repoKey
	serverDetails, err := config.GetSpecificConfig(serverId, true, false)
	if err != nil {
		return fmt.Errorf("unable to get server details : %w", err)
	}
	conf.serverDetails = serverDetails
	err = getRepoDetails(&conf)
	if err != nil {
		return err
	}
	handler, hasHandler := handlers[conf.repoDetails.PackageType]
	if !hasHandler {
		return fmt.Errorf("%s package type is not handled", conf.repoDetails.PackageType)
	}
	log.Info(fmt.Sprintf("Setting up repository %s of type %s on %s", conf.repositoryKey, conf.repoDetails.PackageType, conf.serverDetails.ArtifactoryUrl))
	err = handler(conf)
	if err != nil {
		log.Error(fmt.Sprintf("An error occured : %v", err))
		return err
	}
	return nil
}

func getRepoDetails(conf *SetMeUpConfiguration) error {
	get, body, err := conf.artifactoryHttpGet(fmt.Sprintf("api/repositories/%s", conf.repositoryKey))
	if err != nil {
		return fmt.Errorf("error occured when querying repository details : %w", err)
	}
	if get.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected http status when getting repository details : %d (%s)", get.StatusCode, string(body))
	}
	repoDetails := RepoDetails{}
	err = json.Unmarshal(body, &repoDetails)
	if err != nil {
		return fmt.Errorf("can not read repository details : %w", err)
	}
	conf.repoDetails = repoDetails
	return nil
}
