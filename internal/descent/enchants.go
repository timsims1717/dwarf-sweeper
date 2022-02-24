package descent

import (
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/util"
)

type Enchantment struct {
	OnGain  func(*Dwarf)
	OnLose  func(*Dwarf)
	Key     string
	Title   string
	Desc    string
	Require string
}

func PickEnchantments(enchants []string) []*Enchantment {
	var result []*Enchantment
	var have []string
	list := Enchantments
	for _, i := range enchants {
		have = append(have, i)
	}
outer:
	for i := len(list) - 1; i >= 0; i-- {
		e := list[i]
		req := e.Require != ""
		for _, h := range have {
			if e.Key == h {
				list = append(list[:i], list[i+1:]...)
				continue outer
			} else if e.Require == h {
				req = false
			}
		}
		if req {
			list = append(list[:i], list[i+1:]...)
		}
	}
	choices := util.RandomSampleRange(util.Min(len(list), 3), 0, len(list), random.CaveGen)
	for _, c := range choices {
		result = append(result, list[c])
	}
	return result
}

func AddEnchantment(e1 *Enchantment, d *Dwarf) {
	for _, in := range d.Enchants {
		if in == e1.Key {
			return
		}
	}
	e1.OnGain(d)
	d.Enchants = append(d.Enchants, e1.Key)
	return
}

func RemoveEnchantment(e1 *Enchantment, d *Dwarf) {
	index := -1
	for j, in := range d.Enchants {
		if in == e1.Key {
			index = j
		}
	}
	if index == -1 {
		return
	}
	e1.OnLose(d)
	if len(d.Enchants) > 1 {
		d.Enchants = append(d.Enchants[:index], d.Enchants[index+1:]...)
	} else {
		d.Enchants = []string{}
	}
	return
}