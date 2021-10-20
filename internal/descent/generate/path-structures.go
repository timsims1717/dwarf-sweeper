package generate

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
)

type Path struct {
	Dir   data.Direction
	Count int
}

func SemiStraightPath(c *cave.Cave, start, end world.Coords, dir data.Direction, rb bool) ([]world.Coords, []world.Coords) {
	var path []world.Coords
	var deadends []world.Coords
	type pathDir struct {
		l     bool
		r     bool
		u     bool
		d     bool
		last  data.Direction
		width int
		wLeft bool
	}
	// initialize a starting width
	pDir := pathDir{
		last:  dir,
		width: random.CaveGen.Intn(3) + 1,
		wLeft: random.CaveGen.Intn(2) == 0,
	}
	curr := start
	tile := c.GetTileInt(curr.X, curr.Y)
	// after finding each path tile, we will wall up all the un-pathed tiles around it
	wallUp(tile, pDir.width < 3 && rb)
	done := false
	for !done {
		// check each direction, see if we can go there
		pDir.l = true
		pDir.r = true
		pDir.u = true
		pDir.d = true
		if curr.X < c.Left*constants.ChunkSize+6 {
			pDir.l = false
		}
		if curr.X > (c.Right+1)*constants.ChunkSize-6 {
			pDir.r = false
		}
		if curr.Y < 6 {
			pDir.u = false
		}
		if curr.Y > (c.Bottom+1)*constants.ChunkSize-6 {
			pDir.d = false
		}
		n := pDir.last

		// if we are within 8 of the end in both directions, head straight there
		if util.Abs(curr.X - end.X) < 8 && util.Abs(curr.Y - end.Y) < 8 {
			if curr.Y > end.Y {
				n = data.Up
			} else if curr.Y < end.Y {
				n = data.Down
			} else if curr.X > end.X {
				n = data.Left
			} else {
				n = data.Right
			}
		} else {
			// if we can't go the direction we were going, or in a 1/20 chance ...
			if (n == data.Left && !pDir.l) ||
				(n == data.Right && !pDir.r) ||
				(n == data.Up && !pDir.u) ||
				(n == data.Down && !pDir.d) ||
				random.CaveGen.Intn(20) == 0 {
				// choose a new direction, weighted to get us closer to the end
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
				choice := random.CaveGen.Intn(tC)
				if choice < lC {
					n = data.Left
				} else if choice < rC {
					n = data.Right
				} else if choice < uC {
					n = data.Up
				} else {
					n = data.Down
				}
			}
		}
		// if the new direction is the opposite of the last direction, it's a deadend
		if util.Abs(int(pDir.last - n)) == 2 {
			deadends = append(deadends, curr)
		}
		// move to the next tile
		if n == data.Left {
			curr.X -= 1
		} else if n == data.Right {
			curr.X += 1
		} else if n == data.Up {
			curr.Y -= 1
		} else {
			curr.Y += 1
		}
		pDir.last = n
		// maybe change width
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
		// wall up all tiles surrounding the touched tiles
		tile = c.GetTileInt(curr.X, curr.Y)
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
	return path, deadends
}

func BranchOff(start world.Coords, dir data.Direction) {
	// take off in that direction a random amount or until stopped by the edge
	// return path
}