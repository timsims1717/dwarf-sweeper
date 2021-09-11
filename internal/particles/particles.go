package particles

import (
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/internal/util"
	gween "dwarf-sweeper/pkg/gween64"
	"dwarf-sweeper/pkg/gween64/ease"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"fmt"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image/color"
)

type particle struct {
	Sprite    *pixel.Sprite
	Transform *transform.Transform
	entity    *ecs.Entity
	//ColorEffect transform.ColorEffect
	Frame     bool
	color     color.RGBA
	fader     *gween.Tween
	done      bool
	frame     bool

}

func (p *particle) Update() {
	if !p.done {
		if p.fader != nil {
			a, fin := p.fader.Update(timing.DT)
			if fin {
				p.done = true
			}
			p.color.A = uint8(a)
		}
		if p.Frame && !p.frame {
			p.frame = true
		} else if p.frame {
			p.done = true
		}
	}
}

func Initialize() {
	partSheet, err := img.LoadSpriteSheet("assets/img/particles.json")
	if err != nil {
		panic(err)
	}
	PartBatcher = img.NewBatcher(partSheet)
	// block particles
	for i := 0; i < 8; i++ {
		blocks = append(blocks, fmt.Sprintf("b_%d", i))
	}
}

var PartBatcher *img.Batcher
var particles []*particle

func Update() {
	var drop []int
	for i, p := range particles {
		p.Update()
		if p.done {
			drop = append(drop, i)
		}
	}
	for i := len(drop)-1; i >= 0; i-- {
		particles = append(particles[:drop[i]], particles[drop[i]+1:]...)
	}
}

func Draw(win *pixelgl.Window) {
	PartBatcher.Clear()
	for _, p := range particles {
		//p.Sprite.DrawColorMask(PartBatcher.Batch(), p.Transform.Mat, p.color)
		p.Sprite.Draw(PartBatcher.Batch(), p.Transform.Mat)
	}
	PartBatcher.Draw(win)
	//for _, p := range particles {
	//	debug.AddLine(colornames.Green, imdraw.RoundEndShape, p.Transform.Pos, p.Transform.Pos, 2.)
	//}
}

func Clear() {
	particles = []*particle{}
}

var blocks []string

func BlockParticles(pos pixel.Vec) {
	c := random.Effects.Intn(3) + 4
	for i := 0; i < c; i++ {
		phys, tran := util.RandomVelocity(pos, 1.0, random.Effects)
		if random.Effects.Intn(2) == 0 {
			tran.Flip = true
		}
		if random.Effects.Intn(2) == 0 {
			tran.Flop = true
		}
		particles = append(particles, &particle{
			Sprite:    PartBatcher.Sprites[blocks[random.Effects.Intn(len(blocks))]],
			Transform: tran,
			entity:    myecs.Manager.NewEntity().
				AddComponent(myecs.Transform, tran).
				AddComponent(myecs.Physics, phys),
			color:     colornames.White,
			fader:     gween.New(255., 0., 1.0, ease.Linear),
		})
	}
}
func CreateRandomStaticParticles(min, max int, keys []string, orig pixel.Vec, variance, dur, durVar float64) {
	c := random.Effects.Intn(max - min + 1) + min
	for i := 0; i < c; i++ {
		tran := transform.NewTransform()
		tran.Pos = util.RandomPosition(orig, variance, random.Effects)
		if random.Effects.Intn(2) == 0 {
			tran.Flip = true
		}
		if random.Effects.Intn(2) == 0 {
			tran.Flop = true
		}
		nDur := dur + (random.Effects.Float64() - 0.5) * durVar
		key := keys[random.Effects.Intn(len(keys))]
		particles = append(particles, &particle{
			Sprite:    PartBatcher.Sprites[key],
			Transform: tran,
			entity:    myecs.Manager.NewEntity().
				AddComponent(myecs.Transform, tran),
			color:     colornames.White,
			fader:     gween.New(255., 0., nDur, ease.Linear),
		})
	}
}

func CreateStaticParticle(key string, orig pixel.Vec) {
	tran := transform.NewTransform()
	tran.Pos = orig
	particles = append(particles, &particle{
		Sprite:    PartBatcher.Sprites[key],
		Transform: tran,
		entity:    myecs.Manager.NewEntity().
			AddComponent(myecs.Transform, tran),
		color:     colornames.White,
		Frame:     true,
	})
}