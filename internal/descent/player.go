package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/pkg/camera"
	"github.com/faiface/pixel"
	"math"
)

func UpdatePlayer(d *Dwarf) {
	if d.Player.Puzzle != nil {
		if d.Player.Puzzle.IsClosed() {
			if d.Player.Puzzle.Solved() {
				d.Player.Puzzle.OnSolve()
			}
			d.Player.Puzzle = nil
		} else {
			d.Player.Puzzle.Update(d.Player.Input)
			if d.Player.Puzzle.Solved() && d.Player.Puzzle.IsOpen() {
				d.Player.Puzzle.Close()
			}
		}
	} else {
		if Descent.DisableInput {
			d.Update(nil)
		} else {
			d.Update(d.Player.Input)
		}
		d.Player.Inventory.Update()
	}
}

func UpdateView(d *Dwarf, i, l int) {
	wH := constants.ActualW*0.5*camera.Cam.GetZoomScale()
	wQ := constants.ActualW*0.25*camera.Cam.GetZoomScale()
	hH := constants.BaseH*0.5*camera.Cam.GetZoomScale()
	hQ := constants.BaseH*0.25*camera.Cam.GetZoomScale()
	var canvasWidth, canvasHeight float64
	if l == 1 {
		canvasWidth = wH
		canvasHeight = hH
	} else if l == 2 {
		canvasWidth = wQ
		canvasHeight = hH
	} else {
		canvasWidth = wQ
		canvasHeight = hQ
	}
	if !Descent.FreeCam {
		d.Player.CamPos = d.Transform.Pos
		//dPos := d.Transform.Pos
		//if d.Physics.IsMovingX() {
		//	d.Player.CamTar.X = d.Physics.Velocity.X / d.Speed * world.TileSize * 4.
		//}
		//if d.Physics.IsMovingY() {
		//	d.Player.CamTar.Y = d.Physics.Velocity.Y / d.Speed * world.TileSize * 4.
		//} else if d.Player.Input.Get("up").Pressed() && !d.Player.Input.Get("down").Pressed() {
		//	d.Player.CamTar.Y = world.TileSize * 4.
		//} else if !d.Player.Input.Get("up").Pressed() && d.Player.Input.Get("down").Pressed() {
		//	d.Player.CamTar.Y = world.TileSize * -4.
		//} else {
		//	d.Player.CamTar.Y = 0.
		//}
		//d.Player.CamPos.X = 10. * timing.DT * (dPos.X - d.Player.CamTar.X)
		//d.Player.CamPos.Y = 10. * timing.DT * (dPos.Y - d.Player.CamTar.Y)

		// make sure it still stays with the Dwarf
		//dist := world.TileSize * 5.
		//if d.Player.CamPos.X >= dPos.X +dist {
		//	d.Player.CamPos.X = dPos.X + dist
		//} else if d.Player.CamPos.X <= dPos.X-dist {
		//	d.Player.CamPos.X = dPos.X - dist
		//}
		//if d.Player.CamPos.Y >= dPos.Y +dist {
		//	d.Player.CamPos.Y = dPos.Y + dist
		//} else if d.Player.CamPos.Y <= dPos.Y-dist {
		//	d.Player.CamPos.Y = dPos.Y - dist
		//}
		//bl, tr := Descent.GetCave().CurrentBoundaries()
		//ratio := camera.Cam.Height / constants.BaseH
		//bl.X += camera.Cam.Width * 0.5 / ratio * camera.Cam.GetZoomScale()
		//bl.Y += constants.BaseH * 0.5 * camera.Cam.GetZoomScale()
		//tr.X -= camera.Cam.Width*0.5/ratio*camera.Cam.GetZoomScale() + world.TileSize
		//tr.Y -= constants.BaseH * 0.5 * camera.Cam.GetZoomScale()
		//if bl.X <= tr.X {
		//	if bl.X > d.Player.CamPos.X {
		//		d.Player.CamPos.X = bl.X
		//	} else if tr.X < d.Player.CamPos.X {
		//		d.Player.CamPos.X = tr.X
		//	}
		//}
		//if bl.Y <= tr.Y {
		//	if bl.Y > d.Player.CamPos.Y {
		//		d.Player.CamPos.Y = bl.Y
		//	} else if tr.Y < d.Player.CamPos.Y {
		//		d.Player.CamPos.Y = tr.Y
		//	}
		//}
		d.Player.CamPos.X = math.Round(d.Player.CamPos.X)
		d.Player.CamPos.Y = math.Round(d.Player.CamPos.Y)
	}

	r := pixel.R(math.Round(d.Player.CamPos.X - canvasWidth), math.Round(d.Player.CamPos.Y - canvasHeight), math.Round(d.Player.CamPos.X + canvasWidth), math.Round(d.Player.CamPos.Y + canvasHeight))
	d.Player.Canvas.SetBounds(r)
	var y, x float64
	if l == 1 {
		y = 0.
		x = 0.
	} else if l == 2 {
		y = 0.
		if i == 0 {
			x = -wQ
		} else {
			x = wQ
		}
	} else if l == 3 {
		if i == 0 {
			y = 0.
			x = -wQ
		} else {
			x = wQ
			if i == 1 {
				y = hQ
			} else {
				x = -hQ
			}
		}
	} else if l == 4 {
		if i % 2 == 0 {
			x = -wQ
		} else {
			x = wQ
		}
		if i < 2 {
			y = hQ
		} else {
			y = -hQ
		}
	}
	d.Player.CanvasPos = pixel.V(x, y)
}