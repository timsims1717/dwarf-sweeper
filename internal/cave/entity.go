package cave

import (
	"dwarf-sweeper/pkg/img"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Entity interface {
	Update()
	Draw(pixel.Target)
	Create(pixel.Vec, *img.Batcher)
	Remove() bool
}

var Entities entities

type entities struct{
	set     []Entity
	batcher *img.Batcher
}

func (e *entities) Initialize() {
	sheet, err := img.LoadSpriteSheet("assets/img/entities.json")
	if err != nil {
		panic(err)
	}
	e.batcher = img.NewBatcher(sheet)
}

func (e *entities) Update() {
	var drop []int
	for i, o := range e.set {
		o.Update()
		if o.Remove() {
			drop = append(drop, i)
		}
	}
	for i := len(drop)-1; i >= 0; i-- {
		e.set = append(e.set[:drop[i]], e.set[drop[i]+1:]...)
	}
}

func (e *entities) Draw(win *pixelgl.Window) {
	e.batcher.Clear()
	for _, o := range e.set {
		o.Draw(e.batcher.Batch())
	}
	e.batcher.Draw(win)
}

func (e *entities) Add(entity Entity, vec pixel.Vec) {
	entity.Create(vec, e.batcher)
	e.set = append(e.set, entity)
}

func (e *entities) Clear() {
	e.set = []Entity{}
}