package show

import (
	"context"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/commons"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetDockerReposSelection(t *testing.T) {
	docker := getCurrentDocker(context.Background())
	require.NotEmpty(t, docker)
}

func TestGetDockerReposSelectionSubdomainHttpsImplicit(t *testing.T) {
	selections := getDockerReposSelections(
		dockerConf{
			Auths: map[string]interface{}{
				"repo.docker.com": "",
			}},
		map[*config.ServerDetails]commons.ProxySettings{
			&config.ServerDetails{
				ServerId: "test1",
			}: {
				ServerName:               "docker.com",
				DockerReverseProxyMethod: "SUBDOMAIN",
				UseHttps:                 true,
				HttpsPort:                443,
			},
		},
	)
	require.Equal(t, []repoSelection{
		{
			serverId: "test1",
			repoKey:  "repo",
			unknown:  false,
		},
	}, selections)
}

func TestGetDockerReposSelectionSubdomainPortExplicit(t *testing.T) {
	selections := getDockerReposSelections(
		dockerConf{
			Auths: map[string]interface{}{
				"repo.docker.com:333": "",
			}},
		map[*config.ServerDetails]commons.ProxySettings{
			&config.ServerDetails{
				ServerId: "test1",
			}: {
				ServerName:               "docker.com",
				DockerReverseProxyMethod: "SUBDOMAIN",
				UseHttps:                 true,
				HttpsPort:                333,
			},
		},
	)
	require.Equal(t, []repoSelection{
		{
			serverId: "test1",
			repoKey:  "repo",
			unknown:  false,
		},
	}, selections)
}

func TestGetDockerReposSelectionSubdomainHttp(t *testing.T) {
	selections := getDockerReposSelections(
		dockerConf{
			Auths: map[string]interface{}{
				"repo.docker.com:80": "",
			}},
		map[*config.ServerDetails]commons.ProxySettings{
			&config.ServerDetails{
				ServerId: "test1",
			}: {
				ServerName:               "docker.com",
				DockerReverseProxyMethod: "SUBDOMAIN",
				UseHttp:                  true,
				HttpPort:                 80,
			},
		},
	)
	require.Equal(t, []repoSelection{
		{
			serverId: "test1",
			repoKey:  "repo",
			unknown:  false,
		},
	}, selections)
}

func TestGetDockerReposSelectionRepoPathPrefix(t *testing.T) {
	selections := getDockerReposSelections(
		dockerConf{
			Auths: map[string]interface{}{
				"repo.docker.com": "",
			}},
		map[*config.ServerDetails]commons.ProxySettings{
			&config.ServerDetails{
				ServerId: "test1",
			}: {
				ServerName:               "repo.docker.com",
				DockerReverseProxyMethod: "REPOPATHPREFIX",
				UseHttps:                 true,
				HttpsPort:                443,
			},
		},
	)
	require.Equal(t, []repoSelection{
		{
			serverId: "test1",
			repoKey:  "(ALL REPOSITORIES)",
			unknown:  false,
		},
	}, selections)
}

func TestGetDockerReposSelectionCloud(t *testing.T) {
	selections := getDockerReposSelections(
		dockerConf{
			Auths: map[string]interface{}{
				"repo.docker.com": "",
			}},
		map[*config.ServerDetails]commons.ProxySettings{
			&config.ServerDetails{
				ServerId: "test1",
			}: {
				ServerName:               "repo.docker.com",
				DockerReverseProxyMethod: "CLOUD",
				UseHttps:                 true,
				HttpsPort:                443,
			},
		},
	)
	require.Equal(t, []repoSelection{
		{
			serverId: "test1",
			repoKey:  "(ALL REPOSITORIES)",
			unknown:  false,
		},
	}, selections)
}

func TestGetDockerReposSelectionRepoPortPerRepo(t *testing.T) {
	selections := getDockerReposSelections(
		dockerConf{
			Auths: map[string]interface{}{
				"repo.docker.com:777": "",
			}},
		map[*config.ServerDetails]commons.ProxySettings{
			&config.ServerDetails{
				ServerId: "test1",
			}: {
				ServerName:               "docker.com",
				DockerReverseProxyMethod: "PORTPERREPO",
				UseHttps:                 true,
				HttpsPort:                443,
				ReverseProxyRepositories: commons.ReverseProxyRepositories{
					ReverseProxyRepoConfigs: []commons.ReverseProxyRepoConfigs{
						{RepoRef: "testrepo", Port: 777, ServerName: "repo.docker.com"},
					}},
			},
		},
	)
	require.Equal(t, []repoSelection{
		{
			serverId: "test1",
			repoKey:  "testrepo",
			unknown:  false,
		},
	}, selections)
}

func TestGetDockerReposSelectionNoMatch(t *testing.T) {
	selections := getDockerReposSelections(
		dockerConf{
			Auths: map[string]interface{}{
				"repo.docker.com:777": "",
			}},
		map[*config.ServerDetails]commons.ProxySettings{
			&config.ServerDetails{
				ServerId: "test1",
			}: {
				ServerName:               "docker.com",
				DockerReverseProxyMethod: "PORTPERREPO",
				UseHttps:                 true,
				HttpsPort:                443,
				ReverseProxyRepositories: commons.ReverseProxyRepositories{
					ReverseProxyRepoConfigs: []commons.ReverseProxyRepoConfigs{
						{RepoRef: "testrepo", Port: 999, ServerName: "repo.docker.com"},
					}},
			},
		},
	)
	require.Equal(t, []repoSelection{
		{
			serverId: "",
			repoKey:  "repo.docker.com:777",
			unknown:  true,
		},
	}, selections)
}
