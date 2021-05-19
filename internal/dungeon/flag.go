package dungeon

import (
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/img"
	"github.com/faiface/pixel"
	"math"
)

type Flag struct {
	Transform *transform.Transform
	Tile      *Tile
	created   bool
	done      bool
	animation *img.Instance
}

func (f *Flag) Update() {
	if f.created && !f.done {
		f.Transform.Update()
		f.animation.Update()
		f.animation.SetMatrix(f.Transform.Mat)
		if !f.Tile.Solid || f.Tile.destroyed || !f.Tile.marked {
			f.done = true
			// todo: particles?
		}
	}
}

func (f *Flag) Draw(target pixel.Target) {
	if f.created && !f.done {
		f.animation.Draw(target)
	}
}

func (f *Flag) Create(from pixel.Vec, batcher *img.Batcher) {
	f.Transform = transform.NewTransform()
	f.Transform.Pos = f.Tile.Transform.Pos
	f.created = true
	f.animation = batcher.Animations["flag_hang"].NewInstance()
}

func (f *Flag) Remove() bool {
	return f.done
}

type OldFlag struct {
	Transform *transform.Transform
	Tile      *Tile
	created   bool
	done      bool
	sprite    *pixel.Sprite
}

func (f *OldFlag) Update() {
	if f.created && !f.done {
		f.Transform.Update()
		if !f.Tile.Solid || f.Tile.destroyed || !f.Tile.marked {
			f.done = true
			// todo: particles?
		}
	}
}

func (f *OldFlag) Draw(target pixel.Target) {
	if f.created && !f.done {
		f.sprite.Draw(target, f.Transform.Mat)
	}
}

func (f *OldFlag) Create(from pixel.Vec, batcher *img.Batcher) {
	f.Transform = transform.NewTransform()
	f.created = true
	if f.Tile != nil {
		f.Transform.Pos = f.Tile.Transform.Pos
		mag := from.Sub(f.Tile.Transform.Pos)
		angle := mag.Angle()
		ns := f.Tile.SubCoords.Neighbors()
		var is []int
		if angle > math.Pi * 0.25 && angle < math.Pi * 0.75 {
			// top
			is = []int{4,6,2,0}
		} else if angle < math.Pi * -0.25 && angle > math.Pi * -0.75 {
			// bottom
			is = []int{0,2,6,4}
		} else if angle < math.Pi * 0.25 && angle > math.Pi * -0.25 {
			// right
			is = []int{2,0,4,6}
		} else {
			// left
			is = []int{6,4,0,2}
		}
		if angle > math.Pi * 0.5 && angle < math.Pi {
			// top left
			is = append(is, []int{5,3,7,1}...)
		} else if angle > 0. && angle < math.Pi * 0.5 {
			// top right
			is = append(is, []int{3,5,1,7}...)
		} else if angle < 0. && angle > math.Pi * -0.5 {
			// bottom right
			is = append(is, []int{1,3,7,5}...)
		} else {
			// left
			is = append(is, []int{7,5,1,3}...)
		}
		for _, i := range is {
			t := f.Tile.Chunk.Get(ns[i])
			if t != nil && !t.Solid {
				if i % 2 == 0 {
					f.sprite = batcher.Sprites["flag"]
				} else {
					f.sprite = batcher.Sprites["flag_c"]
				}
				f.Transform.Pos = t.Transform.Pos
				switch i {
				case 0: // bottom
					f.Transform.Rot = math.Pi
				case 2: // right
					f.Transform.Rot = math.Pi * -0.5
				case 6: // left
					f.Transform.Flop = true
					f.Transform.Rot = math.Pi * -0.5
				case 3: // top right
					f.Transform.Flip = true
				case 1: // bottom right
					f.Transform.Flip = true
					f.Transform.Rot = math.Pi * -0.5
				case 7: // bottom left
					f.Transform.Rot = math.Pi * 0.5
				}
				break
			}
		}
	}
}

func (f *OldFlag) Remove() bool {
	return f.done
}