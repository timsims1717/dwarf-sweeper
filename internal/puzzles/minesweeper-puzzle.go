package puzzles

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/minesweeper"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/input"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

const (
	scalar = 1.8
)

type MinePuzzle struct {
	SizeW int
	SizeH int
	Level int

	Board *minesweeper.Board
	Hover world.Coords

	Box *menus.MenuBox

	CellTrans [][]*transform.Transform
}

func NewMinePuzzle(cam *camera.Camera, level int) *MinePuzzle {
	size := 4 + level / 4
	sizeW := util.Min(size, 9)
	sizeH := util.Min(size, 7)
	box := menus.NewBox(cam)
	var cts [][]*transform.Transform
	for y := 0; y < sizeH; y++ {
		cts = append(cts, []*transform.Transform{})
		evenY := sizeH % 2 == 0
		midY := sizeH / 2
		var multY float64
		if evenY {
			multY = float64(y-midY)+0.5
		} else {
			multY = float64(y-midY)
		}
		yPos := multY*(world.TileSize*scalar+1)
		for x := 0; x < sizeW; x++ {
			evenX := sizeW % 2 == 0
			midX := sizeW / 2
			var multX float64
			if evenX {
				multX = float64(x-midX)+0.5
			} else {
				multX = float64(x-midX)
			}
			trans := transform.NewTransform()
			trans.Pos = pixel.V(multX*(world.TileSize*scalar+1), yPos)
			trans.Scalar = pixel.V(scalar, scalar)
			cts[y] = append(cts[y], trans)
		}
	}
	return &MinePuzzle{
		SizeW:     sizeW,
		SizeH:     sizeH,
		Level:     level,
		Box:       box,
		CellTrans: cts,
	}
}

func (mp *MinePuzzle) Open() {
	area := mp.SizeW * mp.SizeH
	amt := area / 6 + random.Effects.Intn(area / 12) + mp.Level / 5 + random.Effects.Intn((mp.Level + 5) / 5)
	board := minesweeper.CreateBoard(mp.SizeW, mp.SizeH, amt, random.Effects)
	minesweeper.RevealTilSolvableP(board, random.Effects)
	mp.Board = minesweeper.CreateBoard(mp.SizeW, mp.SizeH, amt, random.Effects)
	minesweeper.RevealTilSolvableP(board, random.Effects)
	mp.Box.Open()
	mp.Box.SetSize(pixel.R(0., 0., float64(mp.SizeW) * (world.TileSize + 2.) * scalar, float64(mp.SizeH) * (world.TileSize + 2.) * scalar))
}

func (mp *MinePuzzle) Close() {
	mp.Box.Close()
}

func (mp *MinePuzzle) Update(in *input.Input) {
	mp.Box.Update()
	mp.UpdateTransforms()
}

func (mp *MinePuzzle) UpdateTransforms() {
	if mp.Box.Cam != nil {
		for _, row := range mp.CellTrans {
			for _, ct := range row {
				ct.UIZoom = mp.Box.Cam.GetZoomScale()
				ct.UIPos = mp.Box.Cam.APos
				ct.Update()
			}
		}
	}
}

func (mp *MinePuzzle) Draw(target pixel.Target) {
	mp.Box.Draw(target)
	if !mp.Box.IsClosed() && mp.Box.IsOpen() {
		for _, row := range mp.CellTrans {
			for _, ct := range row {
				debug.AddLine(colornames.Wheat, imdraw.RoundEndShape, ct.Pos, ct.Pos, 2.)
				img.Batchers[constants.EntityKey].Sprites["bomb_fuse"].Draw(target, ct.Mat)
			}
		}
	}
}