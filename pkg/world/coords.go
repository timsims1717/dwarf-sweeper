package world

// Coords is a convenience struct for passing tile coordinates
type Coords struct {
	X int
	Y int
}

// Eq checks if a and b are equal.
func (a Coords) Eq(b Coords) bool {
	return a.X == b.X && a.Y == b.Y
}

// Neighbors returns the eight tiles surrounding the Coords, starting at the top and
// moving clockwise
func (a Coords) Neighbors() []Coords {
	return []Coords{
		{a.X, a.Y + 1},
		{a.X + 1, a.Y + 1},
		{a.X + 1, a.Y},
		{a.X + 1, a.Y - 1},
		{a.X, a.Y - 1},
		{a.X - 1, a.Y - 1},
		{a.X - 1, a.Y},
		{a.X - 1, a.Y + 1},
	}
}
