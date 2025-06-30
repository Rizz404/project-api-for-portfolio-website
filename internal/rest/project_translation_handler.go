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

type ProjectTranslationService interface {
	// * CREATE
	CreateProjectTranslation(ctx context.Context, payload *domain.CreateProjectTranslationPayload) (domain.ProjectTranslation, error)
	CreateProjectTranslationsBatch(ctx context.Context, payload *[]domain.CreateProjectTranslationPayload) ([]domain.ProjectTranslation, error)

	// * READ (MANY)
	GetProjectTranslationsPaginated(ctx context.Context, limit int32, offset int32) ([]domain.ProjectTranslation, error)
	GetProjectTranslationsCursorFirst(ctx context.Context, limit int32) ([]domain.ProjectTranslation, error)
	SearchProjectTranslations(ctx context.Context, searchTerm string) ([]domain.ProjectTranslation, error)
	SearchProjectTranslationsPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.ProjectTranslation, error)

	// * READ (ONE & UTILITY)
	GetProjectTranslation(ctx context.Context, id string) (domain.ProjectTranslation, error)
	GetProjectTranslationByName(ctx context.Context, name string) (domain.ProjectTranslation, error)
	CheckProjectTranslationExists(ctx context.Context, id string) (bool, error)
	CountProjectTranslations(ctx context.Context) (int64, error)

	// * UPDATE
	UpdateProjectTranslation(ctx context.Context, payload *domain.ProjectTranslation) (domain.ProjectTranslation, error)

	// * DELETE
	DeleteProjectTranslation(ctx context.Context, id string) error
}

type ProjectTranslationHandler struct {
	Service ProjectTranslationService
}

func NewProjectTranslationHandler(r chi.Router, s ProjectTranslationService) {
	handler := &ProjectTranslationHandler{
		Service: s,
	}

	r.Route("/project-translations", func(r chi.Router) {
		// * CREATE
		r.With(middleware.Auth).Post("/", handler.CreateProjectTranslation)
		r.With(middleware.Auth).Post("/batch", handler.CreateProjectTranslationsBatch)

		// * READ (MANY)
		r.Get("/", handler.GetProjectTranslationsPaginated)
		r.Get("/cursor", handler.GetProjectTranslationsCursorFirst)
		r.Get("/search", handler.SearchProjectTranslations)
		r.Get("/search/paginated", handler.SearchProjectTranslationsPaginated)

		// * READ (ONE & UTILITY)
		r.Get("/{id}", handler.GetProjectTranslation)
		r.Get("/name/{name}", handler.GetProjectTranslationByName)
		r.Get("/exists/{id}", handler.CheckProjectTranslationExists)
		r.Get("/count", handler.CountProjectTranslations)

		// * UPDATE
		r.With(middleware.Auth).Patch("/{id}", handler.UpdateProjectTranslation)

		// * DELETE
		r.With(middleware.Auth).Delete("/{id}", handler.DeleteProjectTranslation)
	})
}

// * Helper untuk pagination
func (h *ProjectTranslationHandler) parsePaginationParams(r *http.Request) (int32, int32, error) {
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
func (h *ProjectTranslationHandler) CreateProjectTranslation(w http.ResponseWriter, r *http.Request) {
	var payload domain.CreateProjectTranslationPayload

	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	projectTranslation, err := h.Service.CreateProjectTranslation(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusCreated, "ProjectTranslation created successfully", projectTranslation)
}
func (h *ProjectTranslationHandler) CreateProjectTranslationsBatch(w http.ResponseWriter, r *http.Request) {
	var payload []domain.CreateProjectTranslationPayload

	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	projectTranslation, err := h.Service.CreateProjectTranslationsBatch(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusCreated, "ProjectTranslation created successfully", projectTranslation)
}

// *===========================READ (MANY)===========================*
func (h *ProjectTranslationHandler) GetProjectTranslationsPaginated(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := h.parsePaginationParams(r)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	totalProjectTranslations, err := h.Service.CountProjectTranslations(r.Context())
	if err != nil {
		web.HandleError(w, err)
		return
	}

	projectTranslations, err := h.Service.GetProjectTranslationsPaginated(r.Context(), limit, offset)
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

	web.SuccessWithPageInfo(w, http.StatusOK, "ProjectTranslations retrieved successfully", projectTranslations, totalProjectTranslations, perPage, currentPage)
}

/*
func (h *ProjectTranslationHandler) GetProjectTranslationsCursorForward(w http.ResponseWriter, r *http.Request) {
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

	projectTranslations, err := h.Service.GetProjectTranslationsCursorForward(r.Context(), cursor, limit)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "ProjectTranslations retrieved successfully", projectTranslations)
}
*/

/*
func (h *ProjectTranslationHandler) GetProjectTranslationsCursorBackward(w http.ResponseWriter, r *http.Request) {
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

	projectTranslations, err := h.Service.GetProjectTranslationsCursorBackward(r.Context(), cursor, limit)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "ProjectTranslations retrieved successfully", projectTranslations)
}
*/

func (h *ProjectTranslationHandler) GetProjectTranslationsCursorFirst(w http.ResponseWriter, r *http.Request) {
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
	projectTranslations, err := h.Service.GetProjectTranslationsCursorFirst(r.Context(), limit)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	// Logika untuk menentukan parameter cursor
	var nextCursor string
	perPage := int(limit)
	// Asumsi: jika jumlah data yang kembali sama dengan limit, kemungkinan ada halaman selanjutnya.
	hasNextPage := len(projectTranslations) == perPage

	if hasNextPage {
		// Ambil cursor dari item terakhir
		nextCursor = projectTranslations[len(projectTranslations)-1].CreatedAt.Format(time.RFC3339Nano)
	}

	web.SuccessWithCursor(w, http.StatusOK, "ProjectTranslations retrieved successfully", projectTranslations, nextCursor, hasNextPage, perPage)
}

// ... (sisa file tidak diubah)

func (h *ProjectTranslationHandler) SearchProjectTranslations(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")
	if searchTerm == "" {
		web.HandleError(w, domain.ErrBadRequest("Search term is required"))
		return
	}

	projectTranslations, err := h.Service.SearchProjectTranslations(r.Context(), searchTerm)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "ProjectTranslations searched successfully", projectTranslations)
}

func (h *ProjectTranslationHandler) SearchProjectTranslationsPaginated(w http.ResponseWriter, r *http.Request) {
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

	projectTranslations, err := h.Service.SearchProjectTranslationsPaginated(r.Context(), searchTerm, limit, offset)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "ProjectTranslations searched successfully", projectTranslations)
}

// *===========================READ (ONE & UTILITY)===========================*
func (h *ProjectTranslationHandler) GetProjectTranslation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("ProjectTranslation ID is required"))
		return
	}

	projectTranslation, err := h.Service.GetProjectTranslation(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "ProjectTranslation retrieved successfully", projectTranslation)
}

func (h *ProjectTranslationHandler) GetProjectTranslationByName(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "name")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("ProjectTranslation name is required"))
		return
	}

	projectTranslation, err := h.Service.GetProjectTranslationByName(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "ProjectTranslation retrieved successfully", projectTranslation)
}

func (h *ProjectTranslationHandler) CheckProjectTranslationExists(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("ProjectTranslation ID is required"))
		return
	}

	exists, err := h.Service.CheckProjectTranslationExists(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := map[string]bool{"exists": exists}
	web.Success(w, http.StatusOK, "ProjectTranslation existence checked successfully", response)
}

func (h *ProjectTranslationHandler) CountProjectTranslations(w http.ResponseWriter, r *http.Request) {
	count, err := h.Service.CountProjectTranslations(r.Context())
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := map[string]int64{"count": count}
	web.Success(w, http.StatusOK, "ProjectTranslations counted successfully", response)
}

// *===========================UPDATE===========================*
func (h *ProjectTranslationHandler) UpdateProjectTranslation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("ProjectTranslation ID is required"))
		return
	}

	var payload domain.ProjectTranslation
	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	payload.ID = id

	projectTranslation, err := h.Service.UpdateProjectTranslation(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "ProjectTranslation updated successfully", projectTranslation)
}

// *===========================DELETE===========================*
func (h *ProjectTranslationHandler) DeleteProjectTranslation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("ProjectTranslation ID is required"))
		return
	}

	err := h.Service.DeleteProjectTranslation(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "ProjectTranslation deleted successfully", nil)
}
