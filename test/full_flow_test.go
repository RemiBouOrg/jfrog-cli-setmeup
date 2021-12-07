package test

import (
	"context"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/environment"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/show"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var serverDetails *config.ServerDetails

func TestMain(m *testing.M) {
	var err error
	serverDetails, err = config.GetDefaultServerConf()
	if err != nil {
		panic(err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestInitShowApplyNpm(t *testing.T) {
	setTempHome(t)
	currDir := cdTempDir(t)
	// default-npm-local is from https://maxim.jfrog.io/
	err := environment.InitEnv(serverDetails, "default-npm-local", "myenv")
	require.NoError(t, err)
	require.FileExists(t, path.Join(currDir, ".jfrog-setmeup.yaml"))

	err = show.ShowCommand(context.Background())
	require.NoError(t, err)

	err = environment.ApplyEnv(context.Background(), "jfrog-set-me-up-test", "myenv")
	require.NoError(t, err)

	err = show.ShowCommand(context.Background())
	require.NoError(t, err)
}

func setTempHome(t *testing.T) string {
	dir, err := ioutil.TempDir(os.TempDir(), "full_flow")
	require.NoError(t, err)
	err = os.Setenv("HOME", dir)
	require.NoError(t, err)
	return dir
}

func cdTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir(os.TempDir(), "jfrogConfig")
	require.NoError(t, err)
	err = os.Chdir(dir)
	require.NoError(t, err)
	return dir
}
