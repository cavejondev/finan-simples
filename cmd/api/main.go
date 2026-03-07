// Package main é o pacote principal da aplicação
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/cavejondev/finan-simples/internal/domain/logger"
	"github.com/cavejondev/finan-simples/internal/domain/person"
	"github.com/cavejondev/finan-simples/internal/infrastructure/database"
	"github.com/cavejondev/finan-simples/internal/infrastructure/handler/middleware"
	personHttp "github.com/cavejondev/finan-simples/internal/infrastructure/handler/person"
	loggerPersistent "github.com/cavejondev/finan-simples/internal/infrastructure/persistence/logger"
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

	db, err := database.NewMainPostgresConnection()
	if err != nil {
		log.Fatal(err)
	}
	dbLog, err := database.NewLogPostgresConnection()
	if err != nil {
		log.Fatal(err)
	}

	database.RunMigrationsMain(db)
	database.RunMigrationsLogs(dbLog)

	repo := personPersistent.NewPersonRepository(db)
	hasher := security.NewBcryptHasher()

	// JWT secret via env
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	tokenGenerator := security.NewJWTGenerator(secret)

	logRepo := loggerPersistent.NewRepository(dbLog)
	logService := logger.NewService(logRepo)

	// ========================
	// Domain Service
	// ========================

	service := person.NewService(repo, hasher, tokenGenerator, logService)

	// ========================
	// HTTP Handler
	// ========================

	handler := personHttp.NewHandler(service)

	r := chi.NewRouter()

	r.Use(middleware.RequestMiddleware(logService))

	personHttp.RegisterRoutes(r, handler)

	r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	fmt.Println("Servidor na porta 8080 🚀")
	log.Fatal(http.ListenAndServe(":8080", r))
}
