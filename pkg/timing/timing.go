package timing

import "time"

var (
	DT   float64
	last time.Time
)

func init() {
	Reset()
}

func Reset() {
	DT = 0.0
	last = time.Now()
}

func Update() {
	DT = time.Since(last).Seconds()
	last = time.Now()
}
