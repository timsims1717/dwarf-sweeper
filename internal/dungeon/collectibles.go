package dungeon

import (
	"dwarf-sweeper/internal/cfg"
	"dwarf-sweeper/internal/data"
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

const (
	GemDiamond = "gem_diamond"
	Beer       = "beer"
	BubbleItem = "bubble_item"
	Heart      = "heart_item"
	XRayItem   = "xray_helmet"
)

var Collectibles = map[string]*data.Collectible{}

func InitCollectibles() {
	gemSpr := img.Batchers[cfg.EntityKey].Sprites["gem_diamond"]
	beerSpr := img.Batchers[cfg.EntityKey].Sprites["beer"]
	bubbleSpr := img.Batchers[cfg.EntityKey].Sprites["bubble_item"]
	heartSpr := img.Batchers[cfg.EntityKey].Sprites["heart_item"]
	xRaySpr := img.Batchers[cfg.EntityKey].Sprites["x-ray-helmet"]
	Collectibles[GemDiamond] = &data.Collectible{
		OnCollect: func(pos pixel.Vec) bool {
			GemsFound++
			particles.CreateRandomStaticParticles(2, 4, []string{"sparkle_1","sparkle_2","sparkle_3","sparkle_4","sparkle_5"}, pos, 1.0, 1.0, 0.5)
			sfx.SoundPlayer.PlaySound("clink", 1.0)
			return true
		},
		Sprite: gemSpr,
	}
	Collectibles[Beer] = &data.Collectible{
		OnCollect: func(pos pixel.Vec) bool {
			return AddToInventory(&InvItem{
				Name:   "beer",
				Sprite: beerSpr,
				OnUse:  func() bool {
					Dungeon.Player.Entity.AddComponent(myecs.Healing, &data.Heal{
						TmpAmount: 1,
					})
					return true
				},
				Count:  1,
				Unique: false,
			})
		},
		Sprite: beerSpr,
	}
	Collectibles[BubbleItem] = &data.Collectible{
		OnCollect: func(pos pixel.Vec) bool {
			return AddToInventory(&InvItem{
				Name:   "bubble",
				Sprite: bubbleSpr,
				OnUse:  func() bool {
					if Dungeon.Player.Bubble == nil {
						bubble := &Bubble{}
						bubble.Create(pixel.Vec{})
						return true
					} else {
						return false
					}
				},
				Count:  1,
				Unique: true,
			})
		},
		Sprite: bubbleSpr,
	}
	Collectibles[Heart] = &data.Collectible{
		OnCollect: func(pos pixel.Vec) bool {
			return AddToInventory(&InvItem{
				Name:   "heart_item",
				Sprite: heartSpr,
				OnUse:  func() bool {
					if Dungeon.Player.Health.Curr < Dungeon.Player.Health.Max {
						Dungeon.Player.Entity.AddComponent(myecs.Healing, &data.Heal{
							Amount: 1,
						})
						return true
					}
					return false
				},
				Count:  1,
				Unique: false,
			})
		},
		Sprite: heartSpr,
	}
	Collectibles[XRayItem] = &data.Collectible{
		OnCollect: func(pos pixel.Vec) bool {
			return AddToInventory(&InvItem{
				Name:   "xray",
				Sprite: xRaySpr,
				OnUse:  func() bool {
					//if Dungeon.Player.Bubble == nil {
						xray := &XRayHelmet{}
						xray.Create(pixel.Vec{})
						return true
					//} else {
					//	return false
					//}
				},
				Count:  1,
				Unique: true,
			})
		},
		Sprite: xRaySpr,
	}
}

type CollectibleItem struct {
	Physics   *physics.Physics
	Transform *transform.Transform
	created   bool
	collect   *data.Collectible
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
	b.Physics, b.Transform = util.RandomVelocity(pos, 1.0, random.Effects)
	b.Transform.Pos = pos
	b.created = true
	b.sprite = b.collect.Sprite
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

func (b *CollectibleItem) Delete() {
	myecs.Manager.DisposeEntity(b.entity)
}