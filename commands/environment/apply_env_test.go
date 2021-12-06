package environment

import (
	"context"
	"github.com/jfrog/jfrog-cli-plugin-template/jfrogconfig"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestApplyEnv(t *testing.T) {
	home := setTempHome(t)
	_ = cdTempDir(t)
	writeConf(t)
	require.NoFileExists(t, path.Join(home, ".npmrc"))

	err := ApplyEnv(context.Background(), "jfrog-set-me-up-test", "env-a")

	require.NoError(t, err)
	require.FileExists(t, path.Join(home, ".npmrc"))
}

func setTempHome(t *testing.T) string {
	dir, err := ioutil.TempDir(os.TempDir(), "handle_npm")
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

func writeConf(t *testing.T) {
	conf := `
{
	"env-a": {
		"npm": "npm-local"
	}
}
`
	err := jfrogconfig.WriteConfigFile([]byte(conf))
	require.NoError(t, err)
}
