package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-todo/internal/handler"
	"go-todo/internal/model"
	"go-todo/internal/service"
)

func TestCreateTodoHandler(t *testing.T) {
	store := newMemoryTodoStore(nil, 1)

	body := bytes.NewBufferString(`{"title":"Goを学習する","description":"POST /todosを作る"}`)
	req := httptest.NewRequest(http.MethodPost, "/todos", body)
	rec := httptest.NewRecorder()

	newTestRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var todo model.Todo
	if err := json.NewDecoder(rec.Body).Decode(&todo); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if todo.ID != 1 {
		t.Fatalf("expected todo ID 1, got %d", todo.ID)
	}

	if todo.Title != "Goを学習する" {
		t.Fatalf("expected title %q, got %q", "Goを学習する", todo.Title)
	}

	if _, ok := store.todos[todo.ID]; !ok {
		t.Fatalf("expected todo to be saved")
	}
}

func TestCreateTodoHandlerTrimsInput(t *testing.T) {
	store := newMemoryTodoStore(nil, 1)

	body := bytes.NewBufferString(`{"title":"  Goを学習する  ","description":"  空白を取り除く  "}`)
	req := httptest.NewRequest(http.MethodPost, "/todos", body)
	rec := httptest.NewRecorder()

	newTestRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var todo model.Todo
	if err := json.NewDecoder(rec.Body).Decode(&todo); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if todo.Title != "Goを学習する" {
		t.Fatalf("expected title %q, got %q", "Goを学習する", todo.Title)
	}

	if todo.Description != "空白を取り除く" {
		t.Fatalf("expected description %q, got %q", "空白を取り除く", todo.Description)
	}
}

func TestCreateTodoHandlerValidationError(t *testing.T) {
	store := newMemoryTodoStore(nil, 1)

	body := bytes.NewBufferString(`{"title":"   ","description":"titleが空です"}`)
	req := httptest.NewRequest(http.MethodPost, "/todos", body)
	rec := httptest.NewRecorder()

	newTestRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var errorResponse handler.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&errorResponse); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errorResponse.Error != "title is required" {
		t.Fatalf("expected error %q, got %q", "title is required", errorResponse.Error)
	}
}

func TestCreateTodoHandlerTitleTooLong(t *testing.T) {
	store := newMemoryTodoStore(nil, 1)

	body := bytes.NewBufferString(`{"title":"` + strings.Repeat("a", service.MaxTitleLength+1) + `","description":"titleが長すぎます"}`)
	req := httptest.NewRequest(http.MethodPost, "/todos", body)
	rec := httptest.NewRecorder()

	newTestRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var errorResponse handler.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&errorResponse); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errorResponse.Error != "title must be 100 characters or less" {
		t.Fatalf("expected error %q, got %q", "title must be 100 characters or less", errorResponse.Error)
	}
}

func TestGetTodoHandler(t *testing.T) {
	now := time.Now()
	store := newMemoryTodoStore(map[int]model.Todo{
		1: {
			ID:          1,
			Title:       "Goを学習する",
			Description: "GET /todos/{id}を作る",
			Completed:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}, 2)
	req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
	rec := httptest.NewRecorder()

	newTestRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var todo model.Todo
	if err := json.NewDecoder(rec.Body).Decode(&todo); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if todo.ID != 1 {
		t.Fatalf("expected todo ID 1, got %d", todo.ID)
	}
}

func TestGetTodoHandlerNotFound(t *testing.T) {
	store := newMemoryTodoStore(nil, 1)

	req := httptest.NewRequest(http.MethodGet, "/todos/999", nil)
	rec := httptest.NewRecorder()

	newTestRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}

	var errorResponse handler.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&errorResponse); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errorResponse.Error != "todo not found" {
		t.Fatalf("expected error %q, got %q", "todo not found", errorResponse.Error)
	}
}

func TestUpdateTodoHandler(t *testing.T) {
	now := time.Now()
	store := newMemoryTodoStore(map[int]model.Todo{
		1: {
			ID:          1,
			Title:       "古いタイトル",
			Description: "古い説明",
			Completed:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}, 2)
	body := bytes.NewBufferString(`{"title":"新しいタイトル","description":"新しい説明","completed":true}`)
	req := httptest.NewRequest(http.MethodPut, "/todos/1", body)
	rec := httptest.NewRecorder()

	newTestRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var todo model.Todo
	if err := json.NewDecoder(rec.Body).Decode(&todo); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if todo.Title != "新しいタイトル" {
		t.Fatalf("expected title %q, got %q", "新しいタイトル", todo.Title)
	}

	if !todo.Completed {
		t.Fatalf("expected completed to be true")
	}

	if store.todos[1].Description != "新しい説明" {
		t.Fatalf("expected todo in map to be updated")
	}
}

func TestUpdateTodoHandlerDescriptionTooLong(t *testing.T) {
	now := time.Now()
	store := newMemoryTodoStore(map[int]model.Todo{
		1: {
			ID:          1,
			Title:       "古いタイトル",
			Description: "古い説明",
			Completed:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}, 2)
	body := bytes.NewBufferString(`{"title":"新しいタイトル","description":"` + strings.Repeat("a", service.MaxDescriptionLength+1) + `","completed":true}`)
	req := httptest.NewRequest(http.MethodPut, "/todos/1", body)
	rec := httptest.NewRecorder()

	newTestRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var errorResponse handler.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&errorResponse); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errorResponse.Error != "description must be 500 characters or less" {
		t.Fatalf("expected error %q, got %q", "description must be 500 characters or less", errorResponse.Error)
	}
}

func TestUpdateTodoHandlerNotFound(t *testing.T) {
	store := newMemoryTodoStore(nil, 1)

	body := bytes.NewBufferString(`{"title":"新しいタイトル","description":"新しい説明","completed":true}`)
	req := httptest.NewRequest(http.MethodPut, "/todos/999", body)
	rec := httptest.NewRecorder()

	newTestRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestDeleteTodoHandler(t *testing.T) {
	now := time.Now()
	store := newMemoryTodoStore(map[int]model.Todo{
		1: {
			ID:          1,
			Title:       "削除するTodo",
			Description: "DELETE /todos/{id}を作る",
			Completed:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}, 2)
	req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
	rec := httptest.NewRecorder()

	newTestRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}

	if _, ok := store.todos[1]; ok {
		t.Fatalf("expected todo to be deleted")
	}
}

func TestDeleteTodoHandlerNotFound(t *testing.T) {
	store := newMemoryTodoStore(nil, 1)

	req := httptest.NewRequest(http.MethodDelete, "/todos/999", nil)
	rec := httptest.NewRecorder()

	newTestRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestLoginHandlerReturnsToken(t *testing.T) {
	store := newMemoryTodoStore(nil, 1)
	body := bytes.NewBufferString(`{"username":"admin","password":"password"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	rec := httptest.NewRecorder()

	newRawTestRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response model.LoginResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Token == "" {
		t.Fatalf("expected token to be returned")
	}
}

func TestLoginHandlerInvalidCredentials(t *testing.T) {
	store := newMemoryTodoStore(nil, 1)
	body := bytes.NewBufferString(`{"username":"admin","password":"wrong"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	rec := httptest.NewRecorder()

	newRawTestRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestTodoRoutesRequireAuthorization(t *testing.T) {
	store := newMemoryTodoStore(nil, 1)
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	rec := httptest.NewRecorder()

	newRawTestRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	var errorResponse handler.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&errorResponse); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errorResponse.Error != "authorization header is required" {
		t.Fatalf("expected error %q, got %q", "authorization header is required", errorResponse.Error)
	}
}

func TestTodoRoutesRejectInvalidToken(t *testing.T) {
	store := newMemoryTodoStore(nil, 1)
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec := httptest.NewRecorder()

	newRawTestRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func newTestRouter(store *memoryTodoStore) http.Handler {
	router := newRawTestRouter(store)
	token, err := newTestAuthService().GenerateToken("admin")
	if err != nil {
		panic(err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, r)
	})
}

func newRawTestRouter(store *memoryTodoStore) http.Handler {
	return handler.NewRouter(service.NewTodoService(store), newTestAuthService())
}

func newTestAuthService() *service.AuthService {
	return service.NewAuthService("admin", "password", "test-secret")
}

func TestCompleteTodoHandler(t *testing.T) {
	now := time.Now()
	store := newMemoryTodoStore(map[int]model.Todo{
		1: {
			ID:          1,
			Title:       "完了にするTodo",
			Description: "PATCH /todos/{id}/completeを作る",
			Completed:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}, 2)
	req := httptest.NewRequest(http.MethodPatch, "/todos/1/complete", nil)
	rec := httptest.NewRecorder()

	newTestRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var todo model.Todo
	if err := json.NewDecoder(rec.Body).Decode(&todo); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !todo.Completed {
		t.Fatalf("expected completed to be true")
	}

	if !store.todos[1].Completed {
		t.Fatalf("expected todo in map to be completed")
	}
}

func TestCompleteTodoHandlerTogglesBackToIncomplete(t *testing.T) {
	now := time.Now()
	store := newMemoryTodoStore(map[int]model.Todo{
		1: {
			ID:          1,
			Title:       "未完了に戻すTodo",
			Description: "PATCH /todos/{id}/completeを作る",
			Completed:   true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}, 2)
	req := httptest.NewRequest(http.MethodPatch, "/todos/1/complete", nil)
	rec := httptest.NewRecorder()

	newTestRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if store.todos[1].Completed {
		t.Fatalf("expected todo to be incomplete")
	}
}

func TestCompleteTodoHandlerNotFound(t *testing.T) {
	store := newMemoryTodoStore(nil, 1)

	req := httptest.NewRequest(http.MethodPatch, "/todos/999/complete", nil)
	rec := httptest.NewRecorder()

	newTestRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
