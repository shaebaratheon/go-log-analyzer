package formatter

import (
	"encoding/json"
	"time"
)

type LogRecord struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

type JsonFormatter struct {
	MaskKeys map[string]bool
}

func NewJsonFormatter(maskKeys []string) *JsonFormatter {
	m := make(map[string]bool)
	for _, k := range maskKeys {
		m[k] = true
	}
	return &JsonFormatter{MaskKeys: m}
}

func (f *JsonFormatter) Format(record LogRecord) ([]byte, error) {
	if f.MaskKeys != nil && record.Fields != nil {
		sanitized := make(map[string]interface{})
		for k, v := range record.Fields {
			if f.MaskKeys[k] {
				sanitized[k] = "[MASKED]"
			} else {
				sanitized[k] = v
			}
		}
		record.Fields = sanitized
	}
	return json.Marshal(record)
}
