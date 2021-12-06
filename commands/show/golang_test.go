package show

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetCurrentGolang(t *testing.T) {
	os.Setenv("GOPROXY", "https://username:fakepassword@fake.url/blah")
	golang := getCurrentGolang(context.Background())
	require.NotEmpty(t, golang)
	require.NotNil(t, golang[0])
	require.Equal(t, "https://fake.url/blah", golang[0].RepoKey)
}
