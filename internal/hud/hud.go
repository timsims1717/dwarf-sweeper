package hud

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data/player"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/typeface"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
	"math"
)

var (
	imd *imdraw.IMDraw

	heartSpr      *pixel.Sprite
	heartEmptySpr *pixel.Sprite
	gemSpr        *pixel.Sprite
	itemBoxSpr    *pixel.Sprite

	timerText *typeface.Text
	ShowTimer bool

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
	imd = imdraw.New(nil)

	heartSpr = img.Batchers[constants.MenuSprites].GetSprite("heart_full")
	heartEmptySpr = img.Batchers[constants.MenuSprites].GetSprite("heart_empty")

	gemSpr = img.Batchers[constants.EntityKey].GetSprite("gem_diamond")

	timerText = typeface.New(&camera.Cam.APos, "main", typeface.NewAlign(typeface.Center, typeface.Center), 1.0, constants.ActualHintSize, 0., 0.)
	timerText.SetColor(hudTextColor)

	itemBoxSpr = img.Batchers[constants.MenuSprites].GetSprite("item_box")

	bombTransform = transform.New()
	bombTransform.Scalar = pixel.V(camera.Cam.GetZoom(), camera.Cam.GetZoom())
	bombCountText = typeface.New(&camera.Cam.APos, "main", typeface.NewAlign(typeface.Center, typeface.Center), 1.0, constants.ActualOneSize, 0., 0.)
	bombCountText.SetColor(hudTextColor)
	bombSpr = img.Batchers[constants.EntityKey].GetSprite("mine_1")
}

func UpdateHUD() {
	hudDistY := world.TileSize*(math.Floor(constants.BaseH*0.5/world.TileSize))
	currY := 0.

	currY = world.TileSize * 0.75

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
	timerText.SetPos(pixel.V(0., hudDistY))

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
		bombTransform.Pos = pixel.V(0., hudDistY - currY - bombSpr.Frame().H()*0.5)
		bombCountText.SetText(fmt.Sprintf("%d", num))
		bombCountText.SetPos(pixel.V(0., hudDistY - currY - bombSpr.Frame().H()*0.5))
	}

	imd.Clear()
	imd.Color = constants.BGColor
	switch len(descent.Descent.Dwarves) {
	case 2:
		if constants.SplitScreenV {
			imd.Push(pixel.V(0., constants.BaseH*0.5), pixel.V(0., constants.BaseH*-0.5))
			imd.Line(1.0)
		} else {
			imd.Push(pixel.V(constants.ActualW*-0.5, 0.), pixel.V(constants.ActualW*0.5, 0.))
			imd.Line(1.0)
		}
	case 3:
		if constants.SplitScreenV {
			imd.Push(pixel.V(0., constants.BaseH*0.5), pixel.V(0., constants.BaseH*-0.5))
			imd.Line(1.0)
			imd.Push(pixel.V(0., 0.), pixel.V(constants.ActualW*0.5, 0.))
			imd.Line(1.0)
		} else {
			imd.Push(pixel.V(constants.ActualW*-0.5, 0.), pixel.V(constants.ActualW*0.5, 0.))
			imd.Line(1.0)
			imd.Push(pixel.V(0., 0.), pixel.V(0., constants.BaseH*-0.5))
			imd.Line(1.0)
		}
	case 4:
		imd.Push(pixel.V(0., constants.BaseH*0.5), pixel.V(0., constants.BaseH*-0.5))
		imd.Line(1.0)
		imd.Push(pixel.V(constants.ActualW*-0.5, 0.), pixel.V(constants.ActualW*0.5, 0.))
		imd.Line(1.0)
	}
}

func DrawHUD(win *pixelgl.Window) {
	imd.Draw(win)

	if ShowTimer {
		timerText.Update()
		timerText.Draw(win)
	}

	if descent.Descent.Type == cave.Minesweeper {
		bombTransform.UIPos = camera.Cam.APos
		bombTransform.UIZoom = camera.Cam.GetZoomScale()
		bombTransform.Update()
		bombCountText.Update()
		bombSpr.Draw(win, bombTransform.Mat)
		bombCountText.Draw(win)
	}
}
