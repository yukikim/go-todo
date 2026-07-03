package repository

import (
	"context"
	"database/sql"
	"errors"

	"go-todo/internal/model"
)

type TodoRepository interface {
	ListTodos(ctx context.Context) ([]model.Todo, error)
	CreateTodo(ctx context.Context, req model.CreateTodoRequest) (model.Todo, error)
	GetTodo(ctx context.Context, id int) (model.Todo, error)
	UpdateTodo(ctx context.Context, id int, req model.UpdateTodoRequest) (model.Todo, error)
	DeleteTodo(ctx context.Context, id int) error
	ToggleTodoComplete(ctx context.Context, id int) (model.Todo, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) ListTodos(ctx context.Context) ([]model.Todo, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []model.Todo{}
	for rows.Next() {
		var todo model.Todo
		if err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *PostgresRepository) CreateTodo(ctx context.Context, req model.CreateTodoRequest) (model.Todo, error) {
	var todo model.Todo
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO todos (title, description)
		VALUES ($1, $2)
		RETURNING id, title, description, completed, created_at, updated_at
	`, req.Title, req.Description).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	return todo, err
}

func (r *PostgresRepository) GetTodo(ctx context.Context, id int) (model.Todo, error) {
	var todo model.Todo
	err := r.db.QueryRowContext(ctx, `
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE id = $1
	`, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Todo{}, model.ErrTodoNotFound
	}
	return todo, err
}

func (r *PostgresRepository) UpdateTodo(ctx context.Context, id int, req model.UpdateTodoRequest) (model.Todo, error) {
	var todo model.Todo
	err := r.db.QueryRowContext(ctx, `
		UPDATE todos
		SET title = $1,
			description = $2,
			completed = $3,
			updated_at = NOW()
		WHERE id = $4
		RETURNING id, title, description, completed, created_at, updated_at
	`, req.Title, req.Description, req.Completed, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Todo{}, model.ErrTodoNotFound
	}
	return todo, err
}

func (r *PostgresRepository) DeleteTodo(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM todos WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return model.ErrTodoNotFound
	}

	return nil
}

func (r *PostgresRepository) ToggleTodoComplete(ctx context.Context, id int) (model.Todo, error) {
	var todo model.Todo
	err := r.db.QueryRowContext(ctx, `
		UPDATE todos
		SET completed = NOT completed,
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, title, description, completed, created_at, updated_at
	`, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Todo{}, model.ErrTodoNotFound
	}
	return todo, err
}
