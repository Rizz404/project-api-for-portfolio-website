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

type ProjectService interface {
	// * CREATE
	CreateProject(ctx context.Context, payload *domain.CreateProjectPayload, userId string) (domain.Project, error)

	// * READ (MANY)
	GetProjectsPaginated(ctx context.Context, limit int32, offset int32) ([]domain.Project, error)
	GetProjectsCursorFirst(ctx context.Context, limit int32) ([]domain.Project, error)
	SearchProjectsPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.Project, error)

	// * READ (ONE & UTILITY)
	GetProject(ctx context.Context, id string) (domain.Project, error)
	GetProjectByTranslatedName(ctx context.Context, name string) (domain.Project, error)
	CheckProjectExists(ctx context.Context, id string) (bool, error)
	CountProjects(ctx context.Context) (int64, error)

	// * UPDATE
	UpdateProject(ctx context.Context, payload *domain.Project) (domain.Project, error)

	// * DELETE
	DeleteProject(ctx context.Context, id string) error
}

type ProjectHandler struct {
	Service ProjectService
}

func NewProjectHandler(r chi.Router, s ProjectService) {
	handler := &ProjectHandler{
		Service: s,
	}

	r.Route("/projects", func(r chi.Router) {
		// * CREATE
		r.With(middleware.Auth).Post("/", handler.CreateProject)

		// * READ (MANY)
		r.Get("/", handler.GetProjectsPaginated)
		r.Get("/cursor", handler.GetProjectsCursorFirst)
		r.Get("/search/paginated", handler.SearchProjectsPaginated)

		// * READ (ONE & UTILITY)
		r.Get("/{id}", handler.GetProject)
		r.Get("/name/{name}", handler.GetProjectByTranslatedName)
		r.Get("/exists/{id}", handler.CheckProjectExists)
		r.Get("/count", handler.CountProjects)

		// * UPDATE
		r.With(middleware.Auth).Patch("/{id}", handler.UpdateProject)

		// * DELETE
		r.With(middleware.Auth).Delete("/{id}", handler.DeleteProject)
	})
}

// * Helper untuk pagination
func (h *ProjectHandler) parsePaginationParams(r *http.Request) (int32, int32, error) {
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
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	claims, _ := web.GetUserFromContext(r.Context())

	var payload domain.CreateProjectPayload

	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	project, err := h.Service.CreateProject(r.Context(), &payload, claims.IDUser)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusCreated, "Project created successfully", project)
}

// *===========================READ (MANY)===========================*
func (h *ProjectHandler) GetProjectsPaginated(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := h.parsePaginationParams(r)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	totalProjects, err := h.Service.CountProjects(r.Context())
	if err != nil {
		web.HandleError(w, err)
		return
	}

	projects, err := h.Service.GetProjectsPaginated(r.Context(), limit, offset)
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

	web.SuccessWithPageInfo(w, http.StatusOK, "Projects retrieved successfully", projects, totalProjects, perPage, currentPage)
}

func (h *ProjectHandler) GetProjectsCursorFirst(w http.ResponseWriter, r *http.Request) {
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
	projects, err := h.Service.GetProjectsCursorFirst(r.Context(), limit)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	// Logika untuk menentukan parameter cursor
	var nextCursor string
	perPage := int(limit)
	// Asumsi: jika jumlah data yang kembali sama dengan limit, kemungkinan ada halaman selanjutnya.
	hasNextPage := len(projects) == perPage

	if hasNextPage {
		// Ambil cursor dari item terakhir
		nextCursor = projects[len(projects)-1].CreatedAt.Format(time.RFC3339Nano)
	}

	web.SuccessWithCursor(w, http.StatusOK, "Projects retrieved successfully", projects, nextCursor, hasNextPage, perPage)
}

func (h *ProjectHandler) SearchProjectsPaginated(w http.ResponseWriter, r *http.Request) {
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

	projects, err := h.Service.SearchProjectsPaginated(r.Context(), searchTerm, limit, offset)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Projects searched successfully", projects)
}

// *===========================READ (ONE & UTILITY)===========================*
func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("Project ID is required"))
		return
	}

	project, err := h.Service.GetProject(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Project retrieved successfully", project)
}

func (h *ProjectHandler) GetProjectByTranslatedName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		web.HandleError(w, domain.ErrBadRequest("Name is required"))
		return
	}

	project, err := h.Service.GetProjectByTranslatedName(r.Context(), name)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Project retrieved successfully", project)
}

func (h *ProjectHandler) CheckProjectExists(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("Project ID is required"))
		return
	}

	exists, err := h.Service.CheckProjectExists(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := map[string]bool{"exists": exists}
	web.Success(w, http.StatusOK, "Project existence checked successfully", response)
}

func (h *ProjectHandler) CountProjects(w http.ResponseWriter, r *http.Request) {
	count, err := h.Service.CountProjects(r.Context())
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := map[string]int64{"count": count}
	web.Success(w, http.StatusOK, "Projects counted successfully", response)
}

// *===========================UPDATE===========================*
func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("Project ID is required"))
		return
	}

	var payload domain.Project
	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	payload.ID = id

	project, err := h.Service.UpdateProject(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Project updated successfully", project)
}

// *===========================DELETE===========================*
func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("Project ID is required"))
		return
	}

	err := h.Service.DeleteProject(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Project deleted successfully", nil)
}
