package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shaebaratheon/go-log-analyzer/pkg/middleware"
)

func TestPanicRecoveryMiddleware(t *testing.T) {
	panickingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("unexpected nil pointer dereference")
	})

	recoveryMw := middleware.PanicRecovery(panickingHandler)

	req := httptest.NewRequest("GET", "/crash", nil)
	rec := httptest.NewRecorder()

	recoveryMw.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected HTTP 500 on panic, got %d", rec.Code)
	}
}

func TestRequestLoggerMiddleware(t *testing.T) {
	normalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	loggedMw := middleware.RequestLogger(normalHandler)

	req := httptest.NewRequest("GET", "/ping", nil)
	rec := httptest.NewRecorder()

	loggedMw.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected HTTP 200, got %d", rec.Code)
	}
	if rec.Body.String() != "pong" {
		t.Errorf("Expected pong, got %s", rec.Body.String())
	}
}
