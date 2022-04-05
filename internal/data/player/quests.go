package player

const (
	timer = 8.
)

type Quest struct {
	Key       string
	Name      string
	Desc      string
	Check     func(*Profile) bool
	OnFinish  func(*Profile)
	Completed bool
	Hidden    bool
}

func (p *Profile) UpdateQuests() {
	for _, q := range p.Quests {
		if !q.Completed && q.Check(p) {
			q.Completed = true
			q.Hidden = false
			if q.OnFinish != nil {
				q.OnFinish(p)
			}
			// add to notifications
		}
	}
}