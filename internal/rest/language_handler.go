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

type LanguageService interface {
	// * CREATE
	CreateLanguage(ctx context.Context, payload *domain.CreateLanguagePayload) (domain.Language, error)

	// * READ (MANY)
	GetLanguagesPaginated(ctx context.Context, limit int32, offset int32) ([]domain.Language, error)
	GetLanguagesCursorFirst(ctx context.Context, limit int32) ([]domain.Language, error)
	SearchLanguages(ctx context.Context, searchTerm string) ([]domain.Language, error)
	SearchLanguagesPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.Language, error)
	SearchLanguagesByName(ctx context.Context, name string) ([]domain.Language, error)
	SearchLanguagesByLangCode(ctx context.Context, langCode string) ([]domain.Language, error)

	// * READ (ONE & UTILITY)
	GetLanguage(ctx context.Context, id string) (domain.Language, error)
	GetLanguageByName(ctx context.Context, name string) (domain.Language, error)
	GetLanguageByLangCode(ctx context.Context, langCode string) (domain.Language, error)
	CheckLanguageExists(ctx context.Context, id string) (bool, error)
	CheckLangCodeExists(ctx context.Context, langCode string) (bool, error)
	CheckNameExists(ctx context.Context, name string) (bool, error)
	CountLanguages(ctx context.Context) (int64, error)

	// * UPDATE
	UpdateLanguage(ctx context.Context, payload *domain.Language) (domain.Language, error)

	// * DELETE
	DeleteLanguage(ctx context.Context, id string) error
}

type LanguageHandler struct {
	Service LanguageService
}

func NewLanguageHandler(r chi.Router, s LanguageService) {
	handler := &LanguageHandler{
		Service: s,
	}

	r.Route("/languages", func(r chi.Router) {
		// * CREATE
		r.With(middleware.Auth).Post("/", handler.CreateLanguage)

		// * READ (MANY)
		r.Get("/", handler.GetLanguagesPaginated)
		r.Get("/cursor", handler.GetLanguagesCursorFirst)
		r.Get("/search", handler.SearchLanguages)
		r.Get("/search/paginated", handler.SearchLanguagesPaginated)
		r.Get("/search/name", handler.SearchLanguagesByName)
		r.Get("/search/langCode", handler.SearchLanguagesByLangCode)

		// * READ (ONE & UTILITY)
		r.Get("/{id}", handler.GetLanguage)
		r.Get("/name/{name}", handler.GetLanguageByName)
		r.Get("/langCode/{langCode}", handler.GetLanguageByLangCode)
		r.Get("/exists/{id}", handler.CheckLanguageExists)
		r.Get("/exists/langCode/{langCode}", handler.CheckLangCodeExists)
		r.Get("/exists/name/{name}", handler.CheckNameExists)
		r.Get("/count", handler.CountLanguages)

		// * UPDATE
		r.With(middleware.Auth).Patch("/{id}", handler.UpdateLanguage)

		// * DELETE
		r.With(middleware.Auth).Delete("/{id}", handler.DeleteLanguage)
	})
}

// * Helper untuk pagination
func (h *LanguageHandler) parsePaginationParams(r *http.Request) (int32, int32, error) {
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
func (h *LanguageHandler) CreateLanguage(w http.ResponseWriter, r *http.Request) {
	var payload domain.CreateLanguagePayload

	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	language, err := h.Service.CreateLanguage(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusCreated, "Language created successfully", language)
}

// *===========================READ (MANY)===========================*
func (h *LanguageHandler) GetLanguagesPaginated(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := h.parsePaginationParams(r)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	totalLanguages, err := h.Service.CountLanguages(r.Context())
	if err != nil {
		web.HandleError(w, err)
		return
	}

	languages, err := h.Service.GetLanguagesPaginated(r.Context(), limit, offset)
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

	web.SuccessWithPageInfo(w, http.StatusOK, "Languages retrieved successfully", languages, totalLanguages, perPage, currentPage)
}

/*
func (h *LanguageHandler) GetLanguagesCursorForward(w http.ResponseWriter, r *http.Request) {
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

	languages, err := h.Service.GetLanguagesCursorForward(r.Context(), cursor, limit)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Languages retrieved successfully", languages)
}
*/

/*
func (h *LanguageHandler) GetLanguagesCursorBackward(w http.ResponseWriter, r *http.Request) {
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

	languages, err := h.Service.GetLanguagesCursorBackward(r.Context(), cursor, limit)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Languages retrieved successfully", languages)
}
*/

func (h *LanguageHandler) GetLanguagesCursorFirst(w http.ResponseWriter, r *http.Request) {
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
	languages, err := h.Service.GetLanguagesCursorFirst(r.Context(), limit)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	// Logika untuk menentukan parameter cursor
	var nextCursor string
	perPage := int(limit)
	// Asumsi: jika jumlah data yang kembali sama dengan limit, kemungkinan ada halaman selanjutnya.
	hasNextPage := len(languages) == perPage

	if hasNextPage {
		// Ambil cursor dari item terakhir
		nextCursor = languages[len(languages)-1].CreatedAt.Format(time.RFC3339Nano)
	}

	web.SuccessWithCursor(w, http.StatusOK, "Languages retrieved successfully", languages, nextCursor, hasNextPage, perPage)
}

// ... (sisa file tidak diubah)

func (h *LanguageHandler) SearchLanguages(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")
	if searchTerm == "" {
		web.HandleError(w, domain.ErrBadRequest("Search term is required"))
		return
	}

	languages, err := h.Service.SearchLanguages(r.Context(), searchTerm)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Languages searched successfully", languages)
}

func (h *LanguageHandler) SearchLanguagesPaginated(w http.ResponseWriter, r *http.Request) {
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

	languages, err := h.Service.SearchLanguagesPaginated(r.Context(), searchTerm, limit, offset)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Languages searched successfully", languages)
}

func (h *LanguageHandler) SearchLanguagesByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		web.HandleError(w, domain.ErrBadRequest("Name parameter is required"))
		return
	}

	languages, err := h.Service.SearchLanguagesByName(r.Context(), name)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Languages searched successfully", languages)
}

func (h *LanguageHandler) SearchLanguagesByLangCode(w http.ResponseWriter, r *http.Request) {
	langCode := r.URL.Query().Get("langCode")
	if langCode == "" {
		web.HandleError(w, domain.ErrBadRequest("LangCode parameter is required"))
		return
	}

	languages, err := h.Service.SearchLanguagesByLangCode(r.Context(), langCode)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Languages searched successfully", languages)
}

// *===========================READ (ONE & UTILITY)===========================*
func (h *LanguageHandler) GetLanguage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("Language ID is required"))
		return
	}

	language, err := h.Service.GetLanguage(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Language retrieved successfully", language)
}

func (h *LanguageHandler) GetLanguageByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		web.HandleError(w, domain.ErrBadRequest("Name is required"))
		return
	}

	language, err := h.Service.GetLanguageByName(r.Context(), name)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Language retrieved successfully", language)
}

func (h *LanguageHandler) GetLanguageByLangCode(w http.ResponseWriter, r *http.Request) {
	langCode := chi.URLParam(r, "langCode")
	if langCode == "" {
		web.HandleError(w, domain.ErrBadRequest("LangCode is required"))
		return
	}

	language, err := h.Service.GetLanguageByLangCode(r.Context(), langCode)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Language retrieved successfully", language)
}

func (h *LanguageHandler) CheckLanguageExists(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("Language ID is required"))
		return
	}

	exists, err := h.Service.CheckLanguageExists(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := map[string]bool{"exists": exists}
	web.Success(w, http.StatusOK, "Language existence checked successfully", response)
}

func (h *LanguageHandler) CheckLangCodeExists(w http.ResponseWriter, r *http.Request) {
	langCode := chi.URLParam(r, "langCode")
	if langCode == "" {
		web.HandleError(w, domain.ErrBadRequest("LangCode is required"))
		return
	}

	exists, err := h.Service.CheckLangCodeExists(r.Context(), langCode)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := map[string]bool{"exists": exists}
	web.Success(w, http.StatusOK, "LangCode existence checked successfully", response)
}

func (h *LanguageHandler) CheckNameExists(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		web.HandleError(w, domain.ErrBadRequest("Name is required"))
		return
	}

	exists, err := h.Service.CheckNameExists(r.Context(), name)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := map[string]bool{"exists": exists}
	web.Success(w, http.StatusOK, "Name existence checked successfully", response)
}

func (h *LanguageHandler) CountLanguages(w http.ResponseWriter, r *http.Request) {
	count, err := h.Service.CountLanguages(r.Context())
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := map[string]int64{"count": count}
	web.Success(w, http.StatusOK, "Languages counted successfully", response)
}

// *===========================UPDATE===========================*
func (h *LanguageHandler) UpdateLanguage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("Language ID is required"))
		return
	}

	var payload domain.Language
	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	payload.ID = id

	language, err := h.Service.UpdateLanguage(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Language updated successfully", language)
}

// *===========================DELETE===========================*
func (h *LanguageHandler) DeleteLanguage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("Language ID is required"))
		return
	}

	err := h.Service.DeleteLanguage(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Language deleted successfully", nil)
}
