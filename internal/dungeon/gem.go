package dungeon

import (
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/internal/util"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/transform"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

type Gem struct {
	Physics   *physics.Physics
	Transform *transform.Transform
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
	g.Physics, g.Transform = util.RandomVelocity(pos, 1.0, random.Effects)
	g.Transform.Pos = pos
	g.created = true
	g.sprite = batcher.Sprites["gem_diamond"]
	g.collect = &myecs.Collectible{}
	g.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, g.Transform).
		AddComponent(myecs.Physics, g.Physics).
		AddComponent(myecs.Collision, myecs.Collider{ GroundOnly: true }).
		AddComponent(myecs.Collect, g.collect)
}

func (g *Gem) Done() bool {
	return g.done
}

func (g *Gem) Delete() {
	if g.done {
		particles.CreateRandomStaticParticles(2, 4, []string{"sparkle_1","sparkle_2","sparkle_3","sparkle_4","sparkle_5"}, g.Transform.Pos, 1.0, 1.0, 0.5)
		sfx.SoundPlayer.PlaySound("clink", 1.0)
	}
	myecs.Manager.DisposeEntity(g.entity)
}