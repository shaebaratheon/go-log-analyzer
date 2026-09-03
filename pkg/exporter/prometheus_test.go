package exporter_test

import (
	"net/http/httptest"
	"testing"
	"go-log-analyzer/pkg/exporter"
)

func TestPrometheusMetricOutput(t *testing.T) {
	c := exporter.NewCounterVec("http_requests_total", "Total processed HTTP logs")
	c.Inc("200")
	c.Inc("500")

	req := httptest.NewRequest("GET", "/metrics", nil)
	rec := httptest.NewRecorder()
	c.Handler()(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
