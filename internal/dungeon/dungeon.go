package dungeon

import "dwarf-sweeper/internal/myecs"

var Dungeon dungeon

type dungeon struct {
	Cave     *Cave
	Level    int
	Player   *Dwarf
	Entities []myecs.AnEntity
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

func (d *dungeon) AddEntity(e myecs.AnEntity) int {
	i := len(d.Entities)
	d.Entities = append(d.Entities, e)
	return i
}

func (d *dungeon) RemoveEntity(id int) {
	if len(d.Entities) == 1 {
		d.Entities = []myecs.AnEntity{}
	} else if len(d.Entities) > 1 {
		d.Entities = append(d.Entities[:id], d.Entities[id+1:]...)
		for i, e := range d.Entities {
			e.SetId(i)
		}
	}
}

func (d *dungeon) RemoveAllEntities() {
	for i := len(d.Entities)-1; i >= 0; i-- {
		d.Entities[i].Delete()
	}
	d.Entities = []myecs.AnEntity{}
}