package exporter

import (
	"sort"
	"strings"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

// SeasonPlayerStatsForPlayer returns season stat lines for a player, sorted by year.
func SeasonPlayerStatsForPlayer(e dynasty.Export, playerID int) []dynasty.PlayerSeasonStatsExport {
	if playerID <= 0 {
		return nil
	}
	var out []dynasty.PlayerSeasonStatsExport
	for _, s := range e.SeasonPlayerStats {
		if s.PlayerID == playerID {
			out = append(out, s)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		yi, yj := seasonYearValue(out[i].SeasonYear), seasonYearValue(out[j].SeasonYear)
		if yi != yj {
			return yi < yj
		}
		return out[i].PlayerID < out[j].PlayerID
	})
	return out
}

// RecordBookForPlayer returns record-book entries matching a player's name.
func RecordBookForPlayer(e dynasty.Export, firstName, lastName string) []dynasty.RecordBookEntry {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	if firstName == "" && lastName == "" {
		return nil
	}
	var out []dynasty.RecordBookEntry
	for _, r := range e.RecordBook {
		if strings.EqualFold(r.FirstName, firstName) && strings.EqualFold(r.LastName, lastName) {
			out = append(out, r)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Period != out[j].Period {
			return out[i].Period < out[j].Period
		}
		if out[i].StatType != out[j].StatType {
			return out[i].StatType < out[j].StatType
		}
		return out[i].Rank < out[j].Rank
	})
	return out
}

// PlayerHasCareerStatsPage reports whether the career stats page would show data.
func PlayerHasCareerStatsPage(p *dynasty.PlayerExport, e dynasty.Export) bool {
	if p == nil {
		return false
	}
	if p.CareerStats != nil {
		if p.CareerStats.GamesPlayed != nil || p.CareerStats.Offense != nil || p.CareerStats.Defense != nil || p.CareerStats.SpecialTeams != nil {
			return true
		}
	}
	if len(SeasonPlayerStatsForPlayer(e, p.ID)) > 0 {
		return true
	}
	return len(RecordBookForPlayer(e, p.FirstName, p.LastName)) > 0
}

func seasonYearValue(year *int) int {
	if year == nil {
		return 0
	}
	return *year
}
