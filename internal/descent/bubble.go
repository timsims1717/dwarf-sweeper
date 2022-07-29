package descent

//const (
//	BubbleVel = 45.
//	BubbleAcc = 5.
//)
//
//var BubbleSec = 12.
//
//type Bubble struct {
//	Dwarf      *Dwarf
//	Physics    *physics.Physics
//	Transform  *transform.Transform
//	Reanimator *reanimator.Tree
//	entity     *ecs.Entity
//	popped     bool
//	created    bool
//	Timer      *timing.Timer
//}
//
//func (b *Bubble) Update() {
//	if b.created && b.Timer.UpdateDone() && !b.popped {
//		b.Pop()
//	}
//}
//
//func (b *Bubble) Create(_ pixel.Vec) {
//	b.Timer = timing.New(BubbleSec)
//	b.Transform = transform.New()
//	b.Physics = physics.New()
//	b.Physics.Gravity = 50.
//	b.Physics.Friction = 200.
//	b.Dwarf.Entity.AddComponent(myecs.Physics, b.Physics)
//	b.Dwarf.Bubble = b
//	b.created = true
//	vibe := img.Batchers[constants.ParticleKey].GetAnimation("bubble_vibe")
//	b.Reanimator = reanimator.New(reanimator.NewSwitch().
//		AddAnimation(reanimator.NewAnimFromSprites("bubble_vibe", []*pixel.Sprite{
//			vibe.S[0], vibe.S[0],
//			vibe.S[1], vibe.S[1],
//			vibe.S[2], vibe.S[2],
//			vibe.S[1], vibe.S[1],
//		}, reanimator.Loop)).
//		AddAnimation(reanimator.NewAnimFromSprites("bubble_pop", []*pixel.Sprite{
//			img.Batchers[constants.ParticleKey].GetSprite("bubble_pop"),
//			img.Batchers[constants.ParticleKey].GetSprite("bubble_pop"),
//		}, reanimator.Tran).
//			SetTrigger(2, func(_ *reanimator.Anim, _ string, _ int) {
//				b.Delete()
//			})).
//		SetChooseFn(func() int {
//			if !b.popped {
//				return 0
//			} else {
//				return 1
//			}
//		}), "bubble_vibe")
//	b.entity = myecs.Manager.NewEntity().
//		AddComponent(myecs.Entity, b).
//		AddComponent(myecs.Transform, b.Transform).
//		AddComponent(myecs.Parent, b.Dwarf.Transform).
//		AddComponent(myecs.Animation, b.Reanimator).
//		AddComponent(myecs.Drawable, b.Reanimator).
//		AddComponent(myecs.Batch, constants.ParticleKey)
//}
//
//func (b *Bubble) Delete() {
//	if !b.popped {
//		b.Pop()
//	}
//	myecs.Manager.DisposeEntity(b.entity)
//}
//
//func (b *Bubble) Pop() {
//	b.popped = true
//	b.entity.RemoveComponent(myecs.Parent)
//	b.Transform.Pos = b.Dwarf.Transform.Pos
//	b.Dwarf.Physics.Velocity = b.Physics.Velocity
//	b.Dwarf.Entity.AddComponent(myecs.Physics, b.Dwarf.Physics)
//	b.Dwarf.Bubble = nil
//}
