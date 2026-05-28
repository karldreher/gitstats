package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func HTTPTestHandler(req *http.Request, handler func(w http.ResponseWriter, r *http.Request), expect int) bool {
	rr := httptest.NewRecorder()
	http.HandlerFunc(handler).ServeHTTP(rr, req)
	return rr.Result().StatusCode == expect
}

func TestGetRoot(t *testing.T) {
	if !HTTPTestHandler(httptest.NewRequest("GET", "/", nil), GetRoot, http.StatusNotFound) {
		t.Errorf("GetRoot failed")
	}
}

func TestHealthz(t *testing.T) {
	if !HTTPTestHandler(httptest.NewRequest("GET", "/healthz", nil), Healthz, http.StatusOK) {
		t.Errorf("Healthz failed")
	}
}

func TestReadyz(t *testing.T) {
	if !HTTPTestHandler(httptest.NewRequest("GET", "/readyz", nil), Readyz, http.StatusOK) {
		t.Errorf("Readyz failed")
	}
}

func TestPostCommit_Unauthorized(t *testing.T) {
	os.Setenv("API_KEY", "test-key")
	req := httptest.NewRequest("POST", "/api/v1/commit",
		bytes.NewBufferString(`{"commit":"feat: add thing","repo":"myrepo","author":"alice"}`))
	req.Header.Set("Content-Type", "application/json")
	// No x-api-key header — should be rejected
	if !HTTPTestHandler(req, NewPostCommit(nil), http.StatusUnauthorized) {
		t.Errorf("expected 401 without API key")
	}
}

func TestPostCommit_Conventional(t *testing.T) {
	os.Setenv("API_KEY", "test-key")
	req := httptest.NewRequest("POST", "/api/v1/commit",
		bytes.NewBufferString(`{"commit":"feat: add awesome feature","repo":"myrepo","author":"alice"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "test-key")
	if !HTTPTestHandler(req, NewPostCommit(nil), http.StatusCreated) {
		t.Errorf("expected 201 for conventional commit")
	}
}

func TestPostCommit_NonConventional(t *testing.T) {
	os.Setenv("API_KEY", "test-key")
	req := httptest.NewRequest("POST", "/api/v1/commit",
		bytes.NewBufferString(`{"commit":"WIP stuff","repo":"myrepo","author":"bob"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "test-key")
	if !HTTPTestHandler(req, NewPostCommit(nil), http.StatusCreated) {
		t.Errorf("expected 201 for non-conventional commit")
	}
}

func TestPostCommit_MissingAuthor(t *testing.T) {
	os.Setenv("API_KEY", "test-key")
	req := httptest.NewRequest("POST", "/api/v1/commit",
		bytes.NewBufferString(`{"commit":"feat: thing","repo":"myrepo"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "test-key")
	if !HTTPTestHandler(req, NewPostCommit(nil), http.StatusBadRequest) {
		t.Errorf("expected 400 for missing author")
	}
}
