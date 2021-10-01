package data

type DigMode int

const (
	Either = iota
	Movement
	Dedicated
)

func (dm DigMode) String() string {
	switch dm {
	case Either:
		return "Either"
	case Movement:
		return "Movement"
	case Dedicated:
		return "Dedicated"
	default:
		return ""
	}
}