package handler

import (
	"net/http"
	"strings"

	"go-todo/internal/service"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token, ok := strings.CutPrefix(authHeader, "Bearer ")
		if !ok || token == "" {
			writeError(c, http.StatusUnauthorized, "authorization header is required")
			c.Abort()
			return
		}

		username, err := authService.ValidateToken(token)
		if err != nil {
			writeError(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}

		c.Set("username", username)
		c.Next()
	}
}
