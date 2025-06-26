package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/Rizz404/project-api-for-portfolio-website/internal"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/utils"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/web"
)

// * Definisikan custom key untuk context agar tidak ada bentrokan.

var accessTokenSecret = []byte(os.Getenv("JWT_ACCESS_SECRET"))

// * Auth adalah middleware untuk memverifikasi JWT dari Authorization header.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			web.Error(w, http.StatusUnauthorized, "Authorization header is required", nil)
			return
		}

		// * Header harus dalam format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			web.Error(w, http.StatusUnauthorized, "Authorization header format must be Bearer {token}", nil)
			return
		}

		tokenString := parts[1]

		// * Validasi token
		claims, err := utils.ValidateToken(tokenString, accessTokenSecret)
		if err != nil {
			web.Error(w, http.StatusUnauthorized, "Invalid or expired token", err.Error())
			return
		}

		// * Jika token valid, simpan claims di dalam context request
		ctx := context.WithValue(r.Context(), internal.UserClaimsKey, claims)

		// * Lanjutkan ke handler berikutnya dengan context yang sudah dimodifikasi
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
