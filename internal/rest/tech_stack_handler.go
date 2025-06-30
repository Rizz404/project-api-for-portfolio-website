package rest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/rest/middleware"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/web"
	"github.com/go-chi/chi/v5"
)

type TechStackService interface {
	// * CREATE
	CreateTechStack(ctx context.Context, payload *domain.CreateTechStackPayload) (domain.TechStack, error)

	// * UPDATE
	UpdateTechStack(ctx context.Context, payload *domain.UpdateTechStackPayload) (domain.TechStack, error)

	// * DELETE
	DeleteTechStack(ctx context.Context, payload *domain.DeleteTechStackPayload) error
}

type TechStackHandler struct {
	Service TechStackService
}

func NewTechStackHandler(r chi.Router, s TechStackService) {
	handler := &TechStackHandler{
		Service: s,
	}

	r.Route("/tech-stacks", func(r chi.Router) {
		// * CREATE
		r.With(middleware.Auth).Post("/", handler.CreateTechStack)

		// * UPDATE
		r.With(middleware.Auth).Patch("/project/{projectId}/tech/{techId}", handler.UpdateTechStack)

		// * DELETE
		r.With(middleware.Auth).Delete("/project/{projectId}/tech/{techId}", handler.DeleteTechStack)
	})
}

// * Helper untuk pagination
func (h *TechStackHandler) parsePaginationParams(r *http.Request) (int32, int32, error) {
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
func (h *TechStackHandler) CreateTechStack(w http.ResponseWriter, r *http.Request) {
	var payload domain.CreateTechStackPayload

	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	techStack, err := h.Service.CreateTechStack(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusCreated, "TechStack created successfully", techStack)
}

// *===========================UPDATE===========================*
func (h *TechStackHandler) UpdateTechStack(w http.ResponseWriter, r *http.Request) {
	projectId := chi.URLParam(r, "projectId")
	if projectId == "" {
		web.HandleError(w, domain.ErrBadRequest("Project ID is required"))
		return
	}
	techId := chi.URLParam(r, "techId")
	if techId == "" {
		web.HandleError(w, domain.ErrBadRequest("Tech ID is required"))
		return
	}

	var payload domain.UpdateTechStackPayload
	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	techStack, err := h.Service.UpdateTechStack(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "TechStack updated successfully", techStack)
}

// *===========================DELETE===========================*
func (h *TechStackHandler) DeleteTechStack(w http.ResponseWriter, r *http.Request) {
	projectId := chi.URLParam(r, "projectId")
	if projectId == "" {
		web.HandleError(w, domain.ErrBadRequest("Project ID is required"))
		return
	}
	techId := chi.URLParam(r, "techId")
	if techId == "" {
		web.HandleError(w, domain.ErrBadRequest("Tech ID is required"))
		return
	}

	var payload domain.DeleteTechStackPayload
	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	err := h.Service.DeleteTechStack(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "TechStack deleted successfully", nil)
}
