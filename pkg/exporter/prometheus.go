package exporter

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type CounterVec struct {
	mu     sync.RWMutex
	name   string
	help   string
	values map[string]float64
}

func NewCounterVec(name, help string) *CounterVec {
	return &CounterVec{
		name:   name,
		help:   help,
		values: make(map[string]float64),
	}
}

func (c *CounterVec) Inc(label string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[label]++
}

func (c *CounterVec) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.mu.RLock()
		defer c.mu.RUnlock()

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("# HELP %s %s\n", c.name, c.help))
		sb.WriteString(fmt.Sprintf("# TYPE %s counter\n", c.name))
		for l, val := range c.values {
			sb.WriteString(fmt.Sprintf("%s{label=\"%s\"} %f\n", c.name, l, val))
		}
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		_, _ = w.Write([]byte(sb.String()))
	}
}
