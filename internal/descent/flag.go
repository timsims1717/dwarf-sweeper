package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/profile"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/transform"
	"fmt"
)

func CreateFlag(p *data.Player, tile *cave.Tile) {
	correct := tile.Bomb
	if correct {
		profile.CurrentProfile.Stats.CorrectFlags++
		p.Stats.CorrectFlags++
	} else {
		profile.CurrentProfile.Stats.WrongFlags++
		p.Stats.WrongFlags++
	}
	e := myecs.Manager.NewEntity()
	trans := transform.New().WithID("flag")
	trans.Pos = tile.Transform.Pos
	anim := reanimator.NewSimple(reanimator.NewAnimFromSprites("flag_hang", img.Batchers[constants.ParticleKey].GetAnimation(fmt.Sprintf("flag_hang_%s", p.Code)).S, reanimator.Loop))
	fn := data.NewFrameFunc(func() bool {
		if !tile.Solid() || tile.Destroyed || !tile.Flagged {
			tile.Flagged = false
			if tile.Solid() {
				if correct {
					profile.CurrentProfile.Stats.CorrectFlags--
					p.Stats.CorrectFlags--
				} else {
					profile.CurrentProfile.Stats.WrongFlags--
					p.Stats.WrongFlags--
				}
			} else if correct {
				profile.CurrentProfile.Stats.CorrectFlags--
				p.Stats.CorrectFlags--
			}
			myecs.Manager.DisposeEntity(e)
		}
		return false
	})
	e.AddComponent(myecs.Transform, trans).
		AddComponent(myecs.Temp, myecs.ClearFlag(false)).
		AddComponent(myecs.Func, fn).
		AddComponent(myecs.Animation, anim).
		AddComponent(myecs.Drawable, anim).
		AddComponent(myecs.Batch, constants.ParticleKey)
}