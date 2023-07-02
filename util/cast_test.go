package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToString(t *testing.T) {
	t.Parallel()
	require.Equal(t, "123456789", ToString(int(123456789)))
	require.Equal(t, "-123456789012345678", ToString(int64(-123456789012345678)))
	require.Equal(t, "123456789012345678", ToString(uint64(123456789012345678)))
	require.Equal(t, "true", ToString(true))
	require.Equal(t, "false", ToString(false))
	require.Equal(t, "0.123000", ToString(float32(0.123)))
	require.Equal(t, "[1 2 3]", ToString([]int{1, 2, 3}))
	require.Equal(t, "1000000000.000000", ToString(float64(1000000000.0)))
}
