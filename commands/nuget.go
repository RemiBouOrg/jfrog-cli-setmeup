package commands

import (
	"bytes"
	"fmt"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/pkg/errors"
	"os/exec"
)

func handleNuget(configuration SetMeUpConfiguration) error {
	authConfig, _ := configuration.serverDetails.CreateArtAuthConfig()
	httpClientDetails := authConfig.CreateHttpClientDetails()
	jfrogHttpClient, _ := jfroghttpclient.JfrogClientBuilder().Build()
	feedUrl := fmt.Sprintf("%sapi/nuget/v3/%s", configuration.serverDetails.ArtifactoryUrl, configuration.repositoryKey)
	get, _, _, err := jfrogHttpClient.SendGet(feedUrl, false, &httpClientDetails)
	if err != nil {
		return err
	}
	if get.StatusCode == 404 {
		log.Info(fmt.Sprintf("%s is not a v3 nuget repository", configuration.repositoryKey))
		feedUrl = fmt.Sprintf("%sapi/nuget/%s", configuration.serverDetails.ArtifactoryUrl, configuration.repositoryKey)
		get, _, _, err = jfrogHttpClient.SendGet(feedUrl, false, &httpClientDetails)
		if err != nil || get.StatusCode != 200 {
			return fmt.Errorf("cannot find nuget repo version %s", configuration.repositoryKey)
		}
	} else {
		log.Info(fmt.Sprintf("%s is a v3 nuget repository", configuration.repositoryKey))
	}

	_ = exec.Command("nuget", "sources", "Remove",
		"-Name", "Artifactory",
		"-NonInteractive",
	).Run()
	command := exec.Command("nuget", "sources", "Add",
		"-Name", "Artifactory",
		"-Source", feedUrl,
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
