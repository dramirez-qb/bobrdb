package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := NewRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestHealthRoute(t *testing.T) {
	router := NewRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "healthy", w.Body.String())
}

func TestAPIGetRoute(t *testing.T) {
	router := NewRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/key/my_value", nil)
	router.ServeHTTP(w, req)
	req, _ = http.NewRequest("GET", "/api/key", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "my_value")
}

func TestAPIPostRoute(t *testing.T) {
	router := NewRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/key/my_value", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "my_value")
}

func TestAPIDeleteRoute(t *testing.T) {
	router := NewRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/key/my_value", nil)
	router.ServeHTTP(w, req)
	req, _ = http.NewRequest("DELETE", "/api/key", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "my_value")
}
