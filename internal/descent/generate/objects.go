package generate

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/typeface"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
)

func addChest(tile *cave.Tile) {
	tile.Type = cave.Empty
	tile.IsChanged = true
	tile.Fillable = false
	popUp := menus.NewPopUp(fmt.Sprintf("%s to open", typeface.SymbolItem), nil)
	popUp.Symbols = []string{data.GameInput.FirstKey("interact")}
	popUp.Dist = world.TileSize
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Transform, tile.Transform).
		AddComponent(myecs.Sprite, img.Batchers[constants.TileEntityKey].Sprites["chest_closed"]).
		AddComponent(myecs.Batch, constants.TileEntityKey).
		AddComponent(myecs.PopUp, popUp).
		AddComponent(myecs.Temp, myecs.ClearFlag(false)).
		AddComponent(myecs.Interact, &data.Interact{
			OnInteract: func(pos pixel.Vec) bool {
				collectible := ""
				switch random.CaveGen.Intn(5) {
				case 0:
					item := &descent.BombItem{}
					item.Create(pos)
				case 1:
					collectible = descent.Apple
				case 2:
					collectible = descent.Beer
				case 3:
					collectible = descent.BubbleItem
				case 4:
					collectible = descent.XRayItem
				}
				if collectible != "" {
					item := &descent.CollectibleItem{
						Collect: descent.Collectibles[collectible],
					}
					item.Create(pos)
				}
				e.AddComponent(myecs.Sprite, img.Batchers[constants.TileEntityKey].Sprites["chest_opened"])
				e.RemoveComponent(myecs.Interact)
				e.RemoveComponent(myecs.PopUp)
				return true
			},
			Distance:   world.TileSize,
			Interacted: false,
			Remove:     false,
		})
}