package cave

import (
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/img"
	"github.com/faiface/pixel"
)

type Gem struct {
	Transform *physics.Physics
	created   bool
	done      bool
	sprite    *pixel.Sprite
}

func (g *Gem) Update() {
	if g.created && !g.done {
		g.Transform.Update()
	}
}

func (g *Gem) Draw(target pixel.Target) {
	if g.created && !g.done {
		g.sprite.Draw(target, g.Transform.Mat)
	}
}

func (g *Gem) Create(pos pixel.Vec, batcher *img.Batcher) {
	g.Transform = physics.RandomVelocity(pos, 1.0)
	g.Transform.Pos = pos
	g.created = true
	g.sprite = batcher.Sprites["gem_diamond"]
}

func (g *Gem) Remove() bool {
	return g.done
}