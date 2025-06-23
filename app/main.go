package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Rizz404/project-api-for-portfolio-website/category"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/repository/postgresql"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/repository/postgresql/sqlc"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/rest"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
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
	categoryRepo := postgresql.NewCategoryRepository(queries)
	categoryService := category.NewService(categoryRepo)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	// * Middleware untuk recover dari panic
	router.Use(middleware.Recoverer)

	apiRouter := chi.NewRouter()
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
