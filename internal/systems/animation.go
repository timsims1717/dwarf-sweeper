package systems

import (
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/reanimator"
)

func AnimationSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasAnimation) {
		if anim, ok := result.Components[myecs.Animation].(*reanimator.Tree); ok {
			anim.Update()
		}
	}
}