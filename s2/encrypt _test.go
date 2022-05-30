package s2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	key  = "1234"
	text = "string to "
)

func TestEncryption(t *testing.T) {

	value, err := Encrypt(text, key)
	require.NoError(t, err)
	require.Greater(t, len(value), 0)

	actual, err := Decrypt(value, key)
	require.NoError(t, err)
	require.Equal(t, string(actual), text)
}
