package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/shaebaratheon/go-log-analyzer/pkg/aggregator"
)

// SummaryReporter handles formatted rendering of analytical stats.
type SummaryReporter struct {
	writer io.Writer
}

func NewSummaryReporter(w io.Writer) *SummaryReporter {
	return &SummaryReporter{writer: w}
}

// RenderText produces a human-readable ASCII summary report.
func (r *SummaryReporter) RenderText(stats *aggregator.Statistics) error {
	var sb strings.Builder
	sb.WriteString("====================================================\n")
	sb.WriteString("               LOG ANALYZER REPORT                  \n")
	sb.WriteString("====================================================\n")
	sb.WriteString(fmt.Sprintf("Total Requests Handled: %d\n", stats.TotalRequests))
	sb.WriteString(fmt.Sprintf("Total Volume:           %.2f MB\n", float64(stats.TotalBytesSent)/(1024*1024)))
	sb.WriteString(fmt.Sprintf("Client Error Rate:      %.2f%%\n", stats.ErrorRate*100))
	sb.WriteString(fmt.Sprintf("Latency P50 / P90 / P99: %v / %v / %v\n", stats.P50Duration, stats.P90Duration, stats.P99Duration))
	sb.WriteString("----------------------------------------------------\n")
	sb.WriteString("Status Code Distribution:\n")
	for code, count := range stats.StatusDistribution {
		pct := float64(count) / float64(stats.TotalRequests) * 100.0
		sb.WriteString(fmt.Sprintf("  [%d]: %d (%.1f%%)\n", code, count, pct))
	}
	sb.WriteString("----------------------------------------------------\n")
	sb.WriteString("Top Requested Paths:\n")
	for idx, p := range stats.TopPaths {
		sb.WriteString(fmt.Sprintf("  %d. %-30s -> %d hits\n", idx+1, p.Path, p.Count))
	}
	sb.WriteString("====================================================\n")

	_, err := r.writer.Write([]byte(sb.String()))
	return err
}

// RenderJSON serializes the analytics into pretty-printed JSON.
func (r *SummaryReporter) RenderJSON(stats *aggregator.Statistics) error {
	encoder := json.NewEncoder(r.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(stats)
}
