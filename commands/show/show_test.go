package show

import (
	"context"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFindServerIdByUrl(t *testing.T) {
	conf, err := config.GetDefaultServerConf()
	require.NoError(t, err)
	serverId := findServerIdByUrl(conf.ArtifactoryUrl + "/my-repo")
	require.Equal(t, conf.ServerId, serverId)
}

func TestFindServerIdByUrlError(t *testing.T) {
	serverId := findServerIdByUrl("http://google.com")
	require.Equal(t, "", serverId)
}

func TestGetShowCommandFlagsAndArgs(t *testing.T) {
	got := GetShowCommand()
	assert.Equal(t, 0, len(got.Flags))
	assert.Equal(t, 0, len(got.Arguments))
}

func Test_showCommand(t *testing.T) {
	err := ShowCommand(context.Background())
	require.NoError(t, err)
}
