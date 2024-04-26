package helper

import (
	"crypto/rand"
	mathRand "math/rand"
)

var (
	chars       = "0123456789"
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

func RandNumber(n int) (string, error) {
	buffer := make([]byte, n)

	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	charsLen := len(chars)
	for i := 0; i < n; i++ {
		buffer[i] = chars[int(buffer[i])%charsLen]
	}

	return string(buffer), nil
}

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[mathRand.Intn(len(letterRunes))]
	}

	return string(b)
}
