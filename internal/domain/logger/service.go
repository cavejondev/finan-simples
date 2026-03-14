package logger

import (
	"context"
	"encoding/json"
	"runtime"
	"time"

	contextutil "github.com/cavejondev/finan-simples/internal/domain/util"
	"github.com/google/uuid"
)

// Service representa o servico de logger
type Service struct {
	repo Repository
}

// NewService cria nova instancia do servico de logger
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Info registra um log de nível INFO
func (s *Service) Info(ctx context.Context, message string) {
	s.save(ctx, LevelInfo, message, nil)
}

// Warn registra um log de nível WARN
func (s *Service) Warn(ctx context.Context, message string) {
	s.save(ctx, LevelWarn, message, nil)
}

// Registra um log de nível ERROR com erro opcional
func (s *Service) Error(ctx context.Context, message string, err error) {
	var errMsg *string

	if err != nil {
		e := err.Error()
		errMsg = &e
	}

	s.save(ctx, LevelError, message, errMsg)
}

// Debug registra um log de nível DEBUG
func (s *Service) Debug(ctx context.Context, message string) {
	s.save(ctx, LevelDebug, message, nil)
}

func (s *Service) save(
	ctx context.Context,
	level Level,
	message string,
	errorMsg *string,
) {
	fn, file, line := caller()

	meta := map[string]any{
		"file": file,
		"line": line,
	}

	metaJSON, _ := json.Marshal(meta)

	log := &Log{
		ID:        uuid.New(),
		Level:     level,
		Message:   message,
		Service:   &fn,
		RequestID: contextutil.GetRequestID(ctx),
		UserID:    contextutil.GetUserID(ctx),
		Method:    contextutil.GetMethod(ctx),
		Path:      contextutil.GetPath(ctx),
		Error:     errorMsg,
		Metadata:  metaJSON,
		CreatedAt: time.Now(),
	}

	/*var err string
	if errorMsg != nil {
		err = *errorMsg
	}

	fmt.Printf(
		"[LOG]\n level=%s\n service=%s\n message=%s\n error=%s\n file=%s\n line=%d\n",
		level,
		fn,
		message,
		err,
		file,
		line,
	)*/

	_ = s.repo.Create(log)
}

func caller() (string, string, int) {
	pc, file, line, ok := runtime.Caller(3)

	if !ok {
		return "unknown", "unknown", 0
	}

	fn := runtime.FuncForPC(pc)

	return fn.Name(), file, line
}
