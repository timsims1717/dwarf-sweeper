package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"fmt"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

const (
	slugSpeed = 22.
	slugAcc   = 16.
)

type Slug struct {
	Transform  *transform.Transform
	Physics    *physics.Physics
	Collider   *data.Collider
	Reanimator *reanimator.Tree
	Entity     *ecs.Entity
	Health     *data.Health
	faceLeft   bool
	hasDir     bool
	floor      data.Direction
	iCorner    bool
	fell       bool
	firstFall  bool
	move       bool
}

func CreateSlug(c *cave.Cave, pos pixel.Vec) *Slug {
	s := &Slug{}
	s.Transform = transform.New().WithID("slug")
	s.Transform.Pos = pos
	s.Physics = physics.New()
	s.Physics.UseWallFriction = true
	s.Physics.Friction *= 0.25
	s.Physics.WallFriction *= 0.25
	s.Physics.GravityOff = true
	s.Health = &data.Health{
		Max:          2,
		Curr:         2,
		TempInvTimer: timing.New(0.5),
		TempInvSec:   0.5,
		Immune:       data.EnemyImmunity,
	}
	s.Reanimator = reanimator.New(reanimator.NewSwitch().
		AddAnimation(reanimator.NewAnimFromSprites("slug_move", []*pixel.Sprite{img.Batchers[constants.EntityKey].GetFrame("slug_move", 0)}, reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("slug_corner", img.Batchers[constants.EntityKey].Animations["slug_corner"].S, reanimator.Tran).
		SetTrigger(8, func() {
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
		AddAnimation(reanimator.NewAnimFromSprites("slug_move", img.Batchers[constants.EntityKey].Animations["slug_move"].S, reanimator.Loop).
			SetTrigger(2, func() {
				s.move = true
				if s.Transform.Load {
					sfx.SoundPlayer.PlaySound(fmt.Sprintf("sludge%d", random.Effects.Intn(4)+1), 3.)
				}
			}).
			SetTrigger(4, func() {
				s.move = false
			})).
		SetChooseFn(func() int {
			if s.Health.Dazed || s.Health.Dead {
				return 0
			} else if s.iCorner {
				return 1
			} else {
				return 2
			}
		}), "slug_move")
	s.Collider = data.NewCollider(pixel.R(0., 0., 16., 12.), data.Critter)
	s.Collider.Debug = true
	s.Entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, s.Transform).
		AddComponent(myecs.Physics, s.Physics).
		AddComponent(myecs.Collision, s.Collider).
		AddComponent(myecs.Health, s.Health).
		AddComponent(myecs.Update, data.NewFrameFunc(s.Update)).
		AddComponent(myecs.Animation, s.Reanimator).
		AddComponent(myecs.Drawable, s.Reanimator).
		AddComponent(myecs.Batch, constants.EntityKey).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
	return s
}

func (s *Slug) Update() bool {
	if s.fell && s.Physics.Grounded {
		s.hasDir = false
		s.floor = data.Down
		s.Transform.Rot = 0.
		s.iCorner = false
		s.fell = false
		if s.Transform.Load && s.firstFall {
			sfx.SoundPlayer.PlaySound("splat", 0.)
		}
		s.firstFall = true
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
					if s.Collider.CDR && !s.Collider.CDL {
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
					if s.Collider.CUL && !s.Collider.CUR {
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
					if s.Collider.CDL && !s.Collider.CUL {
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
					if s.Collider.CUR && !s.Collider.CDR {
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
					if s.Collider.CDL && !s.Collider.CDR {
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
					if s.Collider.CUR && !s.Collider.CUL {
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
					if s.Collider.CUL && !s.Collider.CDL {
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
					if s.Collider.CDR && !s.Collider.CUR {
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
			s.Physics.UseWallFriction = false
			s.fell = true
			s.Collider.Hitbox = pixel.R(0., 0., 16., 12.)
		} else {
			s.Physics.GravityOff = true
			s.Physics.UseWallFriction = true
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
				s.Reanimator.ForceUpdate()
				s.Physics.CancelMovement()
				s.move = false
			} else if s.move {
				s.Move()
			}
			s.Transform.Flip = s.faceLeft
			myecs.Manager.NewEntity().AddComponent(myecs.AreaDmg, &data.AreaDamage{
				SourceID:  s.Transform.ID,
				Center:    s.Transform.Pos,
				Radius:    4.,
				Amount:    1,
				Dazed:     1.,
				Knockback: 8.,
				Type:      data.Enemy,
			})
		}
	} else {
		s.Physics.GravityOff = false
		s.Physics.UseWallFriction = false
		s.hasDir = false
		s.iCorner = false
		s.floor = data.Down
		s.Transform.Rot = 0.
		s.Collider.Hitbox = pixel.R(0., 0., 16., 12.)
	}
	if s.Health.Dead {
		s.Entity.RemoveComponent(myecs.Update)
		s.Entity.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
			myecs.AddEffect(s.Entity, data.NewBlink(2.))
			return true
		}, 2.))
		s.Entity.AddComponent(myecs.Temp, timing.New(4.))
	}
	return false
}

func (s *Slug) Move() {
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