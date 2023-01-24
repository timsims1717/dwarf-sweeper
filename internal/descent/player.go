package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/camera"
	gween "dwarf-sweeper/pkg/gween64"
	"dwarf-sweeper/pkg/gween64/ease"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"math"
)

func UpdatePlayer(d *Dwarf) {
	d.Player.Message.Update()
	msgDone := d.Player.Message.Done
	if msgDone {
		if len(d.Player.Messages) > 0 {
			m := d.Player.Messages[0]
			if len(d.Player.Messages) > 1 {
				d.Player.Messages = d.Player.Messages[1:]
			} else {
				d.Player.Messages = []*data.RawMessage{}
			}
			d.Player.Message.Text.SetText(m.Raw)
			d.Player.Message.OnClose = m.OnClose
			d.Player.Message.Done = false
			msgDone = false
			d.Player.Message.Box.Open()
		}
	}
	if d.Player.Puzzle != nil {
		if d.Player.Puzzle.IsClosed() {
			if d.Player.Puzzle.Solved() {
				d.Player.Puzzle.OnSolve()
			} else if d.Player.Puzzle.Failed() {
				d.Player.Puzzle.OnFail()
			}
			d.Player.Puzzle = nil
		} else {
			if msgDone {
				d.Player.Puzzle.Update(d.Player.Input)
			} else {
				d.Player.Puzzle.Update(nil)
			}
			if d.Player.Puzzle.IsOpen() && (d.Player.Puzzle.Solved() || d.Player.Puzzle.Failed()) {
				d.Player.Puzzle.Close()
			}
		}
	} else {
		if Descent.DisableInput || !msgDone {
			d.Update(nil)
		} else {
			d.Update(d.Player.Input)
		}
		UpdateInventory(d.Player.Inventory)
	}
}

func UpdateView(d *Dwarf, i, l int) {
	wH := constants.ActualW * 0.5 * camera.Cam.GetZoomScale()
	wQ := constants.ActualW * 0.25 * camera.Cam.GetZoomScale()
	hH := constants.BaseH * 0.5 * camera.Cam.GetZoomScale()
	hQ := constants.BaseH * 0.25 * camera.Cam.GetZoomScale()
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
		// Tween if up
		if d.Player.InterX != nil {
			x, finX := d.Player.InterX.Update(timing.DT)
			d.Player.CamPos.X = x
			if finX {
				d.Player.InterX = nil
			}
		}
		if d.Player.InterY != nil {
			y, finY := d.Player.InterY.Update(timing.DT)
			d.Player.CamPos.Y = y
			if finY {
				d.Player.InterY = nil
			}
		}

		// Cam Position
		d.Player.CamPos.X += timing.DT * d.Player.CamVel.X
		d.Player.CamPos.Y += timing.DT * d.Player.CamVel.Y

		// make sure it still stays with the Dwarf
		if !d.Player.Lock {
			distX *= world.TileSize
			distY *= world.TileSize
			dPos := d.Transform.Pos
			if d.Player.CamPos.X >= dPos.X+distX {
				d.Player.CamPos.X = dPos.X + distX
			} else if d.Player.CamPos.X <= dPos.X-distX {
				d.Player.CamPos.X = dPos.X - distX
			}
			if d.Player.CamPos.Y >= dPos.Y+distY {
				d.Player.CamPos.Y = dPos.Y + distY
			} else if d.Player.CamPos.Y <= dPos.Y-distY {
				d.Player.CamPos.Y = dPos.Y - distY
			}
		}
		// make sure it doesn't move outside the map
		bl, tr := Descent.GetCave().CurrentBoundaries()
		bl.X += d.Player.Canvas.Bounds().W() * 0.5
		bl.Y += d.Player.Canvas.Bounds().H() * 0.5
		tr.X -= d.Player.Canvas.Bounds().W()*0.5 + world.TileSize
		tr.Y -= d.Player.Canvas.Bounds().H() * 0.5
		if bl.X <= tr.X {
			if bl.X > d.Player.CamPos.X {
				d.Player.CamPos.X = bl.X
			} else if tr.X < d.Player.CamPos.X {
				d.Player.CamPos.X = tr.X
			}
		}
		if bl.Y <= tr.Y {
			if bl.Y > d.Player.CamPos.Y {
				d.Player.CamPos.Y = bl.Y
			} else if tr.Y < d.Player.CamPos.Y {
				d.Player.CamPos.Y = tr.Y
			}
		}
		// Shake
		d.Player.PostCamPos.X = d.Player.CamPos.X
		d.Player.PostCamPos.Y = d.Player.CamPos.Y
		if d.Player.ShakeX != nil {
			x, finSX := d.Player.ShakeX.Update(timing.DT)
			d.Player.PostCamPos.X += x
			if finSX {
				d.Player.ShakeX = nil
			}
		}
		if d.Player.ShakeY != nil {
			y, finSY := d.Player.ShakeY.Update(timing.DT)
			d.Player.PostCamPos.Y += y
			if finSY {
				d.Player.ShakeY = nil
			}
		}
		// lock to integer
		d.Player.CamPos.X = math.Round(d.Player.CamPos.X)
		d.Player.CamPos.Y = math.Round(d.Player.CamPos.Y)
		d.Player.PostCamPos.X = math.Round(d.Player.PostCamPos.X)
		d.Player.PostCamPos.Y = math.Round(d.Player.PostCamPos.Y)
	}

	r := pixel.R(math.Round(d.Player.PostCamPos.X-canvasWidth), math.Round(d.Player.PostCamPos.Y-canvasHeight), math.Round(d.Player.PostCamPos.X+canvasWidth), math.Round(d.Player.PostCamPos.Y+canvasHeight))
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
		if i%2 == 0 {
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

func MoveCam(d *Dwarf, pos pixel.Vec, dur float64) {
	d.Player.InterX = gween.New(d.Player.CamPos.X, pos.X, dur, ease.InOutQuad)
	d.Player.InterY = gween.New(d.Player.CamPos.Y, pos.Y, dur, ease.InOutQuad)
}

func ShakeCam(d *Dwarf, dur, freq float64) {
	d.Player.ShakeX = gween.New((random.Effects.Float64()-0.5)*8., 0., dur, SetSine(freq))
	d.Player.ShakeY = gween.New((random.Effects.Float64()-0.5)*8., 0., dur, SetSine(freq))
}

func SetSine(freq float64) func(float64, float64, float64, float64) float64 {
	return func(t, b, c, d float64) float64 {
		return b * math.Pow(math.E, -math.Abs(c)*t) * math.Sin(freq*math.Pi*t)
	}
}
