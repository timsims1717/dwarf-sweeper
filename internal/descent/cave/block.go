package cave

import (
	"bytes"
	"encoding/json"
)

type BlockType int

const (
	Unknown = iota
	Empty
	Blank
	Doorway
	Tunnel
	SecretDoor
	SecretOpen
	Pillar
	Growth
	Collapse
	Dig
	Blast
	Wall
)

var toBlockTypeID = map[string]BlockType{
	"Empty":      Empty,
	"Blank":      Blank,
	"Collapse":   Collapse,
	"Dig":        Dig,
	"Blast":      Blast,
	"Wall":       Wall,
	"Doorway":    Doorway,
	"Tunnel":     Tunnel,
	"SecretDoor": SecretDoor,
	"SecretOpen": SecretOpen,
	"Pillar":     Pillar,
	"Growth":     Growth,
}

func (t BlockType) String() string {
	switch t {
	case Empty:
		return "Empty"
	case Blank:
		return "Blank"
	case Collapse:
		return "Collapse"
	case Dig:
		return "Dig"
	case Blast:
		return "Blast"
	case Wall:
		return "Wall"
	case Doorway:
		return "Doorway"
	case Tunnel:
		return "Tunnel"
	case SecretDoor:
		return "SecretDoor"
	case SecretOpen:
		return "SecretOpen"
	case Pillar:
		return "Pillar"
	case Growth:
		return "Growth"
	default:
		return "Unknown"
	}
}

func (t BlockType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(t.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (t *BlockType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	*t = toBlockTypeID[j]
	return nil
}

func DoorType(t BlockType) BlockType {
	if t == Doorway || t == Tunnel || t == SecretDoor || t == SecretOpen {
		return t
	}
	return Doorway
}