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
		Name:        "set-me-up",
		Description: "Set up your environment to use Artifactory repositoryKey",
		Aliases:     []string{"smu"},
		Arguments:   getSetMeUpArguments(),
		Flags:       getSetMeUpFlags(),
		Action: func(c *components.Context) error {
			return setMeUpCommand(c)
		},
	}
}

func getSetMeUpArguments() []components.Argument {
	return []components.Argument{
		{
			Name:        "repositoryKey",
			Description: "The repositoryKey you want to use",
		},
	}
}

func getSetMeUpFlags() []components.Flag {
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

type RepoDetails struct {
	PackageType string `json:"packageType"`
}

var handlers = map[string]func(SetMeUpConfiguration) error{
	repository.Maven: handleMaven,
}

func setMeUpCommand(c *components.Context) error {
	if len(c.Arguments) != 1 {
		return errors.New("Wrong number of arguments. Expected: 1, " + "Received: " + strconv.Itoa(len(c.Arguments)))
	}
	var conf = SetMeUpConfiguration{}
	conf.repositoryKey = c.Arguments[0]
	serverId := c.GetStringFlagValue("server-id")
	serverDetails, err := config.GetSpecificConfig(serverId, true, false)
	if err != nil {
		return fmt.Errorf("unable to get server details : %w", err)
	}
	conf.serverDetails = serverDetails
	authConfig, err := conf.serverDetails.CreateArtAuthConfig()
	if err != nil {
		return fmt.Errorf("unable to get artifactory details : %w", err)
	}
	httpClientDetails := authConfig.CreateHttpClientDetails()
	jfrogHttpClient, err := jfroghttpclient.JfrogClientBuilder().Build()
	if err != nil {
		return fmt.Errorf("cannot create http client : %w", err)
	}
	get, body, _, err := jfrogHttpClient.SendGet(fmt.Sprintf("%sapi/repositories/%s", authConfig.GetUrl(), conf.repositoryKey), false, &httpClientDetails)
	if err != nil {
		return fmt.Errorf("error occured when querying repository details : %w", err)
	}
	if get.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected http status when getting repository details : %d", get.StatusCode)
	}
	conf.repoDetails = RepoDetails{}
	err = json.Unmarshal(body, &conf.repoDetails)
	if err != nil {
		return fmt.Errorf("can not read repository details : %w", err)
	}
	handler, hasHandler := handlers[conf.repoDetails.PackageType]
	if !hasHandler {
		return fmt.Errorf("%s package type is not handled", conf.repoDetails.PackageType)
	}
	log.Info(fmt.Sprintf("Setting up repository %s of type %s on %s", conf.repositoryKey, conf.repoDetails.PackageType, authConfig.GetUrl()))
	return handler(conf)
}
