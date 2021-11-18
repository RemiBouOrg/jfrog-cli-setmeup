package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/assets"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io/ioutil"
	"os"
)

func handleMaven(ctx context.Context, configuration SetMeUpConfiguration) error {
	const settingsFilePath = "%s/.m2/settings.xml"
	templateData := struct {
		UserName string
		Url      string
		RepoKey  string
		ApiKey   string
	}{
		UserName: configuration.serverDetails.User,
		Url:      fmt.Sprintf("%s%s", configuration.serverDetails.ArtifactoryUrl, configuration.repoDetails.Key),
		RepoKey:  configuration.repoDetails.Key,
		ApiKey:   configuration.serverDetails.Password,
	}
	template, err := assets.MavenSettingsTemplate()
	if err != nil {
		return err
	}
	dirname, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	userMavenFile := fmt.Sprintf(settingsFilePath, dirname)
	if _, err := os.Stat(userMavenFile); !errors.Is(err, os.ErrNotExist) {
		data, err := ioutil.ReadFile(userMavenFile)
		if err != nil {
			return err
		}
		backupFile := userMavenFile + ".bak"
		err = ioutil.WriteFile(backupFile, data, 0644)
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf("Settings backuped to %s", backupFile))
	}
	file, err := os.Create(userMavenFile)
	if err != nil {
		return err
	}
	err = template.Execute(file, templateData)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("Repo %s correctly set in %s", configuration.repoDetails.Key, userMavenFile))
	return nil
}
