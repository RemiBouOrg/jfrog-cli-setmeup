package repository

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/artifactory"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/pkg/errors"
	"os/exec"
)

func handleNuget(ctx context.Context, configuration SetMeUpConfiguration) error {
	feedUrl := fmt.Sprintf("api/nuget/v3/%s", configuration.RepoDetails.Key)
	get, _, err := artifactory.ArtifactoryHttpGet(configuration.ServerDetails, feedUrl)
	if err != nil {
		return err
	}
	if get.StatusCode == 404 {
		log.Info(fmt.Sprintf("%s is not a v3 nuget repository", configuration.RepoDetails.Key))
		feedUrl = fmt.Sprintf("api/nuget/%s", configuration.RepoDetails.Key)
		get, _, err := artifactory.ArtifactoryHttpGet(configuration.ServerDetails, feedUrl)
		if err != nil || get.StatusCode != 200 {
			return fmt.Errorf("cannot find nuget repo version %s", configuration.RepoDetails.Key)
		}
	} else {
		log.Info(fmt.Sprintf("%s is a v3 nuget repository", configuration.RepoDetails.Key))
	}

	_ = exec.Command("nuget", "sources", "Remove",
		"-Name", "Artifactory",
		"-NonInteractive",
	).Run()
	authConfig, _ := configuration.ServerDetails.CreateArtAuthConfig()
	command := exec.Command("nuget", "sources", "Add",
		"-Name", "Artifactory",
		"-Source", fmt.Sprintf("%s%s", configuration.ServerDetails.ArtifactoryUrl, feedUrl),
		"-UserName", authConfig.GetUser(),
		"-Password", resolvePassword(authConfig),
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

func resolvePassword(authConfig auth.ServiceDetails) string {
	if authConfig.GetPassword() != "" {
		return authConfig.GetPassword()
	}
	if authConfig.GetApiKey() != "" {
		return authConfig.GetApiKey()
	}
	if authConfig.GetAccessToken() != "" {
		return authConfig.GetAccessToken()
	}
	log.Debug(fmt.Sprintf("Failed to detect credentials, fallback to empty sptring - might fail"))
	return ""
}
