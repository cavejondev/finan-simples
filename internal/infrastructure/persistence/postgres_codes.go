// Package persistence é o pacote onde tem as operações do banco de dados
package persistence

// Constantes de erros do postgres
const (
	PgUniqueViolation     = "23505"
	PgForeignKeyViolation = "23503"
	PgNotNullViolation    = "23502"
)
