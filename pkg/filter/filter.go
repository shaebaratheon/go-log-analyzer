package filter

import (
	"regexp"
	"time"

	"github.com/shaebaratheon/go-log-analyzer/pkg/parser"
)

// Condition defines a predicate matching on a LogRecord.
type Condition func(*parser.LogRecord) bool

// Filter aggregates multiple conditions with logical AND/OR matching.
type Filter struct {
	conditions []Condition
}

func NewFilter() *Filter {
	return &Filter{
		conditions: make([]Condition, 0),
	}
}

func (f *Filter) WhereStatus(minCode, maxCode int) *Filter {
	f.conditions = append(f.conditions, func(r *parser.LogRecord) bool {
		return r.StatusCode >= minCode && r.StatusCode <= maxCode
	})
	return f
}

func (f *Filter) WherePathMatches(regexPattern string) (*Filter, error) {
	re, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, err
	}
	f.conditions = append(f.conditions, func(r *parser.LogRecord) bool {
		return re.MatchString(r.Path)
	})
	return f, nil
}

func (f *Filter) WhereMethod(method string) *Filter {
	f.conditions = append(f.conditions, func(r *parser.LogRecord) bool {
		return r.Method == method
	})
	return f
}

func (f *Filter) WhereDurationExceeds(threshold time.Duration) *Filter {
	f.conditions = append(f.conditions, func(r *parser.LogRecord) bool {
		return r.Duration >= threshold
	})
	return f
}

func (f *Filter) Matches(r *parser.LogRecord) bool {
	for _, c := range f.conditions {
		if !c(r) {
			return false
		}
	}
	return true
}
