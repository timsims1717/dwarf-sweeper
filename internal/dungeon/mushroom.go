package dungeon

//type Mushroom struct {
//	Physics   *physics.Physics
//	Transform *transform.Transform
//	created   bool
//	collect   *data.Collectible
//	sprite    *pixel.Sprite
//	entity    *ecs.Entity
//	health    *data.Health
//}
//
//func (m *Mushroom) Update() {
//	if m.created {
//		if m.collect.CollectedBy {
//			AddToInventory(&InvItem{
//				Name:   "mushroom",
//				Sprite: m.sprite,
//				OnUse:  func() bool {
//					if Dungeon.Player.Health.Curr < Dungeon.Player.Health.Max {
//						Dungeon.Player.Entity.AddComponent(myecs.Healing, &data.Heal{
//							Amount: 1,
//						})
//						return true
//					}
//					return false
//				},
//				Count:  1,
//				Unique: false,
//			})
//			// todo: effects
//			m.Delete()
//		} else if m.health.Dead {
//			m.Delete()
//		}
//	}
//}
//
//func (m *Mushroom) Create(pos pixel.Vec) {
//	m.Physics, m.Transform = util.RandomVelocity(pos, 1.0, random.Effects)
//	m.Transform.Pos = pos
//	m.created = true
//	m.sprite = img.Batchers[entityKey].Sprites["mushroom"]
//	m.collect = &data.Collectible{}
//	m.health = &data.Health{
//		Max:        1,
//		Curr:       1,
//		Override:   true,
//	}
//	m.entity = myecs.Manager.NewEntity().
//		AddComponent(myecs.Entity, m).
//		AddComponent(myecs.Transform, m.Transform).
//		AddComponent(myecs.Physics, m.Physics).
//		AddComponent(myecs.Collision, data.Collider{ GroundOnly: true }).
//		AddComponent(myecs.Collect, m.collect).
//		AddComponent(myecs.Health, m.health).
//		AddComponent(myecs.Sprite, m.sprite).
//		AddComponent(myecs.Batch, entityKey)
//}
//
//func (m *Mushroom) Delete() {
//	m.health.Delete()
//	myecs.Manager.DisposeEntity(m.entity)
//}