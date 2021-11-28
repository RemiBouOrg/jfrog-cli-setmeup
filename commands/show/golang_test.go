package show

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetCurrentGolang(t *testing.T) {
	golang := getCurrentGolang(context.Background())
	require.NotEmpty(t, golang)
}
