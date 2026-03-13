// Package main é o pacote principal da aplicação
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	// DOMAIN
	"github.com/cavejondev/finan-simples/internal/domain/account"
	"github.com/cavejondev/finan-simples/internal/domain/category"
	"github.com/cavejondev/finan-simples/internal/domain/logger"
	"github.com/cavejondev/finan-simples/internal/domain/person"
	"github.com/cavejondev/finan-simples/internal/domain/subcategory"

	// DATABASE
	"github.com/cavejondev/finan-simples/internal/infrastructure/database"

	// HANDLERS
	accountHttp "github.com/cavejondev/finan-simples/internal/infrastructure/handler/account"
	categoryHttp "github.com/cavejondev/finan-simples/internal/infrastructure/handler/category"
	personHttp "github.com/cavejondev/finan-simples/internal/infrastructure/handler/person"
	subcategoryHttp "github.com/cavejondev/finan-simples/internal/infrastructure/handler/subcategory"

	"github.com/cavejondev/finan-simples/internal/infrastructure/handler/middleware"

	// PERSISTENCE
	accountPersistent "github.com/cavejondev/finan-simples/internal/infrastructure/persistence/account"
	categoryPersistent "github.com/cavejondev/finan-simples/internal/infrastructure/persistence/category"
	loggerPersistent "github.com/cavejondev/finan-simples/internal/infrastructure/persistence/logger"
	personPersistent "github.com/cavejondev/finan-simples/internal/infrastructure/persistence/person"
	subcategoryPersistent "github.com/cavejondev/finan-simples/internal/infrastructure/persistence/subcategory"

	// SECURITY
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

	// RODANDO MIGRATIONS
	database.RunMigrationsMain(db)
	database.RunMigrationsLogs(dbLog)

	// BCRYPT
	hasher := security.NewBcryptHasher()

	// LOGGER
	logRepo := loggerPersistent.NewRepository(dbLog)
	logService := logger.NewService(logRepo)

	// JWT
	jwtService := security.NewJWTService()

	// PERSON
	personRepo := personPersistent.NewPersonRepository(db)
	personService := person.NewService(personRepo, hasher, jwtService, logService)
	personHandler := personHttp.NewHandler(personService)

	// ACCOUNT
	accountRepo := accountPersistent.NewAccountRepository(db)
	accountService := account.NewService(accountRepo, logService)
	accountHandler := accountHttp.NewHandler(accountService)

	// CATEGORY
	categoryRepo := categoryPersistent.NewCategoryRepository(db)
	categoryService := category.NewService(categoryRepo, logService)
	categoryHandler := categoryHttp.NewHandler(categoryService)

	// SUBCATEGORY
	subcategoryRepo := subcategoryPersistent.NewSubcategoryRepository(db)
	subcategoryService := subcategory.NewService(subcategoryRepo, categoryService, logService)
	subcategoryHandler := subcategoryHttp.NewHandler(subcategoryService)

	// ROUTER
	r := chi.NewRouter()

	r.Use(middleware.RequestMiddleware(logService))

	// REGISTRO ROUTES
	personHttp.RegisterRoutes(r, personHandler, jwtService)
	accountHttp.RegisterRoutes(r, accountHandler, jwtService)
	categoryHttp.RegisterRoutes(r, categoryHandler, jwtService)
	subcategoryHttp.RegisterRoutes(r, subcategoryHandler, jwtService)

	r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	fmt.Println("Servidor na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
