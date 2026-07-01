package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getTodosHandler(c *gin.Context) {
	todos, err := todoService.ListTodos(c.Request.Context())
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list todos")
		return
	}

	c.JSON(http.StatusOK, todos)
}

func createTodoHandler(c *gin.Context) {
	// リクエストのJSONを格納するための構造体を作成(CreateTodoRequestはTypescriptのinterfaceのようなもの)
	var req CreateTodoRequest
	// ShouldBindJSONでリクエストボディをデコードし、reqに格納する
	if err := c.ShouldBindJSON(&req); err != nil {
		// JSONのデコードに失敗した場合、HTTPステータスコード400(Bad Request)を返す
		writeError(c, http.StatusBadRequest, "invalid JSON")
		return
	}

	todo, err := todoService.CreateTodo(c.Request.Context(), req)
	if errors.Is(err, errTitleRequired) {
		writeError(c, http.StatusBadRequest, "title is required")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to create todo")
		return
	}

	// 作成したTodoをJSON形式でレスポンスとして返す
	c.JSON(http.StatusCreated, todo)
}

func getTodoHandler(c *gin.Context) {
	id, err := getTodoIDFromContext(c)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid todo ID")
		return
	}

	todo, err := todoService.GetTodo(c.Request.Context(), id)
	if errors.Is(err, errTodoNotFound) {
		writeError(c, http.StatusNotFound, "todo not found")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to get todo")
		return
	}

	c.JSON(http.StatusOK, todo)
}

func getTodoIDFromContext(c *gin.Context) (int, error) {
	// *gin.Context.Param() は、URLパラメータを取得するためのメソッドです。例えば、URLが /todos/123 の場合、c.Param("id") は "123" を返します。
	idText := c.Param("id")
	return strconv.Atoi(idText)
}

func updateTodoHandler(c *gin.Context) {
	id, err := getTodoIDFromContext(c)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid todo ID")
		return
	}

	var req UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON")
		return
	}

	todo, err := todoService.UpdateTodo(c.Request.Context(), id, req)
	if errors.Is(err, errTitleRequired) {
		writeError(c, http.StatusBadRequest, "title is required")
		return
	}
	if errors.Is(err, errTodoNotFound) {
		writeError(c, http.StatusNotFound, "todo not found")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to update todo")
		return
	}

	c.JSON(http.StatusOK, todo)
}

func deleteTodoHandler(c *gin.Context) {
	id, err := getTodoIDFromContext(c)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid todo ID")
		return
	}

	err = todoService.DeleteTodo(c.Request.Context(), id)
	if errors.Is(err, errTodoNotFound) {
		writeError(c, http.StatusNotFound, "todo not found")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to delete todo")
		return
	}

	c.Status(http.StatusNoContent)
}

func completeTodoHandler(c *gin.Context) {
	id, err := getTodoIDFromContext(c)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid todo ID")
		return
	}

	todo, err := todoService.ToggleTodoComplete(c.Request.Context(), id)
	if errors.Is(err, errTodoNotFound) {
		writeError(c, http.StatusNotFound, "todo not found")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to update todo")
		return
	}

	c.JSON(http.StatusOK, todo)
}
