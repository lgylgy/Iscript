package s3

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	text = "string to "
)

func TestStega(t *testing.T) {
	err := Encrypt(filepath.Join("testdata", "image.jpg"),
		filepath.Join("testdata", "output.jpg"), []byte(text))
	require.NoError(t, err)

	message, err := Decrypt(filepath.Join("testdata", "output.jpg"))
	require.NoError(t, err)
	require.Equal(t, string(message), text)
}
