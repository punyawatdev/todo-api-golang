package todo

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type createRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type updateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	switch r.Method {
	case http.MethodPost:
		h.create(ctx, w, r)
	case http.MethodGet:
		h.getOrList(ctx, w, r)
	case http.MethodPut:
		h.update(ctx, w, r)
	case http.MethodDelete:
		h.delete(ctx, w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) create(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	todo, err := h.svc.Create(ctx, req.Title, req.Description)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	respondJSON(w, http.StatusCreated, todo)
}

func (h *Handler) getOrList(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	if idStr == "" {
		idStr = r.URL.Query().Get("id")
	}

	// If no id is provided, list all todos
	if idStr == "" {
		todos, err := h.svc.List(ctx)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "internal error")
			return
		}
		respondJSON(w, http.StatusOK, todos)
		return
	}

	// Get specific todo by id
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	todo, err := h.svc.Get(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			respondError(w, http.StatusNotFound, "todo not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	respondJSON(w, http.StatusOK, todo)
}

func (h *Handler) update(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	if idStr == "" {
		idStr = r.URL.Query().Get("id")
	}

	if idStr == "" {
		respondError(w, http.StatusBadRequest, "id required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	todo, err := h.svc.Update(ctx, id, req.Title, req.Description, req.Completed)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			respondError(w, http.StatusNotFound, "todo not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	respondJSON(w, http.StatusOK, todo)
}

func (h *Handler) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	if idStr == "" {
		idStr = r.URL.Query().Get("id")
	}

	if idStr == "" {
		respondError(w, http.StatusBadRequest, "id required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrNotFound) {
			respondError(w, http.StatusNotFound, "todo not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
