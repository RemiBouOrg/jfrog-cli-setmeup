package commands

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"testing"
)

func TestNuget(t *testing.T) {
	testNugetRepoKey := getRepoListFromDefaultServer("nuget")[0].Key
	err := handleNuget(context.Background(),
		SetMeUpConfiguration{
			serverDetails: serverDetails,
			repoDetails: &RepoDetails{
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
	err := handleNuget(context.Background(),
		SetMeUpConfiguration{
			serverDetails: serverDetails,
			repoDetails: &RepoDetails{
				PackageType: "nuget",
				Key:         "do-not-exists",
			},
		},
	)
	require.Error(t, err)
}
