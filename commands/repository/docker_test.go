package repository

import (
	"context"
	"encoding/json"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/artifactory"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/commons"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func Test_FindDockerHostPort(t *testing.T) {
	testData := []struct {
		name                string
		webServerJson       commons.ProxySettings
		webServerStatusCode int
		repoKey             string
		artiUrl             string
		expectedHost        string
		expectedPort        string
		expectErr           bool
	}{
		{
			name:                "On 403 use arti host and https",
			webServerStatusCode: 403,
			artiUrl:             "https://google.com",
			expectedHost:        "google.com",
			expectedPort:        "443",
		},
		{
			name:         "On SUBDOMAIN concat repo key with hostname",
			repoKey:      "test",
			expectedHost: "test.google.com",
			expectedPort: "443",
			webServerJson: commons.ProxySettings{
				ServerName:               "google.com",
				DockerReverseProxyMethod: "SUBDOMAIN",
				HttpsPort:                443,
				UseHttps:                 true,
			},
		},
		{
			name:         "On SUBDOMAIN if no ttps use http port",
			repoKey:      "test",
			expectedHost: "test.google.com",
			expectedPort: "444",
			webServerJson: commons.ProxySettings{
				ServerName:               "google.com",
				DockerReverseProxyMethod: "SUBDOMAIN",
				HttpsPort:                443,
				HttpPort:                 444,
				UseHttps:                 false,
			},
		},
		{
			name:         "On REPOPATHPREFIX use server name",
			artiUrl:      "google.com",
			repoKey:      "test",
			expectedHost: "google.com",
			expectedPort: "443",
			webServerJson: commons.ProxySettings{
				ServerName:               "google.com",
				DockerReverseProxyMethod: "REPOPATHPREFIX",
				HttpsPort:                443,
				UseHttps:                 true,
			},
		},
		{
			name:         "On PORTPERREPO use specific port and host",
			artiUrl:      "google.com",
			repoKey:      "test",
			expectedHost: "testgoogle.com",
			expectedPort: "888",
			webServerJson: commons.ProxySettings{
				ServerName:               "google.com",
				DockerReverseProxyMethod: "PORTPERREPO",
				ReverseProxyRepositories: commons.ReverseProxyRepositories{
					ReverseProxyRepoConfigs: []commons.ReverseProxyRepoConfigs{
						{RepoRef: "test", Port: 888, ServerName: "testgoogle.com"},
						{RepoRef: "test2", Port: 999, ServerName: "test3google.com"},
					}},
				HttpsPort: 443,
				UseHttps:  true,
			},
		},
		{
			name:      "On PORTPERREPO if no port, then err",
			artiUrl:   "google.com",
			repoKey:   "test3",
			expectErr: true,
			webServerJson: commons.ProxySettings{
				ServerName:               "google.com",
				DockerReverseProxyMethod: "PORTPERREPO",
				ReverseProxyRepositories: commons.ReverseProxyRepositories{
					ReverseProxyRepoConfigs: []commons.ReverseProxyRepoConfigs{
						{RepoRef: "test", Port: 888, ServerName: "testgoogle.com"},
						{RepoRef: "test2", Port: 999, ServerName: "test3google.com"},
					}},
				HttpsPort: 443,
				UseHttps:  true,
			},
		},
		{
			name:      "On unexpected method, err",
			artiUrl:   "google.com",
			repoKey:   "test3",
			expectErr: true,
			webServerJson: commons.ProxySettings{
				ServerName:               "google.com",
				DockerReverseProxyMethod: "DONTEXISTS",
			},
		},
	}
	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			webServerJsonStr, _ := json.Marshal(test.webServerJson)
			host, port, err := findDockerHostAndPort(
				SetMeUpConfiguration{
					ServerDetails: &config.ServerDetails{ArtifactoryUrl: test.artiUrl},
					RepoDetails:   &artifactory.RepoDetails{Key: test.repoKey},
				},
				&http.Response{StatusCode: test.webServerStatusCode},
				webServerJsonStr)
			if test.expectErr {
				require.Error(t, err)
			} else {
				require.Equal(t, test.expectedPort, port)
				require.Equal(t, test.expectedHost, host)
			}
		})
	}
}

func Test_handleDocker(t *testing.T) {
	testDockerRepo := getRepoListFromDefaultServer("docker")[0].Key
	err := handleDocker(context.Background(),
		SetMeUpConfiguration{
			ServerDetails: serverDetails,
			RepoDetails: &artifactory.RepoDetails{
				PackageType: "docker",
				Key:         testDockerRepo,
			},
		},
	)
	require.NoError(t, err)
}
