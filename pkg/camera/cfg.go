package camera

var (
	WindowWidth   int
	WindowHeight  int
	WindowWidthF  float64
	WindowHeightF float64
)

func SetWindowSize(width, height int) {
	WindowWidth = width
	WindowHeight = height
	WindowWidthF = float64(width)
	WindowHeightF = float64(height)
}
