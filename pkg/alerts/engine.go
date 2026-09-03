package alerts

import (
	"fmt"
	"time"

	"github.com/shaebaratheon/go-log-analyzer/pkg/aggregator"
)

type AlertSeverity string

const (
	SeverityInfo     AlertSeverity = "INFO"
	SeverityWarning  AlertSeverity = "WARNING"
	SeverityCritical AlertSeverity = "CRITICAL"
)

type AlertRule struct {
	Name             string
	MaxErrorRate     float64
	MaxP99Duration   time.Duration
	MinTotalRequests int64
}

type AlertNotification struct {
	RuleName  string        `json:"rule_name"`
	Severity  AlertSeverity `json:"severity"`
	Message   string        `json:"message"`
	Timestamp time.Time     `json:"timestamp"`
}

type AlertEngine struct {
	rules []AlertRule
}

func NewAlertEngine() *AlertEngine {
	return &AlertEngine{rules: make([]AlertRule, 0)}
}

func (e *AlertEngine) AddRule(rule AlertRule) {
	e.rules = append(e.rules, rule)
}

func (e *AlertEngine) Evaluate(stats *aggregator.Statistics) []AlertNotification {
	var alerts []AlertNotification
	now := time.Now()

	for _, rule := range e.rules {
		if stats.TotalRequests < rule.MinTotalRequests {
			continue
		}

		if rule.MaxErrorRate > 0 && stats.ErrorRate > rule.MaxErrorRate {
			alerts = append(alerts, AlertNotification{
				RuleName:  rule.Name,
				Severity:  SeverityCritical,
				Message:   fmt.Sprintf("Error rate %.2f%% exceeds threshold %.2f%%", stats.ErrorRate*100, rule.MaxErrorRate*100),
				Timestamp: now,
			})
		}

		if rule.MaxP99Duration > 0 && stats.P99Duration > rule.MaxP99Duration {
			alerts = append(alerts, AlertNotification{
				RuleName:  rule.Name,
				Severity:  SeverityWarning,
				Message:   fmt.Sprintf("Latency P99 %v exceeds threshold %v", stats.P99Duration, rule.MaxP99Duration),
				Timestamp: now,
			})
		}
	}

	return alerts
}
