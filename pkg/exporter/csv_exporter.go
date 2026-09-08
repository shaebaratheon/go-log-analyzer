package exporter

import (
	"encoding/csv"
	"io"
)

type CsvExporter struct {
	writer *csv.Writer
}

func NewCsvExporter(w io.Writer, comma rune) *CsvExporter {
	cw := csv.NewWriter(w)
	if comma != 0 {
		cw.Comma = comma
	}
	return &CsvExporter{writer: cw}
}

func (e *CsvExporter) WriteHeader(headers []string) error {
	return e.writer.Write(headers)
}

func (e *CsvExporter) WriteRow(row []string) error {
	return e.writer.Write(row)
}

func (e *CsvExporter) Flush() {
	e.writer.Flush()
}
