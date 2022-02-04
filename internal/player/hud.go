package player

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
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
	"strconv"
)

var (
	heartTransforms    []*transform.Transform
	lastGem            int
	gemTimer           *timing.FrameTimer
	gemTransform       *transform.Transform
	gemNumberItem      *menu.ItemText
	itemTransform      *transform.Transform
	itemCountItem      *menu.ItemText
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
		tran.Pos = pixel.V(constants.BaseW*-0.5+8.+float64(i)*1.2*world.TileSize, constants.BaseH*0.5-world.TileSize)
		heartTransforms = append(heartTransforms, tran)
	}
	gemTransform = transform.NewTransform()
	gemTransform.Anchor.H = transform.Left
	gemTransform.Anchor.V = transform.Center
	gemTransform.Scalar = pixel.V(1.6, 1.6)
	gemTransform.Pos = pixel.V(constants.BaseW*-0.5+8., constants.BaseH*0.5-(4.+2.0*world.TileSize))
	gemNumberItem = menu.NewItemText("", colornames.Aliceblue, pixel.V(1.6, 1.6), menu.Left, menu.Center)
	gemNumberItem.Transform.Pos = pixel.V(constants.BaseW*-0.5+24., constants.BaseH*0.5-(4.+2.0*world.TileSize))
	itemTransform = transform.NewTransform()
	itemTransform.Anchor.H = transform.Right
	itemTransform.Anchor.V = transform.Center
	itemTransform.Scalar = pixel.V(1.6, 1.6)
	itemTransform.Pos = pixel.V(constants.BaseW*0.5-8., constants.BaseH*0.5-world.TileSize)
	itemCountItem = menu.NewItemText("", colornames.Aliceblue, pixel.V(0.8, 0.8), menu.Right, menu.Bottom)
	itemCountItem.Transform.Pos = pixel.V(constants.BaseW*0.5, constants.BaseH*0.5-world.TileSize*1.5)
	bombCountTransform = transform.NewTransform()
	bombCountTransform.Anchor.H = transform.Right
	bombCountTransform.Anchor.V = transform.Center
	bombCountTransform.Scalar = pixel.V(1.6, 1.6)
	bombCountTransform.Pos = pixel.V(constants.BaseW*0.5-8., constants.BaseH*0.5-(4.+2.0*world.TileSize))
	bombCountItem = menu.NewItemText("", colornames.Aliceblue, pixel.V(1.6, 1.6), menu.Right, menu.Center)
	bombCountItem.Transform.Pos = pixel.V(constants.BaseW*0.5+5., constants.BaseH*0.5-(4.+2.0*world.TileSize))
	tmpAnimation = reanimator.New(reanimator.NewSwitch().
		AddAnimation(reanimator.NewAnimFromSprites("heart_temp_1", img.Batchers[constants.MenuSprites].Animations["heart_temp_1"].S, reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("heart_temp_2", img.Batchers[constants.MenuSprites].Animations["heart_temp_2"].S, reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("heart_temp_3", img.Batchers[constants.MenuSprites].Animations["heart_temp_3"].S, reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("heart_temp_4", img.Batchers[constants.MenuSprites].Animations["heart_temp_4"].S, reanimator.Hold)).
		SetChooseFn(func() int {
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
		}), "heart_temp_1")
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
			tran.Pos = pixel.V(constants.BaseW*-0.5+8.+float64(i)*1.2*world.TileSize, constants.BaseH*0.5-world.TileSize)
			heartTransforms = append(heartTransforms, tran)
		}
	}
	if descent.CaveGemsFound != lastGem {
		lastGem = descent.CaveGemsFound
		gemNumberItem.SetText(fmt.Sprintf("x%d", lastGem))
		gemTimer = timing.New(3.0)
	} else if lastGem == 0 || gemTimer == nil {
		gemTimer = timing.New(0.0)
	}
	gemTimer.Update()
	if len(descent.Inventory) > 0 && descent.InvIndex < len(descent.Inventory) {
		itemCountItem.SetText(strconv.Itoa(descent.Inventory[descent.InvIndex].Count))
	} else {
		itemCountItem.SetText("")
	}
	if descent.Descent.Type == cave.Minesweeper {
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
		bombCountTransform.Pos.X = bombCountItem.Transform.Pos.X - (bombCountItem.Text.Bounds().W()+6.)*1.6
	}
}

func DrawHUD(win *pixelgl.Window) {
	d := descent.Descent.GetPlayer()
	if d.Hovered != nil && !d.Health.Dazed {
		if d.Hovered.Solid() && d.SelectLegal {
			img.Batchers[constants.ParticleKey].GetSprite("target").Draw(win, d.Hovered.Transform.Mat)
		} else {
			img.Batchers[constants.ParticleKey].GetSprite("target_blank").Draw(win, d.Hovered.Transform.Mat)
		}
	}
	for _, tran := range heartTransforms {
		tran.UIPos = camera.Cam.APos
		tran.UIZoom = camera.Cam.GetZoomScale()
		tran.Update()
	}
	i := 0
	hp := descent.Descent.Player.Health
	for i < hp.Curr && i < len(heartTransforms) {
		img.Batchers[constants.MenuSprites].Sprites["heart_full"].Draw(win, heartTransforms[i].Mat)
		i++
	}
	for i < hp.TempHP+hp.Curr && i < len(heartTransforms) {
		tmpAnimation.Draw(win, heartTransforms[i].Mat)
		i++
	}
	for i < util.Min(hp.Max+hp.TempHP, hp.Max) && i < len(heartTransforms) {
		img.Batchers[constants.MenuSprites].Sprites["heart_empty"].Draw(win, heartTransforms[i].Mat)
		i++
	}
	if !gemTimer.Done() {
		gemTransform.UIPos = camera.Cam.APos
		gemTransform.UIZoom = camera.Cam.GetZoomScale()
		gemTransform.Update()
		gemNumberItem.Transform.UIPos = camera.Cam.APos
		gemNumberItem.Transform.UIZoom = camera.Cam.GetZoomScale()
		gemNumberItem.Update(pixel.Rect{})
		img.Batchers[constants.EntityKey].Sprites["gem_diamond"].Draw(win, gemTransform.Mat)
		gemNumberItem.Draw(win)
	}
	itemTransform.UIPos = camera.Cam.APos
	itemTransform.UIZoom = camera.Cam.GetZoomScale()
	itemTransform.Update()
	img.Batchers[constants.MenuSprites].Sprites["item_box"].Draw(win, itemTransform.Mat)
	if len(descent.Inventory) > 0 && descent.InvIndex < len(descent.Inventory) {
		item := descent.Inventory[descent.InvIndex]
		item.Sprite.Draw(win, itemTransform.Mat)
		if item.Timer != nil {
			i := 1
			for float64(i)/16. < item.Timer.Perc() {
				i++
			}
			img.Batchers[constants.MenuSprites].Sprites[fmt.Sprintf("item_timer_%d", i)].Draw(win, itemTransform.Mat)
		}
	}
	itemCountItem.Transform.UIPos = camera.Cam.APos
	itemCountItem.Transform.UIZoom = camera.Cam.GetZoomScale()
	itemCountItem.Update(pixel.Rect{})
	itemCountItem.Draw(win)
	if descent.Descent.Type == cave.Minesweeper {
		bombCountTransform.UIPos = camera.Cam.APos
		bombCountTransform.UIZoom = camera.Cam.GetZoomScale()
		bombCountTransform.Update()
		bombCountItem.Transform.UIPos = camera.Cam.APos
		bombCountItem.Transform.UIZoom = camera.Cam.GetZoomScale()
		bombCountItem.Update(pixel.Rect{})
		img.Batchers[constants.EntityKey].Sprites["bomb_fuse"].Draw(win, bombCountTransform.Mat)
		bombCountItem.Draw(win)
	}
}
