package anomaly

import (
	"sync"
	"time"
)

// AnomalyDetector tracks error log frequency across a sliding time window.
type AnomalyDetector struct {
	mu           sync.Mutex
	window       time.Duration
	threshold    int
	errorTimestamps []time.Time
}

// NewAnomalyDetector creates a new sliding window error rate anomaly detector.
func NewAnomalyDetector(window time.Duration, threshold int) *AnomalyDetector {
	return &AnomalyDetector{
		window:    window,
		threshold: threshold,
	}
}

// RecordError adds an error event and returns true if an anomaly threshold is breached.
func (d *AnomalyDetector) RecordError(t time.Time) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	cutoff := t.Add(-d.window)
	filtered := make([]time.Time, 0, len(d.errorTimestamps))
	for _, ts := range d.errorTimestamps {
		if ts.After(cutoff) {
			filtered = append(filtered, ts)
		}
	}
	filtered = append(filtered, t)
	d.errorTimestamps = filtered

	return len(d.errorTimestamps) >= d.threshold
}

// CurrentCount returns active error count within the sliding window.
func (d *AnomalyDetector) CurrentCount(now time.Time) int {
	d.mu.Lock()
	defer d.mu.Unlock()

	cutoff := now.Add(-d.window)
	count := 0
	for _, ts := range d.errorTimestamps {
		if ts.After(cutoff) {
			count++
		}
	}
	return count
}
