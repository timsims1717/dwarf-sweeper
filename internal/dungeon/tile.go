package dungeon

import (
	"bytes"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"math/rand"
	"time"
)

const (
	startSprite = "blank"
)

var (
	one = []byte("1")
	zero = []byte("0")
)

type TileType int

const (
	Block = iota
	Value
	Wall
	Deco
	Empty
)

func (t TileType) String() string {
	switch t {
	case Block:
		return "Block"
	case Value:
		return "Valuable"
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
	SubCoords world.Coords
	RCoords   world.Coords
	BGSprite  *pixel.Sprite
	BGSpriteS string
	BGSMatrix pixel.Matrix
	FGSprite  *pixel.Sprite
	Entities  []Entity
	bomb      bool
	destroyed bool
	breakable bool
	cracked   bool
	Solid     bool
	Transform  *transform.Transform
	Chunk      *Chunk
	revealT    time.Time
	revealing  bool
	destroyT   time.Time
	destroying bool
	reload     bool
	marked     bool

	neverChange bool
	isChanged   bool
}

func NewTile(x, y int, ch world.Coords, bomb bool, chunk *Chunk) *Tile {
	tran := transform.NewTransform()
	tran.Pos = pixel.V(float64(x + ch.X * ChunkSize) * world.TileSize, -(float64(y + ch.Y * ChunkSize) * world.TileSize))
	spr := chunk.Cave.batcher.Sprites[startSprite]
	return &Tile{
		Type:      Block,
		SubCoords: world.Coords{ X: x, Y: y },
		RCoords:   world.Coords{ X: x + ch.X * ChunkSize, Y: y + ch.Y * ChunkSize },
		BGSprite:  spr,
		BGSpriteS: startSprite,
		BGSMatrix: pixel.IM,
		bomb:      bomb,
		breakable: true,
		Solid:     true,
		Transform: tran,
		Chunk:     chunk,
	}
}

func (tile *Tile) AddEntity(e Entity) {
	tile.Entities = append(tile.Entities, e)
}

func (tile *Tile) Update() {
	if tile.reload {
		if tile.SubCoords.X == 0 || tile.SubCoords.X == ChunkSize - 1 || tile.SubCoords.Y == 0 || tile.SubCoords.Y == ChunkSize - 1 {
			for _, n := range tile.SubCoords.Neighbors() {
				t := tile.Chunk.Get(n)
				if t != nil {
					if t.destroyed {
						tile.Reveal(true)
					}
					t.UpdateSprites()
				}
			}
		}
		tile.UpdateSprites()
		tile.reload = false
	}
	if !tile.destroyed && tile.destroying && tile.breakable {
		s := time.Since(tile.destroyT).Seconds()
		if s >= 0.2 {
			tile.Destroy()
		}
	}
	if tile.Solid && !tile.destroyed && tile.revealing && tile.breakable {
		s := time.Since(tile.revealT).Seconds()
		if s >= 0.2 {
			tile.Reveal(false)
		}
	}
	tile.Transform.Update()
}

func (tile *Tile) Draw(target pixel.Target) {
	if !tile.destroyed {
		if tile.BGSprite != nil {
			tile.BGSprite.Draw(target, tile.BGSMatrix.Moved(tile.Transform.Pos))
		}
		if tile.FGSprite != nil {
			tile.FGSprite.Draw(target, tile.Transform.Mat)
		}
	}
}

func (tile *Tile) ToDestroy() {
	if tile != nil && !tile.destroyed && !tile.destroying && tile.breakable {
		tile.destroyT = time.Now()
		tile.destroying = true
	}
}

func (tile *Tile) Destroy() {
	if tile != nil && !tile.destroyed && tile.breakable {
		tile.destroying = false
		tile.Solid = false
		tile.Type = Empty
		ns := tile.SubCoords.Neighbors()
		c := 0
		for _, n := range ns {
			t := tile.Chunk.Get(n)
			if t != nil {
				if t.bomb && t.breakable && t.Type != Wall {
					c++
				}
				t.UpdateSprites()
			}
		}
		if tile.bomb {
			tile.bomb = false
			tile.destroyed = true
			tile.BGSprite = nil
			tile.FGSprite = nil
		} else {
			if c == 0 {
				tile.destroyed = true
				for _, n := range ns {
					tile.Chunk.Get(n).ToReveal()
				}
			}
			tile.UpdateSprites()
			particles.BlockParticles(tile.Transform.Pos)
			sfx.SoundPlayer.PlaySound(fmt.Sprintf("rocks%d", rand.Intn(5) + 1), -1.0)
		}
		for _, e := range tile.Entities {
			Entities.Add(e, tile.Transform.Pos)
		}
	}
}

func (tile *Tile) ToReveal() {
	if tile != nil && !tile.revealing && tile.Solid && tile.breakable {
		tile.revealT = time.Now()
		tile.revealing = true
	}
}

func (tile *Tile) Reveal(instant bool) {
	if tile != nil && !tile.bomb && tile.Solid && tile.breakable {
		tile.revealing = false
		tile.Solid = false
		tile.Type = Empty
		ns := tile.SubCoords.Neighbors()
		c := 0
		for _, n := range ns {
			t := tile.Chunk.Get(n)
			if t != nil {
				if t.bomb && t.breakable && t.Type != Wall {
					c++
				}
				t.UpdateSprites()
			}
		}
		for _, e := range tile.Entities {
			Entities.Add(e, tile.Transform.Pos)
		}
		if c == 0 {
			tile.destroyed = true
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

func (tile *Tile) Unmark() {
	tile.marked = false
}

func (tile *Tile) Mark(from pixel.Vec) {
	if tile != nil && tile.Solid && !tile.destroyed && tile.Type != Wall {
		correct := tile.bomb
		if !tile.marked {
			tile.marked = true
			Entities.Add(&Flag{
				Tile: tile,
			}, from)
			if correct {
				BombsMarked++
			} else {
				WrongMarks++
			}
		} else {
			tile.marked = false
			if correct {
				BombsMarked--
			} else {
				WrongMarks--
			}
		}
	}
}

func (tile *Tile) UpdateSprites() {
	if tile.Type != Deco {
		ns := tile.SubCoords.Neighbors()
		ss := [8]bool{}
		bs := [4]bool{}
		c := 0
		for i, n := range ns {
			t := tile.Chunk.Get(n)
			if t != nil {
				if t.bomb && t.breakable && t.Type != Wall {
					c++
				}
				if t.Solid {
					ss[i] = true
				}
				if i%2 == 0 && !t.destroyed {
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
		if tile.Solid {
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
			tile.BGSprite = tile.Chunk.Cave.batcher.Sprites[s]
		}
		tile.FGSprite = nil
		if !tile.Solid {
			switch c {
			case 0:
				tile.BGSprite = nil
			case 1:
				tile.FGSprite = tile.Chunk.Cave.batcher.Sprites["one"]
			case 2:
				tile.FGSprite = tile.Chunk.Cave.batcher.Sprites["two"]
			case 3:
				tile.FGSprite = tile.Chunk.Cave.batcher.Sprites["three"]
			case 4:
				tile.FGSprite = tile.Chunk.Cave.batcher.Sprites["four"]
			case 5:
				tile.FGSprite = tile.Chunk.Cave.batcher.Sprites["five"]
			case 6:
				tile.FGSprite = tile.Chunk.Cave.batcher.Sprites["six"]
			case 7:
				tile.FGSprite = tile.Chunk.Cave.batcher.Sprites["seven"]
			case 8:
				tile.FGSprite = tile.Chunk.Cave.batcher.Sprites["eight"]
			}
		} else if tile.cracked {
			tile.FGSprite = tile.Chunk.Cave.batcher.Sprites["crack"]
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
			if t.bomb {
				c++
			}
			if t.Solid {
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