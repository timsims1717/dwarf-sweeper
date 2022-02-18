package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/player"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/internal/util"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"math"
)

func CreateBombItem(pos pixel.Vec) {
	e := myecs.Manager.NewEntity()
	spr := img.Batchers[constants.EntityKey].Sprites["bomb_item"]
	i := &data.Interact{
		OnInteract: func(pos pixel.Vec) bool {
			return player.AddToInventory(&player.InvItem{
				Name:   "bomb",
				Sprite: spr,
				OnUse: func() {
					tile := Descent.GetPlayerTile()
					CreateBomb(tile.Transform.Pos)
				},
				Count: 3,
				Limit: 3,
			})
		},
		Distance:   spr.Frame().W() * 0.5,
		Remove:     true,
	}
	popUp := menus.NewPopUp("{symbol:interact}: pick up")
	popUp.Dist = (spr.Frame().W() + world.TileSize) * 0.5
	phys, trans := util.RandomPosAndVel(pos, 0., 0., math.Pi*0.5, math.Pi*0.25, 5., 2., random.Effects)
	coll := data.NewCollider(pixel.R(0., 0., spr.Frame().W(), spr.Frame().H()), true, false)
	hp := &data.SimpleHealth{Immune: data.ItemImmunity}
	e.AddComponent(myecs.Transform, trans).
		AddComponent(myecs.Physics, phys).
		AddComponent(myecs.Collision, coll).
		AddComponent(myecs.Health, hp).
		AddComponent(myecs.Temp, timing.New(10.)).
		AddComponent(myecs.Interact, i).
		AddComponent(myecs.PopUp, popUp).
		AddComponent(myecs.Sprite, spr).
		AddComponent(myecs.Batch, constants.EntityKey).
		AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
			myecs.AddEffect(e, data.NewBlink(2.))
			return true
		}, 8.)).
		AddComponent(myecs.Update, data.NewFrameFunc(func() bool {
			if hp.Dead {
				CreateBomb(trans.Pos)
			}
			return false
		}))
}