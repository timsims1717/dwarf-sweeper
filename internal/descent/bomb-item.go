package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/internal/util"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/typeface"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
	"math"
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
		CreateBomb(b.Transform.Pos)
		b.Delete()
	}
}

func (b *BombItem) Create(pos pixel.Vec) {
	b.Physics, b.Transform = util.RandomPosAndVel(pos, 0., 0., math.Pi * 0.5, math.Pi * 0.25, 5., 2., random.Effects)
	b.Transform.Pos = pos
	b.created = true
	b.sprite = img.Batchers[constants.EntityKey].Sprites["bomb_item"]
	b.collect = &data.Collectible{
		OnCollect: func(pos pixel.Vec) bool {
			return AddToInventory(&InvItem{
				Name:   "bomb",
				Sprite: b.sprite,
				OnUse:  func() {
					tile := Descent.GetPlayerTile()
					CreateBomb(tile.Transform.Pos)
				},
				Count: 1,
				Limit: 3,
			})
		},
		Sprite: b.sprite,
	}
	b.health = &data.SimpleHealth{}
	popUp := menus.NewPopUp(fmt.Sprintf("%s to pick up", typeface.SymbolItem), nil)
	popUp.Symbols = []string{data.GameInput.FirstKey("interact")}
	popUp.Dist = (b.collect.Sprite.Frame().W() + world.TileSize) * 0.5
	b.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Entity, b).
		AddComponent(myecs.Transform, b.Transform).
		AddComponent(myecs.Physics, b.Physics).
		AddComponent(myecs.Collision, &data.Collider{
			Hitbox: b.sprite.Frame(),
			GroundOnly: true,
		}).
		AddComponent(myecs.Collect, b.collect).
		AddComponent(myecs.Health, b.health).
		AddComponent(myecs.Sprite, b.sprite).
		AddComponent(myecs.PopUp, popUp).
		AddComponent(myecs.Batch, constants.EntityKey)
}

func (b *BombItem) Delete() {
	myecs.Manager.DisposeEntity(b.entity)
}