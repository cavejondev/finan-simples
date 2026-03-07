package person

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/cavejondev/finan-simples/internal/domain/logger"
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
	logger         *logger.Service
}

// NewService cria nova instancia do servico de person
func NewService(
	r Repository,
	hasher BcryptHasher,
	tokenGenerator TokenGenerator,
	log *logger.Service,
) *Service {
	return &Service{
		repository:     r,
		bcryptHasher:   hasher,
		tokenGenerator: tokenGenerator,
		logger:         log,
	}
}

// Register cria uma nova pessoa com senha criptografada
func (s *Service) Register(ctx context.Context, name, email, password string) error {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if name == "" {
		return ErrNameRequired
	}

	if email == "" {
		return ErrEmailRequired
	}

	if password == "" {
		return ErrPasswordRequired
	}

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

		s.logger.Error(
			ctx,
			"person repository find email error",
			err,
		)

		return err
	}

	if existing != nil {
		return ErrPersonDuplicated
	}

	hash, err := s.bcryptHasher.Hash(password)
	if err != nil {

		s.logger.Error(
			ctx,
			"bcrypt hash error",
			err,
		)

		return ErrPersonInternal
	}

	person := &Person{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  hash,
		CreatedAt: time.Now(),
	}

	if err := s.repository.Create(person); err != nil {

		if errors.Is(err, ErrPersistenceEmailDuplicated) {
			return ErrPersonDuplicated
		}

		s.logger.Error(
			ctx,
			"person repository create error",
			err,
		)

		return err
	}

	return nil
}

// ForgotPassword ajuda o usuario a recuperar a senha
func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	person, err := s.repository.FindByEmail(email)
	if err != nil {

		s.logger.Error(
			ctx,
			"person repository find email error",
			err,
		)

		return err
	}

	if person == nil || person.ID == uuid.Nil {
		return ErrPersonNotFound
	}

	// aqui iria envio de email no futuro

	return nil
}

// Login autentica usuario e retorna token
func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	person, err := s.repository.FindByEmail(email)
	if err != nil {

		s.logger.Error(
			ctx,
			"person repository find email error",
			err,
		)

		return "", ErrInvalidCredentials
	}

	if person == nil || person.ID == uuid.Nil {
		return "", ErrInvalidCredentials
	}

	if err := s.bcryptHasher.Compare(person.Password, password); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.tokenGenerator.Generate(person.ID.String())
	if err != nil {

		s.logger.Error(
			ctx,
			"token generator error",
			err,
		)

		return "", ErrPersonInternal
	}

	return token, nil
}
