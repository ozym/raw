package raw

import (
	"sort"
)

type records []Reading

func (r records) Len() int           { return len(r) }
func (r records) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r records) Less(i, j int) bool { return r[i].Less(r[j]) }

func Sort(readings []Reading) []Reading {
	s := append([]Reading{}, readings...)
	sort.Sort(records(s))
	return s
}

func Merge(into, from []Reading) []Reading {

	list, check := Sort(into), Sort(from)

	var overflow []Reading
	for _, v := range check {
		i := sort.Search(len(list), func(k int) bool { return list[k].Key() >= v.Key() })
		if i >= len(list) || list[i].Key() != v.Key() {
			j := sort.Search(len(overflow), func(k int) bool { return overflow[k].Key() >= v.Key() })
			if j >= len(overflow) || overflow[j].Key() != v.Key() {
				overflow = append(overflow, v)
			} else {
				overflow[j] = v
			}
		} else {
			list[i] = v
		}
	}

	return Sort(append(list, overflow...))
}
