package repository

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io/ioutil"
	"os"
	"time"
)

// TODO set package.json: "publishConfig":{"registry":"http://my-internal-registry.local"}  https://docs.npmjs.com/cli/v7/using-npm/registry
// TODO consider set project level .npmrc file with credentials, ensure it's gitignored since it includes secrets. For now setting globally
// TODO consider adding support scoped registry: @myscope:registry=https://mycustomregistry.example.org

func handleNpm(configuration SetMeUpConfiguration) error {
	const settingsFilePath = "%s/.npmrc"
	const npmTemplate = "//%s%s/:_authToken=%s" // first param - artifactory url, second param - repo key, third param - token
	if configuration.repoDetails.PackageType != "npm" {
		return fmt.Errorf("unexpected repo type. Expected 'npm' but was: '%v'", configuration.repoDetails.PackageType)
	}
	dirname, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	npmrcFile := fmt.Sprintf(settingsFilePath, dirname)
	npmrcContent := []byte(fmt.Sprintf(npmTemplate, configuration.serverDetails.ArtifactoryUrl, configuration.repoDetails.Key, configuration.serverDetails.Password))
	if _, err := os.Stat(npmrcFile); !errors.Is(err, os.ErrNotExist) { // if exists
		data, err := ioutil.ReadFile(npmrcFile)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to read the content of '%s' with error: %+v", npmrcFile, err))
			return err
		}
		if bytes.Compare(data, npmrcContent) == 0 {
			log.Debug(fmt.Sprintf("content of '%s' is as expected. Skipping", npmrcFile))
			return nil
		}
		timestamp := time.Now().Format(time.RFC3339)
		backupNpmrc := fmt.Sprintf("%s.%s.bak", npmrcFile, timestamp)
		if err := os.Rename(npmrcFile, backupNpmrc); err != nil {
			log.Error(fmt.Sprintf("Failed to backup npmrc file (move from '%s' to '%s') with error: %+v", npmrcFile, backupNpmrc, err))
			return err
		}
	}
	if err = os.WriteFile(npmrcFile, npmrcContent, 0644); err != nil {
		log.Error(fmt.Sprintf("Failed to write npmrc file to '%s' with error: %+v", npmrcFile, err))
		return err
	}
	log.Info(fmt.Sprintf("Npm repo '%v' configured successfully at '%s'", configuration.repoDetails.Key, npmrcFile))
	return nil
}
