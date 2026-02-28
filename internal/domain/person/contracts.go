package person

// Repository é o repositorio de pessoa
type Repository interface {
	Create(person *Person) error
	FindByEmail(email string) (*Person, error)
}

// TokenGenerator é o servico que gera um JWT
type TokenGenerator interface {
	Generate(personID string) (string, error)
}

// BcryptHasher define contrato para criptografia de senha
type BcryptHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}
