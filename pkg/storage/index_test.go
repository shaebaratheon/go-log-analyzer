package storage_test

import (
	"testing"
	"go-log-analyzer/pkg/storage"
)

func TestInvertedIndexBM25Search(t *testing.T) {
	idx := storage.NewInvertedIndex()
	idx.IndexDocument(1, "connection timeout occurred on database pool")
	idx.IndexDocument(2, "user login success on auth service")
	idx.IndexDocument(3, "database connection retry succeeded")

	res := idx.SearchBM25("database connection", 2)
	if len(res) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(res))
	}
}
