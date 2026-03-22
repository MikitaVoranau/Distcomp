package handler

import (
	"Voronov/internal/errors"
	"Voronov/internal/service"
	"Voronov/internal/transport/dto/request"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	userService     service.UserService
	issueService    service.IssueService
	labelService    service.LabelService
	reactionService service.ReactionService
}

func NewHandler(
	userService service.UserService,
	issueService service.IssueService,
	labelService service.LabelService,
	reactionService service.ReactionService,
) *Handler {
	return &Handler{
		userService:     userService,
		issueService:    issueService,
		labelService:    labelService,
		reactionService: reactionService,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	api := "/api/v1.0"

	mux.HandleFunc(api+"/", h.handleAll)
}

func (h *Handler) handleAll(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if strings.HasPrefix(path, "/api/v1.0/users") {
		h.handleUsers(w, r, path)
		return
	}
	if strings.HasPrefix(path, "/api/v1.0/issues") {
		h.handleIssues(w, r, path)
		return
	}
	if strings.HasPrefix(path, "/api/v1.0/labels") {
		h.handleLabels(w, r, path)
		return
	}
	if strings.HasPrefix(path, "/api/v1.0/reactions") {
		h.handleReactions(w, r, path)
		return
	}

	h.writeError(w, errors.ErrNotFound)
}

func (h *Handler) handleUsers(w http.ResponseWriter, r *http.Request, path string) {
	if strings.HasPrefix(path, "/api/v1.0/users/") {
		idStr := strings.TrimPrefix(path, "/api/v1.0/users/")
		if idStr != "" && !strings.Contains(idStr, "/") {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				h.writeError(w, errors.ErrBadRequest)
				return
			}
			switch r.Method {
			case http.MethodGet:
				h.getUser(w, id)
				return
			case http.MethodPut:
				h.updateUser(w, r, id)
				return
			case http.MethodDelete:
				h.deleteUser(w, id)
				return
			}
		}
	}

	if path == "/api/v1.0/users" || path == "/api/v1.0/users/" {
		switch r.Method {
		case http.MethodGet:
			h.getUsers(w)
			return
		case http.MethodPost:
			h.createUser(w, r)
			return
		}
	}

	h.writeError(w, errors.ErrNotFound)
}

func (h *Handler) handleIssues(w http.ResponseWriter, r *http.Request, path string) {
	if strings.HasPrefix(path, "/api/v1.0/issues/") {
		idStr := strings.TrimPrefix(path, "/api/v1.0/issues/")
		if strings.Contains(idStr, "/") {
			parts := strings.SplitN(idStr, "/", 2)
			id, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				h.writeError(w, errors.ErrBadRequest)
				return
			}
			switch parts[1] {
			case "user":
				h.getUserByIssue(w, id)
				return
			case "labels":
				h.getLabelsByIssue(w, id)
				return
			case "reactions":
				h.getReactionsByIssue(w, id)
				return
			}
		} else {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				h.writeError(w, errors.ErrBadRequest)
				return
			}
			switch r.Method {
			case http.MethodGet:
				h.getIssue(w, id)
				return
			case http.MethodPut:
				h.updateIssue(w, r, id)
				return
			case http.MethodDelete:
				h.deleteIssue(w, id)
				return
			}
		}
	}

	if path == "/api/v1.0/issues" || path == "/api/v1.0/issues/" {
		switch r.Method {
		case http.MethodGet:
			h.getIssues(w)
			return
		case http.MethodPost:
			h.createIssue(w, r)
			return
		}
	}

	h.writeError(w, errors.ErrNotFound)
}

func (h *Handler) handleLabels(w http.ResponseWriter, r *http.Request, path string) {
	if strings.HasPrefix(path, "/api/v1.0/labels/") {
		idStr := strings.TrimPrefix(path, "/api/v1.0/labels/")
		if idStr != "" {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				h.writeError(w, errors.ErrBadRequest)
				return
			}
			switch r.Method {
			case http.MethodGet:
				h.getLabel(w, id)
				return
			case http.MethodPut:
				h.updateLabel(w, r, id)
				return
			case http.MethodDelete:
				h.deleteLabel(w, id)
				return
			}
		}
	}

	if path == "/api/v1.0/labels" || path == "/api/v1.0/labels/" {
		switch r.Method {
		case http.MethodGet:
			h.getLabels(w)
			return
		case http.MethodPost:
			h.createLabel(w, r)
			return
		}
	}

	h.writeError(w, errors.ErrNotFound)
}

func (h *Handler) handleReactions(w http.ResponseWriter, r *http.Request, path string) {
	if strings.HasPrefix(path, "/api/v1.0/reactions/") {
		idStr := strings.TrimPrefix(path, "/api/v1.0/reactions/")
		if idStr != "" {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				h.writeError(w, errors.ErrBadRequest)
				return
			}
			switch r.Method {
			case http.MethodGet:
				h.getReaction(w, id)
				return
			case http.MethodPut:
				h.updateReaction(w, r, id)
				return
			case http.MethodDelete:
				h.deleteReaction(w, id)
				return
			}
		}
	}

	if path == "/api/v1.0/reactions" || path == "/api/v1.0/reactions/" {
		switch r.Method {
		case http.MethodGet:
			h.getReactions(w)
			return
		case http.MethodPost:
			h.createReaction(w, r)
			return
		}
	}

	h.writeError(w, errors.ErrNotFound)
}

func (h *Handler) getUser(w http.ResponseWriter, id int64) {
	user, err := h.userService.FindByID(id)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, user)
}

func (h *Handler) getUsers(w http.ResponseWriter) {
	users, err := h.userService.FindAll()
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, users)
}

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
	user, err := h.userService.Create(&req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusCreated, user)
}

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
	user, err := h.userService.Update(id, &req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, user)
}

func (h *Handler) deleteUser(w http.ResponseWriter, id int64) {
	if err := h.userService.Delete(id); err != nil {
		h.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getIssue(w http.ResponseWriter, id int64) {
	issue, err := h.issueService.FindByID(id)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, issue)
}

func (h *Handler) getIssues(w http.ResponseWriter) {
	issues, err := h.issueService.FindAll()
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, issues)
}

func (h *Handler) createIssue(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, errors.ErrBadRequest)
		return
	}
	var req request.IssueRequestTo
	if err := json.Unmarshal(body, &req); err != nil {
		h.writeError(w, errors.ErrBadRequest)
		return
	}
	issue, err := h.issueService.Create(&req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusCreated, issue)
}

func (h *Handler) updateIssue(w http.ResponseWriter, r *http.Request, id int64) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, errors.ErrBadRequest)
		return
	}
	var req request.IssueRequestTo
	if err := json.Unmarshal(body, &req); err != nil {
		h.writeError(w, errors.ErrBadRequest)
		return
	}
	issue, err := h.issueService.Update(id, &req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, issue)
}

func (h *Handler) deleteIssue(w http.ResponseWriter, id int64) {
	if err := h.issueService.Delete(id); err != nil {
		h.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getUserByIssue(w http.ResponseWriter, issueID int64) {
	user, err := h.issueService.FindByUserID(issueID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, user)
}

func (h *Handler) getLabelsByIssue(w http.ResponseWriter, issueID int64) {
	labels, _, err := h.issueService.FindByIssueID(issueID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, labels)
}

func (h *Handler) getReactionsByIssue(w http.ResponseWriter, issueID int64) {
	_, reactions, err := h.issueService.FindByIssueID(issueID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, reactions)
}

func (h *Handler) getLabel(w http.ResponseWriter, id int64) {
	label, err := h.labelService.FindByID(id)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, label)
}

func (h *Handler) getLabels(w http.ResponseWriter) {
	labels, err := h.labelService.FindAll()
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, labels)
}

func (h *Handler) createLabel(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, errors.ErrBadRequest)
		return
	}
	var req request.LabelRequestTo
	if err := json.Unmarshal(body, &req); err != nil {
		h.writeError(w, errors.ErrBadRequest)
		return
	}
	label, err := h.labelService.Create(&req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusCreated, label)
}

func (h *Handler) updateLabel(w http.ResponseWriter, r *http.Request, id int64) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, errors.ErrBadRequest)
		return
	}
	var req request.LabelRequestTo
	if err := json.Unmarshal(body, &req); err != nil {
		h.writeError(w, errors.ErrBadRequest)
		return
	}
	label, err := h.labelService.Update(id, &req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, label)
}

func (h *Handler) deleteLabel(w http.ResponseWriter, id int64) {
	if err := h.labelService.Delete(id); err != nil {
		h.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getReaction(w http.ResponseWriter, id int64) {
	reaction, err := h.reactionService.FindByID(id)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, reaction)
}

func (h *Handler) getReactions(w http.ResponseWriter) {
	reactions, err := h.reactionService.FindAll()
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, reactions)
}

func (h *Handler) createReaction(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, errors.ErrBadRequest)
		return
	}
	var req request.ReactionRequestTo
	if err := json.Unmarshal(body, &req); err != nil {
		h.writeError(w, errors.ErrBadRequest)
		return
	}
	reaction, err := h.reactionService.Create(&req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusCreated, reaction)
}

func (h *Handler) updateReaction(w http.ResponseWriter, r *http.Request, id int64) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, errors.ErrBadRequest)
		return
	}
	var req request.ReactionRequestTo
	if err := json.Unmarshal(body, &req); err != nil {
		h.writeError(w, errors.ErrBadRequest)
		return
	}
	reaction, err := h.reactionService.Update(id, &req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, reaction)
}

func (h *Handler) deleteReaction(w http.ResponseWriter, id int64) {
	if err := h.reactionService.Delete(id); err != nil {
		h.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.HTTPStatus)
		json.NewEncoder(w).Encode(appErr)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(errors.ErrInternal)
}
