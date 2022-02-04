package puzzles

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/minesweeper"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/input"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
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
	solved    bool

	FlagAnim *reanimator.Tree
	ExAnim   *reanimator.Tree
}

func (mp *MinePuzzle) Create(cam *camera.Camera, level int) {
	size := 4 + level / 4
	mp.SizeW = util.Min(size, 9)
	mp.SizeH = util.Min(size, 7)
	mp.Box = menus.NewBox(cam)
	var cts [][]*transform.Transform
	for y := 0; y < mp.SizeH; y++ {
		cts = append(cts, []*transform.Transform{})
		evenY := mp.SizeH % 2 == 0
		midY := mp.SizeH / 2
		var multY float64
		if evenY {
			multY = float64(y-midY)+0.5
		} else {
			multY = float64(y-midY)
		}
		yPos := multY*(world.TileSize*scalar+1)
		for x := 0; x < mp.SizeW; x++ {
			evenX := mp.SizeW % 2 == 0
			midX := mp.SizeW / 2
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
	mp.CellTrans = cts
	mp.FlagAnim = reanimator.NewSimple(reanimator.NewAnimFromSprites("flag", img.Batchers[constants.PuzzleKey].GetAnimation("flag_hang").S, reanimator.Loop))
	mp.ExAnim = reanimator.NewSimple(reanimator.NewAnimFromSprites("ex", img.Batchers[constants.PuzzleKey].GetAnimation("green_x").S, reanimator.Loop))
}

func (mp *MinePuzzle) IsOpen() bool {
	return mp.Box.IsOpen()
}

func (mp *MinePuzzle) IsClosed() bool {
	return mp.Box.IsClosed()
}

func (mp *MinePuzzle) Open() {
	area := mp.SizeW * mp.SizeH
	r := util.Min(mp.SizeW, mp.SizeH) - 1
	amt := area / util.Max(mp.SizeW, mp.SizeH) + random.Effects.Intn(r) - r / 2
	mp.Board = minesweeper.CreateBoard(mp.SizeW, mp.SizeH, amt, random.Effects)
	minesweeper.RevealTilSolvableP(mp.Board, random.Effects)
	minesweeper.UnRevealWhileSolvableP(mp.Board, random.Effects)
	mp.Box.Open()
	mp.Box.SetSize(pixel.R(0., 0., float64(mp.SizeW) * (world.TileSize + 2.) * scalar, float64(mp.SizeH) * (world.TileSize + 2.) * scalar))
	mp.Hover.Y = mp.SizeH-1
}

func (mp *MinePuzzle) Close() {
	mp.Box.Close()
}

func (mp *MinePuzzle) Update(in *input.Input) {
	mp.Box.Update()
	mp.UpdateTransforms()
	mp.FlagAnim.Update()
	mp.ExAnim.Update()
	if !mp.solved {
		if in.MouseMoved {
			for y, row := range mp.CellTrans {
				for x, ct := range row {
					if mp.Hover.X != x || mp.Hover.Y != y {
						point := in.World
						if util.PointInside(point, pixel.R(0., 0., 16., 16.), ct.Mat) {
							mp.Hover.X = x
							mp.Hover.Y = y
						}
					}
				}
			}
		}
		if in.Get("up").JustPressed() && mp.Hover.Y < mp.SizeH-1 {
			in.Get("up").Consume()
			mp.Hover.Y++
			sfx.SoundPlayer.PlaySound("click", 2.0)
		} else if in.Get("down").JustPressed() && mp.Hover.Y > 0 {
			in.Get("down").Consume()
			mp.Hover.Y--
			sfx.SoundPlayer.PlaySound("click", 2.0)
		} else if in.Get("left").JustPressed() && mp.Hover.X > 0 {
			in.Get("left").Consume()
			mp.Hover.X--
			sfx.SoundPlayer.PlaySound("click", 2.0)
		} else if in.Get("right").JustPressed() && mp.Hover.X < mp.SizeW-1 {
			in.Get("right").Consume()
			mp.Hover.X++
			sfx.SoundPlayer.PlaySound("click", 2.0)
		}
		cell := mp.Board.Board[mp.Hover.Y][mp.Hover.X]
		if in.Get("flag").JustPressed() && !cell.Rev && !cell.Ex {
			in.Get("flag").Consume()
			cell.Flag = !cell.Flag
		} else if in.Get("dig").JustPressed() && !cell.Rev && !cell.Flag {
			in.Get("dig").Consume()
			if cell.Bomb {
				cell.Rev = true
			} else {
				cell.Ex = true
			}
		}
		mp.Board.Board[mp.Hover.Y][mp.Hover.X] = cell
		done := true
		for _, row := range mp.Board.Board {
			for _, c := range row {
				if !c.Rev && ((c.Bomb && !c.Flag) || (!c.Bomb && !c.Ex)) {
					done = false
				}
			}
		}
		mp.solved = done
	}
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
		for y, row := range mp.CellTrans {
			for x, ct := range row {
				img.Batchers[constants.PuzzleKey].Sprites["background"].Draw(target, ct.Mat)
				cell := mp.Board.Board[y][x]
				if cell.Rev {
					if cell.Bomb {
						img.Batchers[constants.EntityKey].Sprites["mine_1"].Draw(target, ct.Mat)
					} else {
						var str string
						switch cell.Num {
						case 0:
							str = "zero"
						case 1:
							str = "one"
						case 2:
							str = "two"
						case 3:
							str = "three"
						case 4:
							str = "four"
						case 5:
							str = "five"
						case 6:
							str = "six"
						case 7:
							str = "seven"
						case 8:
							str = "eight"
						}
						img.Batchers[constants.PuzzleKey].Sprites[str].Draw(target, ct.Mat)
					}
				} else if cell.Flag {
					mp.FlagAnim.Draw(target, ct.Mat)
				} else if cell.Ex {
					mp.ExAnim.Draw(target, ct.Mat)
				}
			}
		}
		img.Batchers[constants.ParticleKey].Sprites["target"].Draw(target, mp.CellTrans[mp.Hover.Y][mp.Hover.X].Mat)
	}
}

func (mp *MinePuzzle) Solved() bool {
	return mp.solved
}