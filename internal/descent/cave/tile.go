package cave

import (
	"bytes"
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data/player"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/noise"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
	"image/color"
)

const (
	startSprite = "blank"
	revealTimer = 0.2
)

var (
	one  = []byte("1")
	zero = []byte("0")
)

type BlockType int

const (
	Unknown = iota
	Deco
	Empty
	BlockCollapse
	BlockDig
	BlockBlast
	Wall
)

func (t BlockType) String() string {
	switch t {
	case BlockCollapse:
		return "Block-Collapse"
	case BlockDig:
		return "Block-Dig"
	case BlockBlast:
		return "Block-Blast"
	case Wall:
		return "Indestructible"
	case Deco:
		return "Decoration"
	case Empty:
		return "Empty"
	default:
		return "Unknown"
	}
}

type Tile struct {
	Type      BlockType
	Special   bool
	SubCoords world.Coords
	RCoords   world.Coords
	Transform *transform.Transform
	Chunk     *Chunk
	Biome     string

	Sprites     []img.Sprite
	SmartStr    string
	FogSmartStr string
	FogSpriteS  string
	Surrounded  bool
	DSurrounded bool
	BombCount   int

	//Entity     *ecs.Entity
	XRay       string
	Bomb       bool
	Destroyed  bool
	revealT    *timing.FrameTimer
	revealing  bool
	destroying bool
	destroyer  *player.Player
	reload     bool
	Flagged    bool
	Exit       bool

	NeverChange bool
	IsChanged   bool
	Change      bool
	Path        bool
	Marked      bool
	DeadEnd     bool
	Group       int
	Perlin      float64

	DestroyTrigger func(*player.Player, *Tile)
	GemRate        float64
}

func NewTile(x, y int, coords world.Coords, bomb bool, chunk *Chunk) *Tile {
	tran := transform.New()
	tran.Scalar = pixel.V(1.001, 1.001)
	tran.Pos = pixel.V(float64(x+coords.X*constants.ChunkSize)*world.TileSize, -(float64(y+coords.Y*constants.ChunkSize) * world.TileSize))
	rCoords := world.Coords{X: x + coords.X*constants.ChunkSize, Y: y + coords.Y*constants.ChunkSize}
	p := noise.BlockType(rCoords)
	//e := myecs.Manager.NewEntity()
	//e.AddComponent(myecs.Transform, tran).
	//	AddComponent(myecs.Batch, constants.FogKey)
	return &Tile{
		Type:      BlockCollapse,
		SubCoords: world.Coords{X: x, Y: y},
		RCoords:   rCoords,
		Sprites:   []img.Sprite{},
		Bomb:      bomb,
		Transform: tran,
		Chunk:     chunk,
		Biome:     chunk.Cave.Biome,
		GemRate:   1.,
		Perlin:    p,
	}
}

func (tile *Tile) Update() {
	if tile.reload {
		if tile.SubCoords.X == 0 || tile.SubCoords.X == constants.ChunkSize-1 || tile.SubCoords.Y == 0 || tile.SubCoords.Y == constants.ChunkSize-1 {
			for _, n := range tile.SubCoords.Neighbors() {
				t := tile.Chunk.Get(n)
				if t != nil {
					if t.Destroyed {
						tile.Reveal(true)
					}
				}
			}
		}
		tile.Chunk.Cave.UpdateBatch = true
		tile.reload = false
	}
	if !tile.Destroyed && tile.destroying {
		tile.destroy()
	}
	if tile.Solid() && !tile.Destroyed && tile.revealing && tile.Breakable() {
		if tile.revealT.UpdateDone() {
			tile.Reveal(false)
		}
	}
	tile.Transform.Update()
	if debug.Debug && tile.Group != 0 {
		tile.Transform.Mask = color.RGBA{
			R: uint8((((tile.Group * 7) % 8) * 32) % 256),
			G: uint8((((tile.Group * 13) % 8) * 32) % 256),
			B: uint8((((tile.Group * 11) % 8) * 32) % 256),
			A: 255,
		}
		// 1: Yellow, 2: Pink, 3: Lime, 4: Gray, 5: Purple, 6: Green, 7: Blue
	} else {
		tile.Transform.Mask = colornames.White
	}
}

func (tile *Tile) Draw() {
	for _, spr := range tile.Sprites {
		if spr.S != nil {
			spr.S.DrawColorMask(spr.B.Batch(), spr.M, tile.Transform.Mask)
		}
		//spr.B.DrawSpriteColor(spr.K, spr.M, tile.Transform.Mask)
	}
}

func (tile *Tile) Destroy(p *player.Player, playSound bool) {
	if tile != nil && !tile.Destroyed && !tile.destroying && tile.Breakable() {
		tile.destroying = true
		tile.destroyer = p
		if playSound {
			sfx.SoundPlayer.PlaySound(fmt.Sprintf("rocks%d", random.Effects.Intn(5)+1), -1.0)
		}
	}
}

func (tile *Tile) DestroySpecial(playSound, allowWalls, ignoreTrigger bool) {
	if tile != nil && !tile.Destroyed && !tile.destroying && (tile.Breakable() || allowWalls) {
		if ignoreTrigger {
			tile.DestroyTrigger = nil
			tile.destroyer = nil
		}
		tile.destroying = true
		if playSound {
			sfx.SoundPlayer.PlaySound(fmt.Sprintf("rocks%d", random.Effects.Intn(5)+1), -1.0)
		}
	}
}

func (tile *Tile) destroy() {
	if tile != nil && !tile.Destroyed {
		if tile.DestroyTrigger != nil {
			tile.DestroyTrigger(tile.destroyer, tile)
		}
		wasSolid := tile.Solid()
		tile.Chunk.Cave.UpdateBatch = true
		tile.destroying = false
		tile.Type = Empty
		ns := tile.SubCoords.Neighbors()
		c := 0
		for _, n := range ns {
			t := tile.Chunk.Get(n)
			if t != nil {
				if t.Bomb && t.Breakable() {
					c++
				}
			}
		}
		if tile.Bomb {
			tile.Bomb = false
			tile.Destroyed = true
		} else if c == 0 {
			tile.Destroyed = true
			for _, n := range ns {
				tile.Chunk.Get(n).ToReveal()
			}
		}
		if wasSolid {
			particles.BlockParticles(tile.Transform.Pos, tile.Chunk.Cave.Biome)
		}
	}
}

func (tile *Tile) ToReveal() {
	if tile != nil && !tile.revealing && tile.Solid() && tile.Breakable() && tile.Type == BlockCollapse {
		tile.revealT = timing.New(revealTimer)
		tile.revealing = true
	}
}

func (tile *Tile) Reveal(instant bool) {
	if tile != nil && !tile.Bomb && tile.Solid() && tile.Breakable() && tile.Type == BlockCollapse {
		tile.Chunk.Cave.UpdateBatch = true
		tile.revealing = false
		tile.Type = Empty
		ns := tile.SubCoords.Neighbors()
		c := 0
		for _, n := range ns {
			t := tile.Chunk.Get(n)
			if t != nil {
				if t.Bomb && t.Breakable() {
					c++
				}
			}
		}
		if c == 0 {
			tile.Destroyed = true
			for _, n := range ns {
				if instant {
					tile.Chunk.Get(n).Reveal(true)
				} else {
					tile.Chunk.Get(n).ToReveal()
				}
			}
		}
		if !instant {
			particles.BlockParticles(tile.Transform.Pos, tile.Chunk.Cave.Biome)
		}
	}
}

func (tile *Tile) UpdateDetails() {
	if tile.Type != Deco {
		ns := tile.SubCoords.Neighbors()
		tile.Surrounded = true
		tile.DSurrounded = true
		ss := [8]bool{} // surrounded string code (FG)
		ts := [8]bool{} // tile string code (BG)
		bs := [4]bool{} // empty num string code (BG)
		c := 0
		for i, n := range ns {
			t := tile.Chunk.Get(n)
			if t != nil {
				if t.Solid() {
					if t.Bomb && t.Breakable() {
						c++
					}
					if t.Type == tile.Type {
						ts[i] = true
					}
					if t.Surrounded {
						ss[i] = true
					} else if !t.Surrounded {
						tile.DSurrounded = false
					}
					if i%2 == 0 {
						bs[i/2] = true
					}
				} else {
					tile.Surrounded = false
					tile.DSurrounded = false
					if i%2 == 0 && t.BombCount > 0 {
						bs[i/2] = true
					}
				}
			} else {
				ts[i] = true
				if i%2 == 0 {
					bs[i/2] = true
				}
			}
		}
		tile.BombCount = c
		var smartStr string
		if tile.Solid() {
			// background
			buf := new(bytes.Buffer)
			for _, b := range ts {
				if b {
					buf.Write(one)
				} else {
					buf.Write(zero)
				}
			}
			smartStr = buf.String()
			// foreground
			if tile.Surrounded && !tile.DSurrounded && tile.Chunk.Cave.Fog {
				buf2 := new(bytes.Buffer)
				for _, b := range ss {
					if b {
						buf2.Write(one)
					} else {
						buf2.Write(zero)
					}
				}
				tile.FogSmartStr = buf2.String()
			} else {
				tile.FogSmartStr = ""
			}
		} else if c > 0 {
			// background
			buf := new(bytes.Buffer)
			for _, b := range bs {
				if b {
					buf.Write(one)
				} else {
					buf.Write(zero)
				}
			}
			smartStr = buf.String()
			tile.FogSmartStr = ""
		} else {
			tile.FogSmartStr = ""
			tile.Type = Empty
		}
		tile.SmartStr = smartStr
	} else {
		tile.FogSmartStr = ""
	}
}

func (tile *Tile) UpdateSprites() {
	if tile.Type != Deco {
		tile.ClearSprites()
		if tile.Solid() {
			tile.FogSpriteS = ""
			// main tile
			spr, mat := SmartTileSolid(tile.Type, tile.SmartStr, tile.DSurrounded && tile.Chunk.Cave.Fog, tile.Perlin)
			tile.AddSprite(spr, mat, tile.Biome, false)
			// fog
			var fogSpr string
			var fogMat pixel.Matrix
			if tile.Surrounded && !tile.DSurrounded && tile.Chunk.Cave.Fog {
				fogSpr, fogMat = SmartTileFade(tile.FogSmartStr)
				tile.FogSpriteS = fogSpr
				tile.AddSprite(fogSpr, fogMat, constants.FogKey, false)
			} else if tile.Surrounded && tile.DSurrounded && tile.Chunk.Cave.Fog {
				tile.AddSprite("empty", pixel.IM, constants.FogKey, false)
			}
		} else if tile.BombCount > 0 {
			// main tile
			spr, mat := SmartTileNum(tile.SmartStr, tile.Perlin)
			// number sprite
			var numSpr string
			switch tile.BombCount {
			case 1:
				numSpr = "one"
			case 2:
				numSpr = "two"
			case 3:
				numSpr = "three"
			case 4:
				numSpr = "four"
			case 5:
				numSpr = "five"
			case 6:
				numSpr = "six"
			case 7:
				numSpr = "seven"
			case 8:
				numSpr = "eight"
			}
			tile.AddSprite(spr, mat, tile.Biome, true)
			tile.AddSprite(numSpr, pixel.IM, tile.Biome, true)
		}
	}
	//tile.Entity.AddComponent(myecs.Drawable, tile.Sprites)
}

func (tile *Tile) GetTileCode() string {
	ns := tile.SubCoords.Neighbors()
	bs := [8]bool{}
	c := 0
	for i, n := range ns {
		t := tile.Chunk.Get(n)
		if t != nil {
			if t.Bomb {
				c++
			}
			if t.Solid() {
				bs[i] = true
			}
		}
	}
	buf := new(bytes.Buffer)
	for _, b := range bs {
		if b {
			buf.Write(one)
		} else {
			buf.Write(zero)
		}
	}
	return buf.String()
}

func (tile *Tile) ClearSprites() {
	tile.Sprites = []img.Sprite{}
	tile.Chunk.Cave.UpdateBatch = true
}

func (tile *Tile) AddSprite(key string, mat pixel.Matrix, biome string, bg bool) {
	var batch *img.Batcher
	if bg {
		batch = img.Batchers[fmt.Sprintf(constants.CaveBGFMT, biome)]
	} else {
		batch = img.Batchers[biome]
	}
	tile.Sprites = append(tile.Sprites, img.Sprite{
		K: key,
		S: batch.GetSprite(key),
		M: mat.ScaledXY(pixel.ZV, tile.Transform.Scalar).Moved(tile.Transform.Pos),
		B: batch,
	})
	tile.Chunk.Cave.UpdateBatch = true
}

func (tile *Tile) IsExit() bool {
	return tile != nil && tile.Exit
}

func (tile *Tile) Breakable() bool {
	return !(tile.Type == Wall || tile.Type == Empty || tile.Type == Deco)
}

func (tile *Tile) Solid() bool {
	return !(tile.Type == Deco || tile.Type == Empty)
}

func (tile *Tile) Diggable() bool {
	return tile.Type == BlockDig || tile.Type == BlockCollapse
}

func (tile *Tile) Neighbors() []*Tile {
	var neighbors []*Tile
	for _, n := range tile.RCoords.Neighbors() {
		t := tile.Chunk.Cave.GetTileInt(n.X, n.Y)
		if t != nil {
			neighbors = append(neighbors, t)
		}
	}
	return neighbors
}