package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

func runMigrations(db *sqlx.DB, dir string) {
	goose.SetDialect("postgres")

	if err := goose.Up(db.DB, dir); err != nil {
		log.Fatal("Erro ao rodar migrations:", err)
	}

	log.Println("Migrations executadas com sucesso")
}

// RunMigrationsMain Banco principal
func RunMigrationsMain(db *sqlx.DB) {
	dir := "internal/infrastructure/database/migrations/main"
	runMigrations(db, dir)
}

// RunMigrationsLogs Banco de logs
func RunMigrationsLogs(db *sqlx.DB) {
	dir := "internal/infrastructure/database/migrations/logs"
	runMigrations(db, dir)
}
