package handlers

import (
	"net/http"
	"net/http/httptest"
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

func TestReadyz_NotReady(t *testing.T) {
	// pollReady is false at test start — no poller running in unit tests
	if !HTTPTestHandler(httptest.NewRequest("GET", "/readyz", nil), Readyz, http.StatusServiceUnavailable) {
		t.Errorf("Readyz should return 503 before first poll")
	}
}
