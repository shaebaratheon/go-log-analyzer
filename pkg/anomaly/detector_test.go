package anomaly_test

import (
	"testing"
	"time"

	"go-log-analyzer/pkg/anomaly"
)

func TestAnomalyDetector_ThresholdBreach(t *testing.T) {
	window := 100 * time.Millisecond
	threshold := 3
	detector := anomaly.NewAnomalyDetector(window, threshold)

	now := time.Now()
	if detector.RecordError(now) {
		t.Fatal("expected no breach at 1 error")
	}
	if detector.RecordError(now.Add(20 * time.Millisecond)) {
		t.Fatal("expected no breach at 2 errors")
	}
	if !detector.RecordError(now.Add(40 * time.Millisecond)) {
		t.Fatal("expected anomaly breach at 3 errors")
	}

	// Wait past window
	afterWindow := now.Add(150 * time.Millisecond)
	if count := detector.CurrentCount(afterWindow); count != 0 {
		t.Fatalf("expected 0 errors after window, got %d", count)
	}
}
