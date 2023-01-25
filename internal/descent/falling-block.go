package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

var (
	fbAngle = pixel.V(1., 1.).Angle()
)

type FallingBlock struct {
	Transform *transform.Transform
	Physics   *physics.Physics
	Collider  *data.Collider
	Animation *reanimator.Tree
	Entity    *ecs.Entity
	Health    *data.Health
	Biome     string
}

func CreateFallingBlock(c *cave.Cave, coords world.Coords) *FallingBlock {
	tile := c.GetTileInt(coords.X, coords.Y)
	fb := &FallingBlock{}
	fb.Transform = transform.New().WithID("falling-block")
	fb.Transform.Pos = tile.Transform.Pos
	fb.Physics = physics.New()
	fb.Health = &data.Health{
		Max:    1,
		Curr:   1,
		Immune: data.EnemyImmunity,
	}
	fb.Collider = data.NewCollider(pixel.R(0., 0., world.TileSize, world.TileSize), data.ItemC)
	fb.Collider.Damage = &data.Damage{
		SourceID: fb.Transform.ID,
		//Amount:    1,
		Dazed: 2.,
		//Knockback: 8.,
		//Angle:     &fbAngle,
		Type: data.Projectile,
	}
	fb.Biome = tile.Biome
	fb.Entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, fb.Transform).
		AddComponent(myecs.Physics, fb.Physics).
		AddComponent(myecs.Collision, fb.Collider).
		AddComponent(myecs.Health, fb.Health).
		AddComponent(myecs.Update, data.NewFrameFunc(fb.Update)).
		AddComponent(myecs.Drawable, img.Batchers[constants.ParticleKey].Sprites[fmt.Sprintf("falling_block_%s", fb.Biome)]).
		AddComponent(myecs.Batch, constants.ParticleKey).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
	return fb
}

func (fb *FallingBlock) Update() bool {
	if fb.Health.Dead || fb.Health.Dazed {
		particles.BlockParticles(fb.Transform.Pos, fb.Biome)
		myecs.Manager.DisposeEntity(fb.Entity)
	}
	if fb.Physics.Grounded {
		tile := Descent.Cave.GetTile(fb.Transform.Pos)
		if tile.Type == cave.Empty {
			tile.Type = cave.Collapse
			tile.Destroyed = false
			tile.Biome = fb.Biome
			tile.Bomb = false
			tile.UpdateDetails()
			tile.UpdateSprites()
		}
		particles.BlockParticles(fb.Transform.Pos, fb.Biome)
		myecs.Manager.DisposeEntity(fb.Entity)
	}
	return false
}
