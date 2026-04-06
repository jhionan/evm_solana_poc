package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jhionan/multichain-staking/pkg/middleware"
	"github.com/stretchr/testify/assert"
)

func TestSecurityHeaders_AllHeadersPresent(t *testing.T) {
	// Inner handler that does nothing.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.SecurityHeaders(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	headers := rec.Header()

	assert.Equal(t, "nosniff", headers.Get("X-Content-Type-Options"),
		"X-Content-Type-Options must be set")
	assert.Equal(t, "DENY", headers.Get("X-Frame-Options"),
		"X-Frame-Options must be set")
	assert.Equal(t, "1; mode=block", headers.Get("X-XSS-Protection"),
		"X-XSS-Protection must be set")
	assert.NotEmpty(t, headers.Get("Strict-Transport-Security"),
		"Strict-Transport-Security must be set")
	assert.NotEmpty(t, headers.Get("Content-Security-Policy"),
		"Content-Security-Policy must be set")
	assert.Equal(t, "no-referrer", headers.Get("Referrer-Policy"),
		"Referrer-Policy must be set")
}

func TestSecurityHeaders_NextHandlerCalled(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})

	handler := middleware.SecurityHeaders(next)
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.True(t, called, "next handler must be called")
	assert.Equal(t, http.StatusNoContent, rec.Code)
}
