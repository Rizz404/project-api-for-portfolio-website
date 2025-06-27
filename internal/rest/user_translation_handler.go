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

type UserTranslationService interface {
	// * CREATE
	CreateUserTranslation(ctx context.Context, payload *domain.CreateUserTranslationPayload) (domain.UserTranslation, error)

	// * READ (MANY)
	GetUserTranslationsPaginated(ctx context.Context, limit int32, offset int32) ([]domain.UserTranslation, error)
	GetUserTranslationsByUserIDPaginated(ctx context.Context, idUser string, limit int32, offset int32) ([]domain.UserTranslation, error)
	GetUserTranslationsCursorFirst(ctx context.Context, limit int32) ([]domain.UserTranslation, error)

	// * READ (ONE & UTILITY)
	GetUserTranslation(ctx context.Context, id string) (domain.UserTranslation, error)
	GetUserTranslationByUserIDAndLangID(ctx context.Context, idUser string, idLanguage string) (domain.UserTranslation, error)
	GetUserTranslationByUserIDAndLangName(ctx context.Context, idUser string, langName string) (domain.UserTranslation, error)
	GetUserTranslationByUserIDAndLangCode(ctx context.Context, idUser string, langCode string) (domain.UserTranslation, error)
	CheckUserTranslationExists(ctx context.Context, id string) (bool, error)
	CountUserTranslations(ctx context.Context) (int64, error)

	// * UPDATE
	UpdateUserTranslation(ctx context.Context, id string, payload *domain.UpdateUserTranslationPayload) (domain.UserTranslation, error)

	// * DELETE
	DeleteUserTranslation(ctx context.Context, id string) error
}

type UserTranslationHandler struct {
	Service UserTranslationService
}

func NewUserTranslationHandler(r chi.Router, s UserTranslationService) {
	handler := &UserTranslationHandler{
		Service: s,
	}

	r.Route("/user-translations", func(r chi.Router) {
		// * CREATE
		r.With(middleware.Auth).Post("/", handler.CreateUserTranslation)

		// * READ (MANY)
		r.Get("/", handler.GetUserTranslationsPaginated)
		r.Get("/cursor", handler.GetUserTranslationsCursorFirst)
		r.Get("/user/{userID}", handler.GetUserTranslationsByUserIDPaginated) // Endpoint spesifik untuk user

		// * READ (ONE & UTILITY)
		r.Get("/{id}", handler.GetUserTranslation)
		r.Get("/lookup/by-lang-id", handler.GetUserTranslationByUserIDAndLangID) // Menggunakan query params
		r.Get("/lookup/by-lang-name", handler.GetUserTranslationByUserIDAndLangName)
		r.Get("/lookup/by-lang-code", handler.GetUserTranslationByUserIDAndLangCode)
		r.Get("/exists/{id}", handler.CheckUserTranslationExists)
		r.Get("/count", handler.CountUserTranslations)

		// * UPDATE
		r.With(middleware.Auth).Patch("/{id}", handler.UpdateUserTranslation)

		// * DELETE
		r.With(middleware.Auth).Delete("/{id}", handler.DeleteUserTranslation)
	})
}

// parsePaginationParams adalah helper lokal untuk mengambil parameter limit dan offset.
func (h *UserTranslationHandler) parsePaginationParams(r *http.Request) (int32, int32, error) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := int32(10) // default limit
	offset := int32(0) // default offset

	if limitStr != "" {
		l, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil || l <= 0 {
			return 0, 0, domain.ErrBadRequest("Invalid limit parameter")
		}
		limit = int32(l)
	}

	if offsetStr != "" {
		o, err := strconv.ParseInt(offsetStr, 10, 32)
		if err != nil || o < 0 {
			return 0, 0, domain.ErrBadRequest("Invalid offset parameter")
		}
		offset = int32(o)
	}

	return limit, offset, nil
}

// *===========================CREATE===========================*
func (h *UserTranslationHandler) CreateUserTranslation(w http.ResponseWriter, r *http.Request) {
	var payload domain.CreateUserTranslationPayload

	// Helper ini akan menangani decoding, validasi, dan response error jika gagal
	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	// Mungkin ada baiknya melakukan otorisasi di sini,
	// misalnya memastikan user yang login adalah user yang sama dengan payload.IDUser
	claims, ok := web.GetUserFromContext(r.Context())
	if !ok || claims.IDUser != payload.IDUser {
		web.HandleError(w, domain.ErrForbidden("You can only create translations for your own user profile"))
		return
	}

	translation, err := h.Service.CreateUserTranslation(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusCreated, "User translation created successfully", translation)
}

// *===========================READ (MANY)===========================*
func (h *UserTranslationHandler) GetUserTranslationsPaginated(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := h.parsePaginationParams(r)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	total, err := h.Service.CountUserTranslations(r.Context())
	if err != nil {
		web.HandleError(w, err)
		return
	}

	translations, err := h.Service.GetUserTranslationsPaginated(r.Context(), limit, offset)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	perPage := int(limit)
	currentPage := 1
	if perPage > 0 {
		currentPage = (int(offset) / perPage) + 1
	}

	web.SuccessWithPageInfo(w, http.StatusOK, "User translations retrieved successfully", translations, total, perPage, currentPage)
}

func (h *UserTranslationHandler) GetUserTranslationsByUserIDPaginated(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		web.HandleError(w, domain.ErrBadRequest("User ID is required"))
		return
	}

	limit, offset, err := h.parsePaginationParams(r)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	translations, err := h.Service.GetUserTranslationsByUserIDPaginated(r.Context(), userID, limit, offset)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	// Untuk endpoint ini, kita tidak menyertakan total count agar lebih sederhana
	web.Success(w, http.StatusOK, "User translations for the specified user retrieved successfully", translations)
}

func (h *UserTranslationHandler) GetUserTranslationsCursorFirst(w http.ResponseWriter, r *http.Request) {
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

	translations, err := h.Service.GetUserTranslationsCursorFirst(r.Context(), limit)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	var nextCursor string
	perPage := int(limit)
	hasNextPage := len(translations) == perPage

	if hasNextPage {
		nextCursor = translations[len(translations)-1].CreatedAt.Format(time.RFC3339Nano)
	}

	web.SuccessWithCursor(w, http.StatusOK, "User translations retrieved successfully", translations, nextCursor, hasNextPage, perPage)
}

// *===========================READ (ONE & UTILITY)===========================*
func (h *UserTranslationHandler) GetUserTranslation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("Translation ID is required"))
		return
	}

	translation, err := h.Service.GetUserTranslation(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "User translation retrieved successfully", translation)
}

func (h *UserTranslationHandler) GetUserTranslationByUserIDAndLangID(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	langID := r.URL.Query().Get("lang_id")

	if userID == "" || langID == "" {
		web.HandleError(w, domain.ErrBadRequest("Parameters 'user_id' and 'lang_id' are required"))
		return
	}

	translation, err := h.Service.GetUserTranslationByUserIDAndLangID(r.Context(), userID, langID)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "User translation retrieved successfully", translation)
}

func (h *UserTranslationHandler) GetUserTranslationByUserIDAndLangName(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	langName := r.URL.Query().Get("lang_name")

	if userID == "" || langName == "" {
		web.HandleError(w, domain.ErrBadRequest("Parameters 'user_id' and 'lang_name' are required"))
		return
	}

	translation, err := h.Service.GetUserTranslationByUserIDAndLangName(r.Context(), userID, langName)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "User translation retrieved successfully", translation)
}

func (h *UserTranslationHandler) GetUserTranslationByUserIDAndLangCode(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	langCode := r.URL.Query().Get("lang_code")

	if userID == "" || langCode == "" {
		web.HandleError(w, domain.ErrBadRequest("Parameters 'user_id' and 'lang_code' are required"))
		return
	}

	translation, err := h.Service.GetUserTranslationByUserIDAndLangCode(r.Context(), userID, langCode)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "User translation retrieved successfully", translation)
}

func (h *UserTranslationHandler) CheckUserTranslationExists(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("Translation ID is required"))
		return
	}

	exists, err := h.Service.CheckUserTranslationExists(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "User translation existence checked successfully", map[string]bool{"exists": exists})
}

func (h *UserTranslationHandler) CountUserTranslations(w http.ResponseWriter, r *http.Request) {
	count, err := h.Service.CountUserTranslations(r.Context())
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "User translations counted successfully", map[string]int64{"count": count})
}

// *===========================UPDATE===========================*
func (h *UserTranslationHandler) UpdateUserTranslation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("Translation ID is required"))
		return
	}

	var payload domain.UpdateUserTranslationPayload
	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	// Otorisasi: Pastikan user yang login adalah pemilik terjemahan ini
	translationToUpdate, err := h.Service.GetUserTranslation(r.Context(), id)
	if err != nil {
		web.HandleError(w, err) // ErrNotFound akan ditangani oleh helper
		return
	}
	claims, ok := web.GetUserFromContext(r.Context())
	if !ok || claims.IDUser != translationToUpdate.IDUser {
		web.HandleError(w, domain.ErrForbidden("You can only update your own translations"))
		return
	}

	updatedTranslation, err := h.Service.UpdateUserTranslation(r.Context(), id, &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "User translation updated successfully", updatedTranslation)
}

// *===========================DELETE===========================*
func (h *UserTranslationHandler) DeleteUserTranslation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("Translation ID is required"))
		return
	}

	// Otorisasi: Sama seperti update, cek kepemilikan sebelum menghapus
	translationToDelete, err := h.Service.GetUserTranslation(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}
	claims, ok := web.GetUserFromContext(r.Context())
	if !ok || claims.IDUser != translationToDelete.IDUser {
		web.HandleError(w, domain.ErrForbidden("You can only delete your own translations"))
		return
	}

	err = h.Service.DeleteUserTranslation(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "User translation deleted successfully", nil)
}
