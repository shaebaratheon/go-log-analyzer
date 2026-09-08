package exporter

import (
	"bytes"
	"strings"
	"testing"
)

func TestCsvExporter(t *testing.T) {
	var buf bytes.Buffer
	exp := NewCsvExporter(&buf, ',')

	if err := exp.WriteHeader([]string{"time", "level", "msg"}); err != nil {
		t.Fatal(err)
	}
	if err := exp.WriteRow([]string{"12:00", "ERROR", "disk full, retry"}); err != nil {
		t.Fatal(err)
	}
	exp.Flush()

	out := buf.String()
	if !strings.Contains(out, "time,level,msg") {
		t.Errorf("missing header: %s", out)
	}
	if !strings.Contains(out, "\"disk full, retry\"") {
		t.Errorf("quotes missing for comma in message: %s", out)
	}
}
