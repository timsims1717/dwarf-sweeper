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

type Tile struct {
	Type      BlockType
	Special   bool
	SubCoords world.Coords
	RCoords   world.Coords
	Transform *transform.Transform
	Chunk     *Chunk
	Biome     string

	Sprite      *img.Sprite
	Sprites     []img.Sprite
	SmartStr    string
	BGSmartStr  string
	FogSmartStr string
	FogSpriteS  string
	Surrounded  bool
	DSurrounded bool
	NeedBG      bool
	BombCount   int

	XRay       string
	Bomb       bool
	Destroyed  bool
	revealT    *timing.Timer
	revealing  bool
	destroying bool
	destroyer  *player.Player
	reload     bool
	Flagged    bool
	Exit       bool
	ExitI      int
	DoorI      int

	NeverChange bool
	IsChanged   bool
	Change      bool
	Path        bool
	Marked      bool
	DeadEnd     bool
	Group       int
	BG          bool
	Perlin      float64
	XPerlin     float64
	YPerlin     float64

	DestroyTrigger func(*player.Player, *Tile)
	GemRate        float64
}

func NewTile(x, y int, coords world.Coords, bomb bool, chunk *Chunk) *Tile {
	tran := transform.New()
	tran.Scalar = pixel.V(1.001, 1.001)
	tran.Pos = pixel.V(float64(x+coords.X*constants.ChunkSize)*world.TileSize, -(float64(y+coords.Y*constants.ChunkSize) * world.TileSize))
	rCoords := world.Coords{X: x + coords.X*constants.ChunkSize, Y: y + coords.Y*constants.ChunkSize}
	return &Tile{
		Type:      Collapse,
		SubCoords: world.Coords{X: x, Y: y},
		RCoords:   rCoords,
		Sprites:   []img.Sprite{},
		Bomb:      bomb,
		Transform: tran,
		Chunk:     chunk,
		Biome:     chunk.Cave.Biome,
		GemRate:   1.,
		Perlin:    noise.Perlin2D(rCoords),
		XPerlin:   noise.Perlin1D(rCoords.X),
		YPerlin:   noise.Perlin1D(rCoords.Y),
	}
}

func (tile *Tile) Update() {
	if tile.reload {
		if tile.SubCoords.X == 0 || tile.SubCoords.X == constants.ChunkSize-1 || tile.SubCoords.Y == 0 || tile.SubCoords.Y == constants.ChunkSize-1 {
			for _, n := range tile.RCoords.Neighbors() {
				t := tile.Chunk.Cave.GetTileInt(n.X, n.Y)
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
	if tile != nil && !tile.Destroyed && !tile.destroying && (tile.Breakable() || (p == nil && tile.Type == SecretDoor)) {
		tile.destroying = true
		tile.destroyer = p
		if playSound {
			sfx.SoundPlayer.PlaySound(fmt.Sprintf("rocks%d", random.Effects.Intn(5)+1), -1.0)
		}
	}
}

func (tile *Tile) DestroySpecial(playSound, allowWalls, ignoreTrigger bool) {
	if tile != nil && !tile.Destroyed && !tile.destroying && (tile.Breakable() || (tile.Type == Wall && allowWalls)) {
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
		wasType := tile.Type
		tile.Chunk.Cave.UpdateBatch = true
		tile.destroying = false
		tile.Type = Empty
		ns := tile.RCoords.Neighbors()
		c := 0
		for _, n := range ns {
			t := tile.Chunk.Cave.GetTileInt(n.X, n.Y)
			if t != nil {
				if t.Bomb && t.Breakable() {
					c++
				}
			}
		}
		tile.Destroyed = true
		if tile.Bomb {
			tile.Bomb = false
		} else if c == 0 {
			if wasType == Collapse {
				for _, n := range ns {
					tile.Chunk.Cave.GetTileInt(n.X, n.Y).ToReveal()
				}
			}
		}
		if wasSolid {
			particles.BlockParticles(tile.Transform.Pos, tile.Biome)
		}
		if wasType == SecretDoor {
			particles.BlockParticles(tile.Transform.Pos, tile.Biome)
			tile.Type = SecretOpen
			for _, n := range ns {
				nt := tile.Chunk.Cave.GetTileInt(n.X, n.Y)
				if nt.Type == SecretDoor && nt.DoorI == tile.DoorI {
					nt.destroy()
				}
			}
		}
	}
}

func (tile *Tile) ToReveal() {
	if tile != nil && !tile.revealing && tile.Solid() && tile.Breakable() && tile.Type == Collapse {
		tile.revealT = timing.New(revealTimer)
		tile.revealing = true
	}
}

func (tile *Tile) Reveal(instant bool) {
	if tile != nil && !tile.Bomb && tile.Solid() && tile.Breakable() && tile.Type == Collapse {
		tile.Chunk.Cave.UpdateBatch = true
		tile.revealing = false
		wasType := tile.Type
		tile.Type = Empty
		ns := tile.RCoords.Neighbors()
		c := 0
		for _, n := range ns {
			t := tile.Chunk.Cave.GetTileInt(n.X, n.Y)
			if t != nil {
				if t.Bomb && t.Breakable() {
					c++
				}
			}
		}
		tile.Destroyed = true
		if c == 0 {
			for _, n := range ns {
				if wasType == Collapse {
					if instant {
						tile.Chunk.Cave.GetTileInt(n.X, n.Y).Reveal(true)
					} else {
						tile.Chunk.Cave.GetTileInt(n.X, n.Y).ToReveal()
					}
				}
			}
		}
		if !instant {
			particles.BlockParticles(tile.Transform.Pos, tile.Biome)
		}
	}
}

func (tile *Tile) UpdateDetails() {
	tile.FogSmartStr = ""
	tile.Surrounded = true
	tile.DSurrounded = true
	ss := [8]bool{} // surrounded string code (for fog)
	ts := [8]bool{} // tile string code (for blocks)
	bs := [4]bool{} // background string code (for background)
	c := 0
	for i, n := range tile.RCoords.Neighbors() {
		t := tile.Chunk.Cave.GetTileInt(n.X, n.Y)
		if t != nil {
			if t.Solid() {
				if t.Bomb && t.Breakable() {
					c++
				}
				if t.Type == tile.Type && t.Biome == tile.Biome {
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
	// background
	buf := new(bytes.Buffer)
	for _, b := range bs {
		if b {
			buf.Write(one)
		} else {
			buf.Write(zero)
		}
	}
	tile.BGSmartStr = buf.String()
	if tile.Solid() {
		// tile
		buf = new(bytes.Buffer)
		for _, b := range ts {
			if b {
				buf.Write(one)
			} else {
				buf.Write(zero)
			}
		}
		tile.SmartStr = buf.String()
		// fog
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
		}
	}
}

func (tile *Tile) UpdateSprites() {
	tile.ClearSprites()
	if tile.Type != Blank {
		if tile.Solid() {
			// background
			sprBG, matBG := SmartTileNum(tile.BGSmartStr, tile.Perlin)
			if sprBG != "" {
				tile.AddSprite(sprBG, matBG, tile.Biome, true)
			}
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
		} else {
			if tile.BombCount > 0 {
				// background
				sprBG, matBG := SmartTileNum(tile.BGSmartStr, tile.Perlin)
				if sprBG != "" {
					tile.AddSprite(sprBG, matBG, tile.Biome, true)
				}
				// number
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
				tile.AddSprite(numSpr, pixel.IM, tile.Biome, true)
			}
			if tile.Type == Pillar || tile.Type == Growth {
				st := "pillar"
				if tile.Type == Growth {
					st = "growth"
				}
				var s string
				above := tile.Chunk.Cave.GetTileInt(tile.RCoords.X, tile.RCoords.Y-1)
				below := tile.Chunk.Cave.GetTileInt(tile.RCoords.X, tile.RCoords.Y+1)
				aboveM := above.Type == tile.Type && above.Biome == tile.Biome
				belowM := below.Type == tile.Type && below.Biome == tile.Biome
				if above.Solid() && below.Solid() {
					s = "%s_single"
				} else if above.Solid() && belowM {
					s = "%s_top"
				} else if above.Solid() {
					s = "%s_top_br"
				} else if below.Solid() && aboveM {
					s = "%s_bottom"
				} else if below.Solid() {
					s = "%s_bottom_br"
				} else if aboveM && belowM {
					s = "%s_mid"
				} else if aboveM {
					s = "%s_mid_br_dwn"
				} else if belowM {
					s = "%s_mid_br_up"
				} else {
					s = "%s_mid_br"
				}
				tile.AddSprite(fmt.Sprintf(s, st), pixel.IM, tile.Biome, tile.BG)
			} else if tile.Type == Doorway || tile.Type == Tunnel || tile.Type == SecretDoor || tile.Type == SecretOpen {
				st := "door"
				if tile.Type == Tunnel {
					st = "tunnel"
				} else if tile.Type == SecretDoor {
					st = "secret"
				} else if tile.Type == SecretOpen {
					st = "secret_open"
				}
				var s string
				aboveT := tile.Chunk.Cave.GetTileInt(tile.RCoords.X, tile.RCoords.Y-1)
				belowT := tile.Chunk.Cave.GetTileInt(tile.RCoords.X, tile.RCoords.Y+1)
				leftT := tile.Chunk.Cave.GetTileInt(tile.RCoords.X-1, tile.RCoords.Y)
				rightT := tile.Chunk.Cave.GetTileInt(tile.RCoords.X+1, tile.RCoords.Y)
				above := aboveT.IsDoor() && aboveT.DoorI == tile.DoorI
				below := belowT.IsDoor() && belowT.DoorI == tile.DoorI
				left := leftT.IsDoor() && leftT.DoorI == tile.DoorI
				right := rightT.IsDoor() && rightT.DoorI == tile.DoorI
				if below {
					if right && left {
						s = "%s_t"
					} else if left {
						s = "%s_tr"
					} else {
						s = "%s_tl"
					}
				} else if above {
					if right && left {
						s = "%s"
					} else if left {
						s = "%s_r"
					} else {
						s = "%s_l"
					}
				}
				if s != "" {
					tile.AddSprite(fmt.Sprintf(s, st), pixel.IM, tile.Biome, true)
				}
			}
		}
	}
	if tile.Sprite != nil {
		tile.Sprites = append(tile.Sprites, *tile.Sprite)
	}
}

func (tile *Tile) GetTileCode() string {
	ns := tile.RCoords.Neighbors()
	bs := [8]bool{}
	c := 0
	for i, n := range ns {
		t := tile.Chunk.Cave.GetTileInt(n.X, n.Y)
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

func (tile *Tile) SetSprite(key string, mat pixel.Matrix, biome string, bg bool) {
	var batch *img.Batcher
	if bg {
		batch = img.Batchers[fmt.Sprintf(constants.CaveBGFMT, biome)]
	} else {
		batch = img.Batchers[biome]
	}
	tile.Sprite = &img.Sprite{
		K: key,
		S: batch.GetSprite(key),
		M: mat.ScaledXY(pixel.ZV, tile.Transform.Scalar).Moved(tile.Transform.Pos),
		B: batch,
	}
	tile.Chunk.Cave.UpdateBatch = true
}

func (tile *Tile) IsExit() bool {
	return tile != nil && tile.Exit && tile.Type != SecretDoor
}

func (tile *Tile) Breakable() bool {
	return !(tile.Type == Wall || tile.Type == Empty || tile.Type == Blank || tile.IsDoor())
}

func (tile *Tile) Solid() bool {
	return !(tile.IsDeco() || tile.Type == Empty || tile.Type == Blank)
}

func (tile *Tile) Diggable() bool {
	return tile.Type == Dig || tile.Type == Collapse
}

func (tile *Tile) IsDeco() bool {
	return tile.IsDoor() || tile.Type == Pillar || tile.Type == Growth
}

func (tile *Tile) IsDoor() bool {
	return tile.Type == Doorway || tile.Type == Tunnel || tile.Type == SecretDoor || tile.Type == SecretOpen
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