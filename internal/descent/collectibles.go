package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/internal/util"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/typeface"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
	"math"
)

const (
	GemDiamond = "gem_diamond"
	Beer       = "beer"
	BubbleItem = "bubble_item"
	Apple      = "apple"
	XRayItem   = "xray_helmet"
)

var Collectibles = map[string]*data.Collectible{}

func InitCollectibles() {
	gemSpr := img.Batchers[constants.EntityKey].Sprites["gem_diamond"]
	beerSpr := img.Batchers[constants.EntityKey].Sprites["beer"]
	bubbleSpr := img.Batchers[constants.EntityKey].Sprites["bubble_item"]
	appleSpr := img.Batchers[constants.EntityKey].Sprites["apple"]
	xRaySpr := img.Batchers[constants.EntityKey].Sprites["x-ray-helmet"]
	Collectibles[GemDiamond] = &data.Collectible{
		OnCollect: func(pos pixel.Vec) bool {
			CaveGemsFound++
			particles.CreateRandomStaticParticles(2, 4, []string{"sparkle_1","sparkle_2","sparkle_3","sparkle_4","sparkle_5"}, pos, 10.0, 1.0, 0.5)
			sfx.SoundPlayer.PlaySound("clink", 1.0)
			return true
		},
		Sprite:      gemSpr,
		AutoCollect: true,
	}
	Collectibles[Beer] = &data.Collectible{
		OnCollect: func(pos pixel.Vec) bool {
			return AddToInventory(&InvItem{
				Name:   "beer",
				Sprite: beerSpr,
				OnUse:  func() {
					Descent.Player.Entity.AddComponent(myecs.Healing, &data.Heal{
						TmpAmount: 2,
					})
				},
				Count: 1,
				Limit: 2,
			})
		},
		Sprite: beerSpr,
	}
	Collectibles[BubbleItem] = &data.Collectible{
		OnCollect: func(pos pixel.Vec) bool {
			return AddToInventory(&InvItem{
				Name:   "bubble",
				Sprite: bubbleSpr,
				OnUse:  func() {
					if Descent.Player.Bubble != nil {
						Descent.Player.Bubble.Pop()
					}
					bubble := &Bubble{}
					bubble.Create(pixel.Vec{})
				},
				Count: 1,
				Limit: 1,
				Sec:   BubbleSec,
			})
		},
		Sprite: bubbleSpr,
	}
	Collectibles[Apple] = &data.Collectible{
		OnCollect: func(pos pixel.Vec) bool {
			if Descent.Player.Health.Curr < Descent.Player.Health.Max {
				Descent.Player.Entity.AddComponent(myecs.Healing, &data.Heal{
					Amount: 1,
				})
				return true
			}
			return false
		},
		Sprite:      appleSpr,
		AutoCollect: true,
	}
	Collectibles[XRayItem] = &data.Collectible{
		OnCollect: func(_ pixel.Vec) bool {
			return AddToInventory(&InvItem{
				Name:   "xray",
				Sprite: xRaySpr,
				OnUse:  func() {
					xray := &XRayHelmet{}
					xray.Create(pixel.Vec{})
				},
				Count: 1,
				Limit: 1,
				Sec:   XRaySec,
			})
		},
		Sprite: xRaySpr,
	}
}

type CollectibleItem struct {
	Physics   *physics.Physics
	Transform *transform.Transform
	created   bool
	Collect   *data.Collectible
	sprite    *pixel.Sprite
	entity    *ecs.Entity
	health    *data.SimpleHealth
}

func (b *CollectibleItem) Update() {
	if b.health.Dead {
		b.Delete()
	}
}

func (b *CollectibleItem) Create(pos pixel.Vec) {
	b.Physics, b.Transform = util.RandomPosAndVel(pos, 0., 0., math.Pi * 0.5, math.Pi * 0.25, 5., 2., random.Effects)
	b.Transform.Pos = pos
	b.created = true
	b.sprite = b.Collect.Sprite
	b.health = &data.SimpleHealth{}
	b.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Entity, b).
		AddComponent(myecs.Transform, b.Transform).
		AddComponent(myecs.Physics, b.Physics).
		AddComponent(myecs.Collision, &data.Collider{
			Hitbox: b.Collect.Sprite.Frame(),
			GroundOnly: true,
		}).
		AddComponent(myecs.Temp, timing.New(10.)).
		AddComponent(myecs.Collect, b.Collect).
		AddComponent(myecs.Health, b.health).
		AddComponent(myecs.Sprite, b.sprite).
		AddComponent(myecs.Batch, constants.EntityKey)
	b.entity.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
		myecs.AddEffect(b.entity, data.NewBlink(2.))
		return true
	}, 8.))
	if !b.Collect.AutoCollect {
		popUp := menus.NewPopUp(fmt.Sprintf("%s to pick up", typeface.SymbolItem), nil)
		popUp.Symbols = []string{data.GameInput.FirstKey("interact")}
		popUp.Dist = (b.Collect.Sprite.Frame().W() + world.TileSize) * 0.5
		b.entity.AddComponent(myecs.PopUp, popUp)
	}
}

func (b *CollectibleItem) Delete() {
	myecs.Manager.DisposeEntity(b.entity)
}