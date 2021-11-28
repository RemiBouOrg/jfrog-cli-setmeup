package show

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetCurrentNuget(t *testing.T) {
	repos := getCurrentNuget(context.Background())
	require.NotEmpty(t, repos)
}
