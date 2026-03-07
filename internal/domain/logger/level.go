// Package logger é o pacote que loga coisas da aplicacao
package logger

// Level é o tipo do level
type Level string

// Levels para serem usados
const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
	LevelFatal Level = "FATAL"
)
