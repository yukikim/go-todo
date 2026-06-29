package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTodoHandler(t *testing.T) {
	todos = map[int]Todo{}
	nextID = 1

	body := bytes.NewBufferString(`{"title":"Goを学習する","description":"POST /todosを作る"}`)
	req := httptest.NewRequest(http.MethodPost, "/todos", body)
	rec := httptest.NewRecorder()

	todosHandler(rec, req)

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

	if _, ok := todos[todo.ID]; !ok {
		t.Fatalf("expected todo to be saved")
	}
}
