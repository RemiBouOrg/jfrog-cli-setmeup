package show

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetCurrentMaven(t *testing.T) {
	repos := getCurrentMaven(context.Background())
	require.NotEmpty(t, repos)
}
