package player

var (
	CaveBlocksDug    int
	CaveGemsFound    int
	CaveTotalBombs   int
	CaveBlownUpBombs int
	CaveBombsLeft    int
	CaveBombsMarked  int
	CaveCorrectMarks int
	CaveWrongMarks   int

	BlocksDug         int
	GemsFound         int
	BombsFlagged      int
	CorrectMarks      int
	WrongFlags        int
)

func ResetStats() {
	BlocksDug = 0
	GemsFound = 0
	BombsFlagged = 0
	WrongFlags = 0
	ResetCaveStats()
}

func ResetCaveStats() {
	CaveBlocksDug = 0
	CaveGemsFound = 0
	CaveTotalBombs = 0
	CaveBlownUpBombs = 0
	CaveBombsLeft = 0
	CaveBombsMarked = 0
	CaveWrongMarks = 0
}

func AddStats() {
	BlocksDug += CaveBlocksDug
	GemsFound += CaveGemsFound
	BombsFlagged += CaveBombsMarked
	WrongFlags += CaveWrongMarks
}
