package structures

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

type pathDir struct {
	l     bool
	r     bool
	u     bool
	d     bool
	last  data.Direction
	width int
	wLeft bool
}

func SemiStraightPath(c *cave.Cave, start, end world.Coords, dir data.Direction, rb bool) ([]world.Coords, []world.Coords, []world.Coords) {
	var path, deadends, marked []world.Coords
	toMark := 0
	mark := false
	// initialize a starting width
	pDir := pathDir{
		last:  dir,
		width: random.CaveGen.Intn(3) + 1,
		wLeft: random.CaveGen.Intn(2) == 0,
	}
	curr := start
	done := false
	for !done {
		tile := c.GetTileInt(curr.X, curr.Y)
		if tile != nil {
			tile.Path = true
			WallUpWidth(tile, pDir.width, pDir.wLeft)
			path = append(path, curr)
			if curr == end {
				done = true
			}
		}
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

		mark = world.DistanceOrthogonal(start, curr) > 16 && world.DistanceOrthogonal(end, curr) > 16
		// if we are within 8 of the end in both directions, head straight there
		if util.Abs(curr.X-end.X) < 8 && util.Abs(curr.Y-end.Y) < 8 {
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
					if curr.X-end.X > 0 {
						t = 25 + (curr.X-end.X)/5
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
					if curr.X-end.X < 0 {
						t = 25 + (end.X-curr.X)/5
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
					if curr.Y-end.Y > 0 {
						t = 25 + (curr.Y-end.Y)/5
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
					if curr.Y-end.Y < 0 {
						t = 25 + (end.Y-curr.Y)/5
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
		if tile != nil {
			// if the new direction is the opposite of the last direction, it's a dead end
			if util.Abs(int(pDir.last-n)) == 2 {
				deadends = append(deadends, curr)
				tile.DeadEnd = true
			} else if mark {
				// if it's not a dead end, see if we should mark it
				var ok bool
				if ok, toMark = addToMarked(toMark); ok {
					marked = append(marked, curr)
					tile.Marked = true
				}
			}
		}
		// move to the next tile
		curr = moveToNextTile(curr, n)
		pDir.last = n
		// maybe change width
		pDir.width, pDir.wLeft = changeWidth(pDir.width, pDir.wLeft)
		// wall up all tiles surrounding the touched tiles
	}
	return path, deadends, marked
}

func BranchOff(c *cave.Cave, start world.Coords, min, max int) ([]world.Coords, []world.Coords, []world.Coords) {
	var path, deadends, marked []world.Coords
	toMark := 0
	// take off in that direction a random amount or until stopped by the edge
	dir := RandomDirection()
	for dir == data.Down {
		dir = RandomDirection()
	}
	pDir := pathDir{
		last:  dir,
		width: random.CaveGen.Intn(3) + 1,
		wLeft: random.CaveGen.Intn(2) == 0,
	}
	curr := start
	done := false
	for !done {
		tile := c.GetTileInt(curr.X, curr.Y)
		if tile != nil {
			tile.Path = true
			WallUpWidth(tile, pDir.width, pDir.wLeft)
			path = append(path, curr)
		}
		// check the direction we are going, or the length of path so far
		// if close to the edge, we're done. If the length is at the max, we're done.
		// if the length is greater than min, 1/20 chance it ends each step
		if dir == data.Left && curr.X < c.Left*constants.ChunkSize+6 {
			done = true
		}
		if dir == data.Right && curr.X > (c.Right+1)*constants.ChunkSize-6 {
			done = true
		}
		if dir == data.Up && curr.Y < 6 {
			done = true
		}
		if len(path)+1 >= max {
			done = true
		}
		if len(path)+1 >= min && random.CaveGen.Intn(20) == 0 {
			done = true
		}
		if tile != nil {
			// once we're done, it's a dead end
			if done {
				deadends = append(deadends, curr)
				tile.DeadEnd = true
			}
			// toMark
			var ok bool
			if ok, toMark = addToMarked(toMark); ok && !done {
				marked = append(marked, curr)
				tile.Marked = true
			}
		}
		// move to the next tile
		curr = moveToNextTile(curr, dir)
		// maybe change width
		pDir.width, pDir.wLeft = changeWidth(pDir.width, pDir.wLeft)
	}
	return path, deadends, marked
}

func moveToNextTile(curr world.Coords, n data.Direction) world.Coords {
	if n == data.Left {
		curr.X -= 1
	} else if n == data.Right {
		curr.X += 1
	} else if n == data.Up {
		curr.Y -= 1
	} else {
		curr.Y += 1
	}
	return curr
}

// 5% chance to changeWidth
func changeWidth(w int, l bool) (int, bool) {
	if random.CaveGen.Intn(20) == 0 {
		two := random.CaveGen.Intn(3) // the choice when the path is 2 wide
		if w == 3 || w == 1 {
			w = 2
			l = random.CaveGen.Intn(2) == 0
		} else if two == 0 { // switch to the "other" two wide
			l = !l
		} else if two == 1 { // switch to one wide
			w = 1
		} else { // switch to three wide
			w = 3
		}
	}
	return w, l
}

// randomly see if we should add this tile to "marked"
// if yes: decrease toMark by 25
// if no: increase toMark by 2
func addToMarked(toMark int) (bool, int) {
	if random.CaveGen.Intn(100) <= toMark {
		toMark -= 25
		if toMark < 0 {
			toMark = 0
		}
		return true, toMark
	} else {
		toMark += 2
		if toMark > 100 {
			toMark = 100
		}
		return false, toMark
	}
}
