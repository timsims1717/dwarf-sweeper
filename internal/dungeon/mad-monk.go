package dungeon

import (
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/transform"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
	"time"
)

const (
	mmSpeed  = 40.
	mmAcc    = 5.
	mmDazedT = 3.
)

type MadMonk struct {
	Transform  *transform.Transform
	Physics    *physics.Physics
	Reanimator *reanimator.Tree
	entity     *ecs.Entity
	created    bool
	Dazed      bool
	DazedTimer time.Time
	DazedVFX   *vfx.VFX
	faceLeft   bool
	done       bool
}

func (m *MadMonk) Update() {
	if Dungeon.GetCave().PointLoaded(m.Transform.Pos) {
		if m.Dazed {
			if m.DazedVFX != nil {
				m.DazedVFX.Matrix = pixel.IM.Moved(m.Transform.APos).Moved(pixel.V(0., 9.))
			}
		} else {
			if m.Physics.Grounded {
				ownPos := Dungeon.GetCave().GetTile(m.Transform.Pos).RCoords
				playerPos := Dungeon.GetPlayerTile().RCoords
				if ownPos.X > playerPos.X {
					m.faceLeft = true
					m.Physics.SetVelX(-mmSpeed, mmAcc)
				} else if ownPos.X < playerPos.X {
					m.faceLeft = false
					m.Physics.SetVelX(mmSpeed, mmAcc)
				}
			}
			m.Transform.Flip = m.faceLeft
		}
	}
	if m.Dazed && time.Since(m.DazedTimer).Seconds() > mmDazedT {
		m.Dazed = false
		m.DazedVFX.Animation.Done = true
		m.DazedVFX = nil
	}
}

func (m *MadMonk) Draw(target pixel.Target) {
	if m.created && !m.done && Dungeon.GetCave().PointLoaded(m.Transform.Pos) {
		m.Reanimator.CurrentSprite().Draw(target, m.Transform.Mat)
	}
}

func (m *MadMonk) Create(pos pixel.Vec, batcher *img.Batcher) {
	m.Transform = transform.NewTransform()
	m.Physics = physics.New()
	m.Transform.Pos = pos
	m.created = true
	m.Reanimator = reanimator.New(&reanimator.Switch{
		Elements: reanimator.NewElements(
			reanimator.NewAnimFromSprites("mm_walk", batcher.Animations["mm_walk"].S, reanimator.Loop, nil),
			reanimator.NewAnimFromSprites("mm_idle", batcher.Animations["mm_idle"].S, reanimator.Hold, nil),
		),
		Check: func() int {
			if m.Physics.IsMovingX() {
				return 0
			} else {
				return 1
			}
		},
	}, "mm_idle")
	m.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, m.Transform).
		AddComponent(myecs.Animation, m.Reanimator).
		AddComponent(myecs.Physics, m.Physics).
		AddComponent(myecs.Collision, myecs.Collider{})
	m.DazedVFX = vfx.CreateDazed(m.Transform.APos.Add(pixel.V(0., 9.)))
	m.Dazed = true
	m.DazedTimer = time.Now()
}

func (m *MadMonk) Done() bool {
	return m.done
}

func (m *MadMonk) Delete() {
	myecs.Manager.DisposeEntity(m.entity)
}