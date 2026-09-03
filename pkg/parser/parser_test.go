package parser_test

import (
	"testing"
	"time"

	"github.com/shaebaratheon/go-log-analyzer/pkg/aggregator"
	"github.com/shaebaratheon/go-log-analyzer/pkg/parser"
	"github.com/shaebaratheon/go-log-analyzer/pkg/filter"
)

func TestParseCommonLogFormat(t *testing.T) {
	line := `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326 "http://example.com" "Mozilla/5.0"`
	p := parser.NewParser(parser.FormatNginxCombined)

	rec, err := p.ParseLine(line)
	if err != nil {
		t.Fatalf("Failed to parse valid CLF line: %v", err)
	}

	if rec.RemoteAddr != "127.0.0.1" {
		t.Errorf("Expected RemoteAddr 127.0.0.1, got %s", rec.RemoteAddr)
	}
	if rec.Method != "GET" {
		t.Errorf("Expected Method GET, got %s", rec.Method)
	}
	if rec.Path != "/apache_pb.gif" {
		t.Errorf("Expected Path /apache_pb.gif, got %s", rec.Path)
	}
	if rec.StatusCode != 200 {
		t.Errorf("Expected Status 200, got %d", rec.StatusCode)
	}
	if rec.BodyBytesSent != 2326 {
		t.Errorf("Expected Bytes 2326, got %d", rec.BodyBytesSent)
	}
	if rec.Referer != "http://example.com" {
		t.Errorf("Expected Referer http://example.com, got %s", rec.Referer)
	}
}

func TestParseJSONFormat(t *testing.T) {
	line := `{"remote_addr":"10.0.0.1","time_local":"2026-09-01T12:00:00Z","request":"POST /api/v1/orders HTTP/1.1","status":201,"bytes_sent":512,"request_time":0.045}`
	p := parser.NewParser(parser.FormatJSON)

	rec, err := p.ParseLine(line)
	if err != nil {
		t.Fatalf("Failed to parse JSON line: %v", err)
	}

	if rec.RemoteAddr != "10.0.0.1" || rec.Method != "POST" || rec.StatusCode != 201 {
		t.Errorf("JSON parse mismatch: %+v", rec)
	}
	if rec.Duration != 45*time.Millisecond {
		t.Errorf("Expected 45ms duration, got %v", rec.Duration)
	}
}

func TestAggregatorAndFilter(t *testing.T) {
	agg := aggregator.NewLogAggregator()
	flt := filter.NewFilter().WhereStatus(400, 599)

	r1 := &parser.LogRecord{Path: "/home", StatusCode: 200, BodyBytesSent: 100, Duration: 10 * time.Millisecond}
	r2 := &parser.LogRecord{Path: "/login", StatusCode: 401, BodyBytesSent: 50, Duration: 20 * time.Millisecond}
	r3 := &parser.LogRecord{Path: "/admin", StatusCode: 500, BodyBytesSent: 200, Duration: 80 * time.Millisecond}

	records := []*parser.LogRecord{r1, r2, r3}
	for _, r := range records {
		agg.Add(r)
	}

	stats := agg.Snapshot(5)
	if stats.TotalRequests != 3 {
		t.Errorf("Expected 3 requests, got %d", stats.TotalRequests)
	}
	if len(stats.TopPaths) != 3 {
		t.Errorf("Expected 3 top paths, got %d", len(stats.TopPaths))
	}

	// Filter matches only error requests
	var matchedErrors int
	for _, r := range records {
		if flt.Matches(r) {
			matchedErrors++
		}
	}
	if matchedErrors != 2 {
		t.Errorf("Expected 2 matched errors, got %d", matchedErrors)
	}
}

func BenchmarkParserThroughput(b *testing.B) {
	line := `192.168.1.1 - - [22/Jan/2026:10:00:00 +0000] "GET /api/health HTTP/1.1" 200 45 "https://corp.com" "Go-http-client"`
	p := parser.NewParser(parser.FormatNginxCombined)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = p.ParseLine(line)
	}
}
