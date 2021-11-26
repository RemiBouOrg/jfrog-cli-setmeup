package main

import (
	"github.com/jfrog/jfrog-cli-core/v2/plugins"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/environment"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/repository"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/show"
)

func main() {
	plugins.PluginMain(getApp())
}

func getApp() components.App {
	app := components.App{}
	app.Name = "Set Me Up"
	app.Description = "Easily set up your Artifactory repositories."
	app.Version = "v0.1.0"
	app.Commands = getCommands()
	return app
}

func getCommands() []components.Command {
	return []components.Command{
		repository.GetSetMeUpCommand(),
		environment.GetEnvInitCommand(),
		show.GetShowCommand(),
	}
}
