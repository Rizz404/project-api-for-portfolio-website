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

type TechService interface {
	// * CREATE
	CreateTech(ctx context.Context, payload *domain.CreateTechPayload) (domain.Tech, error)

	// * READ (MANY)
	GetTechsPaginated(ctx context.Context, limit int32, offset int32) ([]domain.Tech, error)
	GetTechsCursorFirst(ctx context.Context, limit int32) ([]domain.Tech, error)
	SearchTechs(ctx context.Context, searchTerm string) ([]domain.Tech, error)
	SearchTechsPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.Tech, error)

	// * READ (ONE & UTILITY)
	GetTech(ctx context.Context, id string) (domain.Tech, error)
	GetTechByName(ctx context.Context, name string) (domain.Tech, error)
	CheckTechExists(ctx context.Context, id string) (bool, error)
	CheckTechNameExists(ctx context.Context, name string) (bool, error)
	CountTechs(ctx context.Context) (int64, error)

	// * UPDATE
	UpdateTech(ctx context.Context, payload *domain.Tech) (domain.Tech, error)

	// * DELETE
	DeleteTech(ctx context.Context, id string) error
}

type TechHandler struct {
	Service TechService
}

func NewTechHandler(r chi.Router, s TechService) {
	handler := &TechHandler{
		Service: s,
	}

	r.Route("/techs", func(r chi.Router) {
		// * CREATE
		r.With(middleware.Auth).Post("/", handler.CreateTech)

		// * READ (MANY)
		r.Get("/", handler.GetTechsPaginated)
		r.Get("/cursor", handler.GetTechsCursorFirst)
		r.Get("/search", handler.SearchTechs)
		r.Get("/search/paginated", handler.SearchTechsPaginated)

		// * READ (ONE & UTILITY)
		r.Get("/{id}", handler.GetTech)
		r.Get("/name/{name}", handler.GetTechByName)
		r.Get("/exists/{id}", handler.CheckTechExists)
		r.Get("/exists/name/{name}", handler.CheckTechNameExists)
		r.Get("/count", handler.CountTechs)

		// * UPDATE
		r.With(middleware.Auth).Patch("/{id}", handler.UpdateTech)

		// * DELETE
		r.With(middleware.Auth).Delete("/{id}", handler.DeleteTech)
	})
}

// * Helper untuk pagination
func (h *TechHandler) parsePaginationParams(r *http.Request) (int32, int32, error) {
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
func (h *TechHandler) CreateTech(w http.ResponseWriter, r *http.Request) {
	var payload domain.CreateTechPayload

	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	tech, err := h.Service.CreateTech(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusCreated, "Tech created successfully", tech)
}

// *===========================READ (MANY)===========================*
func (h *TechHandler) GetTechsPaginated(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := h.parsePaginationParams(r)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	totalTechs, err := h.Service.CountTechs(r.Context())
	if err != nil {
		web.HandleError(w, err)
		return
	}

	techs, err := h.Service.GetTechsPaginated(r.Context(), limit, offset)
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

	web.SuccessWithPageInfo(w, http.StatusOK, "Techs retrieved successfully", techs, totalTechs, perPage, currentPage)
}

/*
func (h *TechHandler) GetTechsCursorForward(w http.ResponseWriter, r *http.Request) {
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

	techs, err := h.Service.GetTechsCursorForward(r.Context(), cursor, limit)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Techs retrieved successfully", techs)
}
*/

/*
func (h *TechHandler) GetTechsCursorBackward(w http.ResponseWriter, r *http.Request) {
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

	techs, err := h.Service.GetTechsCursorBackward(r.Context(), cursor, limit)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Techs retrieved successfully", techs)
}
*/

func (h *TechHandler) GetTechsCursorFirst(w http.ResponseWriter, r *http.Request) {
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
	techs, err := h.Service.GetTechsCursorFirst(r.Context(), limit)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	// Logika untuk menentukan parameter cursor
	var nextCursor string
	perPage := int(limit)
	// Asumsi: jika jumlah data yang kembali sama dengan limit, kemungkinan ada halaman selanjutnya.
	hasNextPage := len(techs) == perPage

	if hasNextPage {
		// Ambil cursor dari item terakhir
		nextCursor = techs[len(techs)-1].CreatedAt.Format(time.RFC3339Nano)
	}

	web.SuccessWithCursor(w, http.StatusOK, "Techs retrieved successfully", techs, nextCursor, hasNextPage, perPage)
}

// ... (sisa file tidak diubah)

func (h *TechHandler) SearchTechs(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")
	if searchTerm == "" {
		web.HandleError(w, domain.ErrBadRequest("Search term is required"))
		return
	}

	techs, err := h.Service.SearchTechs(r.Context(), searchTerm)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Techs searched successfully", techs)
}

func (h *TechHandler) SearchTechsPaginated(w http.ResponseWriter, r *http.Request) {
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

	techs, err := h.Service.SearchTechsPaginated(r.Context(), searchTerm, limit, offset)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Techs searched successfully", techs)
}

// *===========================READ (ONE & UTILITY)===========================*
func (h *TechHandler) GetTech(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("Tech ID is required"))
		return
	}

	tech, err := h.Service.GetTech(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Tech retrieved successfully", tech)
}

func (h *TechHandler) GetTechByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		web.HandleError(w, domain.ErrBadRequest("Name is required"))
		return
	}

	tech, err := h.Service.GetTechByName(r.Context(), name)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Tech retrieved successfully", tech)
}

func (h *TechHandler) CheckTechExists(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("Tech ID is required"))
		return
	}

	exists, err := h.Service.CheckTechExists(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := map[string]bool{"exists": exists}
	web.Success(w, http.StatusOK, "Tech existence checked successfully", response)
}

func (h *TechHandler) CheckTechNameExists(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		web.HandleError(w, domain.ErrBadRequest("Name is required"))
		return
	}

	exists, err := h.Service.CheckTechNameExists(r.Context(), name)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := map[string]bool{"exists": exists}
	web.Success(w, http.StatusOK, "Name existence checked successfully", response)
}

func (h *TechHandler) CountTechs(w http.ResponseWriter, r *http.Request) {
	count, err := h.Service.CountTechs(r.Context())
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := map[string]int64{"count": count}
	web.Success(w, http.StatusOK, "Techs counted successfully", response)
}

// *===========================UPDATE===========================*
func (h *TechHandler) UpdateTech(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("Tech ID is required"))
		return
	}

	var payload domain.Tech
	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	payload.ID = id

	tech, err := h.Service.UpdateTech(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Tech updated successfully", tech)
}

// *===========================DELETE===========================*
func (h *TechHandler) DeleteTech(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.HandleError(w, domain.ErrBadRequest("Tech ID is required"))
		return
	}

	err := h.Service.DeleteTech(r.Context(), id)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Tech deleted successfully", nil)
}
