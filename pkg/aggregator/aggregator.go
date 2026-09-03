package aggregator

import (
	"sort"
	"sync"
	"time"

	"github.com/shaebaratheon/go-log-analyzer/pkg/parser"
)

// Statistics contains aggregated insights across processed records.
type Statistics struct {
	TotalRequests      int64             `json:"total_requests"`
	TotalBytesSent     int64             `json:"total_bytes_sent"`
	StatusDistribution map[int]int64     `json:"status_distribution"`
	MethodDistribution map[string]int64  `json:"method_distribution"`
	TopPaths           []PathMetric      `json:"top_paths"`
	TopIPs             []IPMetric        `json:"top_ips"`
	AverageDuration    time.Duration     `json:"average_duration"`
	P50Duration        time.Duration     `json:"p50_duration"`
	P90Duration        time.Duration     `json:"p90_duration"`
	P99Duration        time.Duration     `json:"p99_duration"`
	ErrorRate          float64           `json:"error_rate"`
}

type PathMetric struct {
	Path  string `json:"path"`
	Count int64  `json:"count"`
}

type IPMetric struct {
	IP    string `json:"ip"`
	Count int64  `json:"count"`
}

// LogAggregator processes streaming records and builds analytical metrics.
type LogAggregator struct {
	mu           sync.Mutex
	totalReqs    int64
	totalBytes   int64
	statusCodes  map[int]int64
	methods      map[string]int64
	paths        map[string]int64
	ips          map[string]int64
	durations    []time.Duration
}

func NewLogAggregator() *LogAggregator {
	return &LogAggregator{
		statusCodes: make(map[int]int64),
		methods:     make(map[string]int64),
		paths:       make(map[string]int64),
		ips:         make(map[string]int64),
		durations:   make([]time.Duration, 0, 1024),
	}
}

// Add ingests a single LogRecord into the aggregations.
func (a *LogAggregator) Add(rec *parser.LogRecord) {
	if rec == nil {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	a.totalReqs++
	a.totalBytes += rec.BodyBytesSent
	a.statusCodes[rec.StatusCode]++
	a.methods[rec.Method]++
	a.paths[rec.Path]++
	a.ips[rec.RemoteAddr]++

	if rec.Duration > 0 {
		a.durations = append(a.durations, rec.Duration)
	}
}

// Snapshot computes and returns current statistical aggregations.
func (a *LogAggregator) Snapshot(topK int) *Statistics {
	a.mu.Lock()
	defer a.mu.Unlock()

	stats := &Statistics{
		TotalRequests:      a.totalReqs,
		TotalBytesSent:     a.totalBytes,
		StatusDistribution: make(map[int]int64, len(a.statusCodes)),
		MethodDistribution: make(map[string]int64, len(a.methods)),
	}

	for k, v := range a.statusCodes {
		stats.StatusDistribution[k] = v
	}
	for k, v := range a.methods {
		stats.MethodDistribution[k] = v
	}

	// Calculate Top Paths
	type pathPair struct {
		p string
		c int64
	}
	pathList := make([]pathPair, 0, len(a.paths))
	for p, c := range a.paths {
		pathList = append(pathList, pathPair{p, c})
	}
	sort.Slice(pathList, func(i, j int) bool { return pathList[i].c > pathList[j].c })

	limitP := topK
	if len(pathList) < limitP {
		limitP = len(pathList)
	}
	for i := 0; i < limitP; i++ {
		stats.TopPaths = append(stats.TopPaths, PathMetric{Path: pathList[i].p, Count: pathList[i].c})
	}

	// Calculate Top IPs
	ipList := make([]pathPair, 0, len(a.ips))
	for ip, c := range a.ips {
		ipList = append(ipList, pathPair{ip, c})
	}
	sort.Slice(ipList, func(i, j int) bool { return ipList[i].c > ipList[j].c })

	limitIP := topK
	if len(ipList) < limitIP {
		limitIP = len(ipList)
	}
	for i := 0; i < limitIP; i++ {
		stats.TopIPs = append(stats.TopIPs, IPMetric{IP: ipList[i].p, Count: ipList[i].c})
	}

	// Percentiles
	if len(a.durations) > 0 {
		durCopy := make([]time.Duration, len(a.durations))
		copy(durCopy, a.durations)
		sort.Slice(durCopy, func(i, j int) bool { return durCopy[i] < durCopy[j] })

		var totalDur time.Duration
		for _, d := range durCopy {
			totalDur += d
		}
		stats.AverageDuration = totalDur / time.Duration(len(durCopy))
		stats.P50Duration = durCopy[len(durCopy)*50/100]
		stats.P90Duration = durCopy[len(durCopy)*90/100]
		stats.P99Duration = durCopy[len(durCopy)*99/100]
	}

	// Error rate calculation
	var errCount int64
	for code, count := range a.statusCodes {
		if code >= 400 {
			errCount += count
		}
	}
	if a.totalReqs > 0 {
		stats.ErrorRate = float64(errCount) / float64(a.totalReqs)
	}

	return stats
}
