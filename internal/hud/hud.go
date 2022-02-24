package hud

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data/player"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/typeface"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
	"math"
	"strconv"
)

var (
	tmpAnim *reanimator.Anim
)

var (
	heartSpr        *pixel.Sprite
	//tmpHeartSpr     *pixel.Sprite
	heartEmptySpr   *pixel.Sprite
	//fullHP          bool
	heartTimer      *timing.FrameTimer
	tmpAnimation    *reanimator.Tree
	//heartTransforms []*transform.Transform

	lastGem       int
	gemSpr        *pixel.Sprite
	gemTimer      *timing.FrameTimer
	gemTransform  *transform.Transform
	gemNumberText *typeface.Text

	timerText *typeface.Text

	itemBoxSpr    *pixel.Sprite
	itemTransform *transform.Transform
	itemCountText *typeface.Text

	bombSpr       *pixel.Sprite
	bombTransform *transform.Transform
	bombCountText *typeface.Text

	hudTextColor = color.RGBA{
		R: 218,
		G: 224,
		B: 234,
		A: 255,
	}
)

func InitHUD() {
	dwarf := descent.Descent.GetPlayers()[0]
	//heartTransforms = []*transform.Transform{}
	//for i := 0; i < dwarf.Health.Max; i++ {
	//	tran := transform.New()
	//	tran.Scalar = pixel.V(camera.Cam.GetZoom(), camera.Cam.GetZoom())
	//	heartTransforms = append(heartTransforms, tran)
	//}
	heartSpr = img.Batchers[constants.MenuSprites].GetSprite("heart_full")
	//tmpHeartSpr = img.Batchers[constants.MenuSprites].GetSprite("heart_temp_1")
	heartEmptySpr = img.Batchers[constants.MenuSprites].GetSprite("heart_empty")

	//gemTransform = transform.New()
	//gemTransform.Scalar = pixel.V(camera.Cam.GetZoom(), camera.Cam.GetZoom())
	//gemNumberText = typeface.New(camera.Cam, "main", typeface.NewAlign(typeface.Left, typeface.Center), 1.0, constants.ActualOneSize, 0., 0.)
	//gemNumberText.SetColor(hudTextColor)
	gemSpr = img.Batchers[constants.EntityKey].GetSprite("gem_diamond")

	timerText = typeface.New(camera.Cam, "main", typeface.NewAlign(typeface.Right, typeface.Center), 1.0, constants.ActualHintSize, 0., 0.)
	timerText.SetColor(hudTextColor)

	itemTransform = transform.New()
	itemTransform.Scalar = pixel.V(camera.Cam.GetZoom(), camera.Cam.GetZoom())
	itemCountText = typeface.New(camera.Cam, "main", typeface.NewAlign(typeface.Center, typeface.Center), 1.0, constants.ActualHintSize, 0., 0.)
	itemCountText.SetColor(hudTextColor)
	itemBoxSpr = img.Batchers[constants.MenuSprites].GetSprite("item_box")

	bombTransform = transform.New()
	bombTransform.Scalar = pixel.V(camera.Cam.GetZoom(), camera.Cam.GetZoom())
	bombCountText = typeface.New(camera.Cam, "main", typeface.NewAlign(typeface.Right, typeface.Center), 1.0, constants.ActualOneSize, 0., 0.)
	bombCountText.SetColor(hudTextColor)
	bombSpr = img.Batchers[constants.EntityKey].GetSprite("mine_1")

	tmpAnimation = reanimator.New(reanimator.NewSwitch().
		AddAnimation(reanimator.NewAnimFromSprites("heart_temp_1", img.Batchers[constants.MenuSprites].Animations["heart_temp_1"].S, reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("heart_temp_2", img.Batchers[constants.MenuSprites].Animations["heart_temp_2"].S, reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("heart_temp_3", img.Batchers[constants.MenuSprites].Animations["heart_temp_3"].S, reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("heart_temp_4", img.Batchers[constants.MenuSprites].Animations["heart_temp_4"].S, reanimator.Hold)).
		SetChooseFn(func() int {
			if dwarf.Health.TempHPTimer == nil {
				return 0
			}
			perc := dwarf.Health.TempHPTimer.Perc()
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
	dwarf := descent.Descent.GetPlayers()[0]
	//hudDistXL := world.TileSize*-(math.Floor(constants.ActualW*0.5/world.TileSize) - 0.5)
	hudDistXR := world.TileSize*(math.Floor(constants.ActualW*0.5/world.TileSize) - 0.5)
	hudDistY := world.TileSize*(math.Floor(constants.BaseH*0.5/world.TileSize))
	currY := 0.

	//hp := dwarf.Health
	//thp := util.Max(hp.Curr + hp.TempHP, hp.Max)
	//if len(heartTransforms) != thp {
	//	heartTransforms = []*transform.Transform{}
	//	for i := 0; i < thp; i++ {
	//		tran := transform.New()
	//		tran.Scalar = pixel.V(camera.Cam.GetZoom(), camera.Cam.GetZoom())
	//		tran.Pos = pixel.V(hudDistXL+heartSpr.Frame().W()*0.5 + float64(i) * heartSpr.Frame().W(), hudDistY - heartSpr.Frame().H()*0.5)
	//		heartTransforms = append(heartTransforms, tran)
	//	}
	//} else {
	//	for i, ht := range heartTransforms {
	//		ht.Pos = pixel.V(hudDistXL+heartSpr.Frame().W()*0.5 + float64(i) * heartSpr.Frame().W(), hudDistY - heartSpr.Frame().H()*0.5)
	//	}
	//}
	//wasFull := fullHP
	//fullHP = hp.Curr == hp.Max && hp.TempHP == 0
	//if (fullHP && !wasFull) || heartTimer == nil {
	//	heartTimer = timing.New(3.)
	//}
	//heartTimer.Update()
	currY = world.TileSize * 1.25

	//gemTransform.Pos = pixel.V(hudDistXL+gemSpr.Frame().W()*0.5, hudDistY - currY - gemSpr.Frame().H()*0.5)
	//gemNumberText.Transform.Pos = gemTransform.Pos
	//gemNumberText.Transform.Pos.X += gemSpr.Frame().W()*0.5 + 4.
	//if dwarf.Player.Stats.CaveGemsFound != lastGem {
	//	lastGem = dwarf.Player.Stats.CaveGemsFound
	//	gemNumberText.SetText(fmt.Sprintf("x%d", lastGem))
	//	gemTimer = timing.New(3.0)
	//} else if lastGem == 0 && gemTimer == nil {
	//	gemTimer = timing.New(0.0)
	//}
	//gemTimer.Update()

	itemTransform.Pos = pixel.V(hudDistXR-itemBoxSpr.Frame().W()*0.5, hudDistY - itemBoxSpr.Frame().H()*0.5)
	itemCountText.SetPos(pixel.V(hudDistXR, hudDistY - itemBoxSpr.Frame().H()))
	if len(dwarf.Player.Inventory.Items) > 0 && dwarf.Player.Inventory.Index < len(dwarf.Player.Inventory.Items) && dwarf.Player.Inventory.Items[dwarf.Player.Inventory.Index].Count > 1 {
		itemCountText.SetText(strconv.Itoa(dwarf.Player.Inventory.Items[dwarf.Player.Inventory.Index].Count))
	} else {
		itemCountText.SetText("")
	}

	if descent.Descent.Timer != nil {
		elapsed := int(descent.Descent.Timer.Elapsed())
		h := elapsed / 3600
		m := elapsed / 60 % 60
		s := elapsed % 60
		if h > 0 {
			timerText.SetText(fmt.Sprintf("%d:%02d:%02d", h, m, s))
		} else {
			timerText.SetText(fmt.Sprintf("%02d:%02d", m, s))
		}
	} else {
		timerText.SetText("")
	}
	timerText.SetPos(pixel.V(hudDistXR - itemBoxSpr.Frame().W() - 8., hudDistY))

	currY = world.TileSize * 1.5
	if descent.Descent.Type == cave.Minesweeper {
		num := player.CaveBombsLeft - (player.OverallStats.CaveWrongFlags + player.OverallStats.CaveBombsFlagged)
		if num == 0 {
			if player.OverallStats.CaveWrongFlags > 0 {
				bombCountText.SetColor(color.RGBA{
					R: 180,
					G: 32,
					B: 42,
					A: 255,
				})
			} else {
				bombCountText.SetColor(color.RGBA{
					R: 20,
					G: 160,
					B: 46,
					A: 255,
				})
			}
		} else {
			bombCountText.SetColor(hudTextColor)
		}
		bombTransform.Pos = pixel.V(hudDistXR-bombSpr.Frame().W()*0.5, hudDistY - currY - bombSpr.Frame().H()*0.5)
		bombCountText.SetText(fmt.Sprintf("%dx", num))
		bombCountText.SetPos(pixel.V(hudDistXR-bombSpr.Frame().W()-4., hudDistY - currY - bombSpr.Frame().H()*0.5))

	}
}

func DrawHUD(win *pixelgl.Window) {
	d := descent.Descent.GetPlayers()[0]
	//for _, tran := range heartTransforms {
	//	tran.UIPos = camera.Cam.APos
	//	tran.UIZoom = camera.Cam.GetZoomScale()
	//	tran.Update()
	//}
	//hp := d.Health
	//if !fullHP || !heartTimer.Done() {
	//	i := 0
	//	for i < hp.Curr && i < len(heartTransforms) {
	//		heartSpr.Draw(win, heartTransforms[i].Mat)
	//		i++
	//	}
	//	for i < hp.TempHP+hp.Curr && i < len(heartTransforms) {
	//		tmpAnimation.Draw(win, heartTransforms[i].Mat)
	//		i++
	//	}
	//	for i < len(heartTransforms) {
	//		heartEmptySpr.Draw(win, heartTransforms[i].Mat)
	//		i++
	//	}
	//}

	if !gemTimer.Done() {
		gemTransform.UIPos = camera.Cam.APos
		gemTransform.UIZoom = camera.Cam.GetZoomScale()
		gemTransform.Update()
		gemNumberText.Update()
		gemSpr.Draw(win, gemTransform.Mat)
		gemNumberText.Draw(win)
	}

	timerText.Update()
	timerText.Draw(win)

	itemTransform.UIPos = camera.Cam.APos
	itemTransform.UIZoom = camera.Cam.GetZoomScale()
	itemTransform.Update()
	itemBoxSpr.Draw(win, itemTransform.Mat)
	if len(d.Player.Inventory.Items) > 0 && d.Player.Inventory.Index < len(d.Player.Inventory.Items) {
		item := d.Player.Inventory.Items[d.Player.Inventory.Index]
		item.Sprite.Draw(win, itemTransform.Mat)
		if item.Timer != nil {
			i := 1
			for float64(i)/16. < item.Timer.Perc() {
				i++
			}
			img.Batchers[constants.MenuSprites].Sprites[fmt.Sprintf("item_timer_%d", i)].Draw(win, itemTransform.Mat)
		}
	}
	itemCountText.Update()
	itemCountText.Draw(win)

	if descent.Descent.Type == cave.Minesweeper {
		bombTransform.UIPos = camera.Cam.APos
		bombTransform.UIZoom = camera.Cam.GetZoomScale()
		bombTransform.Update()
		bombCountText.Update()
		img.Batchers[constants.EntityKey].Sprites["bomb_fuse"].Draw(win, bombTransform.Mat)
		bombCountText.Draw(win)
	}
}
