package enchants

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/util"
)

func PickEnchantments() []*data.Enchantment {
	var result []*data.Enchantment
	var have []string
	pe := descent.Descent.GetPlayer().Enchants
	list := Enchantments
	for _, i := range pe {
		have = append(have, i)
	}
outer:
	for i := len(list)-1; i >= 0; i-- {
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

func AddEnchantment(e1 *data.Enchantment) {
	for _, in := range descent.Descent.GetPlayer().Enchants {
		if in == e1.Key {
			return
		}
	}
	e1.OnGain()
	descent.Descent.GetPlayer().Enchants = append(descent.Descent.GetPlayer().Enchants, e1.Key)
	return
}

func RemoveEnchantment(e1 *data.Enchantment) {
	index := -1
	for j, in := range descent.Descent.GetPlayer().Enchants {
		if in == e1.Key {
			index = j
		}
	}
	if index == -1 {
		return
	}
	e1.OnLose()
	if len(descent.Descent.GetPlayer().Enchants) > 1 {
		descent.Descent.GetPlayer().Enchants = append(descent.Descent.GetPlayer().Enchants[:index], descent.Descent.GetPlayer().Enchants[index+1:]...)
	} else {
		descent.Descent.GetPlayer().Enchants = []string{}
	}
	return
}