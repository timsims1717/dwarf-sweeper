package player

type Stats struct {
	CaveBlocksDug    int
	CaveGemsFound    int
	CaveBombsFlagged int
	CaveCorrectFlags int
	CaveWrongFlags   int

	BlocksDug    int
	GemsFound    int
	BombsFlagged int
	CorrectFlags int
	WrongFlags   int
}

var (
	OverallStats   Stats
	CaveTotalBombs int
	CaveBombsLeft  int
)

func (s *Stats) ResetStats() {
	s.BlocksDug = 0
	s.GemsFound = 0
	s.BombsFlagged = 0
	s.CorrectFlags = 0
	s.WrongFlags = 0
	s.ResetCaveStats()
}

func (s *Stats) ResetCaveStats() {
	s.CaveBlocksDug = 0
	s.CaveGemsFound = 0
	s.CaveBombsFlagged = 0
	s.CaveCorrectFlags = 0
	s.CaveWrongFlags = 0
}

func (s *Stats) AddStats() {
	s.BlocksDug += s.CaveBlocksDug
	s.GemsFound += s.CaveGemsFound
	s.BombsFlagged += s.CaveBombsFlagged
	s.CorrectFlags += s.CaveCorrectFlags
	s.WrongFlags += s.CaveWrongFlags
}
