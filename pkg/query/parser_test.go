package query_test

import (
	"testing"
	"go-log-analyzer/pkg/query"
)

func TestSQLParserWhereClause(t *testing.T) {
	q := "SELECT timestamp message FROM logs WHERE service = 'auth' AND status = 500"
	p := query.NewParser(q)
	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if ast.Table != "logs" {
		t.Errorf("expected table 'logs', got '%s'", ast.Table)
	}
	if len(ast.Where.Conditions) != 2 {
		t.Fatalf("expected 2 conditions, got %d", len(ast.Where.Conditions))
	}
}
