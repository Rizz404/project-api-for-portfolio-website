package middleware

import (
	"net/http"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/web"
)

func AuthorizeRole(allowedRoles ...domain.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// * Jika tidak ada role yang ditentukan, izinkan semua role (selama sudah terotentikasi).
			if len(allowedRoles) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			// * Ambil claims dari context yang sudah di-set oleh middleware Auth.
			claims, ok := web.GetUserFromContext(r.Context())
			if !ok || claims == nil {
				// * Ini seharusnya tidak terjadi jika middleware Auth berjalan dengan benar.
				web.Error(w, http.StatusInternalServerError, "User claims not found in context", nil)
				return
			}

			// * Periksa apakah pengguna memiliki role
			if claims.Role == nil {
				web.Error(w, http.StatusForbidden, "User does not have a role", nil)
				return
			}

			// * Periksa apakah role pengguna ada di dalam daftar yang diizinkan.
			isAllowed := false
			for _, role := range allowedRoles {
				if *claims.Role == role {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				web.Error(w, http.StatusForbidden, "You do not have the required role to access this resource", nil)
				return
			}

			// * Jika role sesuai, lanjutkan ke handler berikutnya.
			next.ServeHTTP(w, r)
		})
	}
}
