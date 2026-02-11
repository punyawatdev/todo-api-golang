package todo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) Create(ctx context.Context, todo *Todo) error {
	query := `
		INSERT INTO todos (title, description, completed)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, todo.Title, todo.Description, todo.Completed).
		Scan(&todo.ID, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return fmt.Errorf("create todo: %w", err)
	}

	return nil
}

func (r *postgresRepo) GetByID(ctx context.Context, id int) (*Todo, error) {
	query := `
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE id = $1
	`

	var todo Todo
	err := r.db.QueryRow(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get todo: %w", err)
	}

	return &todo, nil
}

func (r *postgresRepo) List(ctx context.Context) ([]*Todo, error) {
	query := `
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list todos: %w", err)
	}
	defer rows.Close()

	var todos []*Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan todo: %w", err)
		}
		todos = append(todos, &todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return todos, nil
}

func (r *postgresRepo) Update(ctx context.Context, todo *Todo) error {
	query := `
		UPDATE todos
		SET title = $1, description = $2, completed = $3
		WHERE id = $4
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.ID,
	).Scan(&todo.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("update todo: %w", err)
	}

	return nil
}

func (r *postgresRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM todos WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete todo: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
