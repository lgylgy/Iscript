package s3

import (
	"bufio"
	"bytes"
	"image"
	_ "image/jpeg"
	"os"

	"github.com/auyer/steganography"
)

func Encrypt(input, output string, message []byte) error {
	file, err := os.Open(input)
	if err != nil {
		return err
	}
	defer file.Close()
	img, _, err := image.Decode(bufio.NewReader(file))
	if err != nil {
		return err
	}
	encodedImg := new(bytes.Buffer)
	err = steganography.Encode(encodedImg, img, message)
	if err != nil {
		return err
	}
	outFile, err := os.Create(output)
	if err != nil {
		return err
	}
	_, err = bufio.NewWriter(outFile).Write(encodedImg.Bytes())
	return err
}

func Decrypt(input string) ([]byte, error) {
	file, err := os.Open(input)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(bufio.NewReader(file))
	if err != nil {
		return nil, err
	}
	msg := steganography.Decode(steganography.GetMessageSizeFromImage(img), img)
	return msg, nil
}
