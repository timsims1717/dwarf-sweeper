package data

import (
	"dwarf-sweeper/pkg/timing"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

type SimpleHealth struct {
	Dead   bool
	DigMe  bool
	Immune map[DamageType]Immunity
}

type Health struct {
	Max  int  // maximum HP
	Curr int  // current HP
	Dead bool // convenience boolean (set if curr == 0)
	Inv  bool // invulnerability override (debugging, etc)
	// entity can have temporary hit points
	TempHP      int           // the current amount of temp HP
	TempHPTimer *timing.Timer // the timer
	// entity can be dazed by attacks
	Dazed       bool          // convenience boolean
	DazedTime   float64       // how long will the entity be dazed? If 0., use the dazed value of the attack
	DazedTimer  *timing.Timer // the timer
	DazedEntity *ecs.Entity   // the vfx entity
	// entity can be invulnerable after receiving damage
	TempInvTimer *timing.Timer // the timer
	TempInvSec   float64       // how long (0. would mean no invulnerable frames)
	// entity can be immune to different types of damage
	Immune map[DamageType]Immunity // which damage types is the entity immune to?
}

type DamageType int

const (
	Blast = iota
	Shovel
	Enemy
	Projectile
)

type Damage struct {
	SourceID  string
	Amount    int
	Dazed     float64
	Knockback float64
	Angle     *float64
	Source    pixel.Vec
	Type      DamageType
}

type AreaDamage struct {
	SourceID       string
	Amount         int
	Dazed          float64
	Knockback      float64
	Angle          *float64
	Type           DamageType
	Center         pixel.Vec
	Radius         float64
	Rect           pixel.Rect
	KnockbackDecay bool
}

type Heal struct {
	Amount    int
	TmpAmount int
}

type Immunity struct {
	KB    bool
	DMG   bool
	Dazed bool
}

var (
	NoImmunity = map[DamageType]Immunity{
		Blast:  {},
		Shovel: {},
		Enemy:  {},
	}
	FullImmunity = map[DamageType]Immunity{
		Blast: {
			KB:    true,
			DMG:   true,
			Dazed: true,
		},
		Shovel: {
			KB:    true,
			DMG:   true,
			Dazed: true,
		},
		Enemy: {
			KB:    true,
			DMG:   true,
			Dazed: true,
		},
		Projectile: {
			KB:    true,
			DMG:   true,
			Dazed: true,
		},
	}
	ItemImmunity1 = map[DamageType]Immunity{
		Enemy: {
			KB:    true,
			DMG:   true,
			Dazed: true,
		},
		Shovel: {
			DMG:   true,
			Dazed: true,
		},
	}
	ItemImmunity2 = map[DamageType]Immunity{
		Enemy: {
			KB:    true,
			DMG:   true,
			Dazed: true,
		},
		Shovel: {
			KB:    true,
			DMG:   true,
			Dazed: true,
		},
	}
	EnemyImmunity = map[DamageType]Immunity{
		Enemy: {
			KB:    true,
			DMG:   true,
			Dazed: true,
		},
	}
	ShovelImmunity = map[DamageType]Immunity{
		Shovel: {
			KB:    true,
			DMG:   true,
			Dazed: true,
		},
	}
	UndergroundImmunity = map[DamageType]Immunity{
		Enemy: {
			KB:    true,
			DMG:   true,
			Dazed: true,
		},
		Blast: {
			KB: true,
		},
		Shovel: {
			KB:    true,
			DMG:   true,
			Dazed: true,
		},
	}
)
