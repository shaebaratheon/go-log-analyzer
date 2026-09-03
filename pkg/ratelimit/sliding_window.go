package ratelimit

import (
	"sync"
	"time"
)

type SlidingWindowLimiter struct {
	mu           sync.Mutex
	windowSize   time.Duration
	maxRequests  int
	requestTimes map[string][]time.Time
}

func NewSlidingWindowLimiter(windowSize time.Duration, maxRequests int) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		windowSize:   windowSize,
		maxRequests:  maxRequests,
		requestTimes: make(map[string][]time.Time),
	}
}

func (l *SlidingWindowLimiter) Allow(clientKey string, reqTime time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	cutoff := reqTime.Add(-l.windowSize)
	times := l.requestTimes[clientKey]

	// Prune expired entries
	validIdx := 0
	for i, t := range times {
		if t.After(cutoff) {
			validIdx = i
			break
		}
		validIdx = i + 1
	}
	if validIdx > 0 && validIdx <= len(times) {
		times = times[validIdx:]
	}

	if len(times) >= l.maxRequests {
		l.requestTimes[clientKey] = times
		return false
	}

	times = append(times, reqTime)
	l.requestTimes[clientKey] = times
	return true
}

func (l *SlidingWindowLimiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.requestTimes = make(map[string][]time.Time)
}
