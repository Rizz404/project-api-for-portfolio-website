package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Rizz404/project-api-for-portfolio-website/internal/repository/postgresql"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/repository/postgresql/sqlc"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/rest"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/rest/middleware"
	"github.com/Rizz404/project-api-for-portfolio-website/services/auth"
	"github.com/Rizz404/project-api-for-portfolio-website/services/category"
	"github.com/Rizz404/project-api-for-portfolio-website/services/language"
	"github.com/Rizz404/project-api-for-portfolio-website/services/user"
	"github.com/common-nighthawk/go-figure"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// * Kalo di js namanya figlet
// * Handler untuk root ("/")
func handleRoot(w http.ResponseWriter, r *http.Request) {
	myFigure := figure.NewFigure("API is Running", "standard", true)
	// * Penting untuk set Content-Type ke text/plain agar format ASCII art tidak rusak di browser
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, myFigure.String())
}

// * Handler untuk 404 Not Found
func handleNotFound(w http.ResponseWriter, r *http.Request) {
	myFigure := figure.NewFigure("404 Not Found", "standard", true)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, myFigure.String())
}

func main() {
	// * DATABASE
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	dbConn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("failed to open connection to database: %v", err)
	}
	defer dbConn.Close()

	if err = dbConn.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	log.Printf("successfully connected to database")

	// * DEPENDENCY INJECTION
	queries := sqlc.New(dbConn)

	// * Repo
	authRepo := postgresql.NewUserRepository(queries)
	languageRepo := postgresql.NewLanguageRepository(queries)
	userRepo := postgresql.NewUserRepository(queries)
	categoryRepo := postgresql.NewCategoryRepository(queries)

	// * Service
	authService := auth.NewService(authRepo)
	languageService := language.NewService(languageRepo)
	userService := user.NewService(userRepo)
	categoryService := category.NewService(categoryRepo)

	router := chi.NewRouter()

	router.Use(middleware.Cors)
	router.Use(chiMiddleware.Logger)
	// * Middleware untuk recover dari panic
	router.Use(chiMiddleware.Recoverer)

	router.Get("/", handleRoot)
	router.NotFound(handleNotFound)

	apiRouter := chi.NewRouter()

	// * Handler
	rest.NewAuthHandler(apiRouter, authService)
	rest.NewLanguageHandler(apiRouter, languageService)
	rest.NewUserHandler(apiRouter, userService)
	rest.NewCategoryHandler(apiRouter, categoryService)

	router.Mount("/api/v1", apiRouter)

	// * SERVER
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":5000"
	}

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.Printf("server running on http://localhost%s", addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
