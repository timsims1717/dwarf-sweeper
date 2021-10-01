package dungeon

import (
	"dwarf-sweeper/internal/cfg"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
)

const (
	ChunkSize = 32
	ChunkArea = ChunkSize * ChunkSize
)

type Chunk struct {
	Coords  world.Coords
	Rows    [ChunkSize][ChunkSize]*Tile
	display bool
	reload  bool
	Cave    *Cave
}

func GenerateChunk(coords world.Coords, cave *Cave) *Chunk {
	// Array of 1024 bools
	list := [ChunkArea]bool{}
	// fill first 10-20% with true
	bCount := random.CaveGen.Intn(ChunkArea/ int(100 * (cave.bombPMax - cave.bombPMin))) + ChunkArea/ int(100 * cave.bombPMin)
	for i := 0; i < bCount; i++ {
		list[i] = true
	}
	// randomize list
	for i := len(list) - 1; i > 0; i-- {
		j := random.CaveGen.Intn(i)
		list[i], list[j] = list[j], list[i]
	}
	// create chunk, distribute bombs (trues), build tiles
	chunk := &Chunk{
		Coords:  coords,
		Rows:    [ChunkSize][ChunkSize]*Tile{},
		display: true,
		reload:  true,
		Cave:    cave,
	}
	y := 0
	x := 0
	for _, b := range list {
		var tile *Tile
		if cave.finite &&
			((coords.Y == cave.bottom && y == ChunkSize - 1) ||
			(coords.X == cave.left && x == 0) ||
			(coords.X == cave.right && x == ChunkSize - 1)) {
			tile = NewTile(x, y, coords, false, chunk)
			tile.Type = Wall
			tile.neverChange = true
			tile.breakable = false
		} else if coords.Y == 0 && y == 0 {
			tile = NewTile(x, y, coords, false, chunk)
			tile.Type = Wall
			tile.neverChange = true
			tile.breakable = false
		} else {
			tile = NewTile(x, y, coords, b, chunk)
		}
		if b {
			tile.Entity = &Bomb{
				Tile: tile,
				FuseLength: tile.Chunk.Cave.fuseLen,
			}
			tile.XRay = img.Batchers[cfg.EntityKey].Sprites["bomb_fuse"]
		} else if random.CaveGen.Intn(cave.gemRate) == 0 {
			collect := Collectibles[GemDiamond]
			tile.Entity = &CollectibleItem{
				collect: collect,
			}
			tile.XRay = collect.Sprite
		} else if random.CaveGen.Intn(cave.itemRate) == 0 {
			collectible := ""
			//collectible = XRayItem
			switch random.CaveGen.Intn(5) {
			case 0:
				tile.Entity = &BombItem{}
				tile.XRay = img.Batchers[cfg.EntityKey].Sprites["bomb_unlit"]
			case 1:
				collectible = Heart
			case 2:
				collectible = Beer
			case 3:
				collectible = BubbleItem
			case 4:
				collectible = XRayItem
			}
			if collectible != "" {
				collect := Collectibles[collectible]
				tile.Entity = &CollectibleItem{
					collect: collect,
				}
				tile.XRay = collect.Sprite
			}
		} else if random.CaveGen.Intn(75) == 0 {
			tile.Entity = &MadMonk{}
		}
		chunk.Rows[y][x] = tile
		x++
		if x % ChunkSize == 0 {
			x = 0
			y++
		}
	}
	return chunk
}

func (chunk *Chunk) Update() {
	if chunk.reload {
		for _, row := range chunk.Rows {
			for _, tile := range row {
				tile.reload = true
			}
		}
		chunk.reload = false
	}
	if chunk.display {
		for _, row := range chunk.Rows {
			for _, tile := range row {
				tile.Update()
			}
		}
	}
}

func (chunk *Chunk) Draw(target pixel.Target) {
	if chunk.display {
		for _, row := range chunk.Rows {
			for _, tile := range row {
				if !tile.destroyed {
					tile.Draw(target)
				}
			}
		}
	}
}

func (chunk *Chunk) Get(coords world.Coords) *Tile {
	if chunk == nil {
		return nil
	}
	if coords.X < 0 || coords.Y < 0 || coords.X >= ChunkSize || coords.Y >= ChunkSize {
		ax := coords.X
		ay := coords.Y
		cx := 0
		cy := 0
		if coords.X < 0 {
			cx = -1
			ax = ChunkSize - 1
		} else if coords.X >= ChunkSize {
			cx = 1
			ax = 0
		}
		if coords.Y < 0 {
			cy = -1
			ay = ChunkSize - 1
		} else if coords.Y >= ChunkSize {
			cy = 1
			ay = 0
		}
		cc := chunk.Coords
		cc.X += cx
		cc.Y += cy
		ac := world.Coords{
			X: ax,
			Y: ay,
		}
		return chunk.Cave.GetChunk(cc).Get(ac)
	}
	return chunk.Rows[coords.Y][coords.X]
}