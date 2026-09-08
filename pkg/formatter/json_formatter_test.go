package formatter

import (
	"encoding/json"
	"testing"
	"time"
)

func TestJsonFormatterMasking(t *testing.T) {
	f := NewJsonFormatter([]string{"token", "secret"})
	rec := LogRecord{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "User login",
		Fields: map[string]interface{}{
			"user":  "bob",
			"token": "secret_abc",
		},
	}
	data, err := f.Format(rec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	fields := parsed["fields"].(map[string]interface{})
	if fields["token"] != "[MASKED]" {
		t.Errorf("expected masked token, got %v", fields["token"])
	}
	if fields["user"] != "bob" {
		t.Errorf("expected user bob, got %v", fields["user"])
	}
}
