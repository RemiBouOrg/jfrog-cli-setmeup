package show

import (
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
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
