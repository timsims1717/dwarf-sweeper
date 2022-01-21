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
