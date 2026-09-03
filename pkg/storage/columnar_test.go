package storage_test

import (
	"testing"
	"time"
	"go-log-analyzer/pkg/storage"
)

func TestColumnarChunkScan(t *testing.T) {
	chunk := storage.NewColumnarChunk(100)
	chunk.Append(time.Now(), "ERROR", "svc-1", "msg1")
	chunk.Append(time.Now(), "INFO", "svc-1", "msg2")
	chunk.Append(time.Now(), "ERROR", "svc-2", "msg3")

	if chunk.ScanLevelCount("ERROR") != 2 {
		t.Fatalf("expected 2 errors")
	}
}
