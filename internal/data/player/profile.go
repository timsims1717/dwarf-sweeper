package player

type Profile struct {
	Gems       int
	Flags      *Flags
	Stats      Stats
	Quests     []*Quest
	BiomeExits map[string]map[string]int
	SecretExit map[string]float64
}

func (p *Profile) AddQuest(q *Quest) {
	p.Quests = append(p.Quests, q)
	if !q.Hidden {
		// add to notifications
	}
}

type Flags struct {
	// Message Flags
	MinePuzzSeen bool
	BigBombFail  bool

	// Statistics
	BlocksDug        int
	CorrectFlags     int
	WrongFlags       int
	BombsBlown       int
	BigBombsDisarmed int

	// Explore
	Discovered []string
}

type Stats struct {
	BlocksDug        int
	CorrectFlags     int
	WrongFlags       int
	BombsBlown       int
	BigBombsDisarmed int
}

func AddStats(s1, s2 Stats) Stats {
	return Stats{
		BlocksDug:        s1.BlocksDug + s2.BlocksDug,
		CorrectFlags:     s1.CorrectFlags + s2.CorrectFlags,
		WrongFlags:       s1.WrongFlags + s2.WrongFlags,
		BombsBlown:       s1.BombsBlown + s2.BombsBlown,
		BigBombsDisarmed: s1.BigBombsDisarmed + s2.BigBombsDisarmed,
	}
}