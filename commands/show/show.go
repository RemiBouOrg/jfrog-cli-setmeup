package show

import (
	"context"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/repository"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/manifoldco/promptui"
	"os"
	"strings"
	"text/template"
)

func GetShowCommand() components.Command {
	return components.Command{
		Name:        "show",
		Description: "Display all currently setup repository by reading configuration/environment",
		Aliases:     []string{"s"},
		Action: func(c *components.Context) error {
			return ShowCommand(context.Background())
		},
	}
}

type repoSelection struct {
	ServerId    string
	RepoKey     string
	Description string
	Unknown     bool
}

var showers = []struct {
	PackageType     string
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
	{
		repository.Go,
		getCurrentGolang,
	},
	{
		repository.Docker,
		getCurrentDocker,
	},
	{
		repository.Npm,
		getCurrentNpm,
	},
}

func ShowCommand(ctx context.Context) error {
	repoTypeTmpl, _ := template.New("repoType").
		Funcs(promptui.FuncMap).
		Parse("{{ .PackageType }} :: \n")
	repoLineTmpl, _ := template.New("repoLine").
		Funcs(promptui.FuncMap).
		Parse("\t{{if .Unknown}}" +
			"{{ \"unknown\" | printf \"%-10s\" | cyan  }} {{ .RepoKey | yellow }}" +
			"{{else}}" +
			"{{ .ServerId | printf \"%-10s\" | green  }} {{ .RepoKey | yellow }}" +
			"{{end}} {{.Description}}\n")
	for _, shower := range showers {
		repos := shower.getCurrentSetup(ctx)
		if len(repos) > 0 {
			repoTypeTmpl.Execute(os.Stdout, shower)
			for _, repo := range repos {
				repoLineTmpl.Execute(os.Stdout, repo)
			}
		} else {
			log.Debug(fmt.Sprintf("%s : no repository setup", shower.PackageType))
		}
	}

	return nil
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
