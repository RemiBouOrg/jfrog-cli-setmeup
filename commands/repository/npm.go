package repository

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/artifactory"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// TODO set package.json: "publishConfig":{"registry":"http://my-internal-registry.local"}  https://docs.npmjs.com/cli/v7/using-npm/registry
// TODO consider set project level .npmrc file with credentials, ensure it's gitignored since it includes secrets. For now setting globally
// TODO consider adding support scoped registry: @myscope:registry=https://mycustomregistry.example.org

func handleNpm(ctx context.Context, configuration SetMeUpConfiguration) error {
	if configuration.RepoDetails.PackageType != "npm" {
		return fmt.Errorf("unexpected repo type. Expected 'npm' but was: '%v'", configuration.RepoDetails.PackageType)
	}
	dirname, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	npmrcFile := fmt.Sprintf("%s/.npmrc", dirname)
	npmConfig, err := getNpmAuthConfig(configuration.ServerDetails)
	if err != nil {
		return err
	}
	registryUrl := fmt.Sprintf("%sapi/npm/%s/", configuration.ServerDetails.ArtifactoryUrl, configuration.RepoDetails.Key)
	npmrcContent := []byte(fmt.Sprintf("registry=%s\n%s", registryUrl, npmConfig))
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
	log.Info(fmt.Sprintf("Npm repo '%v' configured successfully at '%s'", configuration.RepoDetails.Key, npmrcFile))
	return nil
}

func getNpmAuthConfig(details *config.ServerDetails) (string, error) {
	res, body, err := artifactory.ArtifactoryHttpGet(details, "api/npm/auth")
	if err != nil {
		return "", fmt.Errorf("failed to fetch npm token: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch npm token, bad status: %v", res.Status)
	}
	return string(body), nil
}
