package util

import (
	"time"
)

func SinceMicros(start time.Time) time.Duration {
	return time.Since(start).Round(time.Microsecond)
}

func SinceMillis(start time.Time) time.Duration {
	return time.Since(start).Round(time.Millisecond)
}

func SinceSecs(start time.Time) time.Duration {
	return time.Since(start).Round(time.Second)
}
