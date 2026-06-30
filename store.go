package main

import "context"

type TodoStore interface {
	ListTodos(ctx context.Context) ([]Todo, error)
	CreateTodo(ctx context.Context, req CreateTodoRequest) (Todo, error)
	GetTodo(ctx context.Context, id int) (Todo, error)
	UpdateTodo(ctx context.Context, id int, req UpdateTodoRequest) (Todo, error)
	DeleteTodo(ctx context.Context, id int) error
	ToggleTodoComplete(ctx context.Context, id int) (Todo, error)
}

var todoStore TodoStore
