package engine_test

import (
	"testing"
	"time"

	"go-log-analyzer/pkg/engine"
	"go-log-analyzer/pkg/pipeline"
	"go-log-analyzer/pkg/query"
)

func TestExecutionEngineQueryAndAggregation(t *testing.T) {
	entries := []*pipeline.LogEntry{
		{Timestamp: time.Now(), Level: "ERROR", Service: "auth", Message: "Password mismatch"},
		{Timestamp: time.Now(), Level: "INFO", Service: "auth", Message: "User login ok"},
		{Timestamp: time.Now(), Level: "ERROR", Service: "payment", Message: "Gateway timeout"},
	}

	exec := engine.NewExecutionEngine(entries)
	q := &query.SQLQuery{
		Fields: []string{"service", "message"},
		Table:  "logs",
		Where: &query.QueryExpression{
			Conditions: []query.Condition{
				{Field: "level", Op: query.OpEqual, Value: "ERROR"},
			},
		},
	}

	res, err := exec.Execute(q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Errorf("expected 2 error records, got %d", len(res))
	}

	agg, err := exec.AggregateByField("service")
	if err != nil || len(agg) == 0 {
		t.Fatalf("aggregation failed")
	}
}
