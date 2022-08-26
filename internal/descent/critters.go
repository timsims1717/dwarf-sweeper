package descent

import (
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/sfx"
	"fmt"
	"github.com/faiface/pixel"
	"math"
	"time"
)

var (
	lastSqueak time.Time
	upAngle    = math.Pi * 0.5
	digTimer   = 0.3
)

func PlaySqueak() {
	if time.Since(lastSqueak) > 1. {
		sfx.SoundPlayer.PlaySound(fmt.Sprintf("squeak%d", random.Effects.Intn(5)+1), 0.)
		lastSqueak = time.Now()
	}
}

func PlayStep(vol float64) {
	sfx.SoundPlayer.PlaySound(fmt.Sprintf("step%d", random.Effects.Intn(4)+1), vol)
}

func PlayRocks(vol float64) {
	sfx.SoundPlayer.PlaySound(fmt.Sprintf("rocks%d", random.Effects.Intn(5)+1), vol)
}

func DigParticle(pos pixel.Vec, biome string) {
	particles.CreateStaticParticle(fmt.Sprintf("dig_thru_%s", biome), pos, 0., 1.5, 0., true)
	PlayStep(-2.)
}