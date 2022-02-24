package hud

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/typeface"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"math"
)

var (
	HUDs []*HUD
)

type HUD struct {
	Dwarf *descent.Dwarf

	MaxHP     int
	DisplayHP bool
	HPTimer   *timing.FrameTimer
	HPTrans   []*transform.Transform

	LastGem    int
	DisplayGem bool
	GemTimer   *timing.FrameTimer
	GemTrans   *transform.Transform
	GemText    *typeface.Text
}

func New(dwarf *descent.Dwarf) *HUD {
	var hpTrans []*transform.Transform
	for i := 0; i < dwarf.Health.Max; i++ {
		trans := transform.New()
		trans.Scalar = pixel.V(camera.Cam.GetZoom(), camera.Cam.GetZoom())
		hpTrans = append(hpTrans, trans)
	}
	gemTrans := transform.New()
	gemTrans.Scalar = pixel.V(camera.Cam.GetZoom(), camera.Cam.GetZoom())
	gemText := typeface.New(camera.Cam, "main", typeface.NewAlign(typeface.Left, typeface.Center), 1.0, constants.ActualOneSize, 0., 0.)
	gemText.SetColor(hudTextColor)
	return &HUD{
		Dwarf:     dwarf,
		DisplayHP: true,
		HPTrans:   hpTrans,
		GemTrans:  gemTrans,
		GemText:   gemText,
	}
}

func (hud *HUD) Update() {
	hudDistXL := math.Floor((hud.Dwarf.Player.CanvasPos.X - world.TileSize*(math.Floor(hud.Dwarf.Player.Canvas.Bounds().W()*0.5/world.TileSize) - 0.5)) * camera.Cam.Zoom)
	//hudDistXR := world.TileSize*(math.Floor(hud.Dwarf.Player.Canvas.Bounds().W()*0.5/world.TileSize) - 0.5)
	hudDistY := math.Floor((hud.Dwarf.Player.CanvasPos.Y + world.TileSize*(math.Floor(hud.Dwarf.Player.Canvas.Bounds().H()*0.5/world.TileSize))) * camera.Cam.Zoom)
	currY := 0.

	hp := hud.Dwarf.Health
	thp := util.Max(hp.Curr + hp.TempHP, hp.Max)
	if len(hud.HPTrans) != thp {
		hud.HPTrans = []*transform.Transform{}
		for i := 0; i < thp; i++ {
			trans := transform.New()
			trans.Scalar = pixel.V(camera.Cam.GetZoom(), camera.Cam.GetZoom())
			trans.Pos = pixel.V(hudDistXL+heartSpr.Frame().W()*0.5 + float64(i) * heartSpr.Frame().W(), hudDistY - heartSpr.Frame().H()*0.5)
			hud.HPTrans = append(hud.HPTrans, trans)
		}
	} else {
		for i, ht := range hud.HPTrans {
			ht.Pos = pixel.V(hudDistXL+heartSpr.Frame().W()*0.5 + float64(i) * heartSpr.Frame().W(), hudDistY - heartSpr.Frame().H()*0.5)
		}
	}

	displayHP := hud.DisplayHP
	hud.DisplayHP = hud.MaxHP != hp.Max || hp.Curr != hp.Max || hp.TempHP != 0
	hud.MaxHP = hp.Max
	if displayHP != hud.DisplayHP || hud.HPTimer == nil {
		hud.HPTimer = timing.New(3.)
	}
	hud.HPTimer.Update()
	currY = world.TileSize * 1.25

	hud.GemTrans.Pos = pixel.V(hudDistXL+gemSpr.Frame().W()*0.5, hudDistY - currY - gemSpr.Frame().H()*0.5)
	gtp := hud.GemTrans.Pos
	gtp.X += gemSpr.Frame().W()*0.5 + 4.
	hud.GemText.SetPos(gtp)
	if hud.Dwarf.Player.Stats.CaveGemsFound != hud.LastGem {
		hud.LastGem = hud.Dwarf.Player.Stats.CaveGemsFound
		hud.GemText.SetText(fmt.Sprintf("x%d", hud.LastGem))
		hud.GemTimer = timing.New(3.0)
	} else if hud.LastGem == 0 && hud.GemTimer == nil {
		hud.GemTimer = timing.New(0.0)
	}
	hud.GemTimer.Update()
}

func (hud *HUD) Draw(target pixel.Target) {
	hp := hud.Dwarf.Health

	for _, ht := range hud.HPTrans {
		ht.UIPos = camera.Cam.APos
		ht.UIZoom = camera.Cam.GetZoomScale()
		ht.Update()
	}
	if hud.DisplayHP || !hud.HPTimer.Done() {
		i := 0
		for i < hp.Curr && i < len(hud.HPTrans) {
			heartSpr.Draw(target, hud.HPTrans[i].Mat)
			i++
		}
		for i < hp.TempHP+hp.Curr && i < len(hud.HPTrans) {
			tmpAnimation.Draw(target, hud.HPTrans[i].Mat)
			i++
		}
		for i < len(hud.HPTrans) {
			heartEmptySpr.Draw(target, hud.HPTrans[i].Mat)
			i++
		}
	}

	if !hud.GemTimer.Done() {
		hud.GemTrans.UIPos = camera.Cam.APos
		hud.GemTrans.UIZoom = camera.Cam.GetZoomScale()
		hud.GemTrans.Update()
		hud.GemText.Update()
		gemSpr.Draw(target, hud.GemTrans.Mat)
		hud.GemText.Draw(target)
	}
}