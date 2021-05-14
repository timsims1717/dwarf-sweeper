package character

import "dwarf-sweeper/pkg/img"

type Character struct {
	RagDoll    bool
	Facing     bool
	Animations map[string]*img.Instance
	currAnim   string
}