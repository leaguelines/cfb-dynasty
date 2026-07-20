package exporter

import "sort"

// positionTabOrder is the usual depth-chart / roster tab order.
var positionTabOrder = []string{
	"QB", "HB", "FB", "WR", "TE",
	"LT", "LG", "C", "RG", "RT",
	"LE", "RE", "DT",
	"LOLB", "MLB", "ROLB",
	"CB", "FS", "SS",
	"K", "P",
}

func sortPositionKeys(seen map[string]struct{}) []string {
	order := make(map[string]int, len(positionTabOrder))
	for i, pos := range positionTabOrder {
		order[pos] = i
	}
	out := make([]string, 0, len(seen))
	for pos := range seen {
		out = append(out, pos)
	}
	sort.Slice(out, func(i, j int) bool {
		oi, oki := order[out[i]]
		oj, okj := order[out[j]]
		switch {
		case oki && okj:
			return oi < oj
		case oki:
			return true
		case okj:
			return false
		default:
			return out[i] < out[j]
		}
	})
	return out
}
