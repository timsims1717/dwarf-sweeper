package cave

import (
	"bytes"
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/random"
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

type Tile struct {
	Type      TileType
	Special   bool
	SubCoords world.Coords
	RCoords   world.Coords
	Transform *transform.Transform
	Chunk     *Chunk

	SmartStr    string
	BGSprite    *pixel.Sprite
	BGSpriteS   string
	BGMatrix    pixel.Matrix
	FGSprite    *pixel.Sprite
	FGSpriteS   string
	FGMatrix    pixel.Matrix
	Surrounded  bool
	DSurrounded bool

	Entity      myecs.AnEntity
	XRay        *pixel.Sprite
	Bomb        bool
	Destroyed   bool
	Fillable    bool
	revealT     *timing.FrameTimer
	revealing   bool
	destroying  bool
	reload      bool
	Marked      bool
	Exit        bool

	NeverChange bool
	IsChanged   bool
	PartOfPath  bool

	DestroyTrigger func(*Tile)
}

func NewTile(x, y int, ch world.Coords, bomb bool, chunk *Chunk) *Tile {
	tran := transform.NewTransform()
	tran.Pos = pixel.V(float64(x + ch.X * constants.ChunkSize) * world.TileSize, -(float64(y + ch.Y * constants.ChunkSize) * world.TileSize))
	spr := chunk.Cave.Batcher.Sprites[startSprite]
	return &Tile{
		Type:      BlockCollapse,
		SubCoords: world.Coords{ X: x, Y: y },
		RCoords:   world.Coords{ X: x + ch.X * constants.ChunkSize, Y: y + ch.Y * constants.ChunkSize},
		BGSprite:  spr,
		BGSpriteS: startSprite,
		BGMatrix:  pixel.IM,
		Bomb:      bomb,
		Transform: tran,
		Chunk:     chunk,
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
					t.UpdateSprites()
				}
			}
		}
		tile.UpdateSprites()
		tile.reload = false
	}
	if !tile.Destroyed && tile.destroying && tile.Breakable() {
		tile.destroy()
	}
	if tile.Solid() && !tile.Destroyed && tile.revealing && tile.Breakable() {
		if tile.revealT.UpdateDone() {
			tile.Reveal(false)
		}
	}
	tile.Transform.Update()
}

func (tile *Tile) Draw(target pixel.Target) {
	if !tile.Destroyed {
		if tile.BGSprite != nil {
			tile.BGSprite.Draw(target, tile.BGMatrix.ScaledXY(pixel.ZV, pixel.V(1.0001, 1.0001)).Moved(tile.Transform.Pos))
		}
		if tile.FGSprite != nil {
			tile.FGSprite.Draw(target, tile.FGMatrix.ScaledXY(pixel.ZV, pixel.V(1.0001, 1.0001)).Moved(tile.Transform.Pos))
		}
	}
}

func (tile *Tile) Destroy(playSound bool) {
	if tile != nil && !tile.Destroyed && !tile.destroying && tile.Breakable() {
		tile.destroying = true
		if playSound {
			sfx.SoundPlayer.PlaySound(fmt.Sprintf("rocks%d", random.Effects.Intn(5)+1), -1.0)
		}
	}
}

func (tile *Tile) destroy() {
	if tile != nil && !tile.Destroyed && tile.Breakable() {
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
				t.UpdateSprites()
			}
		}
		if tile.Bomb {
			tile.Bomb = false
			tile.Destroyed = true
			tile.BGSprite = nil
			tile.FGSprite = nil
		} else {
			if c == 0 {
				tile.Destroyed = true
				for _, n := range ns {
					tile.Chunk.Get(n).ToReveal()
				}
			}
			tile.UpdateSprites()
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
				t.UpdateSprites()
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
		tile.UpdateSprites()
		if !instant {
			particles.BlockParticles(tile.Transform.Pos, tile.Chunk.Cave.Biome)
		}
	}
}

func (tile *Tile) UpdateSprites() {
	if tile.Type != Deco {
		tile.Chunk.Cave.UpdateBatch = true
		ns := tile.SubCoords.Neighbors()
		tile.Surrounded = true
		tile.DSurrounded = true
		ss := [4]bool{}
		ts := [8]bool{}
		bs := [4]bool{}
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
					if t.Surrounded && i%2 == 1 {
						ss[i/2] = true
					} else if !t.Surrounded {
						tile.DSurrounded = false
					}
				} else {
					tile.Surrounded = false
					tile.DSurrounded = false
					if i%2 == 0 && t.BGSpriteS != "" {
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
		var bgs, fgs string
		var bgm, fgm pixel.Matrix
		var smartStr string
		hasSprite := true
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
			bgs, bgm = SmartTileSolid(tile.Type, buf.String(), tile.DSurrounded && tile.Chunk.Cave.Fog)
			smartStr = buf.String()
			// foreground
			if tile.Surrounded && !tile.DSurrounded {
				//buf2 := new(bytes.Buffer)
				//for _, b := range ss {
				//	if b {
				//		buf2.Write(one)
				//	} else {
				//		buf2.Write(zero)
				//	}
				//}
				//fgs, fgm = SmartTileFade(buf2.String())
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
			bgs, bgm = SmartTileNum(buf.String())
			smartStr = buf.String()
			// foreground
			switch c {
			case 1:
				fgs = "one"
			case 2:
				fgs = "two"
			case 3:
				fgs = "three"
			case 4:
				fgs = "four"
			case 5:
				fgs = "five"
			case 6:
				fgs = "six"
			case 7:
				fgs = "seven"
			case 8:
				fgs = "eight"
			}
			fgm = pixel.IM
		} else {
			tile.FGSprite = nil
			tile.BGSprite = nil
			hasSprite = false
		}
		if hasSprite && (tile.SmartStr != smartStr || tile.BGSpriteS != bgs) {
			tile.BGMatrix = bgm
			tile.BGSpriteS = bgs
			tile.BGSprite = tile.Chunk.Cave.Batcher.Sprites[bgs]
			tile.SmartStr = smartStr
		}
		if hasSprite && (tile.FGSpriteS != fgs) {
			tile.FGMatrix = fgm
			tile.FGSpriteS = fgs
			tile.FGSprite = tile.Chunk.Cave.Batcher.Sprites[fgs]
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