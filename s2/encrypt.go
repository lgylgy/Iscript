package s2

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

const (
	nonceSize = 12
)

func createHash(key string) (string, error) {
	hasher := md5.New()
	_, err := hasher.Write([]byte(key))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func Encrypt(value string, key string) ([]byte, error) {
	k, err := createHash(key)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher([]byte(k))
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

func Decrypt(value []byte, key string) ([]byte, error) {
	k, err := createHash(key)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher([]byte(k))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := value[:nonceSize]
	value = value[nonceSize:]
	text, err := gcm.Open(nil, nonce, value, nil)
	if err != nil {
		return nil, err
	}
	return text, nil
}
