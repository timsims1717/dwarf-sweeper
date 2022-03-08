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
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/typeface"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
)

const (
	scalar = 1.8
	timer  = 60.
)

type MinePuzzle struct {
	SizeW int
	SizeH int
	Level int

	Misses int
	Timer  *timing.Timer
	start  bool

	Board *minesweeper.Board
	Hover world.Coords

	Box        *menubox.MenuBox
	InfoText   *typeface.Text
	TimerText  *typeface.Text
	miss1Trans *transform.Transform
	miss2Trans *transform.Transform
	miss3Trans *transform.Transform

	CellTrans [][]*transform.Transform
	solved    bool
	failed    bool
	OnSolveFn func()
	OnFailFn  func()

	FlagAnim *reanimator.Tree
	ExAnim   *reanimator.Tree

	ButtonPressed   bool
	ButtonCancelled bool
}

func (mp *MinePuzzle) Create(cam *camera.Camera, level int) {
	mp.solved = false
	size := 5 + level / 4
	mp.SizeW = util.Min(size, 9)
	mp.SizeH = util.Min(size, 6)
	mp.InfoText = typeface.New(camera.Cam,"main", typeface.NewAlign(typeface.Left, typeface.Center), 2.0, constants.ActualHintSize, 0., 0.)
	mp.InfoText.SetColor(constants.DefaultColor)
	mp.InfoText.SetText("{symbol:flag}:mark bomb\n{symbol:dig}:mark safe")
	mp.TimerText = typeface.New(camera.Cam,"main", typeface.NewAlign(typeface.Right, typeface.Center), 2.0, constants.ActualHintSize, 0., 0.)
	mp.TimerText.Color = constants.DefaultColor
	mp.TimerText.SetText("\n")
	mp.Box = menubox.NewBox(cam, 1.4)
	mp.Box.SetSize(pixel.R(0., 0., float64(mp.SizeW) * (world.TileSize + 2.) * scalar, float64(mp.SizeH) * (world.TileSize + 2.) * scalar + mp.InfoText.Height))
	mp.InfoText.SetPos(pixel.V(mp.Box.Rect.W() * -0.5 + 5., mp.Box.Rect.H() * 0.5 - 15.))
	mp.TimerText.SetPos(pixel.V(mp.Box.Rect.W() * 0.5 - 5., mp.Box.Rect.H() * 0.5 - 15.))
	mp.miss1Trans = transform.New()
	mp.miss1Trans.Pos.X = mp.Box.Rect.W() * 0.5 - 5. - 4. * scalar
	mp.miss1Trans.Pos.Y = mp.Box.Rect.H() * 0.5 - 15. + 2. * scalar
	mp.miss1Trans.Scalar = pixel.V(scalar, scalar)
	mp.miss2Trans = transform.New()
	mp.miss2Trans.Pos.X = mp.Box.Rect.W() * 0.5 - 5. - 13. * scalar
	mp.miss2Trans.Pos.Y = mp.Box.Rect.H() * 0.5 - 15. + 2. * scalar
	mp.miss2Trans.Scalar = pixel.V(scalar, scalar)
	mp.miss3Trans = transform.New()
	mp.miss3Trans.Pos.X = mp.Box.Rect.W() * 0.5 - 5. - 21.5 * scalar
	mp.miss3Trans.Pos.Y = mp.Box.Rect.H() * 0.5 - 15. + 2. * scalar
	mp.miss3Trans.Scalar = pixel.V(scalar, scalar)
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
	mp.ExAnim = reanimator.NewSimple(reanimator.NewAnimFromSprites("ex", img.Batchers[constants.PuzzleKey].GetAnimation("green_x").S, reanimator.Loop))
}

func (mp *MinePuzzle) IsOpen() bool {
	return mp.Box.IsOpen()
}

func (mp *MinePuzzle) IsClosed() bool {
	return mp.Box.IsClosed()
}

func (mp *MinePuzzle) Open(pCode string) {
	mp.InfoText.SetText(fmt.Sprintf("{symbol:%s-flag}:mark bomb\n{symbol:%s-dig}:mark safe", pCode, pCode))
	mp.FlagAnim = reanimator.NewSimple(reanimator.NewAnimFromSprites("flag", img.Batchers[constants.PuzzleKey].GetAnimation(fmt.Sprintf("flag_hang_%s", pCode)).S, reanimator.Loop))
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
	if mp.Box.IsOpen() && !mp.start {
		mp.Timer = timing.New(timer)
		mp.start = true
	}
	if mp.Timer != nil {
		mp.Timer.Update()
		timeLeft := timer - mp.Timer.Elapsed()
		if timeLeft < 0. {
			timeLeft = 0.
		}
		secs := int(timeLeft)
		min := secs / 60
		sec := secs % 60
		mp.TimerText.SetText(fmt.Sprintf("\n%02d:%02d", min, sec))
	}
	mp.UpdateTransforms()
	mp.InfoText.Update()
	mp.TimerText.Update()
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
		} else if !mp.ButtonPressed {
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
		}
		cell := mp.Board.Board[mp.Hover.Y][mp.Hover.X]
		if !mp.ButtonCancelled {
			mp.ButtonCancelled = in.Get("flag").Pressed() && in.Get("dig").Pressed()
			if !mp.ButtonCancelled {
				mp.ButtonPressed = in.Get("flag").Pressed() || in.Get("dig").Pressed()
				if in.Get("flag").JustReleased() && !cell.Rev && !cell.Ex {
					in.Get("flag").Consume()
					cell.Flag = !cell.Flag
				} else if in.Get("dig").JustReleased() && !cell.Rev && !cell.Flag && !cell.Ex {
					in.Get("dig").Consume()
					if cell.Bomb {
						cell.Rev = true
						sfx.SoundPlayer.PlaySound("mpwrong", 1.)
						mp.CellTrans[mp.Hover.Y][mp.Hover.X].Shake(random.Effects)
						mp.Misses++
					} else {
						cell.Ex = true
						sfx.SoundPlayer.PlaySound("mpcorrect", 1.)
					}
				}
			}
		} else if !in.Get("flag").Pressed() && !in.Get("dig").Pressed() {
			mp.ButtonCancelled = false
			mp.ButtonPressed = false
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
		mp.failed = mp.Misses > 2 || mp.Timer.Done()
	}
}

func (mp *MinePuzzle) UpdateTransforms() {
	if mp.Box.Cam != nil {
		mp.miss1Trans.UIZoom = mp.Box.Cam.GetZoomScale()
		mp.miss1Trans.UIPos = mp.Box.Cam.APos
		mp.miss1Trans.Update()
		mp.miss2Trans.UIZoom = mp.Box.Cam.GetZoomScale()
		mp.miss2Trans.UIPos = mp.Box.Cam.APos
		mp.miss2Trans.Update()
		mp.miss3Trans.UIZoom = mp.Box.Cam.GetZoomScale()
		mp.miss3Trans.UIPos = mp.Box.Cam.APos
		mp.miss3Trans.Update()
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
		mp.TimerText.Draw(target)
		img.Batchers[constants.PuzzleKey].GetSprite("miss_box").Draw(target, mp.miss1Trans.Mat)
		if mp.Misses > 0 {
			img.Batchers[constants.PuzzleKey].GetSprite("miss_x").Draw(target, mp.miss1Trans.Mat)
		}
		img.Batchers[constants.PuzzleKey].GetSprite("miss_box").Draw(target, mp.miss2Trans.Mat)
		if mp.Misses > 1 {
			img.Batchers[constants.PuzzleKey].GetSprite("miss_x").Draw(target, mp.miss2Trans.Mat)
		}
		img.Batchers[constants.PuzzleKey].GetSprite("miss_box").Draw(target, mp.miss3Trans.Mat)
		if mp.Misses > 2 {
			img.Batchers[constants.PuzzleKey].GetSprite("miss_x").Draw(target, mp.miss3Trans.Mat)
		}
		for y, row := range mp.CellTrans {
			for x, ct := range row {
				cell := mp.Board.Board[y][x]
				if cell.Rev {
					if cell.Bomb {
						img.Batchers[constants.PuzzleKey].GetSprite("background_error").Draw(target, ct.Mat)
						img.Batchers[constants.EntityKey].GetSprite("mine_1").Draw(target, ct.Mat)
					} else {
						img.Batchers[constants.PuzzleKey].GetSprite("background_num").Draw(target, ct.Mat)
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
						img.Batchers[constants.PuzzleKey].GetSprite(str).Draw(target, ct.Mat)
					}
				} else if cell.Ex {
					img.Batchers[constants.PuzzleKey].GetSprite("background_empty").Draw(target, ct.Mat)
					mp.ExAnim.Draw(target, ct.Mat)
				} else {
					if mp.ButtonPressed && !mp.ButtonCancelled && mp.Hover.X == x && mp.Hover.Y == y {
						img.Batchers[constants.PuzzleKey].GetSprite("background_click").Draw(target, ct.Mat)
					} else {
						img.Batchers[constants.PuzzleKey].GetSprite("background_empty").Draw(target, ct.Mat)
					}
					if cell.Flag {
						mp.FlagAnim.Draw(target, ct.Mat)
					}
				}
			}
		}
		img.Batchers[constants.ParticleKey].GetSprite("target").Draw(target, mp.CellTrans[mp.Hover.Y][mp.Hover.X].Mat)
	}
}

func (mp *MinePuzzle) Solved() bool {
	return mp.solved
}

func (mp *MinePuzzle) Failed() bool {
	return mp.failed
}

func (mp *MinePuzzle) OnSolve() {
	if mp.OnSolveFn != nil {
		mp.OnSolveFn()
	}
}

func (mp *MinePuzzle) OnFail() {
	if mp.OnFailFn != nil {
		mp.OnFailFn()
	}
}

func (mp *MinePuzzle) SetOnSolve(fn func()) {
	mp.OnSolveFn = fn
}

func (mp *MinePuzzle) SetOnFail(fn func()) {
	mp.OnFailFn = fn
}