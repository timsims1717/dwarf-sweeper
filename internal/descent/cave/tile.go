package cave

import (
	"bytes"
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
)

const (
	startSprite = "blank"
	revealTimer = 0.2
)

var (
	one = []byte("1")
	zero = []byte("0")
)

type TileType int

const (
	Deco = iota
	Empty
	BlockCollapse
	BlockDig
	BlockBlast
	Wall
)

func (t TileType) String() string {
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

type TileSpr struct {
	SprKey string
	Matrix pixel.Matrix
	BG     bool
}

type Tile struct {
	Type      TileType
	Special   bool
	SubCoords world.Coords
	RCoords   world.Coords
	Transform *transform.Transform
	Chunk     *Chunk

	Sprites      []TileSpr
	SmartStr     string
	SmartChange  bool
	FogSmartStr  string
	FogSprite    string
	FogSpriteS   string
	FogMatrix    pixel.Matrix
	Surrounded   bool
	DSurrounded  bool
	BombCount    int

	Entity     myecs.AnEntity
	XRay       *pixel.Sprite
	Bomb       bool
	Destroyed  bool
	revealT    *timing.FrameTimer
	revealing  bool
	destroying bool
	reload     bool
	Flagged    bool
	Exit       bool

	NeverChange bool
	IsChanged   bool
	Path        bool
	Marked      bool
	DeadEnd     bool
	Room        bool

	DestroyTrigger func(*Tile)
}

func NewTile(x, y int, ch world.Coords, bomb bool, chunk *Chunk) *Tile {
	tran := transform.NewTransform()
	tran.Pos = pixel.V(float64(x + ch.X * constants.ChunkSize) * world.TileSize, -(float64(y + ch.Y * constants.ChunkSize) * world.TileSize))
	return &Tile{
		Type:        BlockCollapse,
		SubCoords:   world.Coords{ X: x, Y: y },
		RCoords:     world.Coords{ X: x + ch.X * constants.ChunkSize, Y: y + ch.Y * constants.ChunkSize},
		Sprites:     []TileSpr{
			{
				SprKey: startSprite,
			},
		},
		Bomb:        bomb,
		Transform:   tran,
		Chunk:       chunk,
	}
}

func (tile *Tile) Update() {
	if tile.reload {
		if tile.SubCoords.X == 0 || tile.SubCoords.X == constants.ChunkSize- 1 || tile.SubCoords.Y == 0 || tile.SubCoords.Y == constants.ChunkSize- 1 {
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
}

func (tile *Tile) Draw() {
	//if !tile.Destroyed {
	for _, spr := range tile.Sprites {
		if spr.SprKey != "" {
			mat := spr.Matrix.ScaledXY(pixel.ZV, pixel.V(1.001, 1.001)).Moved(tile.Transform.Pos)
			if spr.BG {
				img.Batchers[constants.CaveBGKey].DrawSprite(spr.SprKey, mat)
			} else {
				img.Batchers[constants.CaveKey].DrawSprite(spr.SprKey, mat)
			}
		}
	}
	if tile.FogSprite != "" {
		img.Batchers[constants.FogKey].DrawSprite(tile.FogSprite, tile.FogMatrix.ScaledXY(pixel.ZV, pixel.V(1.001, 1.001)).Moved(tile.Transform.Pos))
	}
	//}
}

func (tile *Tile) Destroy(playSound bool) {
	if tile != nil && !tile.Destroyed && !tile.destroying && tile.Breakable() {
		tile.destroying = true
		if playSound {
			sfx.SoundPlayer.PlaySound(fmt.Sprintf("rocks%d", random.Effects.Intn(5)+1), -1.0)
		}
	}
}

func (tile *Tile) DestroySpecial(playSound, allowWalls, ignoreTrigger bool) {
	if tile != nil && !tile.Destroyed && !tile.destroying && (tile.Breakable() || allowWalls) {
		if ignoreTrigger {
			tile.DestroyTrigger = nil
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
			tile.DestroyTrigger(tile)
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
		} else {
			if c == 0 {
				tile.Destroyed = true
				for _, n := range ns {
					tile.Chunk.Get(n).ToReveal()
				}
			}
		}
		if wasSolid {
			particles.BlockParticles(tile.Transform.Pos, tile.Chunk.Cave.Biome)
			if !util.IsNil(tile.Entity) {
				tile.Entity.Create(tile.Transform.Pos)
			}
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
		if !util.IsNil(tile.Entity) {
			tile.Entity.Create(tile.Transform.Pos)
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
		tile.SmartChange = tile.SmartStr != smartStr
		tile.SmartStr = smartStr
	} else {
		tile.FogSmartStr = ""
		tile.FogSpriteS = ""
		tile.FogSprite = ""
	}
}

func (tile *Tile) UpdateSprites() {
	if tile.Type != Deco {
		if tile.Solid() {
			// main tile
			spr, mat := SmartTileSolid(tile.Type, tile.SmartStr, tile.DSurrounded && tile.Chunk.Cave.Fog)
			if len(tile.Sprites) == 0 || tile.Sprites[0].SprKey != spr || tile.SmartChange {
				tile.Sprites = []TileSpr{}
				tile.AddSprite(spr, mat, false)
			}
			// fog
			var fogSpr string
			var fogMat pixel.Matrix
			if tile.Surrounded && !tile.DSurrounded && tile.Chunk.Cave.Fog {
				fogSpr, fogMat = SmartTileFade(tile.FogSmartStr)
			} else if tile.Surrounded && tile.DSurrounded && tile.Chunk.Cave.Fog {
				fogSpr, fogMat = "empty", pixel.IM
			} else {
				fogSpr = ""
				fogMat = pixel.IM
			}
			if tile.FogSpriteS != fogSpr || tile.FogMatrix != fogMat {
				tile.FogMatrix = fogMat
				tile.FogSpriteS = fogSpr
				if fogSpr != "" {
					tile.FogSprite = fogSpr
				} else {
					tile.FogSprite = ""
				}
			}
		} else if tile.BombCount > 0 {
			// main tile
			spr, mat := SmartTileNum(tile.SmartStr)
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
			if len(tile.Sprites) == 0 || tile.Sprites[0].SprKey != spr || tile.SmartChange {
				tile.Sprites = []TileSpr{}
				tile.AddSprite(spr, mat, true)
				tile.AddSprite(numSpr, pixel.IM, true)
			} else {
				if len(tile.Sprites) > 1 {
					tile.Sprites[1].SprKey = numSpr
				} else {
					tile.AddSprite(numSpr, pixel.IM, true)
				}
			}
			tile.FogSprite = ""
		} else {
			tile.FogSprite = ""
			tile.Sprites = []TileSpr{}
		}
	}
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
	tile.Sprites = []TileSpr{}
}

func (tile *Tile) AddSprite(key string, mat pixel.Matrix, isBG bool) {
	tile.Sprites = append(tile.Sprites, TileSpr{
		SprKey: key,
		Matrix: mat,
		BG:     isBG,
	})
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