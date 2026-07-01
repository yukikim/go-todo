package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// newRouter は、Gin のルーターを作成し、各エンドポイントに対応するハンドラを設定する関数です。
func newRouter() *gin.Engine {
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

	r.GET("/todos", getTodosHandler)
	r.POST("/todos", createTodoHandler)
	r.GET("/todos/:id", getTodoHandler)
	r.PUT("/todos/:id", updateTodoHandler)
	r.DELETE("/todos/:id", deleteTodoHandler)
	r.PATCH("/todos/:id/complete", completeTodoHandler)

	return r
}
