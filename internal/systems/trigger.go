package systems

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/transform"
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
)

func TriggerSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasTrigger) {
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		coll, okC := result.Components[myecs.Collision].(*data.Collider)
		fn, okF := result.Components[myecs.Trigger].(*data.FrameFunc)
		if okT && okC && okF {
			if descent.Descent.Player != nil {
				hb := coll.Hitbox.Moved(tran.Pos).Moved(pixel.V(coll.Hitbox.W()*-0.5, coll.Hitbox.H()*-0.5))
				if coll.Debug {
					debug.AddRect(colornames.Red, tran.Pos, coll.Hitbox, 0.5)
				}
				if hb.Contains(descent.Descent.Player.Transform.Pos) {
					if fn.Func() {
						result.Entity.RemoveComponent(myecs.Trigger)
					}
				}
			}
		}
	}
}