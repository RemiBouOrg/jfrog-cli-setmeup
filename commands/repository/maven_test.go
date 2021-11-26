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

func TestMavenSetMeUp(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "maven")
	require.NoError(t, err)
	mavenDir := fmt.Sprintf("%s/.m2/", dir)
	err = os.Mkdir(mavenDir, 0777)
	require.NoError(t, err)
	_ = os.Setenv("HOME", dir)
	err = handleMaven(context.Background(),
		SetMeUpConfiguration{
			ServerDetails: serverDetails,
			RepoDetails: &artifactory.RepoDetails{
				PackageType: "maven",
				Key:         testMavenRepoKey,
			},
		},
	)
	require.NoError(t, err)
	_, err = os.Stat(fmt.Sprintf("%s/.m2/settings.xml", dir))
	require.NoError(t, err)
}
