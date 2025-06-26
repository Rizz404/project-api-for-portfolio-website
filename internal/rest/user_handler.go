package rest

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/rest/middleware"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/web"
	"github.com/go-chi/chi/v5"
)

type UserService interface {
	// * CREATE
	CreateUser(ctx context.Context, payload *domain.CreateUserPayload) (domain.User, error)

	// * READ (MANY)
	GetUsersPaginated(ctx context.Context, limit int32, offset int32) ([]domain.User, error)
	GetUsersCursorFirst(ctx context.Context, limit int32) ([]domain.User, error)
	SearchUsers(ctx context.Context, searchTerm string) ([]domain.User, error)
	SearchUsersPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.User, error)
	SearchUsersByUsername(ctx context.Context, username string) ([]domain.User, error)
	SearchUsersByEmail(ctx context.Context, email string) ([]domain.User, error)
	SearchUsersByFullName(ctx context.Context, fullName string) ([]domain.User, error)

	// * READ (ONE & UTILITY)
	GetUser(ctx context.Context, id string) (domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	CheckUserExists(ctx context.Context, id string) (bool, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
	CountUsers(ctx context.Context) (int64, error)

	// * UPDATE
	UpdateUser(ctx context.Context, payload *domain.User) (domain.User, error)
	UpdateUserProfile(ctx context.Context, payload *domain.User) (domain.User, error)

	// * DELETE
	DeleteUser(ctx context.Context, id string) error
}

type UserHandler struct {
	Service UserService
}

func NewUserHandler(r chi.Router, s UserService) {
	handler := &UserHandler{
		Service: s,
	}

	r.Route("/users", func(r chi.Router) {
		// * CREATE
		r.With(middleware.Auth).Post("/", handler.CreateUser)

		// * READ (MANY)
		r.Get("/", handler.GetUsersPaginated)
		r.Get("/cursor", handler.GetUsersCursorFirst)
		r.Get("/search", handler.SearchUsers)
		r.Get("/search/paginated", handler.SearchUsersPaginated)
		r.Get("/search/username", handler.SearchUsersByUsername)
		r.Get("/search/email", handler.SearchUsersByEmail)
		r.Get("/search/fullname", handler.SearchUsersByFullName)

		// * READ (ONE & UTILITY)
		r.Get("/{id}", handler.GetUser)
		r.Get("/username/{username}", handler.GetUserByUsername)
		r.Get("/email/{email}", handler.GetUserByEmail)
		r.Get("/exists/{id}", handler.CheckUserExists)
		r.Get("/exists/email/{email}", handler.CheckEmailExists)
		r.Get("/exists/username/{username}", handler.CheckUsernameExists)
		r.Get("/count", handler.CountUsers)

		// * UPDATE
		r.With(middleware.Auth).Patch("/{id}", handler.UpdateUser)
		r.With(middleware.Auth).Patch("/{id}/profile", handler.UpdateUserProfile)

		// * DELETE
		r.With(middleware.Auth).Delete("/{id}", handler.DeleteUser)

		// * CURRENT
		r.Route("/current", func(r chi.Router) {
			r.Use(middleware.Auth)

			r.Get("/", handler.GetCurrentUser)
			r.Patch("/", handler.UpdateCurrentUser)
			r.Patch("/profile", handler.UpdateCurrentUserProfile)
		})
	})
}

// * Helper untuk pagination
func (h *UserHandler) parsePaginationParams(r *http.Request) (int32, int32, error) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := int32(10)
	offset := int32(0)

	if limitStr != "" {
		l, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil {
			return 0, 0, domain.ErrBadRequest("Invalid limit parameter")
		}
		limit = int32(l)
	}

	if offsetStr != "" {
		o, err := strconv.ParseInt(offsetStr, 10, 32)
		if err != nil {
			return 0, 0, domain.ErrBadRequest("Invalid offset parameter")
		}
		offset = int32(o)
	}

	return limit, offset, nil
}

// *===========================CREATE===========================*
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload domain.CreateUserPayload

	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	user, err := h.Service.CreateUser(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusCreated, "User created successfully", user)
}

// *===========================READ (MANY)===========================*
func (h *UserHandler) GetUsersPaginated(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := h.parsePaginationParams(r)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	totalUsers, err := h.Service.CountUsers(r.Context())
	if err != nil {
		web.HandleError(w, err)
		return
	}

	users, err := h.Service.GetUsersPaginated(r.Context(), limit, offset)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	// * Kalkulasi untuk memenuhi parameter helper
	perPage := int(limit)
	currentPage := 1
	if perPage > 0 {
		currentPage = (int(offset) / perPage) + 1
	}

	web.SuccessWithPageInfo(w, http.StatusOK, "Users retrieved successfully", users, totalUsers, perPage, currentPage)
}

/*
func (h *UserHandler) GetUsersCursorForward(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")

	cursor, err := time.Parse(time.RFC3339, cursorStr)
	if err != nil {
		web.HandleError(w, domain.ErrBadRequest("Invalid cursor parameter"))
		return
	}

	limit := int32(10) // default
	if limitStr != "" {
		l, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil {
			web.HandleError(w, domain.ErrBadRequest("Invalid limit parameter"))
			return
		}
		limit = int32(l)
	}

	users, err := h.Service.GetUsersCursorForward(r.Context(), cursor, limit)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Users retrieved successfully", users)
}
*/

/*
func (h *UserHandler) GetUsersCursorBackward(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")

	cursor, err := time.Parse(time.RFC3339, cursorStr)
	if err != nil {
		web.HandleError(w, domain.ErrBadRequest("Invalid cursor parameter"))
		return
	}

	limit := int32(10) // default
	if limitStr != "" {
		l, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil {
			web.HandleError(w, domain.ErrBadRequest("Invalid limit parameter"))
			return
		}
		limit = int32(l)
	}

	users, err := h.Service.GetUsersCursorBackward(r.Context(), cursor, limit)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Users retrieved successfully", users)
}
*/

func (h *UserHandler) GetUsersCursorFirst(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")

	limit := int32(10) // default
	if limitStr != "" {
		l, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil {
			web.HandleError(w, domain.ErrBadRequest("Invalid limit parameter"))
			return
		}
		limit = int32(l)
	}

	// Pro Tip: Untuk implementasi hasNextPage yang akurat, idealnya kita meminta 'limit + 1' data.
	// Namun, untuk saat ini kita jaga agar service tidak perlu diubah.
	users, err := h.Service.GetUsersCursorFirst(r.Context(), limit)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	// Logika untuk menentukan parameter cursor
	var nextCursor string
	perPage := int(limit)
	// Asumsi: jika jumlah data yang kembali sama dengan limit, kemungkinan ada halaman selanjutnya.
	hasNextPage := len(users) == perPage

	if hasNextPage {
		// Ambil cursor dari item terakhir
		nextCursor = users[len(users)-1].CreatedAt.Format(time.RFC3339Nano)
	}

	web.SuccessWithCursor(w, http.StatusOK, "Users retrieved successfully", users, nextCursor, hasNextPage, perPage)
}

// ... (sisa file tidak diubah)

func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")
	if searchTerm == "" {
		web.HandleError(w, domain.ErrBadRequest("Search term is required"))
		return
	}

	users, err := h.Service.SearchUsers(r.Context(), searchTerm)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Users searched successfully", users)
}

func (h *UserHandler) SearchUsersPaginated(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")
	if searchTerm == "" {
		web.HandleError(w, domain.ErrBadRequest("Search term is required"))
		return
	}

	limit, offset, err := h.parsePaginationParams(r)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	users, err := h.Service.SearchUsersPaginated(r.Context(), searchTerm, limit, offset)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Users searched successfully", users)
}

func (h *UserHandler) SearchUsersByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		web.HandleError(w, domain.ErrBadRequest("Username parameter is required"))
		return
	}

	users, err := h.Service.SearchUsersByUsername(r.Context(), username)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Users searched successfully", users)
}

func (h *UserHandler) SearchUsersByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		web.HandleError(w, domain.ErrBadRequest("Email parameter is required"))
		return
	}

	users, err := h.Service.SearchUsersByEmail(r.Context(), email)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Users searched successfully", users)
}

func (h *UserHandler) SearchUsersByFullName(w http.ResponseWriter, r *http.Request) {
	fullName := r.URL.Query().Get("fullname")
	if fullName == "" {
		web.HandleError(w, domain.ErrBadRequest("Full name parameter is required"))
		return
	}

	users, err := h.Service.SearchUsersByFullName(r.Context(), fullName)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Users searched successfully", users)
}

// *===========================READ (ONE & UTILITY)===========================*
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("User ID is required"))
		return
	}

	user, err := h.Service.GetUser(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "User retrieved successfully", user)
}

func (h *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		web.HandleError(w, domain.ErrBadRequest("Username is required"))
		return
	}

	user, err := h.Service.GetUserByUsername(r.Context(), username)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "User retrieved successfully", user)
}

func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if email == "" {
		web.HandleError(w, domain.ErrBadRequest("Email is required"))
		return
	}

	user, err := h.Service.GetUserByEmail(r.Context(), email)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "User retrieved successfully", user)
}

func (h *UserHandler) CheckUserExists(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("User ID is required"))
		return
	}

	exists, err := h.Service.CheckUserExists(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := map[string]bool{"exists": exists}
	web.Success(w, http.StatusOK, "User existence checked successfully", response)
}

func (h *UserHandler) CheckEmailExists(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if email == "" {
		web.HandleError(w, domain.ErrBadRequest("Email is required"))
		return
	}

	exists, err := h.Service.CheckEmailExists(r.Context(), email)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := map[string]bool{"exists": exists}
	web.Success(w, http.StatusOK, "Email existence checked successfully", response)
}

func (h *UserHandler) CheckUsernameExists(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		web.HandleError(w, domain.ErrBadRequest("Username is required"))
		return
	}

	exists, err := h.Service.CheckUsernameExists(r.Context(), username)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := map[string]bool{"exists": exists}
	web.Success(w, http.StatusOK, "Username existence checked successfully", response)
}

func (h *UserHandler) CountUsers(w http.ResponseWriter, r *http.Request) {
	count, err := h.Service.CountUsers(r.Context())
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := map[string]int64{"count": count}
	web.Success(w, http.StatusOK, "Users counted successfully", response)
}

// * Curent user
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	claims, _ := web.GetUserFromContext(r.Context())

	if claims.IDUser == "" {
		web.HandleError(w, domain.ErrBadRequest("User ID is required"))
		return
	}

	user, err := h.Service.GetUser(r.Context(), claims.IDUser)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "User retrieved successfully", user)
}

// *===========================UPDATE===========================*
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("User ID is required"))
		return
	}

	var payload domain.User
	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	payload.ID = id

	user, err := h.Service.UpdateUser(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "User updated successfully", user)
}

func (h *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("User ID is required"))
		return
	}

	var payload domain.User
	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	payload.ID = id // Ensure ID matches URL parameter

	user, err := h.Service.UpdateUserProfile(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "User profile updated successfully", user)
}

func (h *UserHandler) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	claims, _ := web.GetUserFromContext(r.Context())
	if claims.IDUser == "" {
		web.HandleError(w, domain.ErrBadRequest("User ID is required"))
		return
	}

	var payload domain.User
	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	payload.ID = claims.IDUser

	user, err := h.Service.UpdateUser(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "User updated successfully", user)
}

func (h *UserHandler) UpdateCurrentUserProfile(w http.ResponseWriter, r *http.Request) {
	claims, _ := web.GetUserFromContext(r.Context())
	if claims.IDUser == "" {
		web.HandleError(w, domain.ErrBadRequest("User ID is required"))
		return
	}

	var payload domain.User
	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	payload.ID = claims.IDUser

	user, err := h.Service.UpdateUserProfile(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "User profile updated successfully", user)
}

// *===========================DELETE===========================*
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("User ID is required"))
		return
	}

	err := h.Service.DeleteUser(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "User deleted successfully", nil)
}
