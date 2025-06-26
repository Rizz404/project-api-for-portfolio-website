package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

func Cors(next http.Handler) http.Handler {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // * Waktu maksimal (detik) cache preflight request
	})

	return corsMiddleware.Handler(next)
}
