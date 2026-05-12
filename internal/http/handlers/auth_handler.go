package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/joaosimonis/ultimate-fodinha-api/internal/auth"
	"github.com/joaosimonis/ultimate-fodinha-api/internal/http/response"
)

type AuthService interface {
	Login(ctx context.Context, req auth.LoginRequest) (auth.LoginResponse, error)
}

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "dados inválidos")
		return
	}

	login, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, auth.ErrInvalidCredentials.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, "erro ao autenticar")
		return
	}

	response.Data(c, http.StatusOK, login)
}
