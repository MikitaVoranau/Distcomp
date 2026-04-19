package handler

import (
	"Voronov/internal/errors"
	"Voronov/internal/transport/dto/request"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) handleUsers(w http.ResponseWriter, r *http.Request, path string) {
	// Убираем лишние слэши в конце для стабильности
	path = strings.TrimSuffix(path, "/")

	if strings.HasPrefix(path, "/api/v1.0/users/") {
		idStr := strings.TrimPrefix(path, "/api/v1.0/users/")
		// Если после ID нет больше сегментов пути
		if !strings.Contains(idStr, "/") {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				h.writeError(w, errors.ErrBadRequest)
				return
			}
			switch r.Method {
			case http.MethodGet:
				h.getUser(w, r, id)
			case http.MethodPut:
				h.updateUser(w, r, id)
			case http.MethodDelete:
				h.deleteUser(w, r, id)
			default:
				h.writeError(w, errors.ErrNotFound)
			}
			return
		}
	}

	if path == "/api/v1.0/users" {
		switch r.Method {
		case http.MethodGet:
			h.getUsers(w, r)
		case http.MethodPost:
			h.createUser(w, r)
		default:
			h.writeError(w, errors.ErrNotFound)
		}
		return
	}
	h.writeError(w, errors.ErrNotFound)
}

// GET /api/v1.0/users/{id}
func (h *Handler) getUser(w http.ResponseWriter, r *http.Request, id int64) {
	user, err := h.userService.FindByID(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, user)
}

// GET /api/v1.0/users
func (h *Handler) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.FindAll(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, users)
}

// POST /api/v1.0/users
func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, errors.ErrBadRequest)
		return
	}
	var req request.UserRequestTo
	if err := json.Unmarshal(body, &req); err != nil {
		h.writeError(w, errors.ErrBadRequest)
		return
	}
	user, err := h.userService.Create(r.Context(), &req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusCreated, user)
}

// PUT /api/v1.0/users/{id}
func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request, id int64) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, errors.ErrBadRequest)
		return
	}
	var req request.UserRequestTo
	if err := json.Unmarshal(body, &req); err != nil {
		h.writeError(w, errors.ErrBadRequest)
		return
	}
	user, err := h.userService.Update(r.Context(), id, &req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, user)
}

// DELETE /api/v1.0/users/{id}
func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request, id int64) {
	if err := h.userService.Delete(r.Context(), id); err != nil {
		h.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
