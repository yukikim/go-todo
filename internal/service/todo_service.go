package service

import (
	"context"
	"errors"
	"strings"

	"go-todo/internal/model"
	"go-todo/internal/repository"
)

var ErrTitleRequired = errors.New("title is required")
var ErrTitleTooLong = errors.New("title must be 100 characters or less")
var ErrDescriptionTooLong = errors.New("description must be 500 characters or less")

const MaxTitleLength = 100
const MaxDescriptionLength = 500

type TodoService struct {
	repository repository.TodoRepository
}

func NewTodoService(repository repository.TodoRepository) *TodoService {
	return &TodoService{repository: repository}
}

func (s *TodoService) ListTodos(ctx context.Context) ([]model.Todo, error) {
	return s.repository.ListTodos(ctx)
}

func (s *TodoService) CreateTodo(ctx context.Context, req model.CreateTodoRequest) (model.Todo, error) {
	validReq, err := validateCreateTodoRequest(req)
	if err != nil {
		return model.Todo{}, err
	}

	return s.repository.CreateTodo(ctx, validReq)
}

func (s *TodoService) GetTodo(ctx context.Context, id int) (model.Todo, error) {
	return s.repository.GetTodo(ctx, id)
}

func (s *TodoService) UpdateTodo(ctx context.Context, id int, req model.UpdateTodoRequest) (model.Todo, error) {
	validReq, err := validateUpdateTodoRequest(req)
	if err != nil {
		return model.Todo{}, err
	}

	return s.repository.UpdateTodo(ctx, id, validReq)
}

func (s *TodoService) DeleteTodo(ctx context.Context, id int) error {
	return s.repository.DeleteTodo(ctx, id)
}

func (s *TodoService) ToggleTodoComplete(ctx context.Context, id int) (model.Todo, error) {
	return s.repository.ToggleTodoComplete(ctx, id)
}

func IsValidationError(err error) bool {
	return errors.Is(err, ErrTitleRequired) ||
		errors.Is(err, ErrTitleTooLong) ||
		errors.Is(err, ErrDescriptionTooLong)
}

func validateCreateTodoRequest(req model.CreateTodoRequest) (model.CreateTodoRequest, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)

	if req.Title == "" {
		return model.CreateTodoRequest{}, ErrTitleRequired
	}

	if len([]rune(req.Title)) > MaxTitleLength {
		return model.CreateTodoRequest{}, ErrTitleTooLong
	}

	if len([]rune(req.Description)) > MaxDescriptionLength {
		return model.CreateTodoRequest{}, ErrDescriptionTooLong
	}

	return req, nil
}

func validateUpdateTodoRequest(req model.UpdateTodoRequest) (model.UpdateTodoRequest, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)

	if req.Title == "" {
		return model.UpdateTodoRequest{}, ErrTitleRequired
	}

	if len([]rune(req.Title)) > MaxTitleLength {
		return model.UpdateTodoRequest{}, ErrTitleTooLong
	}

	if len([]rune(req.Description)) > MaxDescriptionLength {
		return model.UpdateTodoRequest{}, ErrDescriptionTooLong
	}

	return req, nil
}
