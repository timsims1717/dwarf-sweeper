package timing

import (
	"math"
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

type FrameTimer struct {
	start   time.Time
	sec     float64
	elapsed float64
}

func New(sec float64) *FrameTimer {
	return &FrameTimer{
		start:   time.Now(),
		sec:     sec,
		elapsed: 0.,
	}
}

func (f *FrameTimer) UpdateDone() bool {
	if f == nil {
		return true
	}
	f.Update()
	return f.Done()
}

func (f *FrameTimer) Update() {
	if f == nil {
		return
	}
	f.elapsed += DT
}

func (f *FrameTimer) Done() bool {
	if f == nil {
		return true
	}
	return f.elapsed >= f.sec
}

func (f *FrameTimer) Elapsed() float64 {
	if f == nil {
		return 0.
	}
	return f.elapsed
}

func (f *FrameTimer) Sec() float64 {
	return f.sec
}

func (f *FrameTimer) Perc() float64 {
	if f == nil {
		return 1.
	}
	if f.sec == 0. {
		return 1.
	}
	return math.Min(f.elapsed / f.sec, 1.)
}

func (f *FrameTimer) Reset() {
	f.start = time.Now()
	f.elapsed = 0.
}