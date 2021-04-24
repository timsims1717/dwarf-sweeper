package util

import (
	"github.com/faiface/pixel"
	"math"
	"reflect"
)

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Min returns the smaller number between a and b.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger number between a and b.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func ContainsStr(s string, a []string) bool {
	for _, as := range a {
		if as == s {
			return true
		}
	}
	return false
}

// PointInside returns true if the pixel.Vec is inside the pixel.Rect
// when unprojected by the pixel.Matrix
func PointInside(p pixel.Vec, r pixel.Rect, m pixel.Matrix) bool {
	return r.Moved(pixel.V(-(r.W() / 2.0), -(r.H() / 2.0))).Contains(m.Unproject(p))
}

// Normalize takes a pixel.Vec and returns a normalized vector, or
// one with a magnitude of 1.0
func Normalize(p pixel.Vec) pixel.Vec {
	return p.Scaled(1 / math.Sqrt(p.X*p.X+p.Y*p.Y))
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
