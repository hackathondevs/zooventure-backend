package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

func Encrypt(str string) (string, error) {
	block, err := aes.NewCipher([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	plain := []byte(str)
	cipherTxt := make([]byte, aes.BlockSize+len(plain))
	iv := cipherTxt[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherTxt[aes.BlockSize:], plain)
	return base64.RawStdEncoding.EncodeToString(cipherTxt), nil
}

func Decrypt(secure string) (string, error) {
	cipherTxt, err := base64.RawStdEncoding.DecodeString(secure)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	if len(cipherTxt) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return "", err
	}
	iv := cipherTxt[:aes.BlockSize]
	cipherTxt = cipherTxt[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherTxt, cipherTxt)
	return string(cipherTxt), nil
}
