package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/pkg/errors"
	"net/url"
	"os/exec"
	"strconv"
)

func handleDocker(configuration SetMeUpConfiguration) error {
	get, jsonBytes, err := configuration.artifactoryHttpGet("api/system/configuration/webServer")
	if err != nil {
		return err
	}
	var host string
	var port string
	parseArtiUrl, err := url.Parse(configuration.serverDetails.ArtifactoryUrl)
	if err != nil {
		return err
	}
	if get.StatusCode == 403 {
		// if 403 it's likely a cloud instance, we'll try to login with the hostname
		// if 403 is due to bad creds it will fail anyway
		host = parseArtiUrl.Hostname()
		port = parseArtiUrl.Port()
		log.Info("Artifactory is likely a cloud instance")
	} else {
		proxySetting := proxySettings{}
		err := json.Unmarshal(jsonBytes, &proxySetting)
		if err != nil {
			return err
		}
		if proxySetting.UseHttps {
			port = strconv.Itoa(proxySetting.HttpsPort)
		} else {
			port = strconv.Itoa(proxySetting.HttpPort)
		}
		switch proxySetting.DockerReverseProxyMethod {
		case "SUBDOMAIN":
			host = fmt.Sprintf("%s.%s", configuration.repositoryKey, parseArtiUrl.Hostname())
			log.Info(fmt.Sprintf("Using subdomain per repository technique with %s:%s", host, port))
		case "REPOPATHPREFIX":
			host = parseArtiUrl.Hostname()
			log.Info(fmt.Sprintf("Using path prefix technique on %s:%s", host, port))
		case "PORTPERREPO":
			for _, portConfig := range proxySetting.ReverseProxyRepositories.ReverseProxyRepoConfigs {
				if portConfig.RepoRef == configuration.repositoryKey {
					host = portConfig.ServerName
					port = strconv.Itoa(portConfig.Port)
					break
				}
			}
			if host == "" {
				return fmt.Errorf("unable to find port config for %s", configuration.repositoryKey)
			}
			log.Info(fmt.Sprintf("Using path prefix per repository technique with %s:%s", host, port))
		default:
			return fmt.Errorf("non handled method %s")

		}
	}
	authConfig, _ := configuration.serverDetails.CreateArtAuthConfig()
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

type reverseProxyRepoConfigs struct {
	RepoRef    string
	Port       int
	ServerName string
}

type reverseProxyRepositories struct {
	ReverseProxyRepoConfigs []reverseProxyRepoConfigs
}

type proxySettings struct {
	UseHttp                  bool
	UseHttps                 bool
	HttpPort                 int
	HttpsPort                int
	DockerReverseProxyMethod string
	ReverseProxyRepositories reverseProxyRepositories
}
