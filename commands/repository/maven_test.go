package repository

import (
	"context"
	"fmt"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/artifactory"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestMavenSetMeUpUndefinedHome(t *testing.T) {
	err := os.Setenv("HOME", "")
	require.NoError(t, err)

	err = handleMaven(context.Background(), getMavenConfig())

	require.Error(t, err)
}

func TestMavenSetMeUpFreshM2(t *testing.T) {
	mavenDir := createTempDotM2(t)

	err := handleMaven(context.Background(), getMavenConfig())

	require.NoError(t, err)
	require.FileExists(t, fmt.Sprintf("%s/settings.xml", mavenDir))
}

func TestMavenSetMeUpSettingsFileExists(t *testing.T) {
	mavenDir := createTempDotM2(t)

	err := handleMaven(context.Background(), getMavenConfig())
	require.NoError(t, err)
	statsBefore, err := os.Stat(fmt.Sprintf("%s/settings.xml", mavenDir))
	require.NoError(t, err)

	err = handleMaven(context.Background(), getMavenConfig())
	require.NoError(t, err)
	statsAfter, err := os.Stat(fmt.Sprintf("%s/settings.xml", mavenDir))
	require.NoError(t, err)

	require.Truef(t, statsAfter.ModTime().After(statsBefore.ModTime()), "new file - different time stamp")
}

func getMavenConfig() SetMeUpConfiguration {
	return SetMeUpConfiguration{
		ServerDetails: serverDetails,
		RepoDetails: &artifactory.RepoDetails{
			PackageType: "maven",
			Key:         testMavenRepoKey,
		},
	}
}

func createTempDotM2(t *testing.T) string {
	dir, err := ioutil.TempDir(os.TempDir(), "maven")
	require.NoError(t, err)
	err = os.Setenv("HOME", dir)
	require.NoError(t, err)
	mavenDir := fmt.Sprintf("%s/.m2/", dir)
	err = os.Mkdir(mavenDir, 0777)
	require.NoError(t, err)
	return mavenDir
}
