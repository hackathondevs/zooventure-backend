package pasetok

import (
	"os"
	"time"

	"aidanwoods.dev/go-paseto"
)

func Encode(token paseto.Token) (string, error) {
	key, err := paseto.NewV4AsymmetricSecretKeyFromHex(os.Getenv("PASETO_SECRET_HEX"))
	if err != nil {
		return "", err
	}
	signed := token.V4Sign(key, nil)
	return signed, nil
}

func Decode(signed string) (paseto.Token, error) {
	key, err := paseto.NewV4AsymmetricSecretKeyFromHex(os.Getenv("PASETO_SECRET_HEX"))
	if err != nil {
		return paseto.Token{}, err
	}
	parser := paseto.NewParser()
	parser.AddRule(paseto.ForAudience("*"))
	parser.AddRule(paseto.IssuedBy(os.Getenv("APP_HOST")))
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))
	parsedToken, err := parser.ParseV4Public(key.Public(), signed, nil)
	if err != nil {
		return paseto.Token{}, err
	}
	return *parsedToken, nil
}
