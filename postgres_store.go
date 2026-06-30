package main

import (
	"context"
	"database/sql"
	"errors"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) ListTodos(ctx context.Context) ([]Todo, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var todo Todo
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

func (s *PostgresStore) CreateTodo(ctx context.Context, req CreateTodoRequest) (Todo, error) {
	var todo Todo
	err := s.db.QueryRowContext(ctx, `
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

func (s *PostgresStore) GetTodo(ctx context.Context, id int) (Todo, error) {
	var todo Todo
	err := s.db.QueryRowContext(ctx, `
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
		return Todo{}, errTodoNotFound
	}
	return todo, err
}

func (s *PostgresStore) UpdateTodo(ctx context.Context, id int, req UpdateTodoRequest) (Todo, error) {
	var todo Todo
	err := s.db.QueryRowContext(ctx, `
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
		return Todo{}, errTodoNotFound
	}
	return todo, err
}

func (s *PostgresStore) DeleteTodo(ctx context.Context, id int) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM todos WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errTodoNotFound
	}

	return nil
}

func (s *PostgresStore) ToggleTodoComplete(ctx context.Context, id int) (Todo, error) {
	var todo Todo
	err := s.db.QueryRowContext(ctx, `
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
		return Todo{}, errTodoNotFound
	}
	return todo, err
}
