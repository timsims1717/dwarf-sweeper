package builder

import (
	"bytes"
	"dwarf-sweeper/internal/descent/cave"
	"encoding/json"
)

type Base int

const (
	Roomy = iota
	Blob
	Maze
	Custom
)

type CaveBuilder struct {
	Key        string        `json:"key"`
	Biome      string        `json:"biome"`
	Title      string        `json:"title"`
	Desc       string        `json:"desc"`
	Tracks     []string      `json:"tracks"`
	Width      int           `json:"width"`
	Height     int           `json:"height"`
	Type       cave.CaveType `json:"type"`
	Base       Base          `json:"base"`
	Structures []Structure   `json:"structures"`
	Enemies    []string      `json:"enemies"`
	Exits      []string      `json:"exits"`
	ExitI      []int         `json:"-"`
}

type DigDist int

const (
	Any = iota
	Close // should be within 1.5 chunks of entrance
	Medium // should be between 1 chunk and the longest side length of the cave
	Far // should be longest side length - 1 chunk or further
)

type Structure struct {
	Key      string   `json:"key"`
	Minimum  int      `json:"minimum"`
	Maximum  int      `json:"maximum"`
	MarginL  int      `json:"marginL"`
	MarginR  int      `json:"marginR"`
	MarginT  int      `json:"marginT"`
	MarginB  int      `json:"marginB"`
	DigDist  DigDist  `json:"digDist"`
	Enemies  []string `json:"enemies"`
}

func (s *Structure) Margins() {
	if s.MarginL < 5 {
		s.MarginL = 5
	}
	if s.MarginR < 5 {
		s.MarginR = 5
	}
	if s.MarginT < 5 {
		s.MarginT = 5
	}
	if s.MarginB < 5 {
		s.MarginB = 5
	}
}

var toBaseString = map[Base]string{
	Roomy:  "Roomy",
	Blob:   "Blob",
	Maze:   "Maze",
	Custom: "Custom",
}

var toBaseID = map[string]Base{
	"Roomy":  Roomy,
	"Blob":   Blob,
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

var toDigDistString = map[DigDist]string{
	Any:    "Any",
	Close:  "Close",
	Medium: "Medium",
	Far:    "Far",
}

var toDigDistID = map[string]DigDist{
	"Any":    Any,
	"Close":  Close,
	"Medium": Medium,
	"Far":    Far,
}

func (digDist DigDist) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toDigDistString[digDist])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (digDist *DigDist) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	*digDist = toDigDistID[j]
	return nil
}

func (cb *CaveBuilder) Copy() CaveBuilder {
	newCB := CaveBuilder{
		Key:        cb.Key,
		Biome:      cb.Biome,
		Title:      cb.Title,
		Desc:       cb.Desc,
		Tracks:     cb.Tracks,
		Width:      cb.Width,
		Height:     cb.Height,
		Type:       cb.Type,
		Base:       cb.Base,
		Structures: cb.Structures,
		Exits:      cb.Exits,
	}
	return newCB
}