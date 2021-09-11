package dungeon

//type Gem struct {
//	Physics   *physics.Physics
//	Transform *transform.Transform
//	created   bool
//	collect   *data.Collectible
//	sprite    *pixel.Sprite
//	entity    *ecs.Entity
//	health    *data.Health
//}
//
//func (g *Gem) Update() {
//	if g.created {
//		if g.collect.CollectedBy {
//			GemsFound++
//			particles.CreateRandomStaticParticles(2, 4, []string{"sparkle_1","sparkle_2","sparkle_3","sparkle_4","sparkle_5"}, g.Transform.Pos, 1.0, 1.0, 0.5)
//			sfx.SoundPlayer.PlaySound("clink", 1.0)
//			g.Delete()
//		} else if g.health.Dead {
//			g.Delete()
//		}
//	}
//}
//
//func (g *Gem) Create(pos pixel.Vec) {
//	g.Physics, g.Transform = util.RandomVelocity(pos, 1.0, random.Effects)
//	g.Transform.Pos = pos
//	g.created = true
//	g.sprite = img.Batchers[entityKey].Sprites["gem_diamond"]
//	g.collect = &data.Collectible{}
//	g.health = &data.Health{
//		Max:        1,
//		Curr:       1,
//		Override:   true,
//	}
//	g.entity = myecs.Manager.NewEntity().
//		AddComponent(myecs.Entity, g).
//		AddComponent(myecs.Transform, g.Transform).
//		AddComponent(myecs.Physics, g.Physics).
//		AddComponent(myecs.Collision, data.Collider{ GroundOnly: true }).
//		AddComponent(myecs.Collect, g.collect).
//		AddComponent(myecs.Health, g.health).
//		AddComponent(myecs.Sprite, g.sprite).
//		AddComponent(myecs.Batch, entityKey)
//}
//
//func (g *Gem) Delete() {
//	g.health.Delete()
//	myecs.Manager.DisposeEntity(g.entity)
//}