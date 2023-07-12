package clock

import (
	"time"
)

type Clocker interface {
	// 現在時刻を返す
	Now() time.Time
}

// Wrapper for time.Now
type RealClocker struct{}

func (r RealClocker) Now() time.Time {
	return time.Now()
}

// for test
type FixedClocker struct{}

func (fc FixedClocker) Now() time.Time {
	return time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC)
}
