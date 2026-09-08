package stats

import "testing"

func TestPercentileEstimator(t *testing.T) {
	pe := NewPercentileEstimator()
	for i := 1; i <= 100; i++ {
		pe.Add(float64(i))
	}

	p50, err := pe.Quantile(0.50)
	if err != nil || p50 != 50.0 {
		t.Errorf("expected 50.0, got %v", p50)
	}

	p99, err := pe.Quantile(0.99)
	if err != nil || p99 != 99.0 {
		t.Errorf("expected 99.0, got %v", p99)
	}
}
