package person

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
)

// Erros de servico
var (
	ErrNameTooShort       = errors.New("name too short")
	ErrEmailTooShort      = errors.New("email too short")
	ErrEmailInvalid       = errors.New("email invalid")
	ErrNameRequired       = errors.New("name is required")
	ErrEmailRequired      = errors.New("email is required")
	ErrPasswordRequired   = errors.New("password is required")
	ErrPasswordTooShort   = errors.New("password must be at least 6 characters")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPersonNotFound     = errors.New("person not found")
	ErrPersonDuplicated   = errors.New("person duplicated")
	ErrPersonInternal     = errors.New("person internal error")
)

// Erros do banco
var (
	ErrPersistenceEmailDuplicated = errors.New("person email already exists")
)

// Service representa o servico de person
type Service struct {
	repository     Repository
	bcryptHasher   BcryptHasher
	tokenGenerator TokenGenerator
}

// NewService cria nova instancia do servico de person
func NewService(
	r Repository,
	hasher BcryptHasher,
	tokenGenerator TokenGenerator,
) *Service {
	return &Service{
		repository:     r,
		bcryptHasher:   hasher,
		tokenGenerator: tokenGenerator,
	}
}

// Register cria uma nova pessoa com senha criptografada
func (s *Service) Register(name, email, password string) error {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if len(name) < 3 {
		return ErrNameTooShort
	}

	if len(email) < 5 {
		return ErrEmailTooShort
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return ErrEmailInvalid
	}

	if len(password) < 6 {
		return ErrPasswordTooShort
	}

	existing, err := s.repository.FindByEmail(email)
	if err != nil {
		return err
	}

	if existing != nil && existing.ID > 0 {
		return ErrPersonDuplicated
	}

	hash, err := s.bcryptHasher.Hash(password)
	if err != nil {
		return ErrPersonInternal
	}

	person := &Person{
		Name:      name,
		Email:     email,
		Password:  hash,
		CreatedAt: time.Now(),
	}

	if err := s.repository.Create(person); err != nil {
		if errors.Is(err, ErrPersistenceEmailDuplicated) {
			return ErrPersonDuplicated
		}
		return ErrPersonInternal
	}

	return nil
}

// ForgotPassword e a funcao que ajuda o usuario a recupera a senha
func (s *Service) ForgotPassword(email string) error {
	person, err := s.repository.FindByEmail(email)
	if err != nil {
		return err
	}

	if person == nil || person.ID <= 0 {
		return ErrPersonNotFound
	}

	// Enviar email para recuperar a senha

	return nil
}

// Login autentica usuario e retorna token
func (s *Service) Login(email, password string) (string, error) {
	person, err := s.repository.FindByEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if person == nil || person.ID <= 0 {
		return "", ErrInvalidCredentials
	}

	if err := s.bcryptHasher.Compare(person.Password, password); err != nil {
		return "", ErrInvalidCredentials
	}

	return s.tokenGenerator.Generate(fmt.Sprintf("%d", person.ID))
}
