package jfrogconfig

import (
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
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

func TestFindRepoKeys(t *testing.T) {
	got, err := FindRepoKeys(&components.Context{Arguments: []string{"maven-local"}}, nil, serverDetails)
	require.NoError(t, err)
	require.Equal(t, []string{"maven-local"}, got)
}

func TestReadWriteConfigFile(t *testing.T) {
	cdTempDir(t)
	writeConf(t)
	got, err := ReadCurrentConfFile()
	require.NoError(t, err)
	require.Equal(t, JFrogConfFile{"env-a": RepoTypeToName{"b": "c"}}, got)
}

func Test_extractValues(t *testing.T) {
	got := extractValues(&RepoTypeToName{"env-a": "a1", "b": "b1"})
	require.ElementsMatch(t, []string{"a1", "b1"}, got)
}

func Test_findRepoKeyFromConfFile(t *testing.T) {
	cdTempDir(t)
	writeConf(t)
	got, err := findRepoKeyFromConfFile("env-a")
	require.NoError(t, err)
	require.Equal(t, &RepoTypeToName{"b": "c"}, got)
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
		"b": "c"
	}
}
`
	err := WriteConfigFile([]byte(conf))
	require.NoError(t, err)
}
