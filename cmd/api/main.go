// Package main é o pacote principal da aplicação
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/cavejondev/finan-simples/internal/domain/account"
	"github.com/cavejondev/finan-simples/internal/domain/logger"
	"github.com/cavejondev/finan-simples/internal/domain/person"
	"github.com/cavejondev/finan-simples/internal/infrastructure/database"
	accountHttp "github.com/cavejondev/finan-simples/internal/infrastructure/handler/account"
	"github.com/cavejondev/finan-simples/internal/infrastructure/handler/middleware"
	personHttp "github.com/cavejondev/finan-simples/internal/infrastructure/handler/person"
	accountPersistent "github.com/cavejondev/finan-simples/internal/infrastructure/persistence/account"
	loggerPersistent "github.com/cavejondev/finan-simples/internal/infrastructure/persistence/logger"
	personPersistent "github.com/cavejondev/finan-simples/internal/infrastructure/persistence/person"
	"github.com/cavejondev/finan-simples/internal/infrastructure/security"
)

func main() {
	// ARQUIVO .ENV
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// CONECTANDO NAS DATABASES
	db, err := database.NewMainPostgresConnection()
	if err != nil {
		log.Fatal(err)
	}
	dbLog, err := database.NewLogPostgresConnection()
	if err != nil {
		log.Fatal(err)
	}

	// RODANDO MIGRATIONS NAS DATABASES
	database.RunMigrationsMain(db)
	database.RunMigrationsLogs(dbLog)

	// BCRYPT
	hasher := security.NewBcryptHasher()

	// LOG
	logRepo := loggerPersistent.NewRepository(dbLog)
	logService := logger.NewService(logRepo)

	// JWT
	jwtService := security.NewJWTService()

	// PERSON
	repo := personPersistent.NewPersonRepository(db)
	service := person.NewService(repo, hasher, jwtService, logService)
	handler := personHttp.NewHandler(service)

	// A
	repoAccount := accountPersistent.NewAccountRepository(db)
	serviceAccount := account.NewService(repoAccount, logService)
	handlerAccount := accountHttp.NewHandler(serviceAccount)

	// ROUTES
	r := chi.NewRouter()
	r.Use(middleware.RequestMiddleware(logService))

	// REGISTRO ROUTES
	personHttp.RegisterRoutes(r, handler, jwtService)
	accountHttp.RegisterRoutes(r, handlerAccount, jwtService)

	r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	fmt.Println("Servidor na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
