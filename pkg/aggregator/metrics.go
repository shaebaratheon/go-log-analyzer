package aggregator

import (
	"sync"
	"time"
)

type RateTracker struct {
	mu       sync.Mutex
	window   time.Duration
	buckets  []int64
	lastSlot int64
}

func NewRateTracker(window time.Duration, slots int) *RateTracker {
	return &RateTracker{
		window:  window,
		buckets: make([]int64, slots),
	}
}

func (r *RateTracker) Increment() {
	r.mu.Lock()
	defer r.mu.Unlock()
	slot := time.Now().Unix() % int64(len(r.buckets))
	if slot != r.lastSlot {
		r.buckets[slot] = 0
		r.lastSlot = slot
	}
	r.buckets[slot]++
}

func (r *RateTracker) Count() int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	var total int64
	for _, c := range r.buckets {
		total += c
	}
	return total
}
