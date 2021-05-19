package dungeon

import (
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/internal/util"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/sfx"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

type Gem struct {
	Transform *physics.Physics
	created   bool
	done      bool
	collect   *myecs.Collectible
	sprite    *pixel.Sprite
	entity    *ecs.Entity
}

func (g *Gem) Update() {
	if g.created && !g.done {
		if g.collect.CollectedBy {
			GemsFound++
			g.done = true
		}
	}
}

func (g *Gem) Draw(target pixel.Target) {
	if g.created && !g.done {
		g.sprite.Draw(target, g.Transform.Mat)
	}
}

func (g *Gem) Create(pos pixel.Vec, batcher *img.Batcher) {
	g.Transform = util.RandomVelocity(pos, 1.0)
	g.Transform.Pos = pos
	g.created = true
	g.sprite = batcher.Sprites["gem_diamond"]
	g.collect = &myecs.Collectible{}
	g.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, g.Transform.Transform).
		AddComponent(myecs.Physics, g.Transform).
		AddComponent(myecs.Collision, myecs.Collider{}).
		AddComponent(myecs.Collect, g.collect)
}

func (g *Gem) Remove() bool {
	if g.done {
		myecs.Manager.DisposeEntity(g.entity)
		particles.CreateRandomStaticParticles(2, 4, []string{"sparkle_1","sparkle_2","sparkle_3","sparkle_4","sparkle_5"}, g.Transform.Pos, 1.0, 1.0, 0.5)
		sfx.SoundPlayer.PlaySound("clink", 1.0)
		return true
	}
	return false
}