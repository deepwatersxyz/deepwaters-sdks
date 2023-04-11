package util

import (
	"time"
)

type Timestamp struct {
	Time             time.Time
	MicrosSinceEpoch uint64
}
