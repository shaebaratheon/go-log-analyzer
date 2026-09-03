package query

import (
	"fmt"
	"strings"
)

type Operator string

const (
	OpEqual       Operator = "="
	OpNotEqual    Operator = "!="
	OpContains    Operator = "CONTAINS"
	OpGreaterThan Operator = ">"
	OpLessThan    Operator = "<"
)

type Condition struct {
	Field    string
	Op       Operator
	Value    interface{}
}

type QueryExpression struct {
	Conditions []Condition
	Limit      int
}

func (q *QueryExpression) Match(data map[string]interface{}) bool {
	for _, cond := range q.Conditions {
		val, exists := data[cond.Field]
		if !exists {
			return false
		}

		switch cond.Op {
		case OpEqual:
			if fmt.Sprintf("%v", val) != fmt.Sprintf("%v", cond.Value) {
				return false
			}
		case OpNotEqual:
			if fmt.Sprintf("%v", val) == fmt.Sprintf("%v", cond.Value) {
				return false
			}
		case OpContains:
			valStr := fmt.Sprintf("%v", val)
			target := fmt.Sprintf("%v", cond.Value)
			if !strings.Contains(valStr, target) {
				return false
			}
		}
	}
	return true
}
