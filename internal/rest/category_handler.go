package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/web"
	"github.com/go-chi/chi/v5"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, payload *domain.CreateCategoryPayload) (domain.Category, error)
	GetCategories(ctx context.Context) ([]domain.Category, error)
	GetCategory(ctx context.Context, id string) (domain.Category, error)
	UpdateCategory(ctx context.Context, id string, payload *domain.UpdateCategoryPayload) (domain.Category, error)
	DeleteCategory(ctx context.Context, id string) error
}

type CategoryHandler struct {
	Service CategoryService
}

func NewCategoryHandler(r chi.Router, s CategoryService) {
	handler := &CategoryHandler{
		Service: s,
	}

	r.Route("/categories", func(r chi.Router) {
		r.Post("/", handler.CreateCategory)
		r.Get("/", handler.GetCategories)
		r.Get("/{id}", handler.GetCategory)
		r.Patch("/{id}", handler.UpdateCategory)
		r.Delete("/{id}", handler.DeleteCategory)
	})
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var payload domain.CreateCategoryPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		web.Error(w, http.StatusBadRequest, "Invalid request payload: "+err.Error(), nil)
		return
	}

	category, err := h.Service.CreateCategory(r.Context(), &payload)
	if err != nil {
		web.Error(w, http.StatusConflict, "Failed to create category, perhaps it already exists.", nil)
		return
	}

	web.Success(w, http.StatusCreated, "Category created successfully", category)
}

func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.Service.GetCategories(r.Context())
	if err != nil {
		web.Error(w, http.StatusInternalServerError, "Could not fetch categories.", nil)
		return
	}
	web.Success(w, http.StatusOK, "Categories fetched successfully", categories)
}

func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	category, err := h.Service.GetCategory(r.Context(), id)
	if err != nil {
		web.Error(w, http.StatusNotFound, "Category with ID '"+id+"' not found.", nil)
		return
	}
	web.Success(w, http.StatusOK, "Category fetched successfully", category)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var payload domain.UpdateCategoryPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		web.Error(w, http.StatusBadRequest, "Invalid request payload: "+err.Error(), nil)
		return
	}

	updatedCategory, err := h.Service.UpdateCategory(r.Context(), id, &payload)
	if err != nil {
		web.Error(w, http.StatusNotFound, "Failed to update: category with ID '"+id+"' not found.", nil)
		return
	}

	web.Success(w, http.StatusOK, "Category updated successfully", updatedCategory)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.Service.DeleteCategory(r.Context(), id)
	if err != nil {
		web.Error(w, http.StatusNotFound, "Failed to delete: category with ID '"+id+"' not found.", nil)
		return
	}
	web.Success(w, http.StatusOK, "Category deleted successfully", nil)
}
