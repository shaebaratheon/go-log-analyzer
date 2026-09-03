package pipeline_test

import (
	"sync"
	"testing"
	"time"

	"go-log-analyzer/pkg/pipeline"
)

type MemorySink struct {
	mu      sync.Mutex
	entries []*pipeline.LogEntry
}

func (m *MemorySink) Write(batch []*pipeline.LogEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, batch...)
	return nil
}

func (m *MemorySink) Close() error {
	return nil
}

func TestStreamPipelineEndToEnd(t *testing.T) {
	sink := &MemorySink{}
	pipe := pipeline.NewStreamPipeline(2, 10, 50*time.Millisecond, sink)
	pipe.Start()

	for i := 0; i < 25; i++ {
		pipe.Submit(&pipeline.LogEntry{
			Timestamp: time.Now(),
			Level:     "INFO",
			Service:   "api-gateway",
			Message:   "Processed request",
		})
	}

	time.Sleep(150 * time.Millisecond)
	pipe.Stop()

	if len(sink.entries) != 25 {
		t.Fatalf("expected 25 entries written to sink, got %d", len(sink.entries))
	}
}
