package security

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims representa os claims do JWT
type Claims struct {
	PersonID string `json:"sub"`
	jwt.RegisteredClaims
}

// JWTService é o serviço de JWT, no estilo "service man"
type JWTService struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

// NewJWTService cria um novo JWTService carregando configs do ambiente
func NewJWTService() *JWTService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "finan-simples"
	}

	hours := 24
	if h := os.Getenv("JWT_TTL_HOURS"); h != "" {
		if v, err := strconv.Atoi(h); err == nil {
			hours = v
		}
	}

	return &JWTService{
		secret: []byte(secret),
		issuer: issuer,
		ttl:    time.Duration(hours) * time.Hour,
	}
}

// Generate gera um JWT assinado
func (j *JWTService) Generate(personID uuid.UUID) (string, error) {
	now := time.Now()
	claims := Claims{
		PersonID: personID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   personID.String(),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// Validate valida o JWT e retorna o ID da pessoa
func (j *JWTService) Validate(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	personID, err := uuid.Parse(claims.PersonID)
	if err != nil {
		return uuid.Nil, err
	}

	return personID, nil
}
