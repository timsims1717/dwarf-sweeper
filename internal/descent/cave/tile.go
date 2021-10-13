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
	Block
	Block1
	Block2
	Wall
)

func (t TileType) String() string {
	switch t {
	case Block:
		return "Block-Diggable-Chain"
	case Block1:
		return "Block-Diggable"
	case Block2:
		return "Bombable"
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
	Type       TileType
	Special    bool
	SubCoords  world.Coords
	RCoords    world.Coords
	BGSprite   *pixel.Sprite
	BGSpriteS  string
	BGSMatrix  pixel.Matrix
	FGSprite   *pixel.Sprite
	Entity     myecs.AnEntity
	XRay       *pixel.Sprite
	Bomb       bool
	Destroyed  bool
	Fillable   bool
	Cracked    bool
	Transform  *transform.Transform
	Chunk      *Chunk
	revealT    *timing.FrameTimer
	revealing  bool
	destroyT   *timing.FrameTimer
	destroying bool
	reload     bool
	Marked     bool
	Exit       bool

	NeverChange bool
	IsChanged   bool
	DigTrigger  func(*Tile)
}

func NewTile(x, y int, ch world.Coords, bomb bool, chunk *Chunk) *Tile {
	tran := transform.NewTransform()
	tran.Pos = pixel.V(float64(x + ch.X * constants.ChunkSize) * world.TileSize, -(float64(y + ch.Y * constants.ChunkSize) * world.TileSize))
	spr := chunk.Cave.Batcher.Sprites[startSprite]
	return &Tile{
		Type:      Block,
		SubCoords: world.Coords{ X: x, Y: y },
		RCoords:   world.Coords{ X: x + ch.X * constants.ChunkSize, Y: y + ch.Y * constants.ChunkSize},
		BGSprite:  spr,
		BGSpriteS: startSprite,
		BGSMatrix: pixel.IM,
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
		if tile.destroyT.UpdateDone() {
			tile.Destroy(false)
		}
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
			tile.BGSprite.Draw(target, tile.BGSMatrix.Moved(tile.Transform.Pos))
		}
		if tile.FGSprite != nil {
			tile.FGSprite.Draw(target, tile.Transform.Mat)
		}
	}
}

func (tile *Tile) ToDestroy() {
	if tile != nil && !tile.Destroyed && !tile.destroying && tile.Breakable() {
		tile.destroyT = timing.New(revealTimer)
		tile.destroying = true
	}
}

func (tile *Tile) Destroy(playSound bool) {
	if tile != nil && !tile.Destroyed && tile.Breakable() {
		if tile.DigTrigger != nil {
			tile.DigTrigger(tile)
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
			particles.BlockParticles(tile.Transform.Pos)
			if playSound {
				sfx.SoundPlayer.PlaySound(fmt.Sprintf("rocks%d", random.Effects.Intn(5)+1), -1.0)
			}
			if !util.IsNil(tile.Entity) {
				tile.Entity.Create(tile.Transform.Pos)
			}
		}
	}
}

func (tile *Tile) ToReveal() {
	if tile != nil && !tile.revealing && tile.Solid() && tile.Breakable() {
		tile.revealT = timing.New(revealTimer)
		tile.revealing = true
	}
}

func (tile *Tile) Reveal(instant bool) {
	if tile != nil && !tile.Bomb && tile.Solid() && tile.Breakable() {
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
			particles.BlockParticles(tile.Transform.Pos)
		}
	}
}

func (tile *Tile) UpdateSprites() {
	if tile.Type != Deco {
		tile.Chunk.Cave.UpdateBatch = true
		ns := tile.SubCoords.Neighbors()
		ss := [8]bool{}
		bs := [4]bool{}
		c := 0
		for i, n := range ns {
			t := tile.Chunk.Get(n)
			if t != nil {
				if t.Bomb && t.Breakable() {
					c++
				}
				if t.Solid() {
					ss[i] = true
				}
				if i%2 == 0 && !t.Destroyed {
					bs[i/2] = true
				}
			} else {
				ss[i] = true
				if i%2 == 0 {
					bs[i/2] = true
				}
			}
		}
		var s string
		var m pixel.Matrix
		if tile.Solid() {
			buf := new(bytes.Buffer)
			for _, b := range ss {
				if b {
					buf.Write(one)
				} else {
					buf.Write(zero)
				}
			}
			s, m = tile.Chunk.Cave.SmartTileSolid(tile.Type, buf.String())
		} else if c > 0 {
			buf := new(bytes.Buffer)
			for _, b := range bs {
				if b {
					buf.Write(one)
				} else {
					buf.Write(zero)
				}
			}
			s, m = tile.Chunk.Cave.SmartTileNum(buf.String())
		}
		if tile.BGSpriteS != s {
			tile.BGSMatrix = m
			tile.BGSpriteS = s
			tile.BGSprite = tile.Chunk.Cave.Batcher.Sprites[s]
		}
		tile.FGSprite = nil
		if !tile.Solid() {
			switch c {
			case 0:
				tile.BGSprite = nil
			case 1:
				tile.FGSprite = tile.Chunk.Cave.Batcher.Sprites["one"]
			case 2:
				tile.FGSprite = tile.Chunk.Cave.Batcher.Sprites["two"]
			case 3:
				tile.FGSprite = tile.Chunk.Cave.Batcher.Sprites["three"]
			case 4:
				tile.FGSprite = tile.Chunk.Cave.Batcher.Sprites["four"]
			case 5:
				tile.FGSprite = tile.Chunk.Cave.Batcher.Sprites["five"]
			case 6:
				tile.FGSprite = tile.Chunk.Cave.Batcher.Sprites["six"]
			case 7:
				tile.FGSprite = tile.Chunk.Cave.Batcher.Sprites["seven"]
			case 8:
				tile.FGSprite = tile.Chunk.Cave.Batcher.Sprites["eight"]
			}
		} else if tile.Cracked {
			tile.FGSprite = tile.Chunk.Cave.Batcher.Sprites["crack"]
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
	return !(tile.Type == Wall || tile.Type == Empty)
}

func (tile *Tile) Solid() bool {
	return !(tile.Type == Deco || tile.Type == Empty)
}