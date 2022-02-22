package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

const (
	slugSpeed = 8.
	slugAcc   = 20.
)

type Slug struct {
	Transform  *transform.Transform
	Physics    *physics.Physics
	Collider   *data.Collider
	Reanimator *reanimator.Tree
	Entity     *ecs.Entity
	created    bool
	Health     *data.Health
	faceLeft   bool
	hasDir     bool
	floor      data.Direction
	iCorner    bool
	//oCorner    bool
	fell bool
}

func (s *Slug) Update() {
	if s.fell && s.Physics.Grounded {
		s.hasDir = false
		s.floor = data.Down
		s.Transform.Rot = 0.
		s.iCorner = false
		s.fell = false
	}
	if !s.Health.Dazed && !s.Health.Dead {
		if !s.hasDir {
			s.faceLeft = random.Effects.Intn(2) == 0
			s.hasDir = true
		}
		fall := false
		corner := false
		if s.faceLeft {
			switch s.floor {
			case data.Down:
				if !s.Physics.Grounded {
					if s.Collider.CDR {
						s.floor = data.Right
						s.Transform.Pos.Y -= 3.0
						s.Physics.SetVelX(0., 0.)
					} else {
						fall = true
					}
				} else if s.Physics.LeftBound && !s.iCorner {
					corner = true
					s.floor = data.Left
				}
			case data.Up:
				if !s.Physics.TopBound {
					if s.Collider.CUL {
						s.floor = data.Left
						s.Transform.Pos.Y += 3.0
						s.Physics.SetVelX(0., 0.)
					} else {
						fall = true
					}
				} else if s.Physics.RightBound && !s.iCorner {
					corner = true
					s.floor = data.Right
				}
			case data.Left:
				if !s.Physics.LeftBound {
					if s.Collider.CDL {
						s.floor = data.Down
						s.Transform.Pos.X -= 3.0
						s.Physics.SetVelY(0., 0.)
					} else {
						fall = true
					}
				} else if s.Physics.TopBound && !s.iCorner {
					corner = true
					s.floor = data.Up
				}
			case data.Right:
				if !s.Physics.RightBound {
					if s.Collider.CUR {
						s.floor = data.Up
						s.Transform.Pos.X += 3.0
						s.Physics.SetVelY(0., 0.)
					} else {
						fall = true
					}
				} else if s.Physics.Grounded && !s.iCorner {
					corner = true
					s.floor = data.Down
				}
			}
		} else {
			switch s.floor {
			case data.Down:
				if !s.Physics.Grounded {
					if s.Collider.CDL {
						s.floor = data.Left
						s.Transform.Pos.Y -= 3.0
						s.Physics.SetVelX(0., 0.)
					} else {
						fall = true
					}
				} else if s.Physics.RightBound && !s.iCorner {
					corner = true
					s.floor = data.Right
				}
			case data.Up:
				if !s.Physics.TopBound {
					if s.Collider.CUR {
						s.floor = data.Right
						s.Transform.Pos.Y += 3.0
						s.Physics.SetVelX(0., 0.)
					} else {
						fall = true
					}
				} else if s.Physics.LeftBound && !s.iCorner {
					corner = true
					s.floor = data.Left
				}
			case data.Left:
				if !s.Physics.LeftBound {
					if s.Collider.CUL {
						s.floor = data.Up
						s.Transform.Pos.X -= 3.0
						s.Physics.SetVelY(0., 0.)
					} else {
						fall = true
					}
				} else if s.Physics.Grounded && !s.iCorner {
					corner = true
					s.floor = data.Down
				}
			case data.Right:
				if !s.Physics.RightBound {
					if s.Collider.CDR {
						s.floor = data.Down
						s.Transform.Pos.X += 3.0
						s.Physics.SetVelY(0., 0.)
					} else {
						fall = true
					}
				} else if s.Physics.TopBound && !s.iCorner {
					corner = true
					s.floor = data.Up
				}
			}
		}
		if fall {
			s.Physics.GravityOff = false
			s.fell = true
		} else {
			s.Physics.GravityOff = true
			switch s.floor {
			case data.Down:
				s.Transform.Rot = 0.
			case data.Up:
				s.Transform.Rot = 1.
			case data.Left:
				s.Transform.Rot = -0.5
			case data.Right:
				s.Transform.Rot = 0.5
			}
			if corner {
				s.Collider.Hitbox = pixel.R(0., 0., 16., 16.)
				s.iCorner = true
			} else if !s.iCorner {
				switch s.floor {
				case data.Down:
					if s.faceLeft {
						s.Physics.SetVelX(-slugSpeed, slugAcc)
					} else {
						s.Physics.SetVelX(slugSpeed, slugAcc)
					}
				case data.Up:
					if s.faceLeft {
						s.Physics.SetVelX(slugSpeed, slugAcc)
					} else {
						s.Physics.SetVelX(-slugSpeed, slugAcc)
					}
				case data.Left:
					if s.faceLeft {
						s.Physics.SetVelY(slugSpeed, slugAcc)
					} else {
						s.Physics.SetVelY(-slugSpeed, slugAcc)
					}
				case data.Right:
					if s.faceLeft {
						s.Physics.SetVelY(-slugSpeed, slugAcc)
					} else {
						s.Physics.SetVelY(slugSpeed, slugAcc)
					}
				}
			}
			s.Transform.Flip = s.faceLeft
			myecs.Manager.NewEntity().AddComponent(myecs.AreaDmg, &data.AreaDamage{
				SourceID:  s.Transform.ID,
				Center:    s.Transform.Pos,
				Rect:      s.Collider.Hitbox,
				Amount:    1,
				Dazed:     1.,
				Knockback: 8.,
				Type:      data.Enemy,
			})
		}
	} else {
		s.Physics.GravityOff = false
		s.hasDir = false
		s.iCorner = false
		s.floor = data.Down
		s.Transform.Rot = 0.
	}
	if s.Health.Dead {
		s.Delete()
	}
}

func (s *Slug) Create(pos pixel.Vec) {
	s.Transform = transform.New()
	s.Transform.Pos = pos
	s.Physics = physics.New()
	s.Physics.GravityOff = true
	s.Health = &data.Health{
		Max:          2,
		Curr:         2,
		Dazed:        true,
		DazedTimer:   timing.New(3.),
		TempInvTimer: timing.New(0.5),
		TempInvSec:   0.5,
		Immune:       data.EnemyImmunity,
	}
	s.created = true
	s.Reanimator = reanimator.New(reanimator.NewSwitch().
		AddAnimation(reanimator.NewAnimFromSprites("slug_corner", img.Batchers[constants.EntityKey].Animations["slug_corner"].S, reanimator.Tran).
			SetTrigger(8, func(_ *reanimator.Anim, _ string, _ int) {
				s.iCorner = false
				switch s.floor {
				case data.Up:
					s.Transform.Pos.Y += 2.1
				case data.Down:
					s.Transform.Pos.Y -= 2.1
				case data.Right:
					s.Transform.Pos.X += 2.1
				case data.Left:
					s.Transform.Pos.X -= 2.1
				}
				s.Collider.Hitbox = pixel.R(0., 0., 16., 12.)
			})).
		AddAnimation(reanimator.NewAnimFromSprites("slug_move", img.Batchers[constants.EntityKey].Animations["slug_move"].S, reanimator.Loop)).
		SetChooseFn(func() int {
			if s.iCorner {
				return 0
			} else {
				return 1
			}
		}), "slug_move")
	s.Collider = &data.Collider{
		Hitbox:     pixel.R(0., 0., 16., 12.),
		GroundOnly: true,
	}
	s.Entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Entity, s).
		AddComponent(myecs.Transform, s.Transform).
		AddComponent(myecs.Animation, s.Reanimator).
		AddComponent(myecs.Drawable, s.Reanimator).
		AddComponent(myecs.Physics, s.Physics).
		AddComponent(myecs.Health, s.Health).
		AddComponent(myecs.Collision, s.Collider).
		AddComponent(myecs.Batch, constants.EntityKey)
}

func (s *Slug) Delete() {
	myecs.Manager.DisposeEntity(s.Entity)
}
