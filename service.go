package main

import (
	"context"
	"errors"
)

var errTitleRequired = errors.New("title is required")

type TodoService struct {
	store TodoStore
}

func NewTodoService(store TodoStore) *TodoService {
	return &TodoService{store: store}
}

func (s *TodoService) ListTodos(ctx context.Context) ([]Todo, error) {
	return s.store.ListTodos(ctx)
}

func (s *TodoService) CreateTodo(ctx context.Context, req CreateTodoRequest) (Todo, error) {
	if req.Title == "" {
		return Todo{}, errTitleRequired
	}

	return s.store.CreateTodo(ctx, req)
}

func (s *TodoService) GetTodo(ctx context.Context, id int) (Todo, error) {
	return s.store.GetTodo(ctx, id)
}

func (s *TodoService) UpdateTodo(ctx context.Context, id int, req UpdateTodoRequest) (Todo, error) {
	if req.Title == "" {
		return Todo{}, errTitleRequired
	}

	return s.store.UpdateTodo(ctx, id, req)
}

func (s *TodoService) DeleteTodo(ctx context.Context, id int) error {
	return s.store.DeleteTodo(ctx, id)
}

func (s *TodoService) ToggleTodoComplete(ctx context.Context, id int) (Todo, error) {
	return s.store.ToggleTodoComplete(ctx, id)
}
