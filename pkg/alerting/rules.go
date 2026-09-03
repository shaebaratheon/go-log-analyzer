package alerting

import (
	"fmt"
	"time"
)

type AlertRule struct {
	Name           string
	Metric         string
	Threshold      float64
	WindowDuration time.Duration
	ActionWebhook  string
}

type AlertEvent struct {
	RuleName  string
	TriggerAt time.Time
	Value     float64
	Message   string
}

type AlertEngine struct {
	rules []AlertRule
}

func NewAlertEngine(rules []AlertRule) *AlertEngine {
	return &AlertEngine{rules: rules}
}

func (a *AlertEngine) Evaluate(metricName string, currentValue float64) []*AlertEvent {
	var events []*AlertEvent
	for _, rule := range a.rules {
		if rule.Metric == metricName && currentValue >= rule.Threshold {
			events = append(events, &AlertEvent{
				RuleName:  rule.Name,
				TriggerAt: time.Now(),
				Value:     currentValue,
				Message:   fmt.Sprintf("Metric %s exceeded threshold %.2f (current: %.2f)", metricName, rule.Threshold, currentValue),
			})
		}
	}
	return events
}
