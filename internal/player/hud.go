package player

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/menu"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	heartTransforms    []*transform.Transform
	lastGem            int
	gemTimer           *timing.FrameTimer
	gemTransform       *transform.Transform
	gemNumberItem      *menu.ItemText
	itemTransform      *transform.Transform
	bombCountTransform *transform.Transform
	bombCountItem      *menu.ItemText
	tmpAnimation       *reanimator.Tree
)

func InitHUD() {
	heartTransforms = []*transform.Transform{}
	for i := 0; i < descent.Descent.Player.Health.Max; i++ {
		tran := transform.NewTransform()
		tran.Anchor.H = transform.Left
		tran.Anchor.V = transform.Center
		tran.Scalar = pixel.V(1.6, 1.6)
		tran.Pos = pixel.V(constants.BaseW * -0.5 + 8. + float64(i) * 1.2 * world.TileSize, constants.BaseH * 0.5 - world.TileSize)
		heartTransforms = append(heartTransforms, tran)
	}
	gemTransform = transform.NewTransform()
	gemTransform.Anchor.H = transform.Left
	gemTransform.Anchor.V = transform.Center
	gemTransform.Scalar = pixel.V(1.6, 1.6)
	gemTransform.Pos = pixel.V(constants.BaseW * -0.5 + 8., constants.BaseH * 0.5 - (4. + 2.0 * world.TileSize))
	gemNumberItem = menu.NewItemText("", colornames.Aliceblue, pixel.V(1.6, 1.6), menu.Left, menu.Center)
	gemNumberItem.Transform.Pos = pixel.V(constants.BaseW * -0.5 + 24., constants.BaseH * 0.5 - (4. + 2.0 * world.TileSize))
	itemTransform = transform.NewTransform()
	itemTransform.Anchor.H = transform.Right
	itemTransform.Anchor.V = transform.Center
	itemTransform.Scalar = pixel.V(1.6, 1.6)
	itemTransform.Pos = pixel.V(constants.BaseW * 0.5 - 8., constants.BaseH * 0.5 - world.TileSize)
	bombCountTransform = transform.NewTransform()
	bombCountTransform.Anchor.H = transform.Right
	bombCountTransform.Anchor.V = transform.Center
	bombCountTransform.Scalar = pixel.V(1.6, 1.6)
	bombCountTransform.Pos = pixel.V(constants.BaseW * 0.5 - 8., constants.BaseH * 0.5 - (4. + 2.0 * world.TileSize))
	bombCountItem = menu.NewItemText("", colornames.Aliceblue, pixel.V(1.6, 1.6), menu.Right, menu.Center)
	bombCountItem.Transform.Pos = pixel.V(constants.BaseW * 0.5 + 5., constants.BaseH * 0.5 - (4. + 2.0 * world.TileSize))
	tmpAnimation = reanimator.New(&reanimator.Switch{
		Elements: reanimator.NewElements(
			reanimator.NewAnimFromSprites("heart_temp_1", img.Batchers["entities"].Animations["heart_temp_1"].S, reanimator.Hold, nil),
			reanimator.NewAnimFromSprites("heart_temp_2", img.Batchers["entities"].Animations["heart_temp_2"].S, reanimator.Hold, nil),
			reanimator.NewAnimFromSprites("heart_temp_3", img.Batchers["entities"].Animations["heart_temp_3"].S, reanimator.Hold, nil),
			reanimator.NewAnimFromSprites("heart_temp_4", img.Batchers["entities"].Animations["heart_temp_4"].S, reanimator.Hold, nil),
		),
		Check: func() int {
			if descent.Descent.Player.Health.TempHPTimer == nil {
				return 0
			}
			perc := descent.Descent.Player.Health.TempHPTimer.Perc()
			if perc < 0.25 {
				return 0
			} else if perc < 0.5 {
				return 1
			} else if perc < 0.75 {
				return 2
			} else {
				return 3
			}
		},
	}, "heart_temp_1")
	myecs.Manager.NewEntity().AddComponent(myecs.Animation, tmpAnimation)
}

func UpdateHUD() {
	thp := descent.Descent.Player.Health.Max + descent.Descent.Player.Health.TempHP
	if len(heartTransforms) != thp {
		heartTransforms = []*transform.Transform{}
		for i := 0; i < thp; i++ {
			tran := transform.NewTransform()
			tran.Anchor.H = transform.Left
			tran.Anchor.V = transform.Center
			tran.Scalar = pixel.V(1.6, 1.6)
			tran.Pos = pixel.V(constants.BaseW * -0.5 + 8. + float64(i) * 1.2 * world.TileSize, constants.BaseH * 0.5 - world.TileSize)
			heartTransforms = append(heartTransforms, tran)
		}
	}
	if descent.GemsFound != lastGem {
		lastGem = descent.GemsFound
		gemNumberItem.SetText(fmt.Sprintf("x%d", lastGem))
		gemTimer = timing.New(3.0)
	} else if lastGem == 0 || gemTimer == nil {
		gemTimer = timing.New(0.0)
	}
	gemTimer.Update()
	if descent.Descent.Type == descent.Minesweeper {
		num := descent.CaveBombsLeft - (descent.CaveWrongMarks + descent.CaveBombsMarked)
		if num == 0 {
			if descent.CaveWrongMarks > 0 {
				bombCountItem.TextColor = colornames.Orangered
			} else {
				bombCountItem.TextColor = colornames.Forestgreen
			}
		} else {
			bombCountItem.TextColor = colornames.Aliceblue
		}
		bombCountItem.SetText(fmt.Sprintf("x%d", num))
		bombCountTransform.Pos.X = bombCountItem.Transform.Pos.X - (bombCountItem.Text.Bounds().W() + 6.) * 1.6
	}
}

func DrawHUD(win *pixelgl.Window) {
	for _, tran := range heartTransforms {
		tran.UIPos = camera.Cam.APos
		tran.UIZoom = camera.Cam.GetZoomScale()
		tran.Update()
	}
	i := 0
	hp := descent.Descent.Player.Health
	for i < hp.Curr && i < len(heartTransforms) {
		img.Batchers["entities"].Sprites["heart_full"].Draw(win, heartTransforms[i].Mat)
		i++
	}
	for i < hp.TempHP + hp.Curr && i < len(heartTransforms) {
		tmpAnimation.CurrentSprite().Draw(win, heartTransforms[i].Mat)
		i++
	}
	for i < util.Min(hp.Max + hp.TempHP, hp.Max) && i < len(heartTransforms) {
		img.Batchers["entities"].Sprites["heart_empty"].Draw(win, heartTransforms[i].Mat)
		i++
	}
	if !gemTimer.Done() {
		gemTransform.UIPos = camera.Cam.APos
		gemTransform.UIZoom = camera.Cam.GetZoomScale()
		gemTransform.Update()
		gemNumberItem.Transform.UIPos = camera.Cam.APos
		gemNumberItem.Transform.UIZoom = camera.Cam.GetZoomScale()
		gemNumberItem.Update(pixel.Rect{})
		img.Batchers["entities"].Sprites["gem_diamond"].Draw(win, gemTransform.Mat)
		gemNumberItem.Draw(win)
	}
	itemTransform.UIPos = camera.Cam.APos
	itemTransform.UIZoom = camera.Cam.GetZoomScale()
	itemTransform.Update()
	img.Batchers["entities"].Sprites["item_box"].Draw(win, itemTransform.Mat)
	if len(descent.Inventory) > 0 && descent.InvIndex < len(descent.Inventory) {
		descent.Inventory[descent.InvIndex].Sprite.Draw(win, itemTransform.Mat)
	}
	if descent.Descent.Type == descent.Minesweeper {
		bombCountTransform.UIPos = camera.Cam.APos
		bombCountTransform.UIZoom = camera.Cam.GetZoomScale()
		bombCountTransform.Update()
		bombCountItem.Transform.UIPos = camera.Cam.APos
		bombCountItem.Transform.UIZoom = camera.Cam.GetZoomScale()
		bombCountItem.Update(pixel.Rect{})
		img.Batchers["entities"].Sprites["bomb_fuse"].Draw(win, bombCountTransform.Mat)
		bombCountItem.Draw(win)
	}
}