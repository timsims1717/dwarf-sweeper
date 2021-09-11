package random

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	Global  *rand.Rand
	CaveGen *rand.Rand
	Effects *rand.Rand
)

func init() {
	seed := time.Now().UnixNano()
	//seed := int64(1627363045028136166)
	Global = rand.New(rand.NewSource(seed))
	PrintSeed(seed, "Global")
	effSeed := Global.Int63()
	Effects = rand.New(rand.NewSource(effSeed))
	PrintSeed(effSeed, "Effects")
	caveSeed := Global.Int63()
	CaveGen = rand.New(rand.NewSource(caveSeed))
}

func PrintSeed(seed int64, s string) {
	fmt.Printf("%s Seed: %d\n", s, seed)
}

func RandGlobalSeed() {
	seed := time.Now().UnixNano()
	Global.Seed(seed)
	PrintSeed(seed, "Global")
}

func SetGlobalSeed(seed int64) {
	Global.Seed(seed)
	PrintSeed(seed, "Global")
}

func RandCaveSeed() {
	seed := Global.Int63()
	//seed := int64(4575405318719733359)
	CaveGen.Seed(seed)
	PrintSeed(seed, "CaveGen")
}

func SetCaveSeed(seed int64) {
	CaveGen.Seed(seed)
	PrintSeed(seed, "CaveGen")
}

func RandEffectsSeed() {
	seed := Global.Int63()
	Effects.Seed(seed)
	PrintSeed(seed, "Effects")
}

func SetEffectsSeed(seed int64) {
	Effects.Seed(seed)
	PrintSeed(seed, "Effects")
}