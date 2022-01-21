package data

type Enchantment struct {
	OnGain  func()
	OnLose  func()
	Key     string
	Title   string
	Desc    string
	Require string
}
