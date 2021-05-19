package dungeon

import (
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"math/rand"
)

type Direction int

const (
	Left = iota
	Right
	Up
	Down
)

func SemiStraightPath(cave *Cave, start, end world.Coords, dir Direction) []world.Coords {
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
		width: rand.Intn(3) + 1,
		wLeft: rand.Intn(2) == 0,
	}
	curr := start
	tile := cave.GetTileInt(curr.X, curr.Y)
	if tile != nil && !tile.neverChange && !tile.isChanged {
		tile.Solid = false
		tile.Type = Empty
		tile.BGSprite = nil
		tile.breakable = false
		tile.bomb = false
		tile.neverChange = true
		tile.isChanged = true
		tile.Entities = []Entity{}
		tile.UpdateSprites()
		path = append(path, curr)
	}
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
				rand.Intn(20) == 0 {
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
				c := rand.Intn(tC)
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
		if rand.Intn(20) == 0 {
			two := rand.Intn(3)
			if pDir.width == 3 || pDir.width == 1 {
				pDir.width = 2
				pDir.wLeft = rand.Intn(2) == 0
			} else if two == 0 {
				pDir.wLeft = !pDir.wLeft
			} else if two == 1 {
				pDir.width = 1
			} else {
				pDir.width = 3
			}
		}
		tile = cave.GetTileInt(curr.X, curr.Y)
		wallUp(tile, pDir.width == 1)
		ns := tile.SubCoords.Neighbors()
		if pDir.width == 3 || (pDir.width == 2 && pDir.wLeft) {
			z := tile.Chunk.Get(ns[4])
			wallUp(z, false)
			y := tile.Chunk.Get(ns[5])
			wallUp(y, false)
			x := tile.Chunk.Get(ns[6])
			wallUp(x, false)
			w := tile.Chunk.Get(ns[7])
			wallUp(w, false)
		}
		if pDir.width == 3 || (pDir.width == 2 && !pDir.wLeft) {
			v := tile.Chunk.Get(ns[0])
			wallUp(v, false)
			u := tile.Chunk.Get(ns[1])
			wallUp(u, false)
			t := tile.Chunk.Get(ns[2])
			wallUp(t, false)
			s := tile.Chunk.Get(ns[3])
			wallUp(s, false)
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
			tile.Entities = []Entity{}
		}
		tile.neverChange = true
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