package utils

import (
	"github.com/goware/jwtauth"
	"github.com/pkg/errors"
	"os"
	"time"
)

type JWTTokenGenerator struct {
	*jwtauth.JwtAuth
	claims jwtauth.Claims
}

func NewJWTGenerator() *JWTTokenGenerator {
	c := make(jwtauth.Claims, 10)

	c.SetExpiryIn(timeFrame())

	return &JWTTokenGenerator{jwtauth.New("HS256", []byte(os.Getenv("JWT")), nil), c}
}

func (j *JWTTokenGenerator) Claims(c map[string]interface{}) {
	for k, v := range c {
		j.claims.Set(k, v)
	}
}

func (j *JWTTokenGenerator) Generate() (string, error) {

	if len(j.claims) == 0 {
		return "", errors.New("Jwt claims not set yet")
	}

	_, token, err := j.Encode(j.claims)

	if err != nil {
		return "", errors.Wrap(err, "Could not generate JWT token")
	}

	return token, nil
}

func timeFrame() time.Duration {
	return time.Second * 60 * 5
}
