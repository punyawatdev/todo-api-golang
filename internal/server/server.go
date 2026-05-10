package server

import (
	"net/http"
	"context"

	"todo-api/internal/todo"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	addr string
	db   *pgxpool.Pool
	httpServer *http.Server
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

	s.httpServer = &http.Server{
    Addr:    s.addr,
    Handler: mux,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
  return s.httpServer.Shutdown(ctx)
}