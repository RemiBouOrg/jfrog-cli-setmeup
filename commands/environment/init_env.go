package environment

import (
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/artifactory"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/commons"
	"github.com/jfrog/jfrog-cli-plugin-template/jfrogconfig"
	"gopkg.in/yaml.v2"
	"strconv"
)



func GetEnvInitCommand() components.Command {
	return components.Command{
		Name:        "init",
		Description: "Store repository config on the current dir to share with the team members",
		Aliases:     []string{"i"},
		Arguments:   getInitEnvArguments(),
		Flags:       getEnvFlags(),
		Action: func(c *components.Context) error {
			if len(c.Arguments) != 1 {
				return errors.New("Wrong number of arguments. Expected: 1, " + "Received: " + strconv.Itoa(len(c.Arguments)))
			}

			return InitEnv(c.Arguments[0], c.GetStringFlagValue(commons.ServerIdFlag), c.GetStringFlagValue(commons.EnvNameFlag))
		},
	}
}

func getInitEnvArguments() []components.Argument {
	return []components.Argument{
		{
			Name:        "repoKey",
			Description: "The repositoryKey you want to use",
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

func InitEnv(repoKey string, serverId string, envName string) error {
	serverDetails, err := config.GetSpecificConfig(serverId, true, false)
	if err != nil {
		return fmt.Errorf("unable to get server details : %w", err)
	}
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

	return jfrogconfig.WriteConfigFile(confContent)
}
