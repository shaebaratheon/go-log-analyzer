package stats

import (
	"errors"
	"math"
	"sort"
)

type PercentileEstimator struct {
	samples []float64
}

func NewPercentileEstimator() *PercentileEstimator {
	return &PercentileEstimator{samples: make([]float64, 0)}
}

func (p *PercentileEstimator) Add(sample float64) {
	p.samples = append(p.samples, sample)
}

func (p *PercentileEstimator) Quantile(q float64) (float64, error) {
	if len(p.samples) == 0 {
		return 0, errors.New("empty sample set")
	}
	if q < 0.0 || q > 1.0 {
		return 0, errors.New("quantile must be between 0 and 1")
	}
	sorted := make([]float64, len(p.samples))
	copy(sorted, p.samples)
	sort.Float64s(sorted)

	idx := int(math.Ceil(q*float64(len(sorted)))) - 1
	if idx < 0 {
		idx = 0
	}
	return sorted[idx], nil
}
