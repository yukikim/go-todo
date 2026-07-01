package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// chi.Router を使ってルーティングを設定する関数
func newRouter() http.Handler {
	r := chi.NewRouter()

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "not found")
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	r.Get("/todos", getTodosHandler)
	r.Post("/todos", createTodoHandler)
	r.Get("/todos/{id}", getTodoHandler)
	r.Put("/todos/{id}", updateTodoHandler)
	r.Delete("/todos/{id}", deleteTodoHandler)
	r.Patch("/todos/{id}/complete", completeTodoHandler)

	return r
}
