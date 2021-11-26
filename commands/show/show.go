package show

import (
	"context"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"strings"
)

func GetShowCommand() components.Command {
	return components.Command{
		Name:        "show",
		Description: "Display all currently setup repository by reading configuration/environment",
		Aliases:     []string{"s"},
		Action: func(c *components.Context) error {
			showCommand(context.Background())
			return nil
		},
	}
}

var showers = []struct {
	packageType     string
	getCurrentSetup func(ctx context.Context) (serverId string, repoKey string)
}{
	{
		"maven",
		getCurrentMaven,
	},
}

func showCommand(ctx context.Context) {
	for _, shower := range showers {
		serverId, repoKey := shower.getCurrentSetup(ctx)
		if repoKey != "" {
			log.Info(fmt.Sprintf("%s : %s (%s)", shower.packageType, repoKey, serverId))
		} else {
			log.Debug(fmt.Sprintf("%s : no repository setup", shower.packageType))
		}
	}
}

func findServerIdByUrl(url string) string {
	configs, err := config.GetAllServersConfigs()
	if err != nil {
		log.Debug("error occurred when reading server configs")
		return ""
	}
	for _, serverConf := range configs {
		if strings.HasPrefix(url, serverConf.ArtifactoryUrl) {
			return serverConf.ServerId
		}
	}
	return ""
}
