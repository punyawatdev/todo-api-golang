package todo

import (
	"context"
	"sync"
	"time"
)

type memoryRepo struct {
	mu     sync.RWMutex
	todos  map[int]*Todo
	nextID int
}

func NewMemoryRepo() Repository {
	return &memoryRepo{
		todos:  make(map[int]*Todo),
		nextID: 1,
	}
}

func (r *memoryRepo) Create(ctx context.Context, todo *Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	todo.ID = r.nextID
	todo.CreatedAt = now
	todo.UpdatedAt = now
	r.todos[todo.ID] = todo
	r.nextID++
	return nil
}

func (r *memoryRepo) GetByID(ctx context.Context, id int) (*Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todo, ok := r.todos[id]
	if !ok {
		return nil, ErrNotFound
	}
	return todo, nil
}

func (r *memoryRepo) List(ctx context.Context) ([]*Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todos := make([]*Todo, 0, len(r.todos))
	for _, t := range r.todos {
		todos = append(todos, t)
	}
	return todos, nil
}

func (r *memoryRepo) Update(ctx context.Context, todo *Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.todos[todo.ID]; !ok {
		return ErrNotFound
	}
	todo.UpdatedAt = time.Now()
	r.todos[todo.ID] = todo
	return nil
}

func (r *memoryRepo) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.todos[id]; !ok {
		return ErrNotFound
	}
	delete(r.todos, id)
	return nil
}
