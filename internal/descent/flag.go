package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/data/player"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/transform"
	"fmt"
)

func CreateFlag(p *player.Player, tile *cave.Tile) {
	correct := tile.Bomb
	if correct {
		player.OverallStats.CaveBombsFlagged++
		player.OverallStats.CaveCorrectFlags++
		p.Stats.CaveBombsFlagged++
		p.Stats.CaveCorrectFlags++
	} else {
		player.OverallStats.CaveWrongFlags++
		p.Stats.CaveWrongFlags++
	}
	e := myecs.Manager.NewEntity()
	trans := transform.New()
	trans.Pos = tile.Transform.Pos
	anim := reanimator.NewSimple(reanimator.NewAnimFromSprites("flag_hang", img.Batchers[constants.ParticleKey].GetAnimation(fmt.Sprintf("flag_hang_%s", p.Code)).S, reanimator.Loop))
	fn := data.NewFrameFunc(func() bool {
		if !tile.Solid() || tile.Destroyed || !tile.Flagged {
			tile.Flagged = false
			if tile.Solid() {
				if correct {
					player.OverallStats.CaveBombsFlagged--
					player.OverallStats.CaveCorrectFlags--
					p.Stats.CaveBombsFlagged--
					p.Stats.CaveCorrectFlags--
				} else {
					player.OverallStats.CaveWrongFlags--
					p.Stats.CaveWrongFlags--
				}
			} else if correct {
				player.OverallStats.CaveBombsFlagged--
				p.Stats.CaveBombsFlagged--
			} else {
				player.OverallStats.CaveWrongFlags--
				p.Stats.CaveWrongFlags--
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

//type Flag struct {
//	Transform  *transform.Transform
//	Tile       *cave.Tile
//	created    bool
//	Reanimator *reanimator.Tree
//	entity     *ecs.Entity
//	correct    bool
//}
//
//func (f *Flag) Update() {
//	if f.created {
//		if !f.Tile.Solid() || f.Tile.Destroyed || !f.Tile.Flagged {
//			f.Delete()
//			// todo: particles?
//		}
//	}
//}
//
//func (f *Flag) Create(_ pixel.Vec) {
//	f.Transform = transform.New()
//	f.Transform.Pos = f.Tile.Transform.Pos
//	f.created = true
//	f.correct = f.Tile.Bomb
//	if f.correct {
//		player2.CaveBombsMarked++
//		player2.CaveCorrectMarks++
//	} else {
//		player2.CaveWrongMarks++
//	}
//	f.Reanimator = reanimator.NewSimple(reanimator.NewAnimFromSprites("flag_hang", img.Batchers[constants.ParticleKey].Animations["flag_hang"].S, reanimator.Loop))
//	f.entity = myecs.Manager.NewEntity().
//		AddComponent(myecs.Entity, f).
//		AddComponent(myecs.Transform, f.Transform).
//		AddComponent(myecs.Animation, f.Reanimator).
//		AddComponent(myecs.Drawable, f.Reanimator).
//		AddComponent(myecs.Batch, constants.ParticleKey)
//}
//
//func (f *Flag) Delete() {
//	f.Tile.Flagged = false
//	if f.Tile.Solid() {
//		if f.correct {
//			player2.CaveBombsMarked--
//			player2.CaveCorrectMarks--
//		} else {
//			player2.CaveWrongMarks--
//		}
//	} else if f.correct {
//		player2.CaveBombsMarked--
//	} else {
//		player2.CaveWrongMarks--
//	}
//	myecs.Manager.DisposeEntity(f.entity)
//}
