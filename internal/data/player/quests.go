package player

var Quests map[string]*Quest

type Quest struct {
	Key       string
	Name      string
	Desc      string
	Check     func(*Profile) bool
	OnFinish  func(*Profile)
	Completed bool
	Hidden    bool
}