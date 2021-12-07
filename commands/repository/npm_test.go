package repository

import (
	"context"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/artifactory"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func Test_handleNpmWrongRepoType(t *testing.T) {
	badConfig := SetMeUpConfiguration{
		ServerDetails: serverDetails,
		RepoDetails: &artifactory.RepoDetails{
			PackageType: "maven",
			Key:         testMavenRepoKey,
		},
	}
	err := handleNpm(context.Background(), badConfig)
	require.Error(t, err)
	require.Equal(t, "unexpected repo type. Expected 'npm' but was: 'maven'", err.Error())
}
func Test_handleNpmHomeUndefined(t *testing.T) {
	err := os.Setenv("HOME", "")
	require.NoError(t, err)
	err = handleNpm(context.Background(), getNpmSetMeUpConfig())
	require.Error(t, err)
	require.Equal(t, "$HOME is not defined", err.Error())
}

func Test_handleNpmNoNpmrcFile(t *testing.T) {
	home := setTempHome(t)

	err := handleNpm(context.Background(), getNpmSetMeUpConfig())

	require.NoError(t, err)
	data, err := os.ReadFile(path.Join(home, ".npmrc"))
	require.NoError(t, err)
	expected := `registry=https://maxim.jfrog.io/artifactory/api/npm/default-npm-local/
_auth = .*
always-auth = true
`
	require.Regexp(t, expected, string(data))
}

func Test_handleNpmIdenticalNpmrcFileExists(t *testing.T) {
	home := setTempHome(t)
	err := handleNpm(context.Background(), getNpmSetMeUpConfig())
	require.NoError(t, err)
	statsBefore, err := os.Stat(path.Join(home, ".npmrc"))
	require.NoError(t, err)
	before := statsBefore.ModTime()

	err = handleNpm(context.Background(), getNpmSetMeUpConfig())
	require.NoError(t, err)
	statesAfter, err := os.Stat(path.Join(home, ".npmrc"))
	require.NoError(t, err)
	after := statesAfter.ModTime()

	require.Equalf(t, before, after, "identical content - shouldn't change")
}

func Test_handleNpmDifferentNpmrcFileExists(t *testing.T) {
	home := setTempHome(t)
	npmrcPath := path.Join(home, ".npmrc")
	err := os.WriteFile(npmrcPath, []byte("blah"), 0644)
	require.NoError(t, err)
	statsBefore, err := os.Stat(npmrcPath)
	require.NoError(t, err)
	before := statsBefore.ModTime()

	err = handleNpm(context.Background(), getNpmSetMeUpConfig())
	require.NoError(t, err)
	statsAfter, err := os.Stat(npmrcPath)
	require.NoError(t, err)
	after := statsAfter.ModTime()

	require.Truef(t, before.Before(after), "different content - should change")
	dirEntries, err := os.ReadDir(home)
	require.NoError(t, err)
	require.Equal(t, 2, len(dirEntries)) // npmrc file and bak file
}

func getNpmSetMeUpConfig() SetMeUpConfiguration {
	return SetMeUpConfiguration{
		ServerDetails: serverDetails,
		RepoDetails: &artifactory.RepoDetails{
			PackageType: "npm",
			Key:         testNpmRepoKey,
		},
	}
}

func setTempHome(t *testing.T) string {
	dir, err := ioutil.TempDir(os.TempDir(), "handle_npm")
	require.NoError(t, err)
	err = os.Setenv("HOME", dir)
	require.NoError(t, err)
	return dir
}
