package cave

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
)

type Cave struct {
	RChunks     map[world.Coords]*Chunk
	LChunks     map[world.Coords]*Chunk
	Pivot       pixel.Vec
	Finite      bool
	UpdateBatch bool
	Batcher     *img.Batcher
	Left        int
	Right       int
	Bottom      int
	StartC      world.Coords
	ExitC       world.Coords
	BombPMin    float64
	BombPMax    float64
	FuseLen     float64
	GemRate     int
	ItemRate    int
}

func NewCave(batcher *img.Batcher, finite bool) *Cave {
	return &Cave{
		RChunks:     make(map[world.Coords]*Chunk),
		LChunks:     make(map[world.Coords]*Chunk),
		Batcher:     batcher,
		Finite:      finite,
		UpdateBatch: true,
	}
}

func (c *Cave) Update() {
	p := WorldToChunk(c.Pivot)
	all := append([]world.Coords{p}, p.Neighbors()...)
	for i, chunk := range c.RChunks {
		dis := world.CoordsIn(i, all)
		if dis && !chunk.Display {
			chunk.Reload = true
			c.UpdateBatch = true
		}
		chunk.Display = dis
		chunk.Update()
	}
	for i, chunk := range c.LChunks {
		dis := world.CoordsIn(i, all)
		if dis && !chunk.Display {
			chunk.Reload = true
			c.UpdateBatch = true
		}
		chunk.Display = dis
		chunk.Update()
	}
	//for _, chunk := range c.RChunks {
	//	chunk.display = true
	//	chunk.Update()
	//}
	//for _, chunk := range c.LChunks {
	//	chunk.display = true
	//	chunk.Update()
	//}
}

func (c *Cave) Draw(win *pixelgl.Window) {
	if c.UpdateBatch {
		c.Batcher.Clear()
		for _, chunk := range c.RChunks {
			chunk.Draw(c.Batcher.Batch())
		}
		for _, chunk := range c.LChunks {
			chunk.Draw(c.Batcher.Batch())
		}
	}
	c.Batcher.Draw(win)
	c.UpdateBatch = false
}

func (c *Cave) Dimensions() (int, int) {
	if c.Finite {
		return (c.Right - c.Left + 1) * constants.ChunkSize, (c.Bottom + 1) * constants.ChunkSize
	} else {
		return -1, -1
	}
}

func (c *Cave) PointLoaded(v pixel.Vec) bool {
	return c.GetChunk(WorldToChunk(v)).Display
}

func (c *Cave) CurrentBoundaries() (pixel.Vec, pixel.Vec) {
	p := WorldToChunk(c.Pivot)
	var all []world.Coords
	if c.Finite {
		for _, chunk := range c.RChunks {
			all = append(all, chunk.Coords)
		}
		for _, chunk := range c.LChunks {
			all = append(all, chunk.Coords)
		}
	} else {
		all = append([]world.Coords{p}, p.Neighbors()...)
	}
	x1 := 10000000.
	y1 := 10000000.
	x2 := -10000000.
	y2 := -10000000.
	for _, i := range all {
		if i.X >= 0 && i.Y >= 0 {
			if chunk, ok := c.RChunks[i]; ok {
				tr := chunk.Rows[0][constants.ChunkSize-1].Transform.Pos
				bl := chunk.Rows[constants.ChunkSize-1][0].Transform.Pos
				if bl.X < x1 {
					x1 = bl.X
				}
				if bl.Y < y1 {
					y1 = bl.Y
				}
				if tr.X > x2 {
					x2 = tr.X
				}
				if tr.Y > y2 {
					y2 = tr.Y
				}
			}
		} else if i.X < 0 && i.Y >= 0 {
			if chunk, ok := c.LChunks[i]; ok {
				tr := chunk.Rows[0][constants.ChunkSize-1].Transform.Pos
				bl := chunk.Rows[constants.ChunkSize-1][0].Transform.Pos
				if bl.X < x1 {
					x1 = bl.X
				}
				if bl.Y < y1 {
					y1 = bl.Y
				}
				if tr.X > x2 {
					x2 = tr.X
				}
				if tr.Y > y2 {
					y2 = tr.Y
				}
			}
		}
	}
	return pixel.V(x1, y1), pixel.V(x2, y2)
}

func (c *Cave) GetTileInt(x, y int) *Tile {
	cX := x / constants.ChunkSize
	if x < 0 {
		cX = (x + 1) / constants.ChunkSize
		cX--
	}
	tX := x % constants.ChunkSize
	if tX < 0 {
		tX += constants.ChunkSize
	}
	cY := y / constants.ChunkSize
	tY := y % constants.ChunkSize
	return c.GetChunk(world.Coords{X: cX, Y: cY}).Get(world.Coords{X: tX, Y: tY})
}

func (c *Cave) GetChunk(coords world.Coords) *Chunk {
	if chunkR, okR := c.RChunks[coords]; okR {
		return chunkR
	} else if chunkL, okL := c.LChunks[coords]; okL {
		return chunkL
	} else {
		return nil
	}
}

func (c *Cave) GetTile(v pixel.Vec) *Tile {
	ch := WorldToChunk(v)
	tl := WorldToTile(v, ch.X < 0)
	chunk := c.GetChunk(ch)
	return chunk.Get(tl)
}

func (c *Cave) GetStart() *Tile {
	return c.GetTileInt(c.StartC.X, c.StartC.Y)
}

func (c *Cave) GetExit() *Tile {
	return c.GetTileInt(c.ExitC.X, c.ExitC.Y)
}

func (c *Cave) MarkAsNotChanged() {
	for _, chunk := range c.RChunks {
		for _, row := range chunk.Rows {
			for _, tile := range row {
				tile.IsChanged = false
			}
		}
	}
	for _, chunk := range c.LChunks {
		for _, row := range chunk.Rows {
			for _, tile := range row {
				tile.IsChanged = false
			}
		}
	}
}

func WorldToChunk(v pixel.Vec) world.Coords {
	if v.X >= 0 - world.TileSize * 0.5 {
		return world.Coords{X: int((v.X+world.TileSize*0.5) / constants.ChunkSize / world.TileSize), Y: int(-(v.Y-world.TileSize*0.5) / constants.ChunkSize / world.TileSize)}
	} else {
		return world.Coords{X: int((v.X+world.TileSize*0.5) /constants.ChunkSize/ world.TileSize)-1, Y: int(-(v.Y-world.TileSize*0.5) / constants.ChunkSize / world.TileSize)}
	}
}

func WorldToTile(v pixel.Vec, left bool) world.Coords {
	x, y := world.WorldToMap(v.X+world.TileSize*0.5, -(v.Y-world.TileSize*0.5))
	x = x % constants.ChunkSize
	y = y % constants.ChunkSize
	if left {
		x = (constants.ChunkSize - (util.Abs(x) + 1)) % constants.ChunkSize
	}
	return world.Coords{
		X: x % constants.ChunkSize,
		Y: y % constants.ChunkSize,
	}
}

func TileInTile(a, b pixel.Vec) bool {
	return math.Abs(a.X - b.X) <= world.TileSize * 0.5 && math.Abs(a.Y - b.Y) <= world.TileSize * 0.5
}

func (c *Cave) PrintCaveToTerminal() {
	if c.Finite {
		fmt.Println("Printing cave ... ")
		fmt.Println()
		for y := 0; y < (c.Bottom+ 1) *constants.ChunkSize; y++ {
			for x := c.Left * constants.ChunkSize; x < (c.Right+1)*constants.ChunkSize; x++ {
				tile := c.GetTileInt(x, y)
				if tile != nil {
					switch tile.Type {
					case Block, Value:
						if tile.Bomb {
							fmt.Print("ó")
						} else {
							fmt.Print("□")
						}
					case Wall:
						fmt.Print("▣")
					case Deco:
						fmt.Print("*")
					case Empty:
						fmt.Print(" ")
					}
				}
			}
			fmt.Print("\n")
		}
	}
}