package s1

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	seed = 2
)

func TestSelection(t *testing.T) {
	values := []string{"A", "B", "C", "D"}

	_, err := Select(values, 10, seed)
	require.Error(t, err)

	result, err := Select(values, 2, seed)
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, result, []string{"B", "C"})

}
