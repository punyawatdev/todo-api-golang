package todo

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, todo *Todo) error
	GetByID(ctx context.Context, id int) (*Todo, error)
	List(ctx context.Context) ([]*Todo, error)
	Update(ctx context.Context, todo *Todo) error
	Delete(ctx context.Context, id int) error
}
