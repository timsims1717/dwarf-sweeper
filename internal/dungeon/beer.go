package dungeon

//type Beer struct {
//	Physics   *physics.Physics
//	Transform *transform.Transform
//	created   bool
//	collect   *data.Collectible
//	sprite    *pixel.Sprite
//	entity    *ecs.Entity
//	health    *data.Health
//}
//
//func (b *Beer) Update() {
//	if b.created {
//		if b.collect.CollectedBy {
//			AddToInventory(&InvItem{
//				Name:   "beer",
//				Sprite: b.sprite,
//				OnUse:  func() bool {
//					Dungeon.Player.Entity.AddComponent(myecs.Healing, &data.Heal{
//						TmpAmount: 1,
//					})
//					return true
//				},
//				Count:  1,
//				Unique: false,
//			})
//			// todo: effects
//			b.Delete()
//		} else if b.health.Dead {
//			b.Delete()
//		}
//	}
//}
//
//func (b *Beer) Create(pos pixel.Vec) {
//	b.Physics, b.Transform = util.RandomVelocity(pos, 1.0, random.Effects)
//	b.Transform.Pos = pos
//	b.created = true
//	b.sprite = img.Batchers[entityKey].Sprites["beer"]
//	b.collect = &data.Collectible{}
//	b.health = &data.Health{
//		Max:        1,
//		Curr:       1,
//		Override:   true,
//	}
//	b.entity = myecs.Manager.NewEntity().
//		AddComponent(myecs.Entity, b).
//		AddComponent(myecs.Transform, b.Transform).
//		AddComponent(myecs.Physics, b.Physics).
//		AddComponent(myecs.Collision, data.Collider{ GroundOnly: true }).
//		AddComponent(myecs.Collect, b.collect).
//		AddComponent(myecs.Health, b.health).
//		AddComponent(myecs.Sprite, b.sprite).
//		AddComponent(myecs.Batch, entityKey)
//}
//
//func (b *Beer) Delete() {
//	b.health.Delete()
//	myecs.Manager.DisposeEntity(b.entity)
//}