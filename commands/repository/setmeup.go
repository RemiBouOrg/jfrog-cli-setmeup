package repository

import (
	"context"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/repository"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/artifactory"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/commons"
	"github.com/jfrog/jfrog-cli-plugin-template/jfrogconfig"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

func GetSetMeUpCommand() components.Command {
	return components.Command{
		Name:        "repository",
		Description: "Set up your environment to use Artifactory repository",
		Aliases:     []string{"r"},
		Arguments:   getRepositorySetMeUpArguments(),
		Flags:       getRepositorySetMeUpFlags(),
		Action: func(c *components.Context) error {
			serverDetails, err := config.GetSpecificConfig(c.GetStringFlagValue(commons.ServerIdFlag), true, false)
			if err != nil {
				return fmt.Errorf("unable to get server details : %w", err)
			}

			repoKeys, err := jfrogconfig.FindRepoKeys(c, serverDetails, c.GetStringFlagValue(commons.EnvNameFlag))
			if err != nil {
				return err
			}
			return setMeUpCommand(context.Background(), repoKeys, serverDetails)
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
			Name:         commons.ServerIdFlag,
			Description:  "The Jfrog Platform you want to use, if not set then the default one is used",
			DefaultValue: "",
		},
		components.StringFlag{
			Name:         commons.EnvNameFlag,
			Description:  "The environment you want to use, if not set then the default one is used",
			DefaultValue: "default",
		},
	}
}

type SetMeUpConfiguration struct {
	ServerDetails *config.ServerDetails
	RepoDetails   *artifactory.RepoDetails
}

var Handlers = map[string]func(context.Context, SetMeUpConfiguration) error{
	repository.Maven:  handleMaven,
	repository.Nuget:  handleNuget,
	repository.Docker: handleDocker,
	repository.Go:     handleGo,
}

func setMeUpCommand(ctx context.Context, repoKeys []string, serverDetails *config.ServerDetails) (err error) {
	for _, repoKey := range repoKeys {
		var conf = SetMeUpConfiguration{}
		conf.ServerDetails = serverDetails
		conf.RepoDetails, err = artifactory.GetRepoDetails(conf.ServerDetails, repoKey)
		if err != nil {
			return err
		}
		handler, hasHandler := Handlers[conf.RepoDetails.PackageType]
		if !hasHandler {
			return fmt.Errorf("%s package type is not handled", conf.RepoDetails.PackageType)
		}
		log.Info(fmt.Sprintf("Setting up repository %s of type %s on %s", repoKey, conf.RepoDetails.PackageType, conf.ServerDetails.ArtifactoryUrl))
		err = handler(ctx, conf)
		if err != nil {
			log.Error(fmt.Sprintf("An error occured : %v", err))
			return err
		}
	}

	return nil
}
