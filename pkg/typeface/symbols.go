package typeface

import (
	"dwarf-sweeper/pkg/transform"
	"github.com/faiface/pixel"
)

var (
	theSymbols = map[string]symbol{}
)

type symbol struct {
	spr *pixel.Sprite
	sca float64
}

type symbolHandle struct {
	symbol symbol
	trans  *transform.Transform
}

func RegisterSymbol(key string, spr *pixel.Sprite, scalar float64) {
	theSymbols[key] = symbol{
		spr: spr,
		sca: scalar,
	}
}