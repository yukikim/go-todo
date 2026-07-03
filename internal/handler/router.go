package handler

import (
	"net/http"

	"go-todo/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(todoService *service.TodoService) *gin.Engine {
	todoHandler := NewTodoHandler(todoService)
	r := gin.Default()

	r.NoRoute(func(c *gin.Context) {
		writeError(c, http.StatusNotFound, "not found")
	})
	r.NoMethod(func(c *gin.Context) {
		writeError(c, http.StatusMethodNotAllowed, "method not allowed")
	})

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK\n")
	})

	r.GET("/todos", todoHandler.GetTodos)
	r.POST("/todos", todoHandler.CreateTodo)
	r.GET("/todos/:id", todoHandler.GetTodo)
	r.PUT("/todos/:id", todoHandler.UpdateTodo)
	r.DELETE("/todos/:id", todoHandler.DeleteTodo)
	r.PATCH("/todos/:id/complete", todoHandler.CompleteTodo)

	return r
}
