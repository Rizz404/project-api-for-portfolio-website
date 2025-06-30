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

type ProjectImageService interface {
	// * CREATE
	CreateProjectImage(ctx context.Context, payload *domain.CreateProjectImagePayload) (domain.ProjectImage, error)
	CreateProjectImagesBatch(ctx context.Context, payload *[]domain.CreateProjectImagePayload) ([]domain.ProjectImage, error)

	// * READ (MANY)
	GetProjectImagesPaginated(ctx context.Context, limit int32, offset int32) ([]domain.ProjectImage, error)
	GetProjectImagesCursorFirst(ctx context.Context, limit int32) ([]domain.ProjectImage, error)
	SearchProjectImages(ctx context.Context, searchTerm string) ([]domain.ProjectImage, error)
	SearchProjectImagesPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.ProjectImage, error)

	// * READ (ONE & UTILITY)
	GetProjectImage(ctx context.Context, id string) (domain.ProjectImage, error)
	CheckProjectImageExists(ctx context.Context, id string) (bool, error)
	CountProjectImages(ctx context.Context) (int64, error)

	// * UPDATE
	UpdateProjectImage(ctx context.Context, payload *domain.ProjectImage) (domain.ProjectImage, error)

	// * DELETE
	DeleteProjectImage(ctx context.Context, id string) error
}

type ProjectImageHandler struct {
	Service ProjectImageService
}

func NewProjectImageHandler(r chi.Router, s ProjectImageService) {
	handler := &ProjectImageHandler{
		Service: s,
	}

	r.Route("/project-images", func(r chi.Router) {
		// * CREATE
		r.With(middleware.Auth).Post("/", handler.CreateProjectImage)
		r.With(middleware.Auth).Post("/batch", handler.CreateProjectImagesBatch)

		// * READ (MANY)
		r.Get("/", handler.GetProjectImagesPaginated)
		r.Get("/cursor", handler.GetProjectImagesCursorFirst)
		r.Get("/search", handler.SearchProjectImages)
		r.Get("/search/paginated", handler.SearchProjectImagesPaginated)

		// * READ (ONE & UTILITY)
		r.Get("/{id}", handler.GetProjectImage)
		r.Get("/exists/{id}", handler.CheckProjectImageExists)
		r.Get("/count", handler.CountProjectImages)

		// * UPDATE
		r.With(middleware.Auth).Patch("/{id}", handler.UpdateProjectImage)

		// * DELETE
		r.With(middleware.Auth).Delete("/{id}", handler.DeleteProjectImage)
	})
}

// * Helper untuk pagination
func (h *ProjectImageHandler) parsePaginationParams(r *http.Request) (int32, int32, error) {
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
func (h *ProjectImageHandler) CreateProjectImage(w http.ResponseWriter, r *http.Request) {
	var payload domain.CreateProjectImagePayload

	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	projectImage, err := h.Service.CreateProjectImage(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusCreated, "ProjectImage created successfully", projectImage)
}
func (h *ProjectImageHandler) CreateProjectImagesBatch(w http.ResponseWriter, r *http.Request) {
	var payload []domain.CreateProjectImagePayload

	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	projectImage, err := h.Service.CreateProjectImagesBatch(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusCreated, "ProjectImage created successfully", projectImage)
}

// *===========================READ (MANY)===========================*
func (h *ProjectImageHandler) GetProjectImagesPaginated(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := h.parsePaginationParams(r)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	totalProjectImages, err := h.Service.CountProjectImages(r.Context())
	if err != nil {
		web.HandleError(w, err)
		return
	}

	projectImages, err := h.Service.GetProjectImagesPaginated(r.Context(), limit, offset)
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

	web.SuccessWithPageInfo(w, http.StatusOK, "ProjectImages retrieved successfully", projectImages, totalProjectImages, perPage, currentPage)
}

/*
func (h *ProjectImageHandler) GetProjectImagesCursorForward(w http.ResponseWriter, r *http.Request) {
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

	projectImages, err := h.Service.GetProjectImagesCursorForward(r.Context(), cursor, limit)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "ProjectImages retrieved successfully", projectImages)
}
*/

/*
func (h *ProjectImageHandler) GetProjectImagesCursorBackward(w http.ResponseWriter, r *http.Request) {
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

	projectImages, err := h.Service.GetProjectImagesCursorBackward(r.Context(), cursor, limit)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "ProjectImages retrieved successfully", projectImages)
}
*/

func (h *ProjectImageHandler) GetProjectImagesCursorFirst(w http.ResponseWriter, r *http.Request) {
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
	projectImages, err := h.Service.GetProjectImagesCursorFirst(r.Context(), limit)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	// Logika untuk menentukan parameter cursor
	var nextCursor string
	perPage := int(limit)
	// Asumsi: jika jumlah data yang kembali sama dengan limit, kemungkinan ada halaman selanjutnya.
	hasNextPage := len(projectImages) == perPage

	if hasNextPage {
		// Ambil cursor dari item terakhir
		nextCursor = projectImages[len(projectImages)-1].CreatedAt.Format(time.RFC3339Nano)
	}

	web.SuccessWithCursor(w, http.StatusOK, "ProjectImages retrieved successfully", projectImages, nextCursor, hasNextPage, perPage)
}

// ... (sisa file tidak diubah)

func (h *ProjectImageHandler) SearchProjectImages(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")
	if searchTerm == "" {
		web.HandleError(w, domain.ErrBadRequest("Search term is required"))
		return
	}

	projectImages, err := h.Service.SearchProjectImages(r.Context(), searchTerm)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "ProjectImages searched successfully", projectImages)
}

func (h *ProjectImageHandler) SearchProjectImagesPaginated(w http.ResponseWriter, r *http.Request) {
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

	projectImages, err := h.Service.SearchProjectImagesPaginated(r.Context(), searchTerm, limit, offset)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "ProjectImages searched successfully", projectImages)
}

// *===========================READ (ONE & UTILITY)===========================*
func (h *ProjectImageHandler) GetProjectImage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("ProjectImage ID is required"))
		return
	}

	projectImage, err := h.Service.GetProjectImage(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "ProjectImage retrieved successfully", projectImage)
}

func (h *ProjectImageHandler) CheckProjectImageExists(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("ProjectImage ID is required"))
		return
	}

	exists, err := h.Service.CheckProjectImageExists(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := map[string]bool{"exists": exists}
	web.Success(w, http.StatusOK, "ProjectImage existence checked successfully", response)
}

func (h *ProjectImageHandler) CountProjectImages(w http.ResponseWriter, r *http.Request) {
	count, err := h.Service.CountProjectImages(r.Context())
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := map[string]int64{"count": count}
	web.Success(w, http.StatusOK, "ProjectImages counted successfully", response)
}

// *===========================UPDATE===========================*
func (h *ProjectImageHandler) UpdateProjectImage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("ProjectImage ID is required"))
		return
	}

	var payload domain.ProjectImage
	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	payload.ID = id

	projectImage, err := h.Service.UpdateProjectImage(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "ProjectImage updated successfully", projectImage)
}

// *===========================DELETE===========================*
func (h *ProjectImageHandler) DeleteProjectImage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("ProjectImage ID is required"))
		return
	}

	err := h.Service.DeleteProjectImage(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "ProjectImage deleted successfully", nil)
}
