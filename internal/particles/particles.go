package particles

import (
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/transform"
	gween "dwarf-sweeper/pkg/gween64"
	"dwarf-sweeper/pkg/gween64/ease"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image/color"
	"math/rand"
)

type particle struct {
	Sprite      *pixel.Sprite
	Transform   *physics.Physics
	//ColorEffect transform.ColorEffect
	Frame       bool
	color       color.RGBA
	fader       *gween.Tween
	done        bool
	frame       bool

}

func (p *particle) Update() {
	if !p.done {
		p.Transform.Update()
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
	c := rand.Intn(3) + 4
	for i := 0; i < c; i++ {
		particles = append(particles, &particle{
			Sprite:    PartBatcher.Sprites[blocks[rand.Intn(len(blocks))]],
			Transform: randomParticleLocation(pos, 1.0),
			color:     colornames.White,
			fader:     gween.New(255., 0., 1.0, ease.Linear),
		})
	}
}

func randomParticleLocation(orig pixel.Vec, variance float64) *physics.Physics {
	transform := transform.NewTransform(true)
	physicsT := &physics.Physics{Transform: transform}
	physicsT.Pos = orig
	actVar := variance * world.TileSize
	//if square {
	xVar := (rand.Float64() - 0.5) * actVar
	yVar := (rand.Float64() - 0.5) * actVar
	physicsT.Pos.X += xVar
	physicsT.Pos.Y += yVar
	physicsT.Velocity.X = xVar * 0.02
	physicsT.Velocity.Y = 0.5
	//}
	if rand.Intn(2) == 0 {
		physicsT.Flip = true
	}
	if rand.Intn(2) == 0 {
		physicsT.Flop = true
	}
	return physicsT
}

func CreateStaticParticle(key string, orig pixel.Vec) {
	transform := transform.NewTransform(true)
	transform.Pos = orig
	particles = append(particles, &particle{
		Sprite:    PartBatcher.Sprites[key],
		Transform: &physics.Physics{Transform: transform, Off: true},
		color:     colornames.White,
		Frame:     true,
	})
}