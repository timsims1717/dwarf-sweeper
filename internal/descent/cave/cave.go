package cave

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
)

type Cave struct {
	Chunks      map[world.Coords]*Chunk
	FillChunk   func(chunk *Chunk)
	Pivots      []pixel.Vec
	UpdateBatch bool
	Type        CaveType
	Biome       string
	Biomes      []string

	Left   int
	Right  int
	Bottom int
	Width  int
	Height int
	bl     pixel.Vec
	tr     pixel.Vec
	StartC world.Coords
	ExitC  world.Coords

	Paths     []world.Coords
	DeadEnds  []world.Coords
	Marked    []world.Coords
	Rooms     []world.Coords
	MainGroup int
	FillVar   float64

	BombPMin float64
	BombPMax float64
	GemRate  float64

	BGBatch    *pixel.Batch
	Background *pixel.Sprite
	BGTC       *transform.Transform
	BGTUL      *transform.Transform
	BGTU       *transform.Transform
	BGTUR      *transform.Transform
	BGTL       *transform.Transform
	BGTR       *transform.Transform
	BGTDL      *transform.Transform
	BGTD       *transform.Transform
	BGTDR      *transform.Transform

	Fog     bool
	LoadAll bool

	PathRule PathRule
}

func NewCave(biome string, caveType CaveType) *Cave {
	var bgSpr *pixel.Sprite
	var bgBatch *pixel.Batch
	bg, err := img.LoadImage(fmt.Sprintf("assets/img/the-%s-bg.png", biome))
	if err != nil {
		fmt.Printf("error loading %s biome background: %s\n", biome, err)
	} else {
		bgSpr = pixel.NewSprite(bg, bg.Bounds())
		bgBatch = pixel.NewBatch(&pixel.TrianglesData{}, bg)
	}

	return &Cave{
		Chunks:      make(map[world.Coords]*Chunk),
		Type:        caveType,
		UpdateBatch: true,
		Biome:       biome,
		Biomes:      []string{biome},
		Background:  bgSpr,
		BGBatch:     bgBatch,
		BGTC:        transform.New(),
		BGTUL:       transform.New(),
		BGTU:        transform.New(),
		BGTUR:       transform.New(),
		BGTL:        transform.New(),
		BGTR:        transform.New(),
		BGTDL:       transform.New(),
		BGTD:        transform.New(),
		BGTDR:       transform.New(),
		Fog:         true,
		GemRate:     0.05,
	}
}

func (c *Cave) SetSize(left, right, bottom int) {
	c.Left = left
	c.Right = right
	c.Bottom = bottom
	c.Width = (right - left + 1) * constants.ChunkSize
	c.Height = (bottom + 1) * constants.ChunkSize
	c.bl = pixel.V(float64(left*constants.ChunkSize)*world.TileSize, -float64((bottom+1)*constants.ChunkSize-1)*world.TileSize)
	c.tr = pixel.V(float64((right+1)*constants.ChunkSize)*world.TileSize, 0.)
	// how much we should fill is based on the size of the cave
	// a 32 chunk size at 3x3 gives a value of 72 for fillVar
	// a 16 chunk size at 3x3 gives a value of 18 for fillVar
	c.FillVar = float64(c.Width * c.Height / 128.)
}

func (c *Cave) Update() {
	var all []world.Coords
	for _, pivot := range c.Pivots {
		p := WorldToChunk(pivot)
		all = append(all, p)
		all = append(all, p.Neighbors()...)
	}
	for i, chunk := range c.Chunks {
		dis := world.CoordsIn(i, all)
		if dis && !chunk.Display {
			chunk.Reload = true
			c.UpdateBatch = true
		}
		chunk.Display = dis || c.LoadAll
		chunk.Update()
	}

	if c.UpdateBatch {
		c.UpdateAllTileSprites()
		for _, biome := range c.Biomes {
			img.Batchers[fmt.Sprintf(constants.CaveBGFMT, biome)].Clear()
			img.Batchers[biome].Clear()
		}
		img.Batchers[constants.FogKey].Clear()
	}

	if c.Background != nil {
		offset := camera.Cam.APos.Scaled(-0.5)
		offset.X = util.FMod(offset.X, c.Background.Frame().W())
		offset.Y = util.FMod(offset.Y, c.Background.Frame().H())
		offset.X = offset.X + c.Background.Frame().W()*0.5
		offset.Y = offset.Y - c.Background.Frame().H()*0.5
		c.BGTC.Pos = offset
		c.BGTUL.Pos = pixel.V(-c.Background.Frame().W()+offset.X, c.Background.Frame().H()+offset.Y)
		c.BGTU.Pos = pixel.V(offset.X, c.Background.Frame().H()+offset.Y)
		c.BGTUR.Pos = pixel.V(c.Background.Frame().W()+offset.X, c.Background.Frame().H()+offset.Y)
		c.BGTR.Pos = pixel.V(c.Background.Frame().W()+offset.X, offset.Y)
		c.BGTL.Pos = pixel.V(-c.Background.Frame().W()+offset.X, offset.Y)
		c.BGTDL.Pos = pixel.V(-c.Background.Frame().W()+offset.X, -c.Background.Frame().H()+offset.Y)
		c.BGTD.Pos = pixel.V(offset.X, -c.Background.Frame().H()+offset.Y)
		c.BGTDR.Pos = pixel.V(c.Background.Frame().W()+offset.X, -c.Background.Frame().H()+offset.Y)
	}
}

func (c *Cave) Draw(win *pixelgl.Window) {
	if c.Background != nil && c.BGBatch != nil {
		c.BGTC.UIPos = camera.Cam.APos
		c.BGTC.Update()
		c.BGTUL.UIPos = camera.Cam.APos
		c.BGTUL.Update()
		c.BGTU.UIPos = camera.Cam.APos
		c.BGTU.Update()
		c.BGTUR.UIPos = camera.Cam.APos
		c.BGTUR.Update()
		c.BGTL.UIPos = camera.Cam.APos
		c.BGTL.Update()
		c.BGTR.UIPos = camera.Cam.APos
		c.BGTR.Update()
		c.BGTDL.UIPos = camera.Cam.APos
		c.BGTDL.Update()
		c.BGTD.UIPos = camera.Cam.APos
		c.BGTD.Update()
		c.BGTDR.UIPos = camera.Cam.APos
		c.BGTDR.Update()

		c.BGBatch.Clear()
		c.Background.Draw(c.BGBatch, c.BGTC.Mat)
		c.Background.Draw(c.BGBatch, c.BGTUL.Mat)
		c.Background.Draw(c.BGBatch, c.BGTU.Mat)
		c.Background.Draw(c.BGBatch, c.BGTUR.Mat)
		c.Background.Draw(c.BGBatch, c.BGTL.Mat)
		c.Background.Draw(c.BGBatch, c.BGTR.Mat)
		c.Background.Draw(c.BGBatch, c.BGTDL.Mat)
		c.Background.Draw(c.BGBatch, c.BGTD.Mat)
		c.Background.Draw(c.BGBatch, c.BGTDR.Mat)
		c.BGBatch.Draw(win)
	}
	if c.UpdateBatch {
		for _, chunk := range c.Chunks {
			chunk.Draw()
		}
		c.UpdateBatch = false
	}
}

func (c *Cave) Dimensions() (int, int) {
	if c.Type == Infinite {
		return -1, -1
	} else {
		return (c.Right - c.Left + 1) * constants.ChunkSize, (c.Bottom + 1) * constants.ChunkSize
	}
}

func (c *Cave) PointLoaded(v pixel.Vec) bool {
	return c.GetChunk(WorldToChunk(v)).Display
}

func (c *Cave) CurrentBoundaries() (pixel.Vec, pixel.Vec) {
	if c.Type != Infinite {
		return c.bl, c.tr
	}
	var all []world.Coords
	for _, pivot := range c.Pivots {
		p := WorldToChunk(pivot)
		all = append(all, p)
		all = append(all, p.Neighbors()...)
	}
	x1 := 10000000.
	y1 := 10000000.
	x2 := -10000000.
	y2 := -10000000.
	for _, i := range all {
		if chunk, ok := c.Chunks[i]; ok {
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
	if chunk, ok := c.Chunks[coords]; ok {
		return chunk
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

func (c *Cave) UpdateAllTileSprites() {
	for _, chunk := range c.Chunks {
		if chunk.Display {
			for _, row := range chunk.Rows {
				for _, tile := range row {
					tile.UpdateDetails()
				}
			}
		}
	}
	for _, chunk := range c.Chunks {
		if chunk.Display {
			for _, row := range chunk.Rows {
				for _, tile := range row {
					tile.UpdateDetails()
				}
			}
		}
	}
	for _, chunk := range c.Chunks {
		if chunk.Display {
			for _, row := range chunk.Rows {
				for _, tile := range row {
					tile.UpdateSprites()
				}
			}
		}
	}
}

func (c *Cave) MarkAsNotChanged() {
	c.MapFn(func(tile *Tile) {
		tile.IsChanged = false
		tile.Change = false
	})
}

func WorldToChunk(v pixel.Vec) world.Coords {
	if v.X >= 0-world.TileSize*0.5 {
		return world.Coords{X: int((v.X+world.TileSize*0.5) / constants.ChunkSize / world.TileSize), Y: int(-(v.Y - world.TileSize*0.5) / constants.ChunkSize / world.TileSize)}
	} else {
		return world.Coords{X: int((v.X+world.TileSize*0.5) / constants.ChunkSize / world.TileSize) - 1, Y: int(-(v.Y - world.TileSize*0.5) / constants.ChunkSize / world.TileSize)}
	}
}

func WorldToTile(v pixel.Vec, left bool) world.Coords {
	x, y := world.WorldToMap(v.X+world.TileSize*0.5, -(v.Y - world.TileSize*0.5))
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
	return math.Abs(a.X-b.X) <= world.TileSize*0.5 && math.Abs(a.Y-b.Y) <= world.TileSize*0.5
}

func (c *Cave) MapFn(fn func(*Tile)) {
	for _, ch := range c.Chunks {
		for _, row := range ch.Rows {
			for _, tile := range row {
				fn(tile)
			}
		}
	}
}

func (c *Cave) PrintCaveToTerminal() {
	if c.Type != Infinite {
		fmt.Println("Printing cave ... ")
		fmt.Println()
		for y := 0; y < (c.Bottom+1)*constants.ChunkSize; y++ {
			for x := c.Left * constants.ChunkSize; x < (c.Right+1)*constants.ChunkSize; x++ {
				tile := c.GetTileInt(x, y)
				if tile != nil {
					if tile.Special {
						fmt.Print("s")
						//} else if tile.Path {
						//	fmt.Print("p")
					} else {
						switch tile.Type {
						case BlockCollapse, BlockDig:
							if tile.Bomb {
								fmt.Print("ó")
							} else {
								fmt.Print("□")
							}
						case BlockBlast:
							fmt.Print("▣")
						case Wall:
							fmt.Print("#")
						case Deco:
							fmt.Print("*")
						case Empty:
							fmt.Print(" ")
						}
					}
				}
			}
			fmt.Print("\n")
		}
	}
}
