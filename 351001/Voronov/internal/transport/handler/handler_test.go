package handler

import (
	"Voronov/internal/model"
	"Voronov/internal/repository"
	"Voronov/internal/service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func setupTestHandler() *Handler {
	userRepo := repository.NewInMemoryRepository(
		func(u *model.User) int64 { return u.ID },
		func(u *model.User, id int64) { u.ID = id },
	)
	issueRepo := repository.NewInMemoryRepository(
		func(i *model.Issue) int64 { return i.ID },
		func(i *model.Issue, id int64) { i.ID = id },
	)
	labelRepo := repository.NewInMemoryRepository(
		func(l *model.Label) int64 { return l.ID },
		func(l *model.Label, id int64) { l.ID = id },
	)
	reactionRepo := repository.NewInMemoryRepository(
		func(r *model.Reaction) int64 { return r.ID },
		func(r *model.Reaction, id int64) { r.ID = id },
	)
	issueLabelRepo := repository.NewInMemoryRepository(
		func(il *model.IssueLabel) int64 { return il.IssueID*1000 + il.LabelID },
		func(il *model.IssueLabel, id int64) { il.IssueID = id / 1000; il.LabelID = id % 1000 },
	)

	mapper := service.NewMapper()

	userService := service.NewUserService(userRepo, mapper)
	issueService := service.NewIssueService(issueRepo, userRepo, labelRepo, reactionRepo, issueLabelRepo, mapper)
	labelService := service.NewLabelService(labelRepo, mapper)
	reactionService := service.NewReactionService(reactionRepo, mapper)

	return NewHandler(userService, issueService, labelService, reactionService)
}

func TestCreateUser(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body := `{"login":"testuser","password":"password123","firstname":"Test","lastname":"User"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1.0/users", bytes.NewBufferString(body))
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
	if resp["id"] == nil {
		t.Error("Expected id to be set")
	}
}

func TestGetUsers(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1.0/users", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp []interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(resp) != 0 {
		t.Errorf("Expected empty array, got %d items", len(resp))
	}
}

func TestGetUserByID(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	createBody := `{"login":"testuser","password":"password123","firstname":"Test","lastname":"User"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/users", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRR := httptest.NewRecorder()
	mux.ServeHTTP(createRR, createReq)

	var created map[string]interface{}
	json.Unmarshal(createRR.Body.Bytes(), &created)
	userID := int64(created["id"].(float64))

	req := httptest.NewRequest(http.MethodGet, "/api/v1.0/users/"+string(rune(userID+'0')), nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp["login"] != "testuser" {
		t.Errorf("Expected login 'testuser', got '%v'", resp["login"])
	}
}

func TestUpdateUser(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	createBody := `{"login":"testuser","password":"password123","firstname":"Test","lastname":"User"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/users", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRR := httptest.NewRecorder()
	mux.ServeHTTP(createRR, createReq)

	var created map[string]interface{}
	json.Unmarshal(createRR.Body.Bytes(), &created)
	userID := int64(created["id"].(float64))

	updateBody := `{"login":"updateduser","password":"newpassword","firstname":"Updated","lastname":"Name"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1.0/users/"+formatInt64(userID), bytes.NewBufferString(updateBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp["login"] != "updateduser" {
		t.Errorf("Expected login 'updateduser', got '%v'", resp["login"])
	}
}

func TestDeleteUser(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	createBody := `{"login":"testuser","password":"password123","firstname":"Test","lastname":"User"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/users", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRR := httptest.NewRecorder()
	mux.ServeHTTP(createRR, createReq)

	var created map[string]interface{}
	json.Unmarshal(createRR.Body.Bytes(), &created)
	userID := int64(created["id"].(float64))

	req := httptest.NewRequest(http.MethodDelete, "/api/v1.0/users/"+formatInt64(userID), nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}

func TestCreateLabel(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body := `{"name":"bug"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1.0/labels", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp["name"] != "bug" {
		t.Errorf("Expected name 'bug', got '%v'", resp["name"])
	}
}

func TestGetLabels(t *testing.T) {
	h := setupTestHandler()
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
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	createBody := `{"name":"feature"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/labels", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRR := httptest.NewRecorder()
	mux.ServeHTTP(createRR, createReq)

	var created map[string]interface{}
	json.Unmarshal(createRR.Body.Bytes(), &created)
	labelID := int64(created["id"].(float64))

	req := httptest.NewRequest(http.MethodGet, "/api/v1.0/labels/"+formatInt64(labelID), nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestUpdateLabel(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	createBody := `{"name":"oldname"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/labels", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRR := httptest.NewRecorder()
	mux.ServeHTTP(createRR, createReq)

	var created map[string]interface{}
	json.Unmarshal(createRR.Body.Bytes(), &created)
	labelID := int64(created["id"].(float64))

	updateBody := `{"name":"newname"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1.0/labels/"+formatInt64(labelID), bytes.NewBufferString(updateBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp["name"] != "newname" {
		t.Errorf("Expected name 'newname', got '%v'", resp["name"])
	}
}

func TestDeleteLabel(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	createBody := `{"name":"temp"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/labels", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRR := httptest.NewRecorder()
	mux.ServeHTTP(createRR, createReq)

	var created map[string]interface{}
	json.Unmarshal(createRR.Body.Bytes(), &created)
	labelID := int64(created["id"].(float64))

	req := httptest.NewRequest(http.MethodDelete, "/api/v1.0/labels/"+formatInt64(labelID), nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}

func TestCreateReaction(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	userBody := `{"login":"testuser","password":"password123","firstname":"Test","lastname":"User"}`
	userReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/users", bytes.NewBufferString(userBody))
	userReq.Header.Set("Content-Type", "application/json")
	userRR := httptest.NewRecorder()
	mux.ServeHTTP(userRR, userReq)

	var userResp map[string]interface{}
	json.Unmarshal(userRR.Body.Bytes(), &userResp)
	userID := int64(userResp["id"].(float64))

	issueBody := `{"userId":` + formatInt64(userID) + `,"title":"Test Issue","content":"Test content"}`
	issueReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/issues", bytes.NewBufferString(issueBody))
	issueReq.Header.Set("Content-Type", "application/json")
	issueRR := httptest.NewRecorder()
	mux.ServeHTTP(issueRR, issueReq)

	var issueResp map[string]interface{}
	json.Unmarshal(issueRR.Body.Bytes(), &issueResp)
	issueID := int64(issueResp["id"].(float64))

	body := `{"issueId":` + formatInt64(issueID) + `,"content":"smile"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1.0/reactions", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp["content"] != "smile" {
		t.Errorf("Expected content 'smile', got '%v'", resp["content"])
	}
}

func TestGetReactions(t *testing.T) {
	h := setupTestHandler()
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
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	userBody := `{"login":"testuser","password":"password123","firstname":"Test","lastname":"User"}`
	userReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/users", bytes.NewBufferString(userBody))
	userReq.Header.Set("Content-Type", "application/json")
	userRR := httptest.NewRecorder()
	mux.ServeHTTP(userRR, userReq)

	var userResp map[string]interface{}
	json.Unmarshal(userRR.Body.Bytes(), &userResp)
	userID := int64(userResp["id"].(float64))

	issueBody := `{"userId":` + formatInt64(userID) + `,"title":"Test Issue","content":"Test content"}`
	issueReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/issues", bytes.NewBufferString(issueBody))
	issueReq.Header.Set("Content-Type", "application/json")
	issueRR := httptest.NewRecorder()
	mux.ServeHTTP(issueRR, issueReq)

	var issueResp map[string]interface{}
	json.Unmarshal(issueRR.Body.Bytes(), &issueResp)
	issueID := int64(issueResp["id"].(float64))

	reactionBody := `{"issueId":` + formatInt64(issueID) + `,"content":"smile"}`
	reactionReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/reactions", bytes.NewBufferString(reactionBody))
	reactionReq.Header.Set("Content-Type", "application/json")
	reactionRR := httptest.NewRecorder()
	mux.ServeHTTP(reactionRR, reactionReq)

	var reactionResp map[string]interface{}
	json.Unmarshal(reactionRR.Body.Bytes(), &reactionResp)
	reactionID := int64(reactionResp["id"].(float64))

	req := httptest.NewRequest(http.MethodGet, "/api/v1.0/reactions/"+formatInt64(reactionID), nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestUpdateReaction(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	userBody := `{"login":"testuser","password":"password123","firstname":"Test","lastname":"User"}`
	userReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/users", bytes.NewBufferString(userBody))
	userReq.Header.Set("Content-Type", "application/json")
	userRR := httptest.NewRecorder()
	mux.ServeHTTP(userRR, userReq)

	var userResp map[string]interface{}
	json.Unmarshal(userRR.Body.Bytes(), &userResp)
	userID := int64(userResp["id"].(float64))

	issueBody := `{"userId":` + formatInt64(userID) + `,"title":"Test Issue","content":"Test content"}`
	issueReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/issues", bytes.NewBufferString(issueBody))
	issueReq.Header.Set("Content-Type", "application/json")
	issueRR := httptest.NewRecorder()
	mux.ServeHTTP(issueRR, issueReq)

	var issueResp map[string]interface{}
	json.Unmarshal(issueRR.Body.Bytes(), &issueResp)
	issueID := int64(issueResp["id"].(float64))

	reactionBody := `{"issueId":` + formatInt64(issueID) + `,"content":"smile"}`
	reactionReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/reactions", bytes.NewBufferString(reactionBody))
	reactionReq.Header.Set("Content-Type", "application/json")
	reactionRR := httptest.NewRecorder()
	mux.ServeHTTP(reactionRR, reactionReq)

	var reactionResp map[string]interface{}
	json.Unmarshal(reactionRR.Body.Bytes(), &reactionResp)
	reactionID := int64(reactionResp["id"].(float64))

	updateBody := `{"issueId":` + formatInt64(issueID) + `,"content":"laugh"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1.0/reactions/"+formatInt64(reactionID), bytes.NewBufferString(updateBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp["content"] != "laugh" {
		t.Errorf("Expected content 'laugh', got '%v'", resp["content"])
	}
}

func TestDeleteReaction(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	userBody := `{"login":"testuser","password":"password123","firstname":"Test","lastname":"User"}`
	userReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/users", bytes.NewBufferString(userBody))
	userReq.Header.Set("Content-Type", "application/json")
	userRR := httptest.NewRecorder()
	mux.ServeHTTP(userRR, userReq)

	var userResp map[string]interface{}
	json.Unmarshal(userRR.Body.Bytes(), &userResp)
	userID := int64(userResp["id"].(float64))

	issueBody := `{"userId":` + formatInt64(userID) + `,"title":"Test Issue","content":"Test content"}`
	issueReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/issues", bytes.NewBufferString(issueBody))
	issueReq.Header.Set("Content-Type", "application/json")
	issueRR := httptest.NewRecorder()
	mux.ServeHTTP(issueRR, issueReq)

	var issueResp map[string]interface{}
	json.Unmarshal(issueRR.Body.Bytes(), &issueResp)
	issueID := int64(issueResp["id"].(float64))

	reactionBody := `{"issueId":` + formatInt64(issueID) + `,"content":"smile"}`
	reactionReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/reactions", bytes.NewBufferString(reactionBody))
	reactionReq.Header.Set("Content-Type", "application/json")
	reactionRR := httptest.NewRecorder()
	mux.ServeHTTP(reactionRR, reactionReq)

	var reactionResp map[string]interface{}
	json.Unmarshal(reactionRR.Body.Bytes(), &reactionResp)
	reactionID := int64(reactionResp["id"].(float64))

	req := httptest.NewRequest(http.MethodDelete, "/api/v1.0/reactions/"+formatInt64(reactionID), nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}

func TestCreateIssue(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	userBody := `{"login":"testuser","password":"password123","firstname":"Test","lastname":"User"}`
	userReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/users", bytes.NewBufferString(userBody))
	userReq.Header.Set("Content-Type", "application/json")
	userRR := httptest.NewRecorder()
	mux.ServeHTTP(userRR, userReq)

	var userResp map[string]interface{}
	json.Unmarshal(userRR.Body.Bytes(), &userResp)
	userID := int64(userResp["id"].(float64))

	body := `{"userId":` + formatInt64(userID) + `,"title":"Test Issue","content":"Test content"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1.0/issues", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp["title"] != "Test Issue" {
		t.Errorf("Expected title 'Test Issue', got '%v'", resp["title"])
	}
}

func TestGetIssues(t *testing.T) {
	h := setupTestHandler()
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
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	userBody := `{"login":"testuser","password":"password123","firstname":"Test","lastname":"User"}`
	userReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/users", bytes.NewBufferString(userBody))
	userReq.Header.Set("Content-Type", "application/json")
	userRR := httptest.NewRecorder()
	mux.ServeHTTP(userRR, userReq)

	var userResp map[string]interface{}
	json.Unmarshal(userRR.Body.Bytes(), &userResp)
	userID := int64(userResp["id"].(float64))

	issueBody := `{"userId":` + formatInt64(userID) + `,"title":"Test Issue","content":"Test content"}`
	issueReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/issues", bytes.NewBufferString(issueBody))
	issueReq.Header.Set("Content-Type", "application/json")
	issueRR := httptest.NewRecorder()
	mux.ServeHTTP(issueRR, issueReq)

	var issueResp map[string]interface{}
	json.Unmarshal(issueRR.Body.Bytes(), &issueResp)
	issueID := int64(issueResp["id"].(float64))

	req := httptest.NewRequest(http.MethodGet, "/api/v1.0/issues/"+formatInt64(issueID), nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestUpdateIssue(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	userBody := `{"login":"testuser","password":"password123","firstname":"Test","lastname":"User"}`
	userReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/users", bytes.NewBufferString(userBody))
	userReq.Header.Set("Content-Type", "application/json")
	userRR := httptest.NewRecorder()
	mux.ServeHTTP(userRR, userReq)

	var userResp map[string]interface{}
	json.Unmarshal(userRR.Body.Bytes(), &userResp)
	userID := int64(userResp["id"].(float64))

	issueBody := `{"userId":` + formatInt64(userID) + `,"title":"Test Issue","content":"Test content"}`
	issueReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/issues", bytes.NewBufferString(issueBody))
	issueReq.Header.Set("Content-Type", "application/json")
	issueRR := httptest.NewRecorder()
	mux.ServeHTTP(issueRR, issueReq)

	var issueResp map[string]interface{}
	json.Unmarshal(issueRR.Body.Bytes(), &issueResp)
	issueID := int64(issueResp["id"].(float64))

	updateBody := `{"userId":` + formatInt64(userID) + `,"title":"Updated Issue","content":"Updated content"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1.0/issues/"+formatInt64(issueID), bytes.NewBufferString(updateBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp["title"] != "Updated Issue" {
		t.Errorf("Expected title 'Updated Issue', got '%v'", resp["title"])
	}
}

func TestDeleteIssue(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	userBody := `{"login":"testuser","password":"password123","firstname":"Test","lastname":"User"}`
	userReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/users", bytes.NewBufferString(userBody))
	userReq.Header.Set("Content-Type", "application/json")
	userRR := httptest.NewRecorder()
	mux.ServeHTTP(userRR, userReq)

	var userResp map[string]interface{}
	json.Unmarshal(userRR.Body.Bytes(), &userResp)
	userID := int64(userResp["id"].(float64))

	issueBody := `{"userId":` + formatInt64(userID) + `,"title":"Test Issue","content":"Test content"}`
	issueReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/issues", bytes.NewBufferString(issueBody))
	issueReq.Header.Set("Content-Type", "application/json")
	issueRR := httptest.NewRecorder()
	mux.ServeHTTP(issueRR, issueReq)

	var issueResp map[string]interface{}
	json.Unmarshal(issueRR.Body.Bytes(), &issueResp)
	issueID := int64(issueResp["id"].(float64))

	req := httptest.NewRequest(http.MethodDelete, "/api/v1.0/issues/"+formatInt64(issueID), nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}

func TestGetUserByIssueID(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	userBody := `{"login":"testuser","password":"password123","firstname":"Test","lastname":"User"}`
	userReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/users", bytes.NewBufferString(userBody))
	userReq.Header.Set("Content-Type", "application/json")
	userRR := httptest.NewRecorder()
	mux.ServeHTTP(userRR, userReq)

	var userResp map[string]interface{}
	json.Unmarshal(userRR.Body.Bytes(), &userResp)
	userID := int64(userResp["id"].(float64))

	issueBody := `{"userId":` + formatInt64(userID) + `,"title":"Test Issue","content":"Test content"}`
	issueReq := httptest.NewRequest(http.MethodPost, "/api/v1.0/issues", bytes.NewBufferString(issueBody))
	issueReq.Header.Set("Content-Type", "application/json")
	issueRR := httptest.NewRecorder()
	mux.ServeHTTP(issueRR, issueReq)

	var issueResp map[string]interface{}
	json.Unmarshal(issueRR.Body.Bytes(), &issueResp)
	issueID := int64(issueResp["id"].(float64))

	req := httptest.NewRequest(http.MethodGet, "/api/v1.0/issues/"+formatInt64(issueID)+"/user", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp["login"] != "testuser" {
		t.Errorf("Expected login 'testuser', got '%v'", resp["login"])
	}
}

func TestNotFoundError(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1.0/users/99999", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rr.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp["errorCode"] == nil {
		t.Error("Expected errorCode in response")
	}
}

func formatInt64(n int64) string {
	return strconv.FormatInt(n, 10)
}
