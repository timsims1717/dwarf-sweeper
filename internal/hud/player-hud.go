package hud

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/typeface"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"math"
	"strconv"
)

var (
	HUDs []*HUD
)

type HUD struct {
	Dwarf   *descent.Dwarf
	Refresh bool

	MaxHP     int
	DisplayHP bool
	HPTimer   *timing.Timer
	HPTrans   []*transform.Transform
	TempAnim  *reanimator.Tree

	LastGem    int
	DisplayGem bool
	GemTimer   *timing.Timer
	GemTrans   *transform.Transform
	GemText    *typeface.Text

	ItemTrans *transform.Transform
	ItemText  *typeface.Text
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
	gemText := typeface.New(&camera.Cam.APos, "main", typeface.NewAlign(typeface.Left, typeface.Center), 1.0, constants.ActualOneSize, 0., 0.)
	gemText.SetColor(hudTextColor)

	itemTrans := transform.New()
	itemTrans.Scalar = pixel.V(camera.Cam.GetZoom(), camera.Cam.GetZoom())
	itemText := typeface.New(&camera.Cam.APos, "main", typeface.NewAlign(typeface.Center, typeface.Center), 1.0, constants.ActualHintSize, 0., 0.)
	itemText.SetColor(hudTextColor)

	anim := reanimator.New(reanimator.NewSwitch().
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
	myecs.Manager.NewEntity().AddComponent(myecs.Animation, anim)

	return &HUD{
		Dwarf:     dwarf,
		DisplayHP: true,
		HPTrans:   hpTrans,
		TempAnim:  anim,
		GemTrans:  gemTrans,
		GemText:   gemText,
		ItemTrans: itemTrans,
		ItemText:  itemText,
	}
}

func (hud *HUD) Update() {
	hudDistXL := math.Floor((hud.Dwarf.Player.CanvasPos.X - world.TileSize*(math.Floor(hud.Dwarf.Player.Canvas.Bounds().W()*0.5/world.TileSize) - 0.5)) * camera.Cam.Zoom)
	hudDistXR := math.Floor((hud.Dwarf.Player.CanvasPos.X + world.TileSize*(math.Floor(hud.Dwarf.Player.Canvas.Bounds().W()*0.5/world.TileSize) - 0.5)) * camera.Cam.Zoom)
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
	if displayHP != hud.DisplayHP || hud.HPTimer == nil || hud.Refresh {
		hud.HPTimer = timing.New(3.)
	}
	hud.HPTimer.Update()
	currY = world.TileSize * 1.25

	hud.GemTrans.Pos = pixel.V(hudDistXL+gemSpr.Frame().W()*0.5, hudDistY - currY - gemSpr.Frame().H()*0.5)
	gtp := hud.GemTrans.Pos
	gtp.X += gemSpr.Frame().W()*0.5 + 4.
	hud.GemText.SetPos(gtp)
	if hud.Dwarf.Player.Stats.CaveGemsFound != hud.LastGem || hud.Refresh {
		hud.LastGem = hud.Dwarf.Player.Stats.CaveGemsFound
		hud.GemText.SetText(fmt.Sprintf("x%d", hud.LastGem))
		hud.GemTimer = timing.New(3.0)
	} else if hud.LastGem == 0 && hud.GemTimer == nil {
		hud.GemTimer = timing.New(0.0)
	}
	hud.GemTimer.Update()

	hud.ItemTrans.Pos = pixel.V(hudDistXR-itemBoxSpr.Frame().W()*0.5, hudDistY - itemBoxSpr.Frame().H()*0.5)
	hud.ItemText.SetPos(pixel.V(hudDistXR, hudDistY - itemBoxSpr.Frame().H()))
	inv := hud.Dwarf.Player.Inventory
	if len(inv.Items) > 0 && inv.Index < len(inv.Items) {
		hud.ItemText.SetText(strconv.Itoa(inv.Items[inv.Index].Count))
	} else {
		hud.ItemText.SetText("")
	}
	hud.Refresh = false
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
			hud.TempAnim.Draw(target, hud.HPTrans[i].Mat)
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

	hud.ItemTrans.UIPos = camera.Cam.APos
	hud.ItemTrans.UIZoom = camera.Cam.GetZoomScale()
	hud.ItemTrans.Update()
	hud.ItemText.Update()
	itemBoxSpr.Draw(target, hud.ItemTrans.Mat)
	inv := hud.Dwarf.Player.Inventory
	if len(inv.Items) > 0 && inv.Index < len(inv.Items) {
		item := inv.Items[inv.Index]
		item.Sprite.Draw(target, hud.ItemTrans.Mat)
		if item.Timer != nil {
			i := 1
			for float64(i)/16. < item.Timer.Perc() {
				i++
			}
			img.Batchers[constants.MenuSprites].Sprites[fmt.Sprintf("item_timer_%d", i)].Draw(target, hud.ItemTrans.Mat)
		}
	}
	hud.ItemText.Update()
	hud.ItemText.Draw(target)
}