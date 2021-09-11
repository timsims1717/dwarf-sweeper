package player

import (
	"dwarf-sweeper/internal/cfg"
	"dwarf-sweeper/internal/dungeon"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/menu"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	heartTransforms []*transform.Transform
	lastGem         int
	gemTimer        *timing.FrameTimer
	gemTransform    *transform.Transform
	gemNumberItem   *menu.ItemText
	itemTransform   *transform.Transform
)

func InitHUD() {
	heartTransforms = []*transform.Transform{}
	for i := 0; i < dungeon.Dungeon.Player.Health.Max; i++ {
		tran := transform.NewTransform()
		tran.Anchor.H = transform.Left
		tran.Anchor.V = transform.Center
		tran.Scalar = pixel.V(1.6, 1.6)
		tran.Pos = pixel.V(cfg.BaseW * -0.5 + 8. + float64(i) * 1.2 * world.TileSize, cfg.BaseH * 0.5 - world.TileSize)
		heartTransforms = append(heartTransforms, tran)
	}
	gemTransform = transform.NewTransform()
	gemTransform.Anchor.H = transform.Left
	gemTransform.Anchor.V = transform.Center
	gemTransform.Scalar = pixel.V(1.6, 1.6)
	gemTransform.Pos = pixel.V(cfg.BaseW * -0.5 + 8., cfg.BaseH * 0.5 - (4. + 2.0 * world.TileSize))
	gemNumberItem = menu.NewItemText("", colornames.Aliceblue, pixel.V(1.6, 1.6), menu.Left, menu.Center)
	gemNumberItem.Transform.Pos = pixel.V(cfg.BaseW * -0.5 + 24., cfg.BaseH * 0.5 - (4. + 2.0 * world.TileSize))
	itemTransform = transform.NewTransform()
	itemTransform.Anchor.H = transform.Right
	itemTransform.Anchor.V = transform.Center
	itemTransform.Scalar = pixel.V(1.6, 1.6)
	itemTransform.Pos = pixel.V(cfg.BaseW * 0.5 - 8., cfg.BaseH * 0.5 - world.TileSize)
}

func UpdateHUD() {
	if dungeon.GemsFound != lastGem {
		lastGem = dungeon.GemsFound
		gemNumberItem.SetText(fmt.Sprintf("x%d", lastGem))
		gemTimer = timing.New(3.0)
	} else if lastGem == 0 || gemTimer == nil {
		gemTimer = timing.New(0.0)
	}
	gemTimer.Update()
}

func DrawHUD(win *pixelgl.Window) {
	for _, tran := range heartTransforms {
		tran.UIPos = camera.Cam.APos
		tran.UIZoom = camera.Cam.GetZoomScale()
		tran.Update()
	}
	i := 0
	for i < dungeon.Dungeon.Player.Health.Curr {
		img.Batchers["entities"].Sprites["heart_full"].Draw(win, heartTransforms[i].Mat)
		i++
	}
	for i < dungeon.Dungeon.Player.Health.Max {
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
	if len(dungeon.Inventory) > 0 && dungeon.InvIndex < len(dungeon.Inventory) {
		dungeon.Inventory[dungeon.InvIndex].Sprite.Draw(win, itemTransform.Mat)
	}
}