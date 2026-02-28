package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

// RunMigrationsMaster roda todas as migrations
func RunMigrationsMaster(db *sqlx.DB) {
	migrationsDir := "internal/infrastructure/database/migrations"

	goose.SetDialect("postgres")

	if err := goose.Up(db.DB, migrationsDir); err != nil {
		log.Fatal("Erro ao rodar migrations:", err)
	}

	log.Println("Migrations executadas com sucesso 🚀")
}
