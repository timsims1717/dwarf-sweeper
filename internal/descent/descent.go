package descent

import (
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate/builder"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/profile"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
)

var (
	Difficulty = 1
	Depth      = 9
)

var Descent = &descent{
	Difficulty: 1,
	CoordsMap:  make(map[string]world.Coords),
}

type descent struct {
	Cave       *cave.Cave
	CurrDepth  int
	Depth      int
	Difficulty int
	Type       cave.CaveType
	canExit    bool
	Exited     bool
	ExitI      int
	Exits      []string
	NextBiome  string
	Builder    *builder.CaveBuilder
	BiomeOrder []string

	Dwarves []*Dwarf

	FreeCam      bool
	DisableInput bool
	CoordsMap    map[string]world.Coords

	Builders [][]builder.CaveBuilder
	Timer    *timing.Timer
}

func New() {
	Clear()
	Descent = &descent{
		Difficulty: Difficulty,
		Depth:      Depth,
		CoordsMap:  make(map[string]world.Coords),
		Timer:      timing.New(0.),
	}
}

func Update() {
	if Descent != nil && Descent.Cave != nil {
		Descent.Timer.Update()
		Descent.Cave.Pivots = []pixel.Vec{}
		for _, d := range Descent.Dwarves {
			Descent.Cave.Pivots = append(Descent.Cave.Pivots, d.Transform.Pos)
		}
		if Descent.Cave.Type == cave.Infinite {
			var all []world.Coords
			for _, pivot := range Descent.Cave.Pivots {
				p := cave.WorldToChunk(pivot)
				all = append(all, p)
				all = append(all, p.Neighbors()...)
			}
			for _, i := range all {
				if _, ok := Descent.Cave.Chunks[i]; !ok {
					Descent.Cave.Chunks[i] = cave.NewChunk(i, Descent.Cave, cave.Collapse)
					Descent.Cave.FillChunk(Descent.Cave.Chunks[i])
					Descent.Cave.UpdateBatch = true
					IncreaseLevelInf()
				}
			}
		}
		Descent.Cave.Update()
		exitM := "{symbol:player-interact}:Exit"
		couldExit := Descent.canExit
		switch Descent.Type {
		case cave.Minesweeper:
			Descent.canExit = (profile.CurrentProfile.Stats.CorrectFlags == Descent.Cave.BombsLeft && profile.CurrentProfile.Stats.WrongFlags < 1) || Descent.GetPlayers()[0].Health.Inv
			if !Descent.canExit {
				exitM = "Flag all the remaining mines to exit."
			}
		default:
			Descent.canExit = true
		}
		for _, c := range Descent.Cave.AllRevealed() {
			below := Descent.Cave.GetTileInt(c.X, c.Y+1)
			if below != nil && below.Type == cave.Empty {
				if random.Effects.Intn(3) == 0 {
					CreateFallingBlock(Descent.Cave, c)
				}
			}
		}
		for i, exit := range Descent.Cave.Exits {
			exitTile := Descent.Cave.GetTileInt(exit.Coords.X, exit.Coords.Y)
			if exit.PopUp == nil {
				exitI := exit.ExitI
				Descent.Cave.Exits[i].PopUp = menus.NewPopUp(exitM)
				Descent.Cave.Exits[i].Type = exitTile.Type
				myecs.Manager.NewEntity().
					AddComponent(myecs.PopUp, Descent.Cave.Exits[i].PopUp).
					AddComponent(myecs.Interact, NewInteract(func(_ pixel.Vec, _ *Dwarf) bool {
						if Descent.canExit && exitTile.IsExit() {
							Descent.Exited = true
							Descent.ExitI = exitI
							Descent.NextBiome = exitTile.Biome
							return true
						}
						return false
					}, world.TileSize, true)).
					AddComponent(myecs.Transform, exitTile.Transform).
					AddComponent(myecs.Temp, myecs.ClearFlag(false))
			} else if Descent.Cave.Exits[i].Type == cave.SecretDoor && exitTile.Type == cave.SecretDoor {
				Descent.Cave.Exits[i].PopUp.SetText("A bomb will open this passage.")
			} else if couldExit != Descent.canExit || Descent.Cave.Exits[i].Type != exitTile.Type {
				Descent.Cave.Exits[i].PopUp.SetText(exitM)
				Descent.Cave.Exits[i].Type = exitTile.Type
			}
		}
	}
}

func UpdatePlayers() {
	for _, d := range Descent.Dwarves {
		UpdatePlayer(d)
	}
}

func UpdateViews() {
	for i, d := range Descent.Dwarves {
		UpdateView(d, i, len(Descent.Dwarves))
	}
}

func (d *descent) CanExit() bool {
	return d.canExit
}

func (d *descent) GetCave() *cave.Cave {
	return d.Cave
}

func (d *descent) SetCave(cave *cave.Cave) {
	d.Cave = cave
	d.Type = cave.Type
	d.Cave.UpdateAllTileSprites()
	d.Cave.UpdateBatch = true
	d.Cave.Biomes = append(d.Cave.Biomes, d.Exits...)
	d.BiomeOrder = append(d.BiomeOrder, d.Cave.Biome)
}

func (d *descent) GetPlayers() []*Dwarf {
	return d.Dwarves
}

func (d *descent) GetClosestPlayer(pos pixel.Vec) *Dwarf {
	dist := -1.
	var pick *Dwarf
	for _, p := range d.GetPlayers() {
		newDist := util.Magnitude(p.Transform.Pos.Sub(pos))
		if dist < 0. || newDist < dist {
			dist = newDist
			pick = p
		}
	}
	return pick
}

func (d *descent) GetClosestPlayerTile(pos pixel.Vec) *cave.Tile {
	dwarf := d.GetClosestPlayer(pos)
	if dwarf != nil {
		return d.Cave.GetTile(dwarf.Transform.Pos)
	}
	return nil
}

func (d *descent) GetRandomPlayer() *Dwarf {
	if len(d.Dwarves) > 0 {
		return d.Dwarves[random.Effects.Intn(len(d.Dwarves))]
	}
	return nil
}

func (d *descent) GetRandomPlayerTile() *cave.Tile {
	dwarf := d.GetRandomPlayer()
	if dwarf != nil {
		return d.Cave.GetTile(dwarf.Transform.Pos)
	}
	return nil
}

func (d *descent) GetTile(pos pixel.Vec) *cave.Tile {
	return d.Cave.GetTile(pos)
}

func IncreaseLevelInf() {
	Descent.CurrDepth++
	Descent.Cave.BombPMin += 0.01
	Descent.Cave.BombPMax += 0.01
	if Descent.Cave.BombPMin > 0.3 {
		Descent.Cave.BombPMin = 0.3
	}
	if Descent.Cave.BombPMax > 0.4 {
		Descent.Cave.BombPMax = 0.4
	}
}

func Clear() {
	if Descent != nil {
		for _, d := range Descent.Dwarves {
			d.Delete()
			if d.Player.Puzzle != nil {
				d.Player.Puzzle.Close()
				d.Player.Puzzle = nil
			}
		}
		Descent.Dwarves = []*Dwarf{}
	}
}
