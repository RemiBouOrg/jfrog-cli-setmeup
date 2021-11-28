package show

import (
	"context"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/repository"
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

type repoSelection struct {
	url         string
	serverId    string
	repoKey     string
	description string
	unknown     bool
}

var showers = []struct {
	packageType     string
	getCurrentSetup func(ctx context.Context) []repoSelection
}{
	{
		repository.Maven,
		getCurrentMaven,
	},
	{
		repository.Nuget,
		getCurrentNuget,
	},
}

func showCommand(ctx context.Context) {
	for _, shower := range showers {
		repos := shower.getCurrentSetup(ctx)
		if len(repos) > 0 {
			log.Info(fmt.Sprintf("%s : ", shower.packageType))
			for _, repo := range repos {
				if repo.unknown {
					repo.serverId = "(Unknown)"
				}
				log.Info(fmt.Sprintf("%s - %s - %s", repo.serverId, repo.repoKey, repo.description))
			}
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
