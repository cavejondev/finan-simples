// Package main é o pacote principal da aplicação
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/cavejondev/finan-simples/internal/domain/person"
	"github.com/cavejondev/finan-simples/internal/infrastructure/database"
	personHttp "github.com/cavejondev/finan-simples/internal/infrastructure/handler/person"
	personPersistent "github.com/cavejondev/finan-simples/internal/infrastructure/persistence/person"
	"github.com/cavejondev/finan-simples/internal/infrastructure/security"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// ========================
	// Infra
	// ========================

	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatal(err)
	}

	database.RunMigrationsMaster(db)

	repo := personPersistent.NewPersonRepository(db)
	hasher := security.NewBcryptHasher()

	// JWT secret via env
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	tokenGenerator := security.NewJWTGenerator(secret)

	// ========================
	// Domain Service
	// ========================

	service := person.NewService(repo, hasher, tokenGenerator)

	// ========================
	// HTTP Handler
	// ========================

	handler := personHttp.NewHandler(service)

	r := chi.NewRouter()

	personHttp.RegisterRoutes(r, handler)

	r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	fmt.Println("Servidor na porta 8080 🚀")
	log.Fatal(http.ListenAndServe(":8080", r))
}
