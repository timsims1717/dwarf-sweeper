package dungeon

var Dungeon dungeon

type dungeon struct {
	Cave     *Cave
	Level    int
	Player   *Dwarf
	removing bool
	Start    bool
}

func (d *dungeon) GetCave() *Cave {
	return d.Cave
}

func (d *dungeon) SetCave(cave *Cave) {
	d.Cave = cave
}

func (d *dungeon) GetPlayer() *Dwarf {
	return d.Player
}

func (d *dungeon) SetPlayer(dwarf *Dwarf) {
	d.Player = dwarf
}

func (d *dungeon) GetPlayerTile() *Tile {
	return d.Cave.GetTile(d.Player.Transform.Pos)
}