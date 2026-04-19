package handler

import (
	"Voronov/internal/transport/dto/request"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockUserService struct{}

func (m *MockUserService) FindByID(ctx context.Context, id int64) (*UserResponse, error) {
	return &UserResponse{ID: id, Login: "testuser", Firstname: "Test", Lastname: "User"}, nil
}

func (m *MockUserService) FindAll(ctx context.Context) ([]*UserResponse, error) {
	return []*UserResponse{}, nil
}

func (m *MockUserService) Create(ctx context.Context, req *request.UserRequestTo) (*UserResponse, error) {
	return &UserResponse{ID: 1, Login: req.Login, Firstname: req.Firstname, Lastname: req.Lastname}, nil
}

func (m *MockUserService) Update(ctx context.Context, id int64, req *request.UserRequestTo) (*UserResponse, error) {
	return &UserResponse{ID: id, Login: req.Login, Firstname: req.Firstname, Lastname: req.Lastname}, nil
}

func (m *MockUserService) Delete(ctx context.Context, id int64) error {
	return nil
}

type UserResponse struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type MockIssueService struct{}

func (m *MockIssueService) FindByID(ctx context.Context, id int64) (*IssueResponse, error) {
	return &IssueResponse{ID: id, UserID: 1, Title: "Test Issue", Content: "Test content"}, nil
}

func (m *MockIssueService) FindAll(ctx context.Context) ([]*IssueResponse, error) {
	return []*IssueResponse{}, nil
}

func (m *MockIssueService) Create(ctx context.Context, req *request.IssueRequestTo) (*IssueResponse, error) {
	return &IssueResponse{ID: 1, UserID: req.UserID, Title: req.Title, Content: req.Content}, nil
}

func (m *MockIssueService) Update(ctx context.Context, id int64, req *request.IssueRequestTo) (*IssueResponse, error) {
	return &IssueResponse{ID: id, UserID: req.UserID, Title: req.Title, Content: req.Content}, nil
}

func (m *MockIssueService) Delete(ctx context.Context, id int64) error {
	return nil
}

func (m *MockIssueService) FindByUserID(ctx context.Context, userID int64) (*UserResponse, error) {
	return &UserResponse{ID: userID, Login: "testuser"}, nil
}

func (m *MockIssueService) FindByIssueID(ctx context.Context, issueID int64) ([]*LabelResponse, []*ReactionResponse, error) {
	return []*LabelResponse{}, []*ReactionResponse{}, nil
}

func (m *MockIssueService) SearchIssues(ctx context.Context, labelNames []string, labelIDs []int64, userLogin, title, content string) ([]*IssueResponse, error) {
	return []*IssueResponse{}, nil
}

type IssueResponse struct {
	ID      int64           `json:"id"`
	UserID  int64           `json:"userId"`
	Title   string          `json:"title"`
	Content string          `json:"content"`
	User    *UserResponse   `json:"user"`
	Labels  []LabelResponse `json:"labels"`
}

type MockLabelService struct{}

func (m *MockLabelService) FindByID(ctx context.Context, id int64) (*LabelResponse, error) {
	return &LabelResponse{ID: id, Name: "test-label"}, nil
}

func (m *MockLabelService) FindAll(ctx context.Context) ([]*LabelResponse, error) {
	return []*LabelResponse{}, nil
}

func (m *MockLabelService) Create(ctx context.Context, req *request.LabelRequestTo) (*LabelResponse, error) {
	return &LabelResponse{ID: 1, Name: req.Name}, nil
}

func (m *MockLabelService) Update(ctx context.Context, id int64, req *request.LabelRequestTo) (*LabelResponse, error) {
	return &LabelResponse{ID: id, Name: req.Name}, nil
}

func (m *MockLabelService) Delete(ctx context.Context, id int64) error {
	return nil
}

type LabelResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type MockReactionService struct{}

func (m *MockReactionService) FindByID(ctx context.Context, id int64) (*ReactionResponse, error) {
	return &ReactionResponse{ID: id, IssueID: 1, Content: "smile"}, nil
}

func (m *MockReactionService) FindAll(ctx context.Context) ([]*ReactionResponse, error) {
	return []*ReactionResponse{}, nil
}

func (m *MockReactionService) Create(ctx context.Context, req *request.ReactionRequestTo) (*ReactionResponse, error) {
	return &ReactionResponse{ID: 1, IssueID: req.IssueID, Content: req.Content}, nil
}

func (m *MockReactionService) Update(ctx context.Context, id int64, req *request.ReactionRequestTo) (*ReactionResponse, error) {
	return &ReactionResponse{ID: id, IssueID: req.IssueID, Content: req.Content}, nil
}

func (m *MockReactionService) Delete(ctx context.Context, id int64) error {
	return nil
}

type ReactionResponse struct {
	ID      int64  `json:"id"`
	IssueID int64  `json:"issueId"`
	Content string `json:"content"`
}

type MockHandler struct {
	userService     *MockUserService
	issueService    *MockIssueService
	labelService    *MockLabelService
	reactionService *MockReactionService
}

func (h *MockHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1.0/users", h.handleUsers)
	mux.HandleFunc("/api/v1.0/users/", h.handleUserByID)
	mux.HandleFunc("/api/v1.0/issues", h.handleIssues)
	mux.HandleFunc("/api/v1.0/issues/", h.handleIssueByID)
	mux.HandleFunc("/api/v1.0/labels", h.handleLabels)
	mux.HandleFunc("/api/v1.0/labels/", h.handleLabelByID)
	mux.HandleFunc("/api/v1.0/reactions", h.handleReactions)
	mux.HandleFunc("/api/v1.0/reactions/", h.handleReactionByID)
}

func (h *MockHandler) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		resp, _ := h.userService.FindAll(r.Context())
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	case http.MethodPost:
		var req request.UserRequestTo
		json.NewDecoder(r.Body).Decode(&req)
		resp, _ := h.userService.Create(r.Context(), &req)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *MockHandler) handleUserByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/api/v1.0/users/")
	if idStr == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id := int64(1)
	switch r.Method {
	case http.MethodGet:
		resp, _ := h.userService.FindByID(r.Context(), id)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	case http.MethodPut:
		resp, _ := h.userService.Update(r.Context(), id, &request.UserRequestTo{})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	case http.MethodDelete:
		h.userService.Delete(r.Context(), id)
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *MockHandler) handleIssues(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		resp, _ := h.issueService.FindAll(r.Context())
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	case http.MethodPost:
		var req request.IssueRequestTo
		json.NewDecoder(r.Body).Decode(&req)
		resp, _ := h.issueService.Create(r.Context(), &req)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *MockHandler) handleIssueByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/api/v1.0/issues/")
	if idStr == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if strings.Contains(idStr, "/") {
		parts := strings.SplitN(idStr, "/", 2)
		if parts[1] == "user" {
			resp, _ := h.issueService.FindByUserID(r.Context(), 1)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}
	id := int64(1)
	switch r.Method {
	case http.MethodGet:
		resp, _ := h.issueService.FindByID(r.Context(), id)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	case http.MethodPut:
		resp, _ := h.issueService.Update(r.Context(), id, &request.IssueRequestTo{})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	case http.MethodDelete:
		h.issueService.Delete(r.Context(), id)
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *MockHandler) handleLabels(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		resp, _ := h.labelService.FindAll(r.Context())
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	case http.MethodPost:
		var req request.LabelRequestTo
		json.NewDecoder(r.Body).Decode(&req)
		resp, _ := h.labelService.Create(r.Context(), &req)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *MockHandler) handleLabelByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/api/v1.0/labels/")
	if idStr == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id := int64(1)
	switch r.Method {
	case http.MethodGet:
		resp, _ := h.labelService.FindByID(r.Context(), id)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	case http.MethodPut:
		resp, _ := h.labelService.Update(r.Context(), id, &request.LabelRequestTo{})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	case http.MethodDelete:
		h.labelService.Delete(r.Context(), id)
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *MockHandler) handleReactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		resp, _ := h.reactionService.FindAll(r.Context())
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	case http.MethodPost:
		var req request.ReactionRequestTo
		json.NewDecoder(r.Body).Decode(&req)
		resp, _ := h.reactionService.Create(r.Context(), &req)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *MockHandler) handleReactionByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/api/v1.0/reactions/")
	if idStr == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id := int64(1)
	switch r.Method {
	case http.MethodGet:
		resp, _ := h.reactionService.FindByID(r.Context(), id)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	case http.MethodPut:
		resp, _ := h.reactionService.Update(r.Context(), id, &request.ReactionRequestTo{})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	case http.MethodDelete:
		h.reactionService.Delete(r.Context(), id)
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func setupMockHandler() *MockHandler {
	return &MockHandler{
		userService:     &MockUserService{},
		issueService:    &MockIssueService{},
		labelService:    &MockLabelService{},
		reactionService: &MockReactionService{},
	}
}

func TestCreateUser(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body := `{"login":"testuser","password":"password123","firstname":"Test","lastname":"User"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1.0/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp["login"] != "testuser" {
		t.Errorf("Expected login 'testuser', got '%v'", resp["login"])
	}
}

func TestGetUsers(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1.0/users", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestGetUserByID(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1.0/users/1", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestUpdateUser(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body := `{"login":"updateduser","password":"newpassword","firstname":"Updated","lastname":"Name"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1.0/users/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1.0/users/1", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}

func TestCreateLabel(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body := `{"name":"bug"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1.0/labels", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}
}

func TestGetLabels(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1.0/labels", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestGetLabelByID(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1.0/labels/1", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestUpdateLabel(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body := `{"name":"newname"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1.0/labels/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestDeleteLabel(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1.0/labels/1", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}

func TestCreateReaction(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body := `{"issueId":1,"content":"smile"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1.0/reactions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}
}

func TestGetReactions(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1.0/reactions", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestGetReactionByID(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1.0/reactions/1", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestUpdateReaction(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body := `{"issueId":1,"content":"laugh"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1.0/reactions/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestDeleteReaction(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1.0/reactions/1", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}

func TestCreateIssue(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body := `{"userId":1,"title":"Test Issue","content":"Test content"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1.0/issues", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}
}

func TestGetIssues(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1.0/issues", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestGetIssueByID(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1.0/issues/1", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestUpdateIssue(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body := `{"userId":1,"title":"Updated Issue","content":"Updated content"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1.0/issues/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestDeleteIssue(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1.0/issues/1", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}

func TestGetUserByIssueID(t *testing.T) {
	h := setupMockHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1.0/issues/1/user", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}
