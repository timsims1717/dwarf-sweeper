package credits

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/pkg/camera"
	gween "dwarf-sweeper/pkg/gween64"
	"dwarf-sweeper/pkg/gween64/ease"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/typeface"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"image/color"
)

var (
	CreditString = `DwarfSweeper

 
Gameplay and Graphics
Tim Sims

Sound
Alexander Kopeikin
PMSFX

Music
Ben Reber

Pixel Art
Kenney

Special Thanks:
My wife Kaylan,
Marshall and Clark,
faiface and the Pixel team,
the Ludum Dare LD48 team.

Thank you for playing!`
	Credits *typeface.Text
	overlay *imdraw.IMDraw
	interY  *gween.Tween
	interA  *gween.Tween
	opening bool
	open    bool
	closing bool
)

func Initialize() {
	overlay = imdraw.New(nil)
	Credits = typeface.New(camera.Cam,"main", typeface.Alignment{ H: typeface.Center, V: typeface.Top }, 1.2, constants.ActualHintSize, 0., 0.)
	Credits.Color = color.RGBA{
		R: 218,
		G: 224,
		B: 234,
		A: 0,
	}
	Credits.SetText(CreditString)
	Credits.Transform.Pos.Y = constants.BaseH*-0.5 - 20.
}

func Update() {
	overlay.Clear()
	if opening {
		open = true
		opening = false
	}
	if open || closing {
		if interA != nil {
			a, fin := interA.Update(timing.DT)
			overlay.Color = color.RGBA{
				R: 0,
				G: 0,
				B: 0,
				A: uint8(a),
			}
			Credits.Color = color.RGBA{
				R: 218,
				G: 224,
				B: 234,
				A: uint8(a),
			}
			if fin {
				interA = nil
				closing = false
			}
		}
		if interY != nil {
			y, fin := interY.Update(timing.DT)
			Credits.Transform.Pos.Y = y
			if fin {
				interY = nil
				Close()
			}
		}
		overlay.EndShape = imdraw.SharpEndShape
		w := camera.Cam.Width
		h := camera.Cam.Height
		overlay.Push(pixel.V(w*-0.5, h*-0.5))
		overlay.Push(pixel.V(w*0.5, h*-0.5))
		overlay.Push(pixel.V(w*0.5, h*0.5))
		overlay.Push(pixel.V(w*-0.5, h*0.5))
		overlay.Polygon(0.)
		Credits.Transform.UIPos = camera.Cam.APos
		Credits.Transform.UIZoom = camera.Cam.GetZoomScale()
		Credits.Update()
	}
}

func Draw(target pixel.Target) {
	if open || closing {
		overlay.Draw(target)
		Credits.Draw(target)
	}
}

func Opened() bool {
	return opening || open || closing
}

func Open() {
	if !opening && !open {
		closing = false
		opening = true
		s := constants.BaseH*-0.5 - 20.
		e := (Credits.Text.BoundsOf(CreditString).H() * Credits.RelativeSize + constants.BaseH*0.5) + 20.
		Credits.Transform.Pos.Y = s
		overlay.Color = color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 0,
		}
		interY = gween.New(s, e, Credits.Text.BoundsOf(CreditString).H()/Credits.Text.LineHeight*0.65, ease.Linear)
		interA = gween.New(0., 210., 2., ease.Linear)
	}
}

func Close() {
	if !closing {
		opening = false
		closing = true
		open = false
		interA = gween.New(210., 0., 0.5, ease.Linear)
	}
}
