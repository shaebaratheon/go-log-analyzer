package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

// ResponseCapture wraps http.ResponseWriter to intercept status and size.
type ResponseCapture struct {
	http.ResponseWriter
	StatusCode int
	BytesSent  int64
}

func (w *ResponseCapture) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *ResponseCapture) Write(b []byte) (int, error) {
	if w.StatusCode == 0 {
		w.StatusCode = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.BytesSent += int64(n)
	return n, err
}

// RequestLogger logs incoming HTTP requests with timing and status codes.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := fmt.Sprintf("%d-%d", time.Now().UnixNano(), r.ContentLength)
		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)

		rec := &ResponseCapture{ResponseWriter: w}
		next.ServeHTTP(rec, r.WithContext(ctx))

		duration := time.Since(start)
		log.Printf("[HTTP] id=%s method=%s path=%s status=%d duration=%v bytes=%d",
			reqID, r.Method, r.URL.Path, rec.StatusCode, duration, rec.BytesSent)
	})
}

// PanicRecovery catches unhandled panics and returns an HTTP 500 error cleanly.
func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				stack := string(debug.Stack())
				log.Printf("[PANIC RECOVERED] err=%v stack=\n%s", rec, stack)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
