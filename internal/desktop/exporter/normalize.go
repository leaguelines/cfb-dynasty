package exporter

import (
	"math"
	"sort"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

// Normalize sorts collections into a user-friendly order once, after parsing.
// Applying it at ingest keeps the table, JSON, and CSV outputs consistent.
func Normalize(e *dynasty.Export) {
	FilterRosters(e)
	FilterLeavingPlayers(e)
	sortRecruits(e.Recruits)
	for i := range e.Rosters {
		sortPlayers(e.Rosters[i].Players)
	}
}

// sortRecruits orders recruits by national rank (best first, unranked last),
// then by overall (high to low), then by name.
func sortRecruits(rs []dynasty.RecruitExport) {
	sort.SliceStable(rs, func(i, j int) bool {
		ri, rj := rankOrLast(rs[i].NationalRank), rankOrLast(rs[j].NationalRank)
		if ri != rj {
			return ri < rj
		}
		oi, oj := playerOverall(rs[i].Player), playerOverall(rs[j].Player)
		if oi != oj {
			return oi > oj
		}
		return playerLastName(rs[i].Player) < playerLastName(rs[j].Player)
	})
}

// sortPlayers orders a roster by overall (high to low), then position, then name.
func sortPlayers(ps []dynasty.PlayerExport) {
	sort.SliceStable(ps, func(i, j int) bool {
		oi, oj := overallOrLast(ps[i].Overall), overallOrLast(ps[j].Overall)
		if oi != oj {
			return oi > oj
		}
		if ps[i].Position != ps[j].Position {
			return ps[i].Position < ps[j].Position
		}
		return ps[i].LastName < ps[j].LastName
	})
}

// rankOrLast returns the rank value, treating nil/<=0 as "unranked" so those
// entries sort to the end of an ascending order.
func rankOrLast(v *int) int {
	if v == nil || *v <= 0 {
		return math.MaxInt
	}
	return *v
}

// overallOrLast returns the overall, treating nil as lowest for descending sort.
func overallOrLast(v *int) int {
	if v == nil {
		return math.MinInt
	}
	return *v
}

func playerOverall(p *dynasty.PlayerExport) int {
	if p == nil {
		return math.MinInt
	}
	return overallOrLast(p.Overall)
}

func playerLastName(p *dynasty.PlayerExport) string {
	if p == nil {
		return ""
	}
	return p.LastName
}
