// Package security implementa serviços relacionados à segurança, como criptografia e verificação de senhas.
package security

import "golang.org/x/crypto/bcrypt"

// BcryptHasher é a implementação de hash de senha utilizando o algoritmo bcrypt.
type BcryptHasher struct{}

// NewBcryptHasher cria uma nova instância de BcryptHasher.
func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

// Hash recebe uma senha em texto puro e retorna versão criptografada utilizando bcrypt.
func (b *BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// Compare verifica se a senha informada corresponde
func (b *BcryptHasher) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}
