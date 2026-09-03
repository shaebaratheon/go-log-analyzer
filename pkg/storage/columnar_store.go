package storage

import (
	"sync"
	"time"
)

type ColumnarChunk struct {
	mu         sync.RWMutex
	Timestamps []int64
	Levels     []string
	Services   []string
	Messages   []string
}

func NewColumnarChunk(capacity int) *ColumnarChunk {
	return &ColumnarChunk{
		Timestamps: make([]int64, 0, capacity),
		Levels:     make([]string, 0, capacity),
		Services:   make([]string, 0, capacity),
		Messages:   make([]string, 0, capacity),
	}
}

func (c *ColumnarChunk) Append(ts time.Time, level, service, msg string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Timestamps = append(c.Timestamps, ts.UnixNano())
	c.Levels = append(c.Levels, level)
	c.Services = append(c.Services, service)
	c.Messages = append(c.Messages, msg)
}

func (c *ColumnarChunk) ScanLevelCount(targetLevel string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	count := 0
	for _, l := range c.Levels {
		if l == targetLevel {
			count++
		}
	}
	return count
}
