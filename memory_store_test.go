package main

import (
	"context"
	"sync"
	"time"

	"go-todo/internal/model"
)

type memoryTodoStore struct {
	mu     sync.Mutex
	todos  map[int]model.Todo
	nextID int
}

func newMemoryTodoStore(initialTodos map[int]model.Todo, nextID int) *memoryTodoStore {
	if initialTodos == nil {
		initialTodos = map[int]model.Todo{}
	}

	return &memoryTodoStore{
		todos:  initialTodos,
		nextID: nextID,
	}
}

func (s *memoryTodoStore) ListTodos(ctx context.Context) ([]model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todoList := make([]model.Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		todoList = append(todoList, todo)
	}

	return todoList, nil
}

func (s *memoryTodoStore) CreateTodo(ctx context.Context, req model.CreateTodoRequest) (model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	todo := model.Todo{
		ID:          s.nextID,
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.todos[todo.ID] = todo
	s.nextID++

	return todo, nil
}

func (s *memoryTodoStore) GetTodo(ctx context.Context, id int) (model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.todos[id]
	if !ok {
		return model.Todo{}, model.ErrTodoNotFound
	}

	return todo, nil
}

func (s *memoryTodoStore) UpdateTodo(ctx context.Context, id int, req model.UpdateTodoRequest) (model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.todos[id]
	if !ok {
		return model.Todo{}, model.ErrTodoNotFound
	}

	todo.Title = req.Title
	todo.Description = req.Description
	todo.Completed = req.Completed
	todo.UpdatedAt = time.Now()
	s.todos[id] = todo

	return todo, nil
}

func (s *memoryTodoStore) DeleteTodo(ctx context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.todos[id]; !ok {
		return model.ErrTodoNotFound
	}

	delete(s.todos, id)
	return nil
}

func (s *memoryTodoStore) ToggleTodoComplete(ctx context.Context, id int) (model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.todos[id]
	if !ok {
		return model.Todo{}, model.ErrTodoNotFound
	}

	todo.Completed = !todo.Completed
	todo.UpdatedAt = time.Now()
	s.todos[id] = todo

	return todo, nil
}
