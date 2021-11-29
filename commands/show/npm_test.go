package show

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetCurrentNpm(t *testing.T) {
	npm := getCurrentNpm(context.Background())
	require.NotEmpty(t, npm)
}
