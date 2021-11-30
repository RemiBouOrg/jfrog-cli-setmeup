package environment

import (
	"context"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/artifactory"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/commons"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/repository"
	"github.com/jfrog/jfrog-cli-plugin-template/jfrogconfig"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

func GetEnvApplyCommand() components.Command {
	return components.Command{
		Name:        "apply",
		Description: "Apply repository from config in the current",
		Aliases:     []string{"a"},
		Arguments:   []components.Argument{},
		Flags:       getEnvFlags(),
		Action: func(c *components.Context) (err error) {
			return ApplyEnv(context.Background(), c.GetStringFlagValue(commons.ServerIdFlag), c.GetStringFlagValue(commons.EnvNameFlag))
		},
	}
}

func ApplyEnv(ctx context.Context, serverId, envName string) error {
	if envName == "" {
		log.Debug("env name not provided, fallback to 'default'")
		envName = "default"
	}

	// empty serverId will result in default, see jfrog-cli-core/v2@v2.3.0/utils/config/config.go:48
	serverDetails, err := config.GetSpecificConfig(serverId, true, false)
	if err != nil {
		return fmt.Errorf("unable to get server details : %w", err)
	}

	confFile, err := jfrogconfig.ReadCurrentConfFile()
	if err != nil {
		return err
	}
	if confFile == nil {
		return errors.New("failed to apply - empty config file")
	}

	repoTypeToName, ok := confFile[envName]
	if !ok {
		return fmt.Errorf("env '%v' not configured", repoTypeToName)
	}

	for repoType, repoKey := range repoTypeToName {
		handler, ok := repository.Handlers[repoType]
		if !ok {
			return fmt.Errorf("unsupported repo type '%v'", repoType)
		}
		setMeUpConfiguration := repository.SetMeUpConfiguration{
			ServerDetails: serverDetails,
			RepoDetails: &artifactory.RepoDetails{
				PackageType: repoType,
				Key:         repoKey,
			},
		}
		err := handler(ctx, setMeUpConfiguration)
		if err != nil {
			return err
		}
	}
	return nil
}
