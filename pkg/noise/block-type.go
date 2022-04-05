package noise

import (
	"dwarf-sweeper/pkg/world"
	"github.com/aquilax/go-perlin"
	"math/rand"
)

const (
	alpha = 2.
	beta  = 2.
	n     = 3
)

var (
	p *perlin.Perlin
)

func Seed(rando *rand.Rand) {
	p = perlin.NewPerlin(alpha, beta, n, rando.Int63())
}

func Perlin2D(coords world.Coords) float64 {
	r := p.Noise2D(float64(coords.X)/10, float64(coords.Y)/10)
	return r
}

func Perlin1D(z int) float64 {
	r := p.Noise1D(float64(z)/10)
	return r
}