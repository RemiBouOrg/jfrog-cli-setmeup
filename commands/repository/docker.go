package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/artifactory"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/commons"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
)

func handleDocker(ctx context.Context, configuration SetMeUpConfiguration) error {
	get, jsonBytes, err := artifactory.ArtifactoryHttpGet(configuration.ServerDetails, "api/system/configuration/webServer")
	if err != nil {
		return err
	}
	host, port, err := findDockerHostAndPort(configuration, get, jsonBytes)
	if err != nil {
		return err
	}
	authConfig, _ := configuration.ServerDetails.CreateArtAuthConfig()
	command := exec.Command("docker",
		"login",
		"-u", authConfig.GetUser(),
		"-p", authConfig.GetPassword(),
		fmt.Sprintf("%s:%s", host, port))
	bufferString := bytes.NewBufferString("")
	command.Stderr = bufferString
	err = command.Run()
	if err != nil {
		return errors.Wrap(err, bufferString.String())
	}
	log.Info(fmt.Sprintf("Docker login to %s successful", host))
	return nil
}

func findDockerHostAndPort(configuration SetMeUpConfiguration, webServerResponse *http.Response, webServerJson []byte) (string, string, error) {

	if webServerResponse.StatusCode == 403 {
		parseArtiUrl, err := url.Parse(configuration.ServerDetails.ArtifactoryUrl)
		if err != nil {
			return "", "", err
		}
		// if 403 it's likely a cloud instance, we'll try to login with the hostname
		// if 403 is due to bad creds it will fail anyway
		log.Info("Artifactory is likely a cloud instance")
		return parseArtiUrl.Hostname(), "443", nil
	} else {
		var port string
		proxySetting := commons.ProxySettings{}
		err := json.Unmarshal(webServerJson, &proxySetting)
		if err != nil {
			return "", "", err
		}
		if proxySetting.UseHttps {
			port = strconv.Itoa(proxySetting.HttpsPort)
		} else {
			port = strconv.Itoa(proxySetting.HttpPort)
		}
		switch proxySetting.DockerReverseProxyMethod {
		case "SUBDOMAIN":
			host := fmt.Sprintf("%s.%s", configuration.RepoDetails.Key, proxySetting.ServerName)
			log.Info(fmt.Sprintf("Using subdomain per repository technique with %s:%s", host, port))
			return host, port, nil
		case "REPOPATHPREFIX":
			host := proxySetting.ServerName
			log.Info(fmt.Sprintf("Using path prefix technique on %s:%s", host, port))
			return host, port, nil
		case "PORTPERREPO":
			var host string
			for _, portConfig := range proxySetting.ReverseProxyRepositories.ReverseProxyRepoConfigs {
				if portConfig.RepoRef == configuration.RepoDetails.Key {
					host = portConfig.ServerName
					port = strconv.Itoa(portConfig.Port)
					break
				}
			}
			if host == "" {
				return "", "", fmt.Errorf("unable to find port config for %s", configuration.RepoDetails.Key)
			}
			log.Info(fmt.Sprintf("Using path prefix per repository technique with %s:%s", host, port))
			return host, port, nil
		default:
			return "", "", fmt.Errorf("non handled method %s", proxySetting.DockerReverseProxyMethod)

		}
	}
}
