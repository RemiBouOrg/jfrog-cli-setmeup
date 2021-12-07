package environment

import (
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/artifactory"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/commons"
	"github.com/jfrog/jfrog-cli-plugin-template/jfrogconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
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

func TestGetEnvInitCommandFlags(t *testing.T) {
	got := GetEnvInitCommand(nil)
	var flagNames []string
	for _, flag := range got.Flags {
		flagNames = append(flagNames, flag.GetName())
	}
	assert.ElementsMatch(t, []string{commons.ServerIdFlag, commons.EnvNameFlag}, flagNames)
}

type dummyFindRepo struct {
}

func (r dummyFindRepo) FindRepo(serverDetails *config.ServerDetails) (*artifactory.RepoDetails, error) {
	return &artifactory.RepoDetails{Key: "default-npm-local", PackageType: "npm"}, nil
}

func TestGetEnvInitCommandFuncNoArgs(t *testing.T) {
	findRepoService := &dummyFindRepo{}
	got := GetEnvInitCommand(findRepoService)

	err := got.Action(&components.Context{})
	require.NoError(t, err)
	require.FileExists(t, jfrogconfig.JfrogConfFilePath)
	_ = os.Remove(jfrogconfig.JfrogConfFilePath)
}
