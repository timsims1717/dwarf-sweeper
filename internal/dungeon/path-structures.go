package dungeon

import (
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
)

type Direction int

const (
	Left = iota
	Right
	Up
	Down
)

func SemiStraightPath(cave *Cave, start, end world.Coords, dir Direction, rb bool) []world.Coords {
	var path []world.Coords
	type pathDir struct {
		l     bool
		r     bool
		u     bool
		d     bool
		last  Direction
		width int
		wLeft bool
	}
	pDir := pathDir{
		last:  dir,
		width: random.CaveGen.Intn(3) + 1,
		wLeft: random.CaveGen.Intn(2) == 0,
	}
	curr := start
	tile := cave.GetTileInt(curr.X, curr.Y)
	wallUp(tile, pDir.width < 3 && rb)
	done := false
	for !done {
		pDir.l = true
		pDir.r = true
		pDir.u = true
		pDir.d = true
		if curr.X < cave.left*ChunkSize+6 {
			pDir.l = false
		}
		if curr.X > (cave.right+1)*ChunkSize-6 {
			pDir.r = false
		}
		if curr.Y < 6 {
			pDir.u = false
		}
		if curr.Y > (cave.bottom+1)*ChunkSize-6 {
			pDir.d = false
		}
		n := pDir.last

		if util.Abs(curr.X - end.X) < 8 && util.Abs(curr.Y - end.Y) < 8 {
			if curr.Y > end.Y {
				n = Up
			} else if curr.Y < end.Y {
				n = Down
			} else if curr.X > end.X {
				n = Left
			} else {
				n = Right
			}
		} else {
			if (n == Left && !pDir.l) ||
				(n == Right && !pDir.r) ||
				(n == Up && !pDir.u) ||
				(n == Down && !pDir.d) ||
				random.CaveGen.Intn(20) == 0 {
				tC := 0
				lC := 0
				rC := 0
				uC := 0
				dC := 0
				var t int
				if pDir.l {
					if curr.X - end.X > 0 {
						t = 25 + (curr.X - end.X) / 5
					} else {
						t = 5 + curr.X - end.X
					}
					if t > 0 {
						lC += t
						rC += t
						uC += t
						dC += t
						tC += t
					}
				}
				if pDir.r {
					if curr.X - end.X < 0 {
						t = 25 + (end.X - curr.X) / 5
					} else {
						t = 5 + end.X - curr.X
					}
					if t > 0 {
						rC += t
						uC += t
						dC += t
						tC += t
					}
				}
				if pDir.u {
					if curr.Y - end.Y > 0 {
						t = 25 + (curr.Y - end.Y) / 5
					} else {
						t = 5 + curr.Y - end.Y
					}
					if t > 0 {
						uC += t
						dC += t
						tC += t
					}
				}
				if pDir.d {
					if curr.Y - end.Y < 0 {
						t = 25 + (end.Y - curr.Y) / 5
					} else {
						t = 5 + end.Y - curr.Y
					}
					if t > 0 {
						dC += t
						tC += t
					}
				}
				c := random.CaveGen.Intn(tC)
				if c < lC {
					n = Left
				} else if c < rC {
					n = Right
				} else if c < uC {
					n = Up
				} else {
					n = Down
				}
			}
		}
		if n == Left {
			curr.X -= 1
		} else if n == Right {
			curr.X += 1
		} else if n == Up {
			curr.Y -= 1
		} else {
			curr.Y += 1
		}
		pDir.last = n
		if random.CaveGen.Intn(20) == 0 {
			two := random.CaveGen.Intn(3)
			if pDir.width == 3 || pDir.width == 1 {
				pDir.width = 2
				pDir.wLeft = random.CaveGen.Intn(2) == 0
			} else if two == 0 {
				pDir.wLeft = !pDir.wLeft
			} else if two == 1 {
				pDir.width = 1
			} else {
				pDir.width = 3
			}
		}
		tile = cave.GetTileInt(curr.X, curr.Y)
		wallUp(tile, pDir.width < 3 && rb)
		ns := tile.SubCoords.Neighbors()
		if pDir.width == 3 || (pDir.width == 2 && pDir.wLeft) {
			z := tile.Chunk.Get(ns[4])
			wallUp(z, pDir.width < 3 && rb)
			y := tile.Chunk.Get(ns[5])
			wallUp(y, pDir.width < 3 && rb)
			x := tile.Chunk.Get(ns[6])
			wallUp(x, pDir.width < 3 && rb)
			w := tile.Chunk.Get(ns[7])
			wallUp(w, pDir.width < 3 && rb)
		}
		if pDir.width == 3 || (pDir.width == 2 && !pDir.wLeft) {
			v := tile.Chunk.Get(ns[0])
			wallUp(v, pDir.width < 3 && rb)
			u := tile.Chunk.Get(ns[1])
			wallUp(u, pDir.width < 3 && rb)
			t := tile.Chunk.Get(ns[2])
			wallUp(t, pDir.width < 3 && rb)
			s := tile.Chunk.Get(ns[3])
			wallUp(s, pDir.width < 3 && rb)
		}
		path = append(path, curr)
		if curr == end {
			done = true
		}
	}
	return path
}

func BranchOff(start world.Coords, dir Direction) {
	// take off in that direction a random amount or until stopped by the edge
	// return path
}

func wallUp(tile *Tile, noBomb bool) {
	if tile != nil && !tile.neverChange && !tile.isChanged {
		tile.Solid = true
		tile.Type = Block
		tile.breakable = true
		if noBomb {
			tile.bomb = false
			tile.Entity = nil
		}
		tile.isChanged = true
		tile.UpdateSprites()
		for _, n := range tile.SubCoords.Neighbors() {
			t := tile.Chunk.Get(n)
			if t != nil && !t.neverChange && !t.isChanged {
				t.Solid = true
				t.Type = Wall
				t.breakable = false
				t.UpdateSprites()
			}
		}
	}
}