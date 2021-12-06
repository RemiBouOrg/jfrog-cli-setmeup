package environment

import (
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/stretchr/testify/require"
	"path"
	"testing"
)

func TestInitEnv(t *testing.T) {
	currDir := cdTempDir(t)
	serverDetails, err := config.GetDefaultServerConf()
	require.NoError(t, err)
	// default-npm-local is from https://maxim.jfrog.io/
	err = InitEnv(serverDetails, "default-npm-local", "myenv")

	require.NoError(t, err)
	require.FileExists(t, path.Join(currDir, ".jfrog-setmeup.yaml"))

}
