package handler

import (
	"Voronov/internal/errors"
	"Voronov/internal/service"
	"encoding/json"
	std_errors "errors"
	"fmt"
	"net/http"
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

	switch {
	case strings.HasPrefix(path, "/api/v1.0/users"):
		h.handleUsers(w, r, path)
	case strings.HasPrefix(path, "/api/v1.0/issues"):
		h.handleIssues(w, r, path)
	case strings.HasPrefix(path, "/api/v1.0/labels"):
		h.handleLabels(w, r, path)
	case strings.HasPrefix(path, "/api/v1.0/reactions"):
		h.handleReactions(w, r, path)
	default:
		h.writeError(w, errors.ErrNotFound)
	}
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	// 1. Выводим ошибку в консоль сервера, чтобы не работать вслепую
	fmt.Println("❌ ОШИБКА НА СЕРВЕРЕ:", err.Error())

	var appErr *errors.AppError

	// 2. ИСПОЛЬЗУЕМ errors.As ВМЕСТО ПРЯМОГО ПРИВЕДЕНИЯ ТИПОВ
	// Это спасет, если слой service оборачивает ошибку через fmt.Errorf
	if std_errors.As(err, &appErr) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.HTTPStatus)
		json.NewEncoder(w).Encode(appErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	// 3. Возвращаем текст ошибки прямо в тест, если это не AppError
	json.NewEncoder(w).Encode(map[string]interface{}{
		"errorCode":    50001,
		"errorMessage": "Unknown Error: " + err.Error(),
	})
}
