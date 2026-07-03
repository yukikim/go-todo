package handler

import (
	"net/http"

	"go-todo/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(todoService *service.TodoService, authService *service.AuthService, allowedOrigins []string) *gin.Engine {
	todoHandler := NewTodoHandler(todoService)
	authHandler := NewAuthHandler(authService)
	r := gin.Default()

	r.Use(CORSMiddleware(allowedOrigins))

	r.NoRoute(func(c *gin.Context) {
		writeError(c, http.StatusNotFound, "not found")
	})
	r.NoMethod(func(c *gin.Context) {
		writeError(c, http.StatusMethodNotAllowed, "method not allowed")
	})

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK\n")
	})

	r.POST("/login", authHandler.Login)

	todos := r.Group("/todos")
	todos.Use(AuthMiddleware(authService))
	todos.GET("", todoHandler.GetTodos)
	todos.POST("", todoHandler.CreateTodo)
	todos.GET("/:id", todoHandler.GetTodo)
	todos.PUT("/:id", todoHandler.UpdateTodo)
	todos.DELETE("/:id", todoHandler.DeleteTodo)
	todos.PATCH("/:id/complete", todoHandler.CompleteTodo)

	return r
}
