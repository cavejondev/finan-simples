package person

import "github.com/google/uuid"

// Repository é o repositorio de pessoa
type Repository interface {
	Create(person *Person) error
	FindByEmail(email string) (*Person, error)
	FindByID(id uuid.UUID) (*Person, error)
}

// JWTService é a interface do serviço que gera JWT
type JWTService interface {
	Generate(personID uuid.UUID) (string, error)
}

// BcryptHasher define contrato para criptografia de senha
type BcryptHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}
