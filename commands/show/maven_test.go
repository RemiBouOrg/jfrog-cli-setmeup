package show

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetCurrentMaven(t *testing.T) {
	serverId, repoKey := getCurrentMaven(context.Background())
	require.NotEmpty(t, serverId)
	require.NotEmpty(t, repoKey)
}
