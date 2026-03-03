package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims representa os claims do JWT
type Claims struct {
	PersonID string `json:"sub"`
	jwt.RegisteredClaims
}

// JWTGenerator é a estrutura do serviço de JWT
type JWTGenerator struct {
	secret []byte
	ttl    time.Duration
	issuer string
}

// NewJWTGenerator cria uma nova instancia do gerador de JWT
func NewJWTGenerator(secret string) *JWTGenerator {
	return &JWTGenerator{
		secret: []byte(secret),
		ttl:    24 * time.Hour,
		issuer: "finan-simples",
	}
}

// Generate gera um JWT assinado
func (j *JWTGenerator) Generate(personID string) (string, error) {
	now := time.Now()

	claims := Claims{
		PersonID: personID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   personID,
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(j.secret)
}
