package cave

import (
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
	"math/rand"
)

var CurrCave *Cave

type Cave struct {
	RChunks map[world.Coords]*Chunk
	LChunks map[world.Coords]*Chunk
	batcher *img.Batcher
}

func NewCave(spriteSheet *img.SpriteSheet) *Cave {
	batcher := img.NewBatcher(spriteSheet)
	cave := &Cave{
		RChunks: nil,
		LChunks: nil,
		batcher: batcher,
	}
	chunk0 := GenerateChunk(world.Coords{X: 0, Y: 0}, cave)
	CarveEntrance(chunk0)

	chunkr1 := GenerateChunk(world.Coords{X: 1, Y: 0}, cave)
	chunkr2 := GenerateChunk(world.Coords{X: 1, Y: 1}, cave)
	chunkr3 := GenerateChunk(world.Coords{X: 0, Y: 1}, cave)

	chunkl1 := GenerateChunk(world.Coords{X: -1, Y: 0}, cave)
	chunkl2 := GenerateChunk(world.Coords{X: -1, Y: 1}, cave)

	cave.RChunks = make(map[world.Coords]*Chunk)
	cave.RChunks[chunk0.Coords] = chunk0
	cave.RChunks[chunkr1.Coords] = chunkr1
	cave.RChunks[chunkr2.Coords] = chunkr2
	cave.RChunks[chunkr3.Coords] = chunkr3

	cave.LChunks = make(map[world.Coords]*Chunk)
	cave.LChunks[chunkl1.Coords] = chunkl1
	cave.LChunks[chunkl2.Coords] = chunkl2
	return cave
}

func (cave *Cave) Update(pos pixel.Vec) {
	cave.batcher.Clear()
	p := WorldToChunk(pos)
	all := append([]world.Coords{p}, p.Neighbors()...)
	for _, i := range all {
		if i.X >= 0 && i.Y >= 0 {
			if _, ok := cave.RChunks[i]; !ok {
				cave.RChunks[i] = GenerateChunk(i, cave)
			}
		} else if i.X < 0 && i.Y >= 0 {
			if _, ok := cave.LChunks[i]; !ok {
				cave.LChunks[i] = GenerateChunk(i, cave)
			}
		}
	}
	for i, chunk := range cave.RChunks {
		dis := world.CoordsIn(i, all)
		if dis && !chunk.display {
			chunk.reload = true
		}
		chunk.display = dis
		chunk.Update()
	}
	for i, chunk := range cave.LChunks {
		dis := world.CoordsIn(i, all)
		if dis && !chunk.display {
			chunk.reload = true
		}
		chunk.display = dis
		chunk.Update()
	}
}

func (cave *Cave) Draw(win *pixelgl.Window) {
	for _, chunk := range cave.RChunks {
		chunk.Draw(cave.batcher.Batch())
	}
	for _, chunk := range cave.LChunks {
		chunk.Draw(cave.batcher.Batch())
	}
	cave.batcher.Draw(win)
}

func (cave *Cave) Get(coords world.Coords) *Chunk {
	if chunkR, okR := cave.RChunks[coords]; okR {
		return chunkR
	} else if chunkL, okL := cave.LChunks[coords]; okL {
		return chunkL
	} else {
		return nil
	}
}

func (cave *Cave) GetTile(v pixel.Vec) *Tile {
	ch := WorldToChunk(v)
	tl := WorldToTile(v, ch.X < 0)
	chunk := cave.Get(ch)
	return chunk.Get(tl)
}

func (cave *Cave) SmartTileNum(list string) (string, pixel.Matrix) {
	switch list {
	case "1111":
		return "num_0", pixel.IM
	case "1101":
		mat := img.IM
		if rand.Intn(2) == 1 {
			mat = img.Flip
		}
		return "num_1", mat
	case "1011":
		mat := img.IM
		if rand.Intn(2) == 1 {
			mat = img.Flip
		}
		return "num_1", mat.Rotated(pixel.ZV, math.Pi * -0.5)
	case "0111":
		mat := img.Flop
		if rand.Intn(2) == 1 {
			mat = img.FlipFlop
		}
		return "num_1", mat
	case "1110":
		mat := img.IM
		if rand.Intn(2) == 1 {
			mat = img.Flip
		}
		return "num_1", mat.Rotated(pixel.ZV, math.Pi * 0.5)
	case "1001":
		mat := img.IM
		if rand.Intn(2) == 1 {
			mat = img.FlipFlop.Rotated(pixel.ZV, math.Pi)
		}
		return "num_2", mat
	case "0011":
		mat := img.Flop
		if rand.Intn(2) == 1 {
			mat = img.IM.Rotated(pixel.ZV, math.Pi * -0.5)
		}
		return "num_2", mat
	case "0110":
		mat := img.FlipFlop
		if rand.Intn(2) == 1 {
			mat = img.IM.Rotated(pixel.ZV, math.Pi)
		}
		return "num_2", mat
	case "1100":
		mat := img.Flip
		if rand.Intn(2) == 1 {
			mat = img.IM.Rotated(pixel.ZV, math.Pi * 0.5)
		}
		return "num_2", mat
	default:
		return "num_0", img.IM
	}
}

func (cave *Cave) SmartTileSolid(list string) (string, pixel.Matrix) {
	switch list {
	case "11111111":
		return "full_0", pixel.IM
	case "11100011", "11110011", "11100111", "11110111":
		mat := pixel.IM
		if rand.Intn(2) == 1 {
			mat = img.Flip
		}
		return "full_1", mat
	case "00111110", "00111111", "01111110", "01111111":
		mat := img.Flop
		if rand.Intn(2) == 1 {
			mat = img.FlipFlop
		}
		return "full_1", mat
	case "10001111", "11001111", "10011111", "11011111":
		mat := pixel.IM
		if rand.Intn(2) == 1 {
			mat = img.Flip
		}
		return "full_1", mat.Rotated(pixel.ZV, math.Pi * -0.5)
	case "11111000", "11111100", "11111001", "11111101":
		mat := pixel.IM
		if rand.Intn(2) == 1 {
			mat = img.Flip
		}
		return "full_1", mat.Rotated(pixel.ZV, math.Pi * 0.5)
	case "10000011", "10000111", "11000011", "11000111",
		"10010011", "10010111", "11010011", "11010111":
		mat := pixel.IM
		if rand.Intn(2) == 1 {
			mat = img.FlipFlop.Rotated(pixel.ZV, math.Pi)
		}
		return "full_2", mat
	case "00001110", "00011110", "00001111", "00011111",
		"01001110", "01011110", "01001111", "01011111":
		mat := img.Flop
		if rand.Intn(2) == 1 {
			mat = img.Flip.Rotated(pixel.ZV, math.Pi)
		}
		return "full_2", mat
	case "00111000", "01111000", "00111100", "01111100",
		"00111001", "01111001", "00111101", "01111101":
		mat := img.FlipFlop
		if rand.Intn(2) == 1 {
			mat = img.IM.Rotated(pixel.ZV, math.Pi)
		}
		return "full_2", mat
	case "11100000", "11110000", "11100001", "11110001",
		"11100100", "11110100", "11100101", "11110101":
		mat := img.Flip
		if rand.Intn(2) == 1 {
			mat = img.Flop.Rotated(pixel.ZV, math.Pi)
		}
		return "full_2", mat
	case "00000010", "00000011", "00000110", "00000111",
		"00010010", "00010011", "00010110", "00010111",
		"01000010", "01000011", "01000110", "01000111",
		"01010010", "01010011", "01010110", "01010111":
		mat := pixel.IM
		if rand.Intn(2) == 1 {
			mat = img.Flop
		}
		return "full_3", mat
	case "10000000", "11000000", "10000001", "11000001",
		"10010000", "11010000", "10010001", "11010001",
		"10000100", "11000100", "10000101", "11000101",
		"10010100", "11010100", "10010101", "11010101":
		mat := pixel.IM
		if rand.Intn(2) == 1 {
			mat = img.Flop
		}
		return "full_3", mat.Rotated(pixel.ZV, math.Pi * 0.5)
	case "00100000", "00110000", "01100000", "01110000",
		"00100100", "00110100", "01100100", "01110100",
		"00100001", "00110001", "01100001", "01110001",
		"00100101", "00110101", "01100101", "01110101":
		mat := img.Flip
		if rand.Intn(2) == 1 {
			mat = img.FlipFlop
		}
		return "full_3", mat
	case "00001000", "00001100", "00011000", "00011100",
		"01001000", "01001100", "01011000", "01011100",
		"00001001", "00001101", "00011001", "00011101",
		"01001001", "01001101", "01011001", "01011101":
		mat := pixel.IM
		if rand.Intn(2) == 1 {
			mat = img.Flop
		}
		return "full_3", mat.Rotated(pixel.ZV, math.Pi * -0.5)
	case "01110111", "00110111", "01100111", "00100111",
		"01110011", "00110011", "01100011", "00100011",
		"01110110", "00110110", "01100110", "00100110",
		"01110010", "00110010", "01100010", "00100010":
		mat := pixel.IM
		if rand.Intn(2) == 1 {
			mat = img.Flop
		}
		return "full_4", mat
	case "11011101", "11011100", "10011101", "10011100",
		"11001101", "11001100", "10001101", "10001100",
		"11011001", "11011000", "10011001", "10011000",
		"11001001", "11001000", "10001001", "10001000":
		mat := pixel.IM
		if rand.Intn(2) == 1 {
			mat = img.Flop
		}
		return "full_4", mat.Rotated(pixel.ZV, math.Pi * -0.5)
	case "00000000", "01010101",
		"01010100", "01010001", "01000101", "00010101",
		"01010000", "01000100", "01000001",
		"00010100", "00010001", "00000101",
		"01000000", "00010000", "00000100", "00000001":
		mat := pixel.IM
		if rand.Intn(2) == 1 {
			mat = img.Flip
		}
		return "full_5", mat
	case "11111011":
		mat := pixel.IM
		if rand.Intn(2) == 1 {
			mat = img.FlipFlop.Rotated(pixel.ZV, math.Pi)
		}
		return "full_6", mat
	case "11111110":
		mat := img.Flop
		if rand.Intn(2) == 1 {
			mat = img.Flip.Rotated(pixel.ZV, math.Pi)
		}
		return "full_6", mat
	case "11101111":
		mat := img.Flip
		if rand.Intn(2) == 1 {
			mat = img.Flop.Rotated(pixel.ZV, math.Pi)
		}
		return "full_6", mat
	case "10111111":
		mat := img.FlipFlop
		if rand.Intn(2) == 1 {
			mat = img.IM.Rotated(pixel.ZV, math.Pi)
		}
		return "full_6", mat
	case "11101011":
		mat := pixel.IM
		if rand.Intn(2) == 1 {
			mat = img.Flip
		}
		return "full_7", mat
	case "11111010":
		mat := img.IM
		if rand.Intn(2) == 1 {
			mat = img.Flip
		}
		return "full_7", mat.Rotated(pixel.ZV, math.Pi * 0.5)
	case "10101111":
		mat := img.IM
		if rand.Intn(2) == 1 {
			mat = img.Flip
		}
		return "full_7", mat.Rotated(pixel.ZV, math.Pi * -0.5)
	case "10111110":
		mat := img.Flop
		if rand.Intn(2) == 1 {
			mat = img.FlipFlop
		}
		return "full_7", mat
	case "10111011":
		mat := pixel.IM
		if rand.Intn(2) == 1 {
			mat = img.FlipFlop
		}
		return "full_8", mat
	case "11101110":
		mat := img.Flop
		if rand.Intn(2) == 1 {
			mat = img.Flip
		}
		return "full_8", mat
	case "10101011":
		mat := pixel.IM
		if rand.Intn(2) == 1 {
			mat = img.FlipFlop.Rotated(pixel.ZV, math.Pi)
		}
		return "full_9", mat
	case "11101010":
		mat := img.Flip
		if rand.Intn(2) == 1 {
			mat = img.Flop.Rotated(pixel.ZV, math.Pi)
		}
		return "full_9", mat
	case "10101110":
		mat := img.FlipFlop
		if rand.Intn(2) == 1 {
			mat = img.IM.Rotated(pixel.ZV, math.Pi)
		}
		return "full_9", mat
	case "10111010":
		mat := img.Flop
		if rand.Intn(2) == 1 {
			mat = img.Flip.Rotated(pixel.ZV, math.Pi)
		}
		return "full_9", mat
	case "10101010":
		mat := img.IM
		c := rand.Intn(4)
		if c == 1 {
			mat = img.Flip
		} else if c == 2 {
			mat = img.Flop
		} else if c == 3 {
			mat = img.FlipFlop
		}
		return "full_10", mat
	case "11110110", "11110010", "11100110", "11100010":
		mat := img.IM
		if rand.Intn(2) == 1 {
			mat = img.FlipFlop.Rotated(pixel.ZV, math.Pi)
		}
		return "full_11", mat
	case "10111101", "10111100", "10111001", "10111000":
		mat := img.IM.Rotated(pixel.ZV, math.Pi * 0.5)
		if rand.Intn(2) == 1 {
			mat = img.FlipFlop.Rotated(pixel.ZV, math.Pi * -0.5)
		}
		return "full_11", mat
	case "01101111", "00101111", "01101110", "00101110":
		mat := img.FlipFlop
		if rand.Intn(2) == 1 {
			mat = img.IM.Rotated(pixel.ZV, math.Pi)
		}
		return "full_11", mat
	case "11011011", "11001011", "10011011", "10001011":
		mat := img.IM.Rotated(pixel.ZV, math.Pi * -0.5)
		if rand.Intn(2) == 1 {
			mat = img.FlipFlop.Rotated(pixel.ZV, math.Pi * 0.5)
		}
		return "full_11", mat
	case "10110111", "10110011", "10100111", "10100011":
		mat := img.Flip
		if rand.Intn(2) == 1 {
			mat = img.Flop.Rotated(pixel.ZV, math.Pi)
		}
		return "full_11", mat
	case "11101101", "11101100", "11101001", "11101000":
		mat := img.Flip.Rotated(pixel.ZV, math.Pi * 0.5)
		if rand.Intn(2) == 1 {
			mat = img.Flop.Rotated(pixel.ZV, math.Pi * -0.5)
		}
		return "full_11", mat
	case "01111011", "00111011", "01111010", "00111010":
		mat := img.Flop
		if rand.Intn(2) == 1 {
			mat = img.Flip.Rotated(pixel.ZV, math.Pi)
		}
		return "full_11", mat
	case "11011110", "11001110", "10011110", "10001110":
		mat := img.Flop.Rotated(pixel.ZV, math.Pi * 0.5)
		if rand.Intn(2) == 1 {
			mat = img.Flip.Rotated(pixel.ZV, math.Pi * -0.5)
		}
		return "full_11", mat
	case "10110110", "10110010", "10100110", "10100010":
		mat := img.IM
		if rand.Intn(2) == 1 {
			mat = img.Flip
		}
		return "full_12", mat
	case "10101101", "10101100", "10101001", "10101000":
		mat := img.IM
		if rand.Intn(2) == 1 {
			mat = img.Flip
		}
		return "full_12", mat.Rotated(pixel.ZV, math.Pi * 0.5)
	case "01101011", "00101011", "01101010", "00101010":
		mat := img.Flop
		if rand.Intn(2) == 1 {
			mat = img.FlipFlop
		}
		return "full_12", mat
	case "11011010", "11001010", "10011010", "10001010":
		mat := img.IM
		if rand.Intn(2) == 1 {
			mat = img.Flip
		}
		return "full_12", mat.Rotated(pixel.ZV, math.Pi * -0.5)
	case "11010110", "11010010", "11000110", "11000010",
		"10010110", "10010010", "10000110", "10000010":
		mat := pixel.IM
		if rand.Intn(2) == 1 {
			mat = img.FlipFlop.Rotated(pixel.ZV, math.Pi)
		}
		return "full_13", mat
	case "01011011", "01001011", "00011011", "00001011",
		"01011010", "01001010", "00011010", "00001010":
		mat := img.Flop
		if rand.Intn(2) == 1 {
			mat = img.Flip.Rotated(pixel.ZV, math.Pi)
		}
		return "full_13", mat
	case "01101101", "00101101", "01101100", "00101100",
		"01101001", "00101001", "01101000", "00101000":
		mat := img.FlipFlop
		if rand.Intn(2) == 1 {
			mat = img.IM.Rotated(pixel.ZV, math.Pi)
		}
		return "full_13", mat
	case "10110101", "10110100", "10110001", "10110000",
		"10100101", "10100100", "10100001", "10100000":
		mat := img.Flip
		if rand.Intn(2) == 1 {
			mat = img.Flop.Rotated(pixel.ZV, math.Pi)
		}
		return "full_13", mat
	default:
		return "full_1", pixel.IM
	}
}

func WorldToChunk(v pixel.Vec) world.Coords {
	if v.X >= 0 - world.TileSize * 0.5 {
		return world.Coords{X: int((v.X+world.TileSize*0.5) / ChunkSize / world.TileSize), Y: int(-(v.Y-world.TileSize*0.5) / ChunkSize / world.TileSize)}
	} else {
		return world.Coords{X: int((v.X+world.TileSize*0.5) / ChunkSize / world.TileSize)-1, Y: int(-(v.Y-world.TileSize*0.5) / ChunkSize / world.TileSize)}
	}
}

func WorldToTile(v pixel.Vec, left bool) world.Coords {
	x, y := world.WorldToMap(v.X+world.TileSize*0.5, -(v.Y-world.TileSize*0.5))
	x = x % ChunkSize
	y = y % ChunkSize
	if left {
		x = (ChunkSize - (util.Abs(x) + 1)) % ChunkSize
	}
	return world.Coords{
		X: x % ChunkSize,
		Y: y % ChunkSize,
	}
}