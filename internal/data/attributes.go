package data

const (
	GemRate = 1.0
)

type Attributes struct {
	GemRate float64
}

func DefaultAttr() Attributes {
	return Attributes{
		GemRate: GemRate,
	}
}