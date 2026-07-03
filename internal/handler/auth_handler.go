package handler

import (
	"errors"
	"net/http"

	"go-todo/internal/model"
	"go-todo/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON")
		return
	}

	res, err := h.service.Login(req)
	if errors.Is(err, service.ErrInvalidCredentials) {
		writeError(c, http.StatusUnauthorized, "invalid username or password")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to login")
		return
	}

	c.JSON(http.StatusOK, res)
}
