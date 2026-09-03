package alerts_test

import (
	"testing"
	"time"

	"github.com/shaebaratheon/go-log-analyzer/pkg/aggregator"
	"github.com/shaebaratheon/go-log-analyzer/pkg/alerts"
	"github.com/shaebaratheon/go-log-analyzer/pkg/geodb"
	"github.com/shaebaratheon/go-log-analyzer/pkg/ratelimit"
)

func TestAlertRuleEvaluation(t *testing.T) {
	engine := alerts.NewAlertEngine()
	engine.AddRule(alerts.AlertRule{
		Name:             "HighErrorRateAlert",
		MaxErrorRate:     0.05, // 5%
		MinTotalRequests: 10,
	})

	// Case 1: Below threshold
	statsGood := &aggregator.Statistics{
		TotalRequests: 100,
		ErrorRate:     0.01,
	}
	alertsFired := engine.Evaluate(statsGood)
	if len(alertsFired) != 0 {
		t.Errorf("Expected 0 alerts, got %d", len(alertsFired))
	}

	// Case 2: Above threshold
	statsBad := &aggregator.Statistics{
		TotalRequests: 100,
		ErrorRate:     0.12,
	}
	alertsFired = engine.Evaluate(statsBad)
	if len(alertsFired) != 1 {
		t.Fatalf("Expected 1 alert, got %d", len(alertsFired))
	}
	if alertsFired[0].Severity != alerts.SeverityCritical {
		t.Errorf("Expected CRITICAL severity, got %v", alertsFired[0].Severity)
	}
}

func TestGeoDBAndRateLimit(t *testing.T) {
	trie := geodb.NewIPPrefixTrie()
	err := trie.InsertCIDR("192.168.1.0/24", &geodb.Location{CountryCode: "US", CountryName: "United States", City: "Mountain View"})
	if err != nil {
		t.Fatalf("InsertCIDR failed: %v", err)
	}

	loc := trie.Lookup("192.168.1.50")
	if loc == nil || loc.City != "Mountain View" {
		t.Errorf("Expected Mountain View, got %+v", loc)
	}

	limiter := ratelimit.NewSlidingWindowLimiter(1*time.Minute, 2)
	now := time.Now()
	if !limiter.Allow("user1", now) {
		t.Errorf("First request should be allowed")
	}
	if !limiter.Allow("user1", now.Add(1*time.Second)) {
		t.Errorf("Second request should be allowed")
	}
	if limiter.Allow("user1", now.Add(2*time.Second)) {
		t.Errorf("Third request should be denied (limit 2)")
	}
}
