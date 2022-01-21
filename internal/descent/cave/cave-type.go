package cave

import (
	"bytes"
	"encoding/json"
)

type CaveType int

const (
	Normal = iota
	Minesweeper
	Infinite
)

var toString = map[CaveType]string{
	Normal:      "Normal",
	Minesweeper: "Minesweeper",
	Infinite:    "Infinite",
}

var toID = map[string]CaveType{
	"Normal":      Normal,
	"Minesweeper": Minesweeper,
	"Infinite":    Infinite,
}

func (t CaveType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toString[t])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (t *CaveType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	*t = toID[j]
	return nil
}
