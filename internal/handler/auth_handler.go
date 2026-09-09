package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/TakedaB/akademi-api/internal/middleware"
	"github.com/TakedaB/akademi-api/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			respondError(w, http.StatusUnauthorized, "email ou senha inválidos")
			return
		}
		respondError(w, http.StatusInternalServerError, "erro interno")
		return
	}

	respondJSON(w, http.StatusOK, loginResponse{Token: token})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*service.Claims)
	if !ok {
		respondError(w, http.StatusUnauthorized, "não autenticado")
		return
	}

	user, err := h.service.GetProfile(claims.UserID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			respondError(w, http.StatusNotFound, "usuário não encontrado")
			return
		}
		respondError(w, http.StatusInternalServerError, "erro interno")
		return
	}

	respondJSON(w, http.StatusOK, user)
}
