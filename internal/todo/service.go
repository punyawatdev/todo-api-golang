package todo

import (
	"context"
	"fmt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, title, description string) (*Todo, error) {
	if title == "" {
		return nil, fmt.Errorf("%w: title required", ErrInvalidInput)
	}

	todo := &Todo{
		Title:       title,
		Description: description,
		Completed:   false,
	}

	if err := s.repo.Create(ctx, todo); err != nil {
		return nil, fmt.Errorf("create todo: %w", err)
	}
	return todo, nil
}

func (s *Service) Get(ctx context.Context, id int) (*Todo, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]*Todo, error) {
	return s.repo.List(ctx)
}

func (s *Service) Update(ctx context.Context, id int, title, desc string, completed bool) (*Todo, error) {
	todo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if title != "" {
		todo.Title = title
	}
	if desc != "" {
		todo.Description = desc
	}
	todo.Completed = completed

	if err := s.repo.Update(ctx, todo); err != nil {
		return nil, fmt.Errorf("update todo: %w", err)
	}
	return todo, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
