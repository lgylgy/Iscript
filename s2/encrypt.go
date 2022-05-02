package s2

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

const (
	nonceSize = 12
)

func Encrypt(value string, key string) ([]byte, error) {
	k, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, nonceSize)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	encrypted := gcm.Seal(nil, nonce, []byte(value), nil)
	return append(nonce, encrypted...), nil
}

func Decrypt(value []byte, key string) (string, error) {
	k, err := hex.DecodeString(key)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := value[:nonceSize]
	value = value[nonceSize:]
	text, err := gcm.Open(nil, nonce, value, nil)
	if err != nil {
		return "", err
	}
	return string(text), nil
}
