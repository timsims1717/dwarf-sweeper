package timing

import (
	"strconv"
	"time"
)

var (
	DT     float64
	last   time.Time
	frames int
	second = time.Tick(time.Second)
	FPS    = "0"
)

func init() {
	Reset()
}

func Reset() {
	DT = 0.0
	last = time.Now()
	frames = 0
}

func Update() {
	DT = time.Since(last).Seconds()
	last = time.Now()

	frames++
	select {
	case <-second:
		FPS = strconv.Itoa(frames)
		frames = 0
	default:
	}
}