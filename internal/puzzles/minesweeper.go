package puzzles

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/menubox"
	"dwarf-sweeper/internal/minesweeper"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/input"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/typeface"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
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

	Box        *menubox.MenuBox
	InfoText   *typeface.Text
	InfoSyms   []pixel.Vec
	CountText  *text.Text
	CountTrans *transform.Transform

	CellTrans [][]*transform.Transform
	solved    bool
	OnSolveFn func()

	FlagAnim *reanimator.Tree
	ExAnim   *reanimator.Tree

	ButtonPressed bool
}

func (mp *MinePuzzle) Create(cam *camera.Camera, level int) {
	mp.solved = false
	size := 5 + level / 4
	mp.SizeW = util.Min(size, 9)
	mp.SizeH = util.Min(size, 6)
	mp.InfoText = typeface.New(camera.Cam,"main", typeface.NewAlign(typeface.Left, typeface.Center), 2.0, constants.ActualHintSize, 0., 0.)
	mp.InfoText.SetColor(constants.DefaultColor)
	mp.InfoText.SetText("{symbol:flag}:mark a bomb tile\n{symbol:dig}:mark a safe tile")
	mp.CountText = text.New(pixel.ZV, typeface.Atlases["main"])
	mp.CountText.LineHeight *= 1.2
	mp.CountText.Color = constants.DefaultColor
	mp.Box = menubox.NewBox(cam, 1.4)
	mp.Box.SetSize(pixel.R(0., 0., float64(mp.SizeW) * (world.TileSize + 2.) * scalar, float64(mp.SizeH) * (world.TileSize + 2.) * scalar + mp.InfoText.Height))
	mp.InfoText.SetPos(pixel.V(mp.Box.Rect.W() * -0.5 + 5., mp.Box.Rect.H() * 0.5 - 15.))
	mp.CountTrans = transform.New()
	mp.CountTrans.Pos.X = mp.Box.Rect.W() * 0.5
	mp.CountTrans.Pos.Y = mp.Box.Rect.H() * 0.5 - mp.CountText.LineHeight * constants.ActualHintSize
	mp.CountTrans.Scalar = pixel.V(constants.ActualHintSize, constants.ActualHintSize)
	topPadding := mp.InfoText.Height * 0.5
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
		yPos := multY*(world.TileSize*scalar+1)-topPadding
		for x := 0; x < mp.SizeW; x++ {
			evenX := mp.SizeW % 2 == 0
			midX := mp.SizeW / 2
			var multX float64
			if evenX {
				multX = float64(x-midX)+0.5
			} else {
				multX = float64(x-midX)
			}
			trans := transform.New()
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
	minesweeper.RevealTilSolvableP(mp.Board, random.Effects, false)
	minesweeper.UnRevealWhileSolvableP(mp.Board, random.Effects, false)
	mp.Box.Open()
	mp.Hover.Y = mp.SizeH-1
}

func (mp *MinePuzzle) Close() {
	mp.Box.Close()
}

func (mp *MinePuzzle) Update(in *input.Input) {
	mp.Box.Update()
	mp.UpdateTransforms()
	mp.InfoText.Update()
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
							sfx.SoundPlayer.PlaySound("click", 2.0)
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
		mp.ButtonPressed = in.Get("flag").Pressed() || in.Get("dig").Pressed()
		if in.Get("flag").JustReleased() && !cell.Rev && !cell.Ex {
			in.Get("flag").Consume()
			cell.Flag = !cell.Flag
		} else if in.Get("dig").JustReleased() && !cell.Rev && !cell.Flag {
			in.Get("dig").Consume()
			if cell.Bomb {
				cell.Rev = true
				sfx.SoundPlayer.PlaySound("mpwrong", 1.)
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
		mp.CountTrans.UIZoom = mp.Box.Cam.GetZoomScale()
		mp.CountTrans.UIPos = mp.Box.Cam.APos
		mp.CountTrans.Update()
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
		mp.InfoText.Draw(target)
		mp.CountText.Draw(target, mp.CountTrans.Mat)
		for y, row := range mp.CellTrans {
			for x, ct := range row {
				cell := mp.Board.Board[y][x]
				if cell.Rev {
					img.Batchers[constants.PuzzleKey].Sprites["background_num"].Draw(target, ct.Mat)
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
				} else {
					if mp.ButtonPressed && mp.Hover.X == x && mp.Hover.Y == y {
						img.Batchers[constants.PuzzleKey].Sprites["background_click"].Draw(target, ct.Mat)
					} else {
						img.Batchers[constants.PuzzleKey].Sprites["background_empty"].Draw(target, ct.Mat)
					}
					if cell.Flag {
						//img.Batchers[constants.EntityKey].Sprites["mine_1"].Draw(target, ct.Mat)
						mp.FlagAnim.Draw(target, ct.Mat)
					} else if cell.Ex {
						mp.ExAnim.Draw(target, ct.Mat)
					}
				}
			}
		}
		img.Batchers[constants.ParticleKey].Sprites["target"].Draw(target, mp.CellTrans[mp.Hover.Y][mp.Hover.X].Mat)
	}
}

func (mp *MinePuzzle) Solved() bool {
	return mp.solved
}

func (mp *MinePuzzle) OnSolve() {
	if mp.OnSolveFn != nil {
		mp.OnSolveFn()
	}
}