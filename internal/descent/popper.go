package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
	"math"
)

const (
	popperSpeed = 5.
	seekDist    = 8.
	fireDist    = 5.
)

type PopperAction int

const (
	Wait = iota
	Seek
	Pop
	Unpop
	Aim
	Fire
)

type Popper struct {
	Transform  *transform.Transform
	Physics    *physics.Physics
	Reanimator *reanimator.Tree
	Entity     *ecs.Entity
	created    bool
	Health     *data.Health
	faceLeft   bool
	angle      pixel.Vec
	tile       *cave.Tile
	action     PopperAction
	target     *cave.Tile
}

func (p *Popper) Update() {
	if !p.Health.Dazed && !p.Health.Dead {
		//action := p.action
		//goUnder := false
		ray := Descent.GetPlayer().Transform.Pos.Sub(p.Transform.Pos)
		distance := util.Magnitude(ray)
		//angle := ray.Angle()
		// todo: blocked
		if distance > world.TileSize * seekDist || distance > world.TileSize * fireDist {
			//goUnder = true
		}
		switch p.action {
		case Wait:
			if distance < world.TileSize * seekDist {
				//action = Seek
			}
		case Seek:
			// todo: find target
		}
	} else {

	}
}

func (p *Popper) Create(pos pixel.Vec) {
	p.Transform = transform.NewTransform()
	p.Transform.Pos = pos
	p.Physics = physics.New()
	p.Physics.GravityOff = true
	p.tile = Descent.Cave.GetTile(pos)
	p.Health = &data.Health{
		Max:          2,
		Curr:         2,
		TempInv:      true,
		TempInvTimer: timing.New(1.),
		TempInvSec:   1.,
		Immune:       []data.DamageType{data.Enemy},
	}
	p.created = true
	p.Reanimator = reanimator.New(&reanimator.Switch{
		Elements: reanimator.NewElements(
			reanimator.NewAnimFromSprites("popper_wait", []*pixel.Sprite{img.Batchers[constants.EntityKey].GetFrame("popper_out", 0)}, reanimator.Hold, nil),
			reanimator.NewAnimFromSprites("popper_out", img.Batchers[constants.EntityKey].Animations["popper_out"].S, reanimator.Tran, map[int]func(){
				5: func() {
					p.action = Aim
				},
			}),
			reanimator.NewAnimFromSprites("popper_in", img.Reverse(img.Batchers[constants.EntityKey].Animations["popper_out"].S), reanimator.Tran, map[int]func(){
				5: func() {
					p.action = Seek
				},
			}),
			reanimator.NewAnimFromSprites("popper_side", []*pixel.Sprite{img.Batchers[constants.EntityKey].GetSprite("popper_side")}, reanimator.Hold, nil),
			reanimator.NewAnimFromSprites("popper_diag", []*pixel.Sprite{img.Batchers[constants.EntityKey].GetSprite("popper_diag")}, reanimator.Hold, nil),
			reanimator.NewAnimFromSprites("popper_up", []*pixel.Sprite{img.Batchers[constants.EntityKey].GetSprite("popper_up")}, reanimator.Hold, nil),
		),
		Check: func() int {
			if p.Health.Dazed {
				return 5
			} else if p.action == Wait || p.action == Seek {
				return 0
			} else if p.action == Pop {
				return 1
			} else if p.action == Unpop {
				return 2
			} else if p.angle.Angle() > math.Pi * 0.8 || p.angle.Angle() < math.Pi * 0.2 {
				return 3
			} else if p.angle.Angle() > math.Pi * 0.6 || p.angle.Angle() < math.Pi * 0.4 {
				return 4
			} else {
				return 5
			}
		},
	}, "popper_wait")
	p.Entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Entity, p).
		AddComponent(myecs.Transform, p.Transform).
		AddComponent(myecs.Animation, p.Reanimator).
		AddComponent(myecs.Physics, p.Physics).
		AddComponent(myecs.Health, p.Health).
		AddComponent(myecs.Batch, constants.EntityKey)
}

func (p *Popper) Delete() {
	p.Health.Delete()
	myecs.Manager.DisposeEntity(p.Entity)
}