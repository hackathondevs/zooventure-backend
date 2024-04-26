package helper

import "golang.org/x/crypto/bcrypt"

func BcryptHash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func BcryptCompare(hashed, raw string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw)); err != nil {
		return err
	}

	return nil
}
