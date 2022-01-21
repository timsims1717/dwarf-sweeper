package builder

import (
	"bytes"
	"dwarf-sweeper/internal/descent/cave"
	"encoding/json"
)

type Base int

const (
	Roomy = iota
	Maze
	Custom
)

type Seed int

const (
	Path = iota
	Marked
	DeadEnd
	Room
	Random
)

type CaveBuilder struct {
	Key        string        `json:"key"`
	Biome      string        `json:"biome"`
	Title      string        `json:"title"`
	Desc       string        `json:"desc"`
	Tracks     []string      `json:"tracks"`
	Type       cave.CaveType `json:"type"`
	Base       Base          `json:"base"`
	Structures []Structure   `json:"structures"`
}

type Structure struct {
	Key      string   `json:"key"`
	Seed     Seed     `json:"seed"`
	Minimum  int      `json:"minimum"`
	Maximum  int      `json:"maximum"`
	MinMult  float64  `json:"minMult"`
	RandMult float64  `json:"randMult"`
	Enemies  []string `json:"enemies"`
}

var toBaseString = map[Base]string{
	Roomy:  "Roomy",
	Maze:   "Maze",
	Custom: "Custom",
}

var toBaseID = map[string]Base{
	"Roomy":  Roomy,
	"Maze":   Maze,
	"Custom": Custom,
}

func (base Base) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toBaseString[base])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (base *Base) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	*base = toBaseID[j]
	return nil
}

var toSeedString = map[Seed]string{
	Path:    "Path",
	Marked:  "Marked",
	DeadEnd: "DeadEnd",
	Room:    "Room",
	Random:  "Random",
}

var toSeedID = map[string]Seed{
	"Path":    Path,
	"Marked":  Marked,
	"DeadEnd": DeadEnd,
	"Room":    Room,
	"Random":  Random,
}

func (seed Seed) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toSeedString[seed])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (seed *Seed) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	*seed = toSeedID[j]
	return nil
}
