package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"go-todo/internal/model"
)

type fakeTodoRepository struct {
	todos        map[int]model.Todo
	nextID       int
	createCalled bool
	updateCalled bool
}

func newFakeTodoRepository(initialTodos map[int]model.Todo, nextID int) *fakeTodoRepository {
	if initialTodos == nil {
		initialTodos = map[int]model.Todo{}
	}

	return &fakeTodoRepository{
		todos:  initialTodos,
		nextID: nextID,
	}
}

func (r *fakeTodoRepository) ListTodos(ctx context.Context) ([]model.Todo, error) {
	todos := make([]model.Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		todos = append(todos, todo)
	}

	return todos, nil
}

func (r *fakeTodoRepository) CreateTodo(ctx context.Context, req model.CreateTodoRequest) (model.Todo, error) {
	r.createCalled = true

	now := time.Now()
	todo := model.Todo{
		ID:          r.nextID,
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	r.todos[todo.ID] = todo
	r.nextID++

	return todo, nil
}

func (r *fakeTodoRepository) GetTodo(ctx context.Context, id int) (model.Todo, error) {
	todo, ok := r.todos[id]
	if !ok {
		return model.Todo{}, model.ErrTodoNotFound
	}

	return todo, nil
}

func (r *fakeTodoRepository) UpdateTodo(ctx context.Context, id int, req model.UpdateTodoRequest) (model.Todo, error) {
	r.updateCalled = true

	todo, ok := r.todos[id]
	if !ok {
		return model.Todo{}, model.ErrTodoNotFound
	}

	todo.Title = req.Title
	todo.Description = req.Description
	todo.Completed = req.Completed
	todo.UpdatedAt = time.Now()
	r.todos[id] = todo

	return todo, nil
}

func (r *fakeTodoRepository) DeleteTodo(ctx context.Context, id int) error {
	if _, ok := r.todos[id]; !ok {
		return model.ErrTodoNotFound
	}

	delete(r.todos, id)
	return nil
}

func (r *fakeTodoRepository) ToggleTodoComplete(ctx context.Context, id int) (model.Todo, error) {
	todo, ok := r.todos[id]
	if !ok {
		return model.Todo{}, model.ErrTodoNotFound
	}

	todo.Completed = !todo.Completed
	todo.UpdatedAt = time.Now()
	r.todos[id] = todo

	return todo, nil
}

func TestTodoServiceCreateTodoTrimsInput(t *testing.T) {
	repo := newFakeTodoRepository(nil, 1)
	service := NewTodoService(repo)

	todo, err := service.CreateTodo(context.Background(), model.CreateTodoRequest{
		Title:       "  Goを学習する  ",
		Description: "  service層をテストする  ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if todo.Title != "Goを学習する" {
		t.Fatalf("expected title %q, got %q", "Goを学習する", todo.Title)
	}

	if todo.Description != "service層をテストする" {
		t.Fatalf("expected description %q, got %q", "service層をテストする", todo.Description)
	}

	if !repo.createCalled {
		t.Fatalf("expected repository CreateTodo to be called")
	}
}

func TestTodoServiceCreateTodoValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		req  model.CreateTodoRequest
		want error
	}{
		{
			name: "title required",
			req: model.CreateTodoRequest{
				Title:       "   ",
				Description: "titleが空です",
			},
			want: ErrTitleRequired,
		},
		{
			name: "title too long",
			req: model.CreateTodoRequest{
				Title:       strings.Repeat("a", MaxTitleLength+1),
				Description: "titleが長すぎます",
			},
			want: ErrTitleTooLong,
		},
		{
			name: "description too long",
			req: model.CreateTodoRequest{
				Title:       "説明が長すぎるTodo",
				Description: strings.Repeat("a", MaxDescriptionLength+1),
			},
			want: ErrDescriptionTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeTodoRepository(nil, 1)
			service := NewTodoService(repo)

			_, err := service.CreateTodo(context.Background(), tt.req)
			if !errors.Is(err, tt.want) {
				t.Fatalf("expected error %v, got %v", tt.want, err)
			}

			if repo.createCalled {
				t.Fatalf("expected repository CreateTodo not to be called")
			}
		})
	}
}

func TestTodoServiceUpdateTodoValidationErrors(t *testing.T) {
	repo := newFakeTodoRepository(map[int]model.Todo{
		1: {
			ID:          1,
			Title:       "古いタイトル",
			Description: "古い説明",
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}, 2)
	service := NewTodoService(repo)

	_, err := service.UpdateTodo(context.Background(), 1, model.UpdateTodoRequest{
		Title:       "   ",
		Description: "titleが空です",
		Completed:   true,
	})
	if !errors.Is(err, ErrTitleRequired) {
		t.Fatalf("expected error %v, got %v", ErrTitleRequired, err)
	}

	if repo.updateCalled {
		t.Fatalf("expected repository UpdateTodo not to be called")
	}
}

func TestTodoServiceUpdateTodoNotFound(t *testing.T) {
	repo := newFakeTodoRepository(nil, 1)
	service := NewTodoService(repo)

	_, err := service.UpdateTodo(context.Background(), 999, model.UpdateTodoRequest{
		Title:       "新しいタイトル",
		Description: "存在しないTodo",
		Completed:   true,
	})
	if !errors.Is(err, model.ErrTodoNotFound) {
		t.Fatalf("expected error %v, got %v", model.ErrTodoNotFound, err)
	}
}

func TestTodoServiceToggleTodoComplete(t *testing.T) {
	repo := newFakeTodoRepository(map[int]model.Todo{
		1: {
			ID:          1,
			Title:       "完了にするTodo",
			Description: "service層をテストする",
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}, 2)
	service := NewTodoService(repo)

	todo, err := service.ToggleTodoComplete(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !todo.Completed {
		t.Fatalf("expected completed to be true")
	}
}
