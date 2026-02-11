package todo

import (
	"errors"
	"time"
)

// Domain model
type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Domain errors
var (
	ErrNotFound     = errors.New("todo not found")
	ErrInvalidInput = errors.New("invalid input")
)
