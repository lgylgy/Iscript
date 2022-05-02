package s2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	key  = "6368616e676520746869732070617373776f726420746f206120736563726574"
	text = "string to "
)

func TestEncryption(t *testing.T) {

	value, err := Encrypt(text, key)
	require.NoError(t, err)
	require.Greater(t, len(value), 0)

	actual, err := Decrypt(value, key)
	require.NoError(t, err)
	require.Equal(t, actual, text)
}
