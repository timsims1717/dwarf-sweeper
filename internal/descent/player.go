package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data/player"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"math"
)

func UpdatePlayer(d *Dwarf) {
	d.Player.Message.Update()
	if d.Player.Message.Done {
		if len(d.Player.Messages) > 0 {
			m := d.Player.Messages[0]
			if len(d.Player.Messages) > 1 {
				d.Player.Messages = d.Player.Messages[1:]
			} else {
				d.Player.Messages = []*player.RawMessage{}
			}
			d.Player.Message.Text.SetText(m.Raw)
			d.Player.Message.OnClose = m.OnClose
			d.Player.Message.Done = false
			d.Player.Message.Box.Open()
		} else if d.Player.Puzzle != nil {
			if d.Player.Puzzle.IsClosed() {
				if d.Player.Puzzle.Solved() {
					d.Player.Puzzle.OnSolve()
				} else if d.Player.Puzzle.Failed() {
					d.Player.Puzzle.OnFail()
				}
				d.Player.Puzzle = nil
			} else {
				d.Player.Puzzle.Update(d.Player.Input)
				if d.Player.Puzzle.IsOpen() && (d.Player.Puzzle.Solved() || d.Player.Puzzle.Failed()) {
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
}

func UpdateView(d *Dwarf, i, l int) {
	wH := constants.ActualW*0.5*camera.Cam.GetZoomScale()
	wQ := constants.ActualW*0.25*camera.Cam.GetZoomScale()
	hH := constants.BaseH*0.5*camera.Cam.GetZoomScale()
	hQ := constants.BaseH*0.25*camera.Cam.GetZoomScale()
	var canvasWidth, canvasHeight float64
	var distX, distY float64
	if l == 1 {
		canvasWidth = wH
		canvasHeight = hH
		distX = 3.
		distY = 2.
	} else if l == 2 {
		if constants.SplitScreenV {
			canvasWidth = wQ
			canvasHeight = hH
			distX = 2.
			distY = 2.
		} else {
			canvasWidth = wH
			canvasHeight = hQ
			distX = 3.
			distY = 1.
		}
	} else {
		if l == 3 && i == 0 {
			if constants.SplitScreenV {
				canvasWidth = wQ
				canvasHeight = hH
				distX = 2.
				distY = 2.
			} else {
				canvasWidth = wH
				canvasHeight = hQ
				distX = 3.
				distY = 1.
			}
		} else {
			canvasWidth = wQ
			canvasHeight = hQ
			distX = 2.
			distY = 1.
		}
	}
	if !Descent.FreeCam {
		//d.Player.CamPos = d.Transform.Pos
		d.Player.CamPos.X += timing.DT * d.Player.CamVel.X
		d.Player.CamPos.Y += timing.DT * d.Player.CamVel.Y

		// make sure it still stays with the Dwarf
		distX *= world.TileSize
		distY *= world.TileSize
		dPos := d.Transform.Pos
		if d.Player.CamPos.X >= dPos.X + distX {
			d.Player.CamPos.X = dPos.X + distX
		} else if d.Player.CamPos.X <= dPos.X - distX {
			d.Player.CamPos.X = dPos.X - distX
		}
		if d.Player.CamPos.Y >= dPos.Y + distY {
			d.Player.CamPos.Y = dPos.Y + distY
		} else if d.Player.CamPos.Y <= dPos.Y - distY {
			d.Player.CamPos.Y = dPos.Y - distY
		}
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
		if constants.SplitScreenV {
			y = 0.
			if i == 0 {
				x = -wQ
			} else {
				x = wQ
			}
		} else {
			x = 0.
			if i == 0 {
				y = hQ
			} else {
				x = -hQ
			}
		}
	} else if l == 3 {
		if constants.SplitScreenV {
			if i == 0 {
				y = 0.
				x = -wQ
			} else {
				x = wQ
				if i == 1 {
					y = hQ
				} else {
					y = -hQ
				}
			}
		} else {
			if i == 0 {
				x = 0.
				y = hQ
			} else {
				y = -hQ
				if i == 1 {
					x = -wQ
				} else {
					x = wQ
				}
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