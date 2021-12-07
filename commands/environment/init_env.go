package environment

import (
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/artifactory"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/commons"
	"github.com/jfrog/jfrog-cli-plugin-template/jfrogconfig"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"gopkg.in/yaml.v2"
)

func GetEnvInitCommand(findRepoService commons.FindRepoService) components.Command {
	return components.Command{
		Name:        "environment",
		Description: "Store repository config on the current dir to share with the team members",
		Aliases:     []string{"e", "env", "environ"},
		Arguments:   getInitEnvArguments(),
		Flags:       getEnvFlags(),
		Action: func(c *components.Context) error {
			serverId := c.GetStringFlagValue(commons.ServerIdFlag)
			serverDetails, err := config.GetSpecificConfig(serverId, true, false)
			if err != nil {
				return fmt.Errorf("unable to get server details : %w", err)
			}

			var repoKey string
			if len(c.Arguments) == 1 {
				repoKey = c.Arguments[0]
			} else {
				repo, err := findRepoService.FindRepo(serverDetails)
				if err != nil {
					return err
				}
				if repo != nil {
					repoKey = repo.Key
				}
			}

			return InitEnv(serverDetails, repoKey, c.GetStringFlagValue(commons.EnvNameFlag))
		},
	}
}

func getInitEnvArguments() []components.Argument {
	return []components.Argument{
		{
			Name:        "repoKey",
			Description: "The repository key you want to use, if not presented, a selection dropdown will be shown",
		},
	}
}

func getEnvFlags() []components.Flag {
	return []components.Flag{
		components.StringFlag{
			Name:         commons.ServerIdFlag,
			Description:  "The JFrog Platform you want to use, if not set then the default one is used",
			DefaultValue: "",
		},
		components.StringFlag{
			Name:         commons.EnvNameFlag,
			Description:  "The environment you want to use, if not set then the default one is used",
			DefaultValue: "default",
		},
	}
}

func InitEnv(serverDetails *config.ServerDetails, repoKey string, envName string) error {

	repoDetails, err := artifactory.GetRepoDetails(serverDetails, repoKey)
	if err != nil {
		return err
	}

	confFile, err := jfrogconfig.ReadCurrentConfFile()
	if err != nil {
		return err
	}
	if confFile == nil {
		confFile = jfrogconfig.JFrogConfFile{}
	}

	if confFile[envName] == nil {
		confFile[envName] = jfrogconfig.RepoTypeToName{}
	}

	confFile[envName][repoDetails.PackageType] = repoKey
	confContent, err := yaml.Marshal(&confFile)
	if err != nil {
		return err
	}

	err = jfrogconfig.WriteConfigFile(confContent)
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("%s of type %s added to the configuration", repoKey,
		repoDetails.PackageType))

	return nil
}
