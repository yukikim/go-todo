package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateTodoHandler(t *testing.T) {
	store := newMemoryTodoStore(nil, 1)
	todoStore = store

	body := bytes.NewBufferString(`{"title":"Goを学習する","description":"POST /todosを作る"}`)
	req := httptest.NewRequest(http.MethodPost, "/todos", body)
	rec := httptest.NewRecorder()

	newRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var todo Todo
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

func TestGetTodoHandler(t *testing.T) {
	now := time.Now()
	todoStore = newMemoryTodoStore(map[int]Todo{
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

	newRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var todo Todo
	if err := json.NewDecoder(rec.Body).Decode(&todo); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if todo.ID != 1 {
		t.Fatalf("expected todo ID 1, got %d", todo.ID)
	}
}

func TestGetTodoHandlerNotFound(t *testing.T) {
	todoStore = newMemoryTodoStore(nil, 1)

	req := httptest.NewRequest(http.MethodGet, "/todos/999", nil)
	rec := httptest.NewRecorder()

	newRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}

	var errorResponse ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&errorResponse); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errorResponse.Error != "todo not found" {
		t.Fatalf("expected error %q, got %q", "todo not found", errorResponse.Error)
	}
}

func TestUpdateTodoHandler(t *testing.T) {
	now := time.Now()
	store := newMemoryTodoStore(map[int]Todo{
		1: {
			ID:          1,
			Title:       "古いタイトル",
			Description: "古い説明",
			Completed:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}, 2)
	todoStore = store

	body := bytes.NewBufferString(`{"title":"新しいタイトル","description":"新しい説明","completed":true}`)
	req := httptest.NewRequest(http.MethodPut, "/todos/1", body)
	rec := httptest.NewRecorder()

	newRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var todo Todo
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

func TestUpdateTodoHandlerNotFound(t *testing.T) {
	todoStore = newMemoryTodoStore(nil, 1)

	body := bytes.NewBufferString(`{"title":"新しいタイトル","description":"新しい説明","completed":true}`)
	req := httptest.NewRequest(http.MethodPut, "/todos/999", body)
	rec := httptest.NewRecorder()

	newRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestDeleteTodoHandler(t *testing.T) {
	now := time.Now()
	store := newMemoryTodoStore(map[int]Todo{
		1: {
			ID:          1,
			Title:       "削除するTodo",
			Description: "DELETE /todos/{id}を作る",
			Completed:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}, 2)
	todoStore = store

	req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
	rec := httptest.NewRecorder()

	newRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}

	if _, ok := store.todos[1]; ok {
		t.Fatalf("expected todo to be deleted")
	}
}

func TestDeleteTodoHandlerNotFound(t *testing.T) {
	todoStore = newMemoryTodoStore(nil, 1)

	req := httptest.NewRequest(http.MethodDelete, "/todos/999", nil)
	rec := httptest.NewRecorder()

	newRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestCompleteTodoHandler(t *testing.T) {
	now := time.Now()
	store := newMemoryTodoStore(map[int]Todo{
		1: {
			ID:          1,
			Title:       "完了にするTodo",
			Description: "PATCH /todos/{id}/completeを作る",
			Completed:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}, 2)
	todoStore = store

	req := httptest.NewRequest(http.MethodPatch, "/todos/1/complete", nil)
	rec := httptest.NewRecorder()

	newRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var todo Todo
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
	store := newMemoryTodoStore(map[int]Todo{
		1: {
			ID:          1,
			Title:       "未完了に戻すTodo",
			Description: "PATCH /todos/{id}/completeを作る",
			Completed:   true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}, 2)
	todoStore = store

	req := httptest.NewRequest(http.MethodPatch, "/todos/1/complete", nil)
	rec := httptest.NewRecorder()

	newRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if store.todos[1].Completed {
		t.Fatalf("expected todo to be incomplete")
	}
}

func TestCompleteTodoHandlerNotFound(t *testing.T) {
	todoStore = newMemoryTodoStore(nil, 1)

	req := httptest.NewRequest(http.MethodPatch, "/todos/999/complete", nil)
	rec := httptest.NewRecorder()

	newRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
