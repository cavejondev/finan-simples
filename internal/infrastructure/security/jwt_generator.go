package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTGenerator é a estrutura do serivco de JWT
type JWTGenerator struct {
	secret []byte
}

// NewJWTGenerator retorna uma nova instancia de JWT
func NewJWTGenerator(secret string) *JWTGenerator {
	return &JWTGenerator{
		secret: []byte(secret),
	}
}

// Generate gera um JWT
func (j *JWTGenerator) Generate(personID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"person_id": personID,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(j.secret)
}
