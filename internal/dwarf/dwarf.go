package dwarf

import (
	"dwarf-sweeper/internal/cave"
	"dwarf-sweeper/internal/input"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/animation"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Dwarf struct {
	Transform  *physics.Physics
	Animations map[string]*animation.AnimationInstance
	currAnim   string
}

func NewDwarf() *Dwarf {
	dwarfSheet, err := img.LoadSpriteSheet("assets/img/dwarf.json")
	if err != nil {
		panic(err)
	}
	idle := animation.NewAnimation(dwarfSheet, 0, 2,true, false, 0.5)
	run := animation.NewAnimation(dwarfSheet, 2, 6,true, false, 0.4)
	animations := make(map[string]*animation.AnimationInstance)
	animations["idle"] = idle.NewInstance()
	animations["run"] = run.NewInstance()
	transform := animation.NewTransform(true)
	physicsT := &physics.Physics{
		Transform: transform,
	}
	physicsT.Pos = pixel.V(16 * world.TileSize, -8 * world.TileSize)
	return &Dwarf{
		Transform:  physicsT,
		Animations: animations,
		currAnim:   "idle",
	}
}

func (d *Dwarf) Update(input *input.Input) {
	d.Transform.Update()
	in := cave.CurrCave.GetTile(d.Transform.Pos)
	if in != nil {
		up := cave.CurrCave.GetTile(pixel.V(d.Transform.Pos.X, d.Transform.Pos.Y+world.TileSize))
		dwn := cave.CurrCave.GetTile(pixel.V(d.Transform.Pos.X, d.Transform.Pos.Y-world.TileSize))
		right := cave.CurrCave.GetTile(pixel.V(d.Transform.Pos.X+world.TileSize, d.Transform.Pos.Y))
		left := cave.CurrCave.GetTile(pixel.V(d.Transform.Pos.X-world.TileSize, d.Transform.Pos.Y))
		if up != nil && up.Solid && d.Transform.Pos.Y > in.Transform.Pos.Y {
			d.Transform.Pos.Y = in.Transform.Pos.Y
			if d.Transform.Velocity.Y > 0 {
				d.Transform.Velocity.Y = 0
			}
		}
		if dwn != nil && dwn.Solid && d.Transform.Pos.Y < in.Transform.Pos.Y {
			d.Transform.Pos.Y = in.Transform.Pos.Y
			if d.Transform.Velocity.Y < 0 {
				d.Transform.Velocity.Y = 0
			}
		}
		if right != nil && right.Solid && d.Transform.Pos.X > in.Transform.Pos.X {
			d.Transform.Pos.X = in.Transform.Pos.X
			if d.Transform.Velocity.X > 0 {
				d.Transform.Velocity.X = 0
			}
		}
		if left != nil && left.Solid && d.Transform.Pos.X < in.Transform.Pos.X {
			d.Transform.Pos.X = in.Transform.Pos.X
			if d.Transform.Velocity.X < 0 {
				d.Transform.Velocity.X = 0
			}
		}
	}
	d.Animations[d.currAnim].Update()
	d.Animations[d.currAnim].SetMatrix(d.Transform.Mat)
	camera.Cam.Pos = d.Transform.Pos
}

func (d *Dwarf) Draw(win *pixelgl.Window) {
	d.Animations[d.currAnim].Draw(win)
}