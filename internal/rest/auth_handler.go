package rest

import (
	"context"
	"net/http"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/web"
	"github.com/go-chi/chi/v5"
)

type AuthService interface {
	Register(ctx context.Context, payload *domain.RegisterPayload) (domain.User, error)
	Login(ctx context.Context, payload *domain.LoginPayload) (domain.LoginResponse, error)
}

type AuthHandler struct {
	Service AuthService
}

func NewAuthHandler(r chi.Router, s AuthService) {
	handler := &AuthHandler{
		Service: s,
	}

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var payload domain.RegisterPayload

	// * Menggunakan utility untuk decode dan validate sekaligus
	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return // * Error sudah di-handle otomatis
	}

	user, err := h.Service.Register(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusCreated, "Register successfull", user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload domain.LoginPayload
	if !web.DecodeValidateAndRespond(w, r, &payload) {
		return
	}

	user, err := h.Service.Login(r.Context(), &payload)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	web.Success(w, http.StatusOK, "Login successfull", user)
}
