package handler

import (
	"errors"
	"net/http"
	"strconv"

	"go-todo/internal/model"
	"go-todo/internal/service"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	service *service.TodoService
}

func NewTodoHandler(service *service.TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

func (h *TodoHandler) GetTodos(c *gin.Context) {
	todos, err := h.service.ListTodos(c.Request.Context())
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list todos")
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req model.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON")
		return
	}

	todo, err := h.service.CreateTodo(c.Request.Context(), req)
	if service.IsValidationError(err) {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to create todo")
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) GetTodo(c *gin.Context) {
	id, err := getTodoIDFromContext(c)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid todo ID")
		return
	}

	todo, err := h.service.GetTodo(c.Request.Context(), id)
	if errors.Is(err, model.ErrTodoNotFound) {
		writeError(c, http.StatusNotFound, "todo not found")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to get todo")
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id, err := getTodoIDFromContext(c)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid todo ID")
		return
	}

	var req model.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON")
		return
	}

	todo, err := h.service.UpdateTodo(c.Request.Context(), id, req)
	if service.IsValidationError(err) {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	if errors.Is(err, model.ErrTodoNotFound) {
		writeError(c, http.StatusNotFound, "todo not found")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to update todo")
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id, err := getTodoIDFromContext(c)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid todo ID")
		return
	}

	err = h.service.DeleteTodo(c.Request.Context(), id)
	if errors.Is(err, model.ErrTodoNotFound) {
		writeError(c, http.StatusNotFound, "todo not found")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to delete todo")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TodoHandler) CompleteTodo(c *gin.Context) {
	id, err := getTodoIDFromContext(c)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid todo ID")
		return
	}

	todo, err := h.service.ToggleTodoComplete(c.Request.Context(), id)
	if errors.Is(err, model.ErrTodoNotFound) {
		writeError(c, http.StatusNotFound, "todo not found")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to update todo")
		return
	}

	c.JSON(http.StatusOK, todo)
}

func getTodoIDFromContext(c *gin.Context) (int, error) {
	idText := c.Param("id")
	return strconv.Atoi(idText)
}
