package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsSupportedCurrency(t *testing.T) {
	res := IsSupportedCurrency(USD)
	require.True(t, res)

	res1 := IsSupportedCurrency("XYZ")
	require.False(t, res1)
}
