package main

import (
	"context"
	"errors"
	"strings"
)

var errTitleRequired = errors.New("title is required")
var errTitleTooLong = errors.New("title must be 100 characters or less")
var errDescriptionTooLong = errors.New("description must be 500 characters or less")

const maxTitleLength = 100
const maxDescriptionLength = 500

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
	validReq, err := validateCreateTodoRequest(req)
	if err != nil {
		return Todo{}, err
	}

	return s.store.CreateTodo(ctx, validReq)
}

func (s *TodoService) GetTodo(ctx context.Context, id int) (Todo, error) {
	return s.store.GetTodo(ctx, id)
}

func (s *TodoService) UpdateTodo(ctx context.Context, id int, req UpdateTodoRequest) (Todo, error) {
	validReq, err := validateUpdateTodoRequest(req)
	if err != nil {
		return Todo{}, err
	}

	return s.store.UpdateTodo(ctx, id, validReq)
}

func (s *TodoService) DeleteTodo(ctx context.Context, id int) error {
	return s.store.DeleteTodo(ctx, id)
}

func (s *TodoService) ToggleTodoComplete(ctx context.Context, id int) (Todo, error) {
	return s.store.ToggleTodoComplete(ctx, id)
}

func validateCreateTodoRequest(req CreateTodoRequest) (CreateTodoRequest, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)

	if req.Title == "" {
		return CreateTodoRequest{}, errTitleRequired
	}

	if len([]rune(req.Title)) > maxTitleLength {
		return CreateTodoRequest{}, errTitleTooLong
	}

	if len([]rune(req.Description)) > maxDescriptionLength {
		return CreateTodoRequest{}, errDescriptionTooLong
	}

	return req, nil
}

func validateUpdateTodoRequest(req UpdateTodoRequest) (UpdateTodoRequest, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)

	if req.Title == "" {
		return UpdateTodoRequest{}, errTitleRequired
	}

	if len([]rune(req.Title)) > maxTitleLength {
		return UpdateTodoRequest{}, errTitleTooLong
	}

	if len([]rune(req.Description)) > maxDescriptionLength {
		return UpdateTodoRequest{}, errDescriptionTooLong
	}

	return req, nil
}
