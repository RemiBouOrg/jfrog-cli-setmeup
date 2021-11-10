package commands

import (
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"testing"
)

func TestNuget(t *testing.T) {
	testNugetRepoKey := getRepoListFromDefaultServer("nuget")[0].Key
	err := handleNuget(
		SetMeUpConfiguration{
			repositoryKey: testNugetRepoKey,
			serverDetails: serverDetails,
			repoDetails: RepoDetails{
				PackageType: "nuget",
				Key:         testNugetRepoKey,
			},
		},
	)
	require.NoError(t, err)
	command := exec.Command("nuget", "search",
		"-Source", "Artifactory",
		"BrowserInterop",
	)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err = command.Run()
	require.NoError(t, err)
}

func TestNugetErr(t *testing.T) {
	err := handleNuget(
		SetMeUpConfiguration{
			repositoryKey: "do-not-exists",
			serverDetails: serverDetails,
			repoDetails: RepoDetails{
				PackageType: "nuget",
				Key:         "do-not-exists",
			},
		},
	)
	require.Error(t, err)
}
