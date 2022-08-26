package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/data/player"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

var throwShovelItem = &player.Item{
	Key:     "throw_shovel",
	Name:    "Throwing Shovel",
	Temp:    true,
	OnUseFn: func(dPos, tPos pixel.Vec, e *ecs.Entity, _ float64) bool {
		var dT *transform.Transform
		if t, ok := e.GetComponentData(myecs.Transform); ok {
			dT = t.(*transform.Transform)
		} else {
			return false
		}
		ts := myecs.Manager.NewEntity()
		trans := transform.New().WithID("throwing-shovel")
		trans.Pos = dT.Pos
		facing := util.Cardinal(dPos, tPos)
		sec := 0.5
		if facing.X < 0 {
			trans.Flip = true
		} else if facing.X > 0 {
			trans.Flip = false
		} else {
			trans.Flip = dT.Flip
		}
		if facing.Y > 0 {
			if facing.X == 0 {
				sec = 0.5
			} else {
				sec = 0.2
			}
		}
		facing.X += (random.Effects.Float64() - 0.5) * 0.05
		facing.Y += (random.Effects.Float64()) * 0.05
		trans.Pos.X += facing.X * world.TileSize
		trans.Pos.Y += facing.Y * world.TileSize
		phys := physics.New()
		phys.SetVelX(facing.X * 200., 0.)
		phys.SetVelY(facing.Y * 200., 0.)
		accel := true
		moving := true
		coll := data.NewCollider(pixel.R(0., 0., throwShovelSpr.Frame().W()-2., throwShovelSpr.Frame().H()-2.), data.Item)
		coll.Damage = &data.Damage{
			SourceID:  dT.ID,
			Dazed:     2.,
			Knockback: 8.,
			Type:      data.Shovel,
		}
		var t *cave.Tile
		hp := &data.SimpleHealth{Immune: data.ItemImmunity2}
		ts.AddComponent(myecs.Transform, trans).
			AddComponent(myecs.Physics, phys).
			AddComponent(myecs.Collision, coll).
			AddComponent(myecs.Health, hp).
			AddComponent(myecs.Drawable, throwShovelSpr).
			AddComponent(myecs.Batch, constants.EntityKey).
			AddComponent(myecs.Temp, myecs.ClearFlag(false)).
			AddComponent(myecs.Update, data.NewFrameFunc(func() bool {
				if moving {
					if reanimator.FrameSwitch {
						sfx.SoundPlayer.PlaySound("shovel", 0.)
						trans.Rot += 0.5
						if trans.Rot > 1. {
							trans.Rot = -0.5
						}
					}
					coll.Damage.Source = trans.Pos
					if coll.Collided {
						accel = false
						moving = false
						phys.CancelMovement()
						coll.Damage = nil
						ts.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
							myecs.AddEffect(ts, data.NewBlink(2.))
							return true
						}, 2.))
						ts.AddComponent(myecs.Temp, timing.New(4.))
						var p *player.Player
						if pl, ok := e.GetComponentData(myecs.Player); ok {
							p, _ = pl.(*player.Player)
						}
						if phys.RightBound {
							t = DigRight(trans.Pos, p)
							trans.Rot = 0.
						} else if phys.LeftBound {
							t = DigLeft(trans.Pos, p)
							trans.Rot = 0.
						} else if phys.TopBound {
							t = DigUp(trans.Pos, p)
							if trans.Flip {
								trans.Rot = -0.5
							} else {
								trans.Rot = 0.5
							}
						} else if phys.Grounded {
							t = DigDown(trans.Pos, p)
							if trans.Flip {
								trans.Rot = 0.5
							} else {
								trans.Rot = -0.5
							}
						}
					}
				} else {
					if t != nil && t.Solid() {
						phys.CancelMovement()
					}
				}
				if accel {
					phys.SetVelX(facing.X * 200., 0.)
					phys.SetVelY(facing.Y * 200., 0.)
				}
				if hp.Dead {
					myecs.Manager.DisposeEntity(ts)
				}
				return false
			})).
			AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
				accel = false
				return true
			}, sec))
		return true
	},
}

func DigRight(pos pixel.Vec, p *player.Player) *cave.Tile {
	rP := pos
	rP.X += world.TileSize
	t1 := Descent.GetTile(rP)
	if t1 != nil && t1.Solid() && !t1.IsDeco() {
		Dig(t1, p)
		return t1
	}
	up := false
	if t1 != nil && rP.Y > t1.Transform.Pos.Y {
		up = true
	}
	if up {
		rP.Y += world.TileSize
	} else {
		rP.Y -= world.TileSize
	}
	t2 := Descent.GetTile(rP)
	if t2 != nil && t2.Solid() && !t2.IsDeco() {
		if Dig(t2, p) {
			return nil
		}
		return t2
	}
	if up {
		rP.Y -= world.TileSize * 2.
	} else {
		rP.Y += world.TileSize * 2.
	}
	t3 := Descent.GetTile(rP)
	if t3 != nil && t3.Solid() && !t3.IsDeco() {
		if Dig(t3, p) {
			return nil
		}
		return t3
	}
	return nil
}

func DigLeft(pos pixel.Vec, p *player.Player) *cave.Tile {
	lP := pos
	lP.X -= world.TileSize
	t1 := Descent.GetTile(lP)
	if t1 != nil && t1.Solid() && !t1.IsDeco() {
		Dig(t1, p)
		return t1
	}
	up := false
	if t1 != nil && lP.Y > t1.Transform.Pos.Y {
		up = true
	}
	if up {
		lP.Y += world.TileSize
	} else {
		lP.Y -= world.TileSize
	}
	t2 := Descent.GetTile(lP)
	if t2 != nil && t2.Solid() && !t2.IsDeco() {
		if Dig(t2, p) {
			return nil
		}
		return t2
	}
	if up {
		lP.Y -= world.TileSize * 2.
	} else {
		lP.Y += world.TileSize * 2.
	}
	t3 := Descent.GetTile(lP)
	if t3 != nil && t3.Solid() && !t3.IsDeco() {
		if Dig(t3, p) {
			return nil
		}
		return t3
	}
	return nil
}

func DigUp(pos pixel.Vec, p *player.Player) *cave.Tile {
	uP := pos
	uP.Y += world.TileSize
	t1 := Descent.GetTile(uP)
	if t1 != nil && t1.Solid() && !t1.IsDeco() {
		Dig(t1, p)
		return t1
	}
	right := false
	if t1 != nil && uP.X > t1.Transform.Pos.X {
		right = true
	}
	if right {
		uP.X += world.TileSize
	} else {
		uP.X -= world.TileSize
	}
	t2 := Descent.GetTile(uP)
	if t2 != nil && t2.Solid() && !t2.IsDeco() {
		if Dig(t2, p) {
			return nil
		}
		return t2
	}
	if right {
		uP.X -= world.TileSize * 2.
	} else {
		uP.X += world.TileSize * 2.
	}
	t3 := Descent.GetTile(uP)
	if t3 != nil && t3.Solid() && !t3.IsDeco() {
		if Dig(t3, p) {
			return nil
		}
		return t3
	}
	return nil
}

func DigDown(pos pixel.Vec, p *player.Player) *cave.Tile {
	uP := pos
	uP.Y -= world.TileSize
	t1 := Descent.GetTile(uP)
	if t1 != nil && t1.Solid() && !t1.IsDeco() {
		Dig(t1, p)
		return t1
	}
	right := false
	if t1 != nil && uP.X > t1.Transform.Pos.X {
		right = true
	}
	if right {
		uP.X += world.TileSize
	} else {
		uP.X -= world.TileSize
	}
	t2 := Descent.GetTile(uP)
	if t2 != nil && t2.Solid() && !t2.IsDeco() {
		if Dig(t2, p) {
			return nil
		}
		return t2
	}
	if right {
		uP.X -= world.TileSize * 2.
	} else {
		uP.X += world.TileSize * 2.
	}
	t3 := Descent.GetTile(uP)
	if t3 != nil && t3.Solid() && !t3.IsDeco() {
		if Dig(t3, p) {
			return nil
		}
		return t3
	}
	return nil
}