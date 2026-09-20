package domain

import "time"

// PostTestMetrics ...
type PostTestMetrics struct {
	TestID             *TestID
	TestDuration       time.Duration
	GracePeriodPercent uint8
}