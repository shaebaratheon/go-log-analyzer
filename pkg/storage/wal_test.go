package storage_test

import (
	"os"
	"testing"
	"time"
	"go-log-analyzer/pkg/storage"
)

func TestJournalWriteAppend(t *testing.T) {
	tmp := "/tmp/test_journal.wal"
	defer os.Remove(tmp)

	jw, err := storage.NewJournalWriter(tmp)
	if err != nil {
		t.Fatalf("journal open failed: %v", err)
	}
	defer jw.Close()

	err = jw.Append(time.Now().Unix(), []byte("sample log payload entry"))
	if err != nil {
		t.Fatalf("append failed: %v", err)
	}
}
