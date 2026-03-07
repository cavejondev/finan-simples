package logger

// Repository é a interface do repositorio do logger
type Repository interface {
	Create(log *Log) error
}
