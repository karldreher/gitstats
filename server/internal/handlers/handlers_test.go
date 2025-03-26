package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
A helper function to test HTTP handlers
*/
func HTTPTestHandler(req *http.Request, HandlerFunction func(w http.ResponseWriter, r *http.Request), expect int) bool {

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerFunction)

	handler.ServeHTTP(rr, req)
	if status := rr.Result().StatusCode; status == expect {
		return true
	}
	return false
}

func TestGetRoot(t *testing.T) {
	test := HTTPTestHandler(httptest.NewRequest("GET", "/", nil), GetRoot, http.StatusNotFound)
	if !test {
		t.Errorf("GetRoot failed")
	}
}
