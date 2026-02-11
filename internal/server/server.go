package server

import (
	"net/http"

	"todo-api/internal/todo"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	addr string
	db   *pgxpool.Pool
}

func New(addr string, db *pgxpool.Pool) *Server {
	return &Server{addr: addr, db: db}
}

func (s *Server) Start() error {
	// Setup dependencies
	// repo := todo.NewMemoryRepo() // In-memory repository for testing
	repo := todo.NewPostgresRepo(s.db)
	svc := todo.NewService(repo)
	handler := todo.NewHandler(svc)

	mux := http.NewServeMux()

	// Register routes
	mux.Handle("/todos", handler)
	mux.Handle("/todos/{id}", handler)

	return http.ListenAndServe(s.addr, mux)
}
