package utils

import (
	"crypto/rand"
	"encoding/base64"
)

type TokenGenerator struct{}

//elithar.githhub.com
func (t TokenGenerator) Generate() (string, error) {
	b := make([]byte, 36)

	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func NewTokenGenerator() TokenGenerator {
	return TokenGenerator{}
}
