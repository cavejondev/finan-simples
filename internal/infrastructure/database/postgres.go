// Package database é o pacote que contem a definição do banco de dados
package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // lib para o banco de dados se comunicar com application
)

// newPostgresConnection cria conexão genérica
func newPostgresConnection(
	host, port, user, password, dbname string,
) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// NewMainPostgresConnection Banco principal da aplicação
func NewMainPostgresConnection() (*sqlx.DB, error) {
	return newPostgresConnection(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
}

// NewLogPostgresConnection Banco exclusivo de logs
func NewLogPostgresConnection() (*sqlx.DB, error) {
	return newPostgresConnection(
		os.Getenv("LOG_DB_HOST"),
		os.Getenv("LOG_DB_PORT"),
		os.Getenv("LOG_DB_USER"),
		os.Getenv("LOG_DB_PASSWORD"),
		os.Getenv("LOG_DB_NAME"),
	)
}
