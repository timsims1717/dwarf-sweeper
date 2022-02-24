package systems

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/transform"
)

func CollisionBoundSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasCollision) {
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		coll, okC := result.Components[myecs.Collision].(*data.Collider)
		phys, okP := result.Components[myecs.Physics].(*physics.Physics)
		if okT && okC && okP && tran.Load {
			phys.TopBound = coll.TopBound
			phys.BottomBound = coll.BottomBound
			phys.LeftBound = coll.LeftBound
			phys.RightBound = coll.RightBound
		}
	}
}
