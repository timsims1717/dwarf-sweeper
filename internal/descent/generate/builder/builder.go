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
	Type       cave.CaveType `json:"type"`
	Base       Base          `json:"base"`
	Structures []Structure   `json:"structures"`
	Enemies    []string      `json:"enemies"`
	Exits      []string      `json:"exits"`
	ExitI      []int         `json:"-"`
}

type Structure struct {
	Key      string   `json:"key"`
	Minimum  int      `json:"minimum"`
	Maximum  int      `json:"maximum"`
	MarginL  int      `json:"marginL"`
	MarginR  int      `json:"marginR"`
	MarginT  int      `json:"marginT"`
	MarginB  int      `json:"marginB"`
	Enemies  []string `json:"enemies"`
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

func (cb *CaveBuilder) Copy() CaveBuilder {
	newCB := CaveBuilder{
		Key:        cb.Key,
		Biome:      cb.Biome,
		Title:      cb.Title,
		Desc:       cb.Desc,
		Tracks:     cb.Tracks,
		Type:       cb.Type,
		Base:       cb.Base,
		Structures: cb.Structures,
		Exits:      cb.Exits,
	}
	return newCB
}