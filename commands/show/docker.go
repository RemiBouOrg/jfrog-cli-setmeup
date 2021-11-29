package show

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/artifactory"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/commons"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type dockerConf struct {
	Auths map[string]interface{}
}

func getCurrentDocker(ctx context.Context) []repoSelection {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Error(fmt.Sprintf("error reading home dir : %v", err))
		return nil
	}
	dockerAuthFile := fmt.Sprintf("%s/.docker/config.json", dirname)
	bytes, err := ioutil.ReadFile(dockerAuthFile)
	if err != nil {
		log.Debug(fmt.Sprintf("error reading docker config : %v", err))
		return nil
	}
	dockerAuths := dockerConf{}
	err = json.Unmarshal(bytes, &dockerAuths)
	if err != nil {
		log.Error(fmt.Sprintf("error unmarshalling docker json config : %v", err))
		return nil
	}
	serverConfigs, err := config.GetAllServersConfigs()
	if err != nil {
		log.Error(fmt.Sprintf("error reading all servers config : %v", err))
		return nil
	}
	var proxySettingsPerServer = make(map[*config.ServerDetails]commons.ProxySettings)
	for _, serverConfig := range serverConfigs {
		get, jsonBytes, err := artifactory.ArtifactoryHttpGet(serverConfig, "api/system/configuration/webServer")
		if err != nil {
			log.Debug(fmt.Sprintf("error getting proxy settings : %v", err))
			continue
		}
		if get.StatusCode != 403 {
			proxySetting := commons.ProxySettings{}
			err = json.Unmarshal(jsonBytes, &proxySetting)
			if err != nil {
				log.Error(fmt.Sprintf("error unmarshallling proxy settings : %v", err))
				return nil
			}
			proxySettingsPerServer[serverConfig] = proxySetting
		} else {
			url, err := url.Parse(serverConfig.ArtifactoryUrl)
			if err != nil {
				log.Error(fmt.Sprintf("error parsing url : %v", err))
				return nil
			}
			proxySettingsPerServer[serverConfig] = commons.ProxySettings{
				DockerReverseProxyMethod: "CLOUD",
				ServerName:               url.Host,
				UseHttps:                 true,
				HttpsPort:                443,
			}
		}
	}
	return getDockerReposSelections(dockerAuths, proxySettingsPerServer)
}

func getDockerReposSelections(auths dockerConf, configs map[*config.ServerDetails]commons.ProxySettings) []repoSelection {
	res := make([]repoSelection, 0)
	for auth, _ := range auths.Auths {
		res = append(res, findAuthMatchingRepoSelection(auth, configs))
	}
	return res
}

func findAuthMatchingRepoSelection(auth string, configs map[*config.ServerDetails]commons.ProxySettings) repoSelection {
	authSplit := strings.Split(auth, ":")
	host := authSplit[0]
	port := ""
	if len(authSplit) == 1 {
		port = "443"
	} else {
		port = authSplit[1]
	}
	portInt, _ := strconv.Atoi(port)
	for details, settings := range configs {
		if settings.DockerReverseProxyMethod == "SUBDOMAIN" {
			if strings.HasSuffix(host, settings.ServerName) && matchesServerPort(settings, portInt) {
				return repoSelection{
					serverId:    details.ServerId,
					repoKey:     strings.TrimSuffix(host, fmt.Sprintf(".%s", settings.ServerName)),
					description: "",
					unknown:     false,
				}
			}
		}
		if settings.DockerReverseProxyMethod == "REPOPATHPREFIX" || settings.DockerReverseProxyMethod == "CLOUD" {
			if settings.ServerName == host && matchesServerPort(settings, portInt) {
				return repoSelection{
					serverId:    details.ServerId,
					repoKey:     "(ALL REPOSITORIES)",
					description: "",
					unknown:     false,
				}
			}
		}
		if settings.DockerReverseProxyMethod == "PORTPERREPO" {
			for _, repoConfig := range settings.ReverseProxyRepositories.ReverseProxyRepoConfigs {
				if repoConfig.ServerName == host && repoConfig.Port == portInt {
					return repoSelection{
						serverId:    details.ServerId,
						repoKey:     repoConfig.RepoRef,
						description: "",
						unknown:     false,
					}
				}
			}
		}
	}
	return repoSelection{
		serverId:    "",
		repoKey:     auth,
		description: "",
		unknown:     true,
	}
}

func matchesServerPort(settings commons.ProxySettings, portInt int) bool {
	return (settings.UseHttp && settings.HttpPort == portInt) ||
		(settings.UseHttps && settings.HttpsPort == portInt)
}
