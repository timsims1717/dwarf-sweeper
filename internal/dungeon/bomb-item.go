package dungeon

import (
	"dwarf-sweeper/internal/character"
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
	EID       int
	Physics   *physics.Physics
	Transform *transform.Transform
	created   bool
	collect   *myecs.Collectible
	sprite    *pixel.Sprite
	entity    *ecs.Entity
	health    *character.Health
}

func (b *BombItem) Update() {
	if b.created {
		if b.collect.CollectedBy {
			AddToInventory(&InvItem{
				Name:   "bomb",
				Sprite: b.sprite,
				OnUse:  func() {
					tile := Dungeon.GetCave().GetTile(Dungeon.Player.Transform.Pos)
					bomb := Bomb{
						Tile: tile,
						FuseLength: BaseFuse,
					}
					bomb.Create(tile.Transform.Pos)
				},
				Count:  1,
				Unique: false,
			})
			// todo: tool/item collect sound
			//sfx.SoundPlayer.PlaySound("clink", 1.0)
			b.Delete()
		} else if b.health.Dead {
			tile := Dungeon.GetCave().GetTile(b.Transform.Pos)
			bomb := Bomb{
				Tile: tile,
				FuseLength: BaseFuse,
			}
			bomb.Create(tile.Transform.Pos)
			b.Delete()
		}
	}
}

func (b *BombItem) Create(pos pixel.Vec) {
	b.Physics, b.Transform = util.RandomVelocity(pos, 1.0, random.Effects)
	b.Transform.Pos = pos
	b.created = true
	b.sprite = img.Batchers[entityBKey].Sprites["bomb_unlit"]
	b.collect = &myecs.Collectible{}
	b.health = &character.Health{
		Max:        1,
		Curr:       1,
		Override:   true,
	}
	b.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Entity, b).
		AddComponent(myecs.Transform, b.Transform).
		AddComponent(myecs.Physics, b.Physics).
		AddComponent(myecs.Collision, myecs.Collider{ GroundOnly: true }).
		AddComponent(myecs.Collect, b.collect).
		AddComponent(myecs.Health, b.health).
		AddComponent(myecs.Sprite, b.sprite).
		AddComponent(myecs.Batch, entityBKey)
	Dungeon.AddEntity(b)
}

func (b *BombItem) Delete() {
	myecs.Manager.DisposeEntity(b.entity)
	Dungeon.RemoveEntity(b.EID)
}

func (b *BombItem) SetId(i int) {
	b.EID = i
}