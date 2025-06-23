package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	defaultAddr = ":5000"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// * Database
	dbUrl := os.Getenv("DATABASE_URL")

	dbConn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("failed to open connection to database", err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal("failed to ping database", err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal("failed to close database", err)
		}
	}()

	// * Server
	addr := os.Getenv("ADDR")
	server := http.Server{Addr: addr}

	log.Printf("connected to database at %s", dbUrl)
	log.Printf("server running on http://localhost%s", addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("server error")
	}

}
