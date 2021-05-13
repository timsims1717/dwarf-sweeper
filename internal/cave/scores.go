package cave

import (
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/typeface"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"image/color"
	"time"
)

var(
	LowestLevel      int
	LowestLevelItem  ScoreItem
	BlocksDug        int
	BlocksDugItem    ScoreItem
	BombsMarked      int
	BombsMarkedItem  ScoreItem
	BlocksMarked     int
	BlocksMarkedItem ScoreItem
	TotalScore       ScoreItem
	ScoreTimer       time.Time
)

type ScoreItem struct {
	Raw       string
	text      *text.Text
	transform *transform.Transform
	color     color.RGBA
	timer     float64
}

func NewScore(raw string, pos pixel.Vec, timer float64) ScoreItem {
	return ScoreItem{
		Raw:       raw,
		text:      text.New(pixel.ZV, typeface.BasicAtlas),
		transform: &transform.Transform{
			Pos:    pos,
			Scalar: pixel.V(3., 3.),
		},
		color: colornames.Aliceblue,
		timer: timer,
	}
}

func (i *ScoreItem) Update() {
	if time.Since(ScoreTimer).Seconds() > i.timer {
		i.text.Clear()
		i.text.Color = i.color
		fmt.Fprintf(i.text, i.Raw)
		i.transform.Update()
		i.transform.Mat = camera.Cam.UITransform(i.transform.APos, i.transform.Scalar, i.transform.Rot)
	}
}

func (i *ScoreItem) Draw(win *pixelgl.Window) {
	if time.Since(ScoreTimer).Seconds() > i.timer {
		i.text.Draw(win, i.transform.Mat)
	}
}