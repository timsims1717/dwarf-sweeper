package typeface

type Align int

const (
	Left = iota
	Center
	Right
	Top    = Left
	Bottom = Right
)

type Alignment struct {
	H Align
	V Align
}

var DefaultAlign = Alignment{
	H: Left,
	V: Bottom,
}

func NewAlign(h, v Align) Alignment {
	return Alignment{
		H: h,
		V: v,
	}
}