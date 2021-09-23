package dungeon

import (
	"dwarf-sweeper/internal/cfg"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/internal/util"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/transform"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

type BombItem struct {
	Physics   *physics.Physics
	Transform *transform.Transform
	created   bool
	collect   *data.Collectible
	sprite    *pixel.Sprite
	entity    *ecs.Entity
	health    *data.SimpleHealth
}

func (b *BombItem) Update() {
	if b.health.Dead {
		tile := Dungeon.GetCave().GetTile(b.Transform.Pos)
		bomb := Bomb{
			Tile: tile,
			FuseLength: BaseFuse,
		}
		bomb.Create(tile.Transform.Pos)
		b.Delete()
	}
}

func (b *BombItem) Create(pos pixel.Vec) {
	b.Physics, b.Transform = util.RandomVelocity(pos, 1.0, random.Effects)
	b.Transform.Pos = pos
	b.created = true
	b.sprite = img.Batchers[cfg.EntityKey].Sprites["bomb_unlit"]
	b.collect = &data.Collectible{
		OnCollect: func(pos pixel.Vec) bool {
			return AddToInventory(&InvItem{
				Name:   "bomb",
				Sprite: b.sprite,
				OnUse:  func() bool {
					tile := Dungeon.GetCave().GetTile(Dungeon.Player.Transform.Pos)
					bomb := Bomb{
						Tile: tile,
						FuseLength: BaseFuse,
					}
					bomb.Create(tile.Transform.Pos)
					return true
				},
				Count:  1,
				Unique: false,
			})
		},
		Sprite:    b.sprite,
	}
	b.health = &data.SimpleHealth{}
	b.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Entity, b).
		AddComponent(myecs.Transform, b.Transform).
		AddComponent(myecs.Physics, b.Physics).
		AddComponent(myecs.Collision, data.Collider{ GroundOnly: true }).
		AddComponent(myecs.Collect, b.collect).
		AddComponent(myecs.Health, b.health).
		AddComponent(myecs.Sprite, b.sprite).
		AddComponent(myecs.Batch, cfg.EntityKey)
}

func (b *BombItem) Delete() {
	myecs.Manager.DisposeEntity(b.entity)
}