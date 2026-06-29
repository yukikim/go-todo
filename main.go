package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

var todos = map[int]Todo{}
var nextID = 1

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})
	http.HandleFunc("/todos", todosHandler)
	http.HandleFunc("/todos/", todoHandler)

	fmt.Println("server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("server error:", err)
	}
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTodosHandler(w, r)
	case http.MethodPost:
		createTodoHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTodoHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	todoList := make([]Todo, 0, len(todos))
	for _, todo := range todos {
		todoList = append(todoList, todo)
	}

	writeJSON(w, http.StatusOK, todoList)
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	now := time.Now()
	todo := Todo{
		ID:          nextID,
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	todos[todo.ID] = todo
	nextID++

	writeJSON(w, http.StatusCreated, todo)
}

func getTodoHandler(w http.ResponseWriter, r *http.Request) {
	idText := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idText)
	if err != nil {
		http.Error(w, "invalid todo ID", http.StatusBadRequest)
		return
	}

	todo, ok := todos[id]
	if !ok {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
