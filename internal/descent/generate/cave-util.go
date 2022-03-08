package generate

func BombLevel(level int) (float64, float64) {
	bombPMin := 0.12
	bombPMax := 0.22
	for i := 1; i < level; i++ {
		bombPMin += 0.02
		bombPMax += 0.02
	}
	if bombPMin > 0.2 {
		bombPMin = 0.2
	}
	if bombPMax > 0.3 {
		bombPMax = 0.3
	}
	return bombPMin, bombPMax
}