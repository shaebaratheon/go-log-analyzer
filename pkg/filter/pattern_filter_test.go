package filter

import "testing"

func TestPatternFilter(t *testing.T) {
	pf, err := NewPatternFilter([]string{"ERROR", "WARN"}, []string{"healthcheck"})
	if err != nil {
		t.Fatal(err)
	}

	if !pf.Match("2026-01-01 ERROR disk failure") {
		t.Error("expected match for ERROR")
	}
	if pf.Match("2026-01-01 ERROR healthcheck ping") {
		t.Error("expected exclude for healthcheck")
	}
	if pf.Match("2026-01-01 INFO normal operation") {
		t.Error("expected no match for INFO")
	}
}
