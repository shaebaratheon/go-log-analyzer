package engine

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"go-log-analyzer/pkg/pipeline"
	"go-log-analyzer/pkg/query"
)

type ExecutionEngine struct {
	records []*pipeline.LogEntry
}

func NewExecutionEngine(records []*pipeline.LogEntry) *ExecutionEngine {
	return &ExecutionEngine{records: records}
}

type GroupByResult struct {
	GroupKey string
	Count    int64
	Entries  []*pipeline.LogEntry
}

func (e *ExecutionEngine) Execute(sql *query.SQLQuery) ([]*pipeline.LogEntry, error) {
	var result []*pipeline.LogEntry
	for _, entry := range e.records {
		recordMap := map[string]interface{}{
			"timestamp": entry.Timestamp.Format(time.RFC3339),
			"level":     entry.Level,
			"service":   entry.Service,
			"message":   entry.Message,
		}
		for k, v := range entry.Fields {
			recordMap[k] = v
		}

		if sql.Where == nil || sql.Where.Match(recordMap) {
			result = append(result, entry)
		}
	}
	return result, nil
}

func (e *ExecutionEngine) AggregateByField(field string) ([]GroupByResult, error) {
	groupMap := make(map[string][]*pipeline.LogEntry)
	for _, entry := range e.records {
		var valStr string
		switch field {
		case "service":
			valStr = entry.Service
		case "level":
			valStr = entry.Level
		default:
			if v, ok := entry.Fields[field]; ok {
				valStr = fmt.Sprintf("%v", v)
			} else {
				valStr = "<unknown>"
			}
		}
		groupMap[valStr] = append(groupMap[valStr], entry)
	}

	var results []GroupByResult
	for k, v := range groupMap {
		results = append(results, GroupByResult{
			GroupKey: k,
			Count:    int64(len(v)),
			Entries:  v,
		})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Count > results[j].Count
	})
	return results, nil
}
