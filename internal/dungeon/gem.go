package dungeon

import (
	"dwarf-sweeper/internal/character"
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
	EID       int
	Physics   *physics.Physics
	Transform *transform.Transform
	created   bool
	collect   *myecs.Collectible
	sprite    *pixel.Sprite
	entity    *ecs.Entity
	health    *character.Health
}

func (g *Gem) Update() {
	if g.created {
		if g.collect.CollectedBy {
			GemsFound++
			particles.CreateRandomStaticParticles(2, 4, []string{"sparkle_1","sparkle_2","sparkle_3","sparkle_4","sparkle_5"}, g.Transform.Pos, 1.0, 1.0, 0.5)
			sfx.SoundPlayer.PlaySound("clink", 1.0)
			g.Delete()
		} else if g.health.Dead {
			g.Delete()
		}
	}
}

func (g *Gem) Create(pos pixel.Vec) {
	g.Physics, g.Transform = util.RandomVelocity(pos, 1.0, random.Effects)
	g.Transform.Pos = pos
	g.created = true
	g.sprite = img.Batchers[entityBKey].Sprites["gem_diamond"]
	g.collect = &myecs.Collectible{}
	g.health = &character.Health{
		Max:        1,
		Curr:       1,
		Override:   true,
	}
	g.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Entity, g).
		AddComponent(myecs.Transform, g.Transform).
		AddComponent(myecs.Physics, g.Physics).
		AddComponent(myecs.Collision, myecs.Collider{ GroundOnly: true }).
		AddComponent(myecs.Collect, g.collect).
		AddComponent(myecs.Health, g.health).
		AddComponent(myecs.Sprite, g.sprite).
		AddComponent(myecs.Batch, entityBKey)
	Dungeon.AddEntity(g)
}

func (g *Gem) Delete() {
	myecs.Manager.DisposeEntity(g.entity)
	Dungeon.RemoveEntity(g.EID)
}

func (g *Gem) SetId(i int) {
	g.EID = i
}