package repository

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/artifactory"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/pkg/errors"
	"os/exec"
)

func handleNuget(ctx context.Context, configuration SetMeUpConfiguration) error {
	feedUrl := fmt.Sprintf("api/nuget/v3/%s", configuration.repoDetails.Key)
	get, _, err := artifactory.ArtifactoryHttpGet(configuration.serverDetails, feedUrl)
	if err != nil {
		return err
	}
	if get.StatusCode == 404 {
		log.Info(fmt.Sprintf("%s is not a v3 nuget repository", configuration.repoDetails.Key))
		feedUrl = fmt.Sprintf("api/nuget/%s", configuration.repoDetails.Key)
		get, _, err := artifactory.ArtifactoryHttpGet(configuration.serverDetails, feedUrl)
		if err != nil || get.StatusCode != 200 {
			return fmt.Errorf("cannot find nuget repo version %s", configuration.repoDetails.Key)
		}
	} else {
		log.Info(fmt.Sprintf("%s is a v3 nuget repository", configuration.repoDetails.Key))
	}

	_ = exec.Command("nuget", "sources", "Remove",
		"-Name", "Artifactory",
		"-NonInteractive",
	).Run()
	authConfig, _ := configuration.serverDetails.CreateArtAuthConfig()
	command := exec.Command("nuget", "sources", "Add",
		"-Name", "Artifactory",
		"-Source", fmt.Sprintf("%s%s", configuration.serverDetails.ArtifactoryUrl, feedUrl),
		"-UserName", authConfig.GetUser(),
		"-Password", authConfig.GetPassword(),
		"-NonInteractive",
	)
	bufferString := bytes.NewBufferString("")
	command.Stderr = bufferString
	err = command.Run()
	if err != nil {
		return errors.Wrap(err, bufferString.String())
	}
	log.Info(fmt.Sprintf("Nuget feed named 'Artifactory' succesfuly set to %s", feedUrl))
	return nil
}
