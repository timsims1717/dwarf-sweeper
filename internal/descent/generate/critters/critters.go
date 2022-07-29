package critters

import (
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/random"
	"github.com/faiface/pixel"
)

func AddRandomCritter(c *cave.Cave, keys []string, pos pixel.Vec) {
	if len(keys) > 0 {
		AddCritter(c, keys[random.CaveGen.Intn(len(keys))], pos)
	}
}

func AddCritter(c *cave.Cave, key string, pos pixel.Vec) {
	switch key {
	case "bat":
		descent.CreateBat(c, pos)
	case "slug":
		descent.CreateSlug(c, pos)
	}
}