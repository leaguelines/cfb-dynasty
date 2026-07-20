package exporter

import (
	"sort"
	"strings"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

// RecordBookPeriods are the record-book periods, each shown on its own page.
var RecordBookPeriods = []string{"career", "season", "game"}

// RecordBookPeriodTitle returns a display label for a period.
func RecordBookPeriodTitle(period string) string {
	switch period {
	case "career":
		return "Career"
	case "season":
		return "Season"
	case "game":
		return "Single Game"
	default:
		return period
	}
}

// ValidRecordBookPeriod reports whether period is a known record-book period.
func ValidRecordBookPeriod(period string) bool {
	for _, p := range RecordBookPeriods {
		if p == period {
			return true
		}
	}
	return false
}

// RecordBookHeader is the column header for the per-period record book view.
var RecordBookHeader = []string{
	"statType", "rank", "statValue", "player", "position", "team", "year", "scope", "scopeName",
}

// RecordBookFilter selects a slice of the record book.
type RecordBookFilter struct {
	Period   string
	Team     string
	StatType string
}

// RecordBookStatTypeLabel returns a short display label for a stat type code.
func RecordBookStatTypeLabel(statType string) string {
	switch statType {
	case "DefensiveInts":
		return "INTs"
	case "DefensiveSacks":
		return "Sacks"
	case "PassTds":
		return "Pass TDs"
	case "PassYards":
		return "Pass Yards"
	case "ReceiveCatches":
		return "Receptions"
	case "ReceiveTDs":
		return "Receiving TDs"
	case "ReceiveYards":
		return "Receiving Yards"
	case "RushTds":
		return "Rush TDs"
	case "RushYards":
		return "Rush Yards"
	default:
		return statType
	}
}

// scopeOrder ranks scopes so league boards sort before conference and team.
var scopeOrder = map[string]int{"league": 0, "conference": 1, "team": 2}

// filterRecordBook returns entries matching the filter, sorted by stat, scope, and rank.
func filterRecordBook(e dynasty.Export, f RecordBookFilter) []dynasty.RecordBookEntry {
	var out []dynasty.RecordBookEntry
	for _, r := range e.RecordBook {
		if r.Period != f.Period {
			continue
		}
		if f.Team != "" && r.TeamName != f.Team {
			continue
		}
		if f.StatType != "" && r.StatType != f.StatType {
			continue
		}
		out = append(out, r)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].StatType != out[j].StatType {
			return out[i].StatType < out[j].StatType
		}
		if scopeOrder[out[i].Scope] != scopeOrder[out[j].Scope] {
			return scopeOrder[out[i].Scope] < scopeOrder[out[j].Scope]
		}
		if out[i].ScopeName != out[j].ScopeName {
			return out[i].ScopeName < out[j].ScopeName
		}
		return out[i].Rank < out[j].Rank
	})
	return out
}

// RecordBookRows returns the flattened, filtered rows for a period.
func RecordBookRows(e dynasty.Export, f RecordBookFilter) [][]string {
	entries := filterRecordBook(e, f)
	rows := make([][]string, 0, len(entries))
	for _, r := range entries {
		rows = append(rows, []string{
			r.StatType, fmtInt(r.Rank), fmtInt(r.StatValue),
			strings.TrimSpace(r.FirstName + " " + r.LastName), r.Position, r.TeamName,
			fmtIntPtr(r.CalendarYear), r.Scope, r.ScopeName,
		})
	}
	return rows
}

// RecordBookJSON returns the filtered entries for JSON export.
func RecordBookJSON(e dynasty.Export, f RecordBookFilter) []dynasty.RecordBookEntry {
	return filterRecordBook(e, f)
}

// RecordBookCount returns the number of entries for a period (ignoring team/stat).
func RecordBookCount(e dynasty.Export, period string) int {
	n := 0
	for _, r := range e.RecordBook {
		if r.Period == period {
			n++
		}
	}
	return n
}

// RecordBookTeams returns the distinct record-holder teams within a period,
// sorted alphabetically, for the filter dropdown.
func RecordBookTeams(e dynasty.Export, period string) []string {
	seen := make(map[string]struct{})
	for _, r := range e.RecordBook {
		if r.Period == period && r.TeamName != "" {
			seen[r.TeamName] = struct{}{}
		}
	}
	teams := make([]string, 0, len(seen))
	for t := range seen {
		teams = append(teams, t)
	}
	sort.Strings(teams)
	return teams
}

// RecordBookStatOption is one stat-type choice for the filter dropdown.
type RecordBookStatOption struct {
	Value string
	Label string
}

// RecordBookStatOptions returns stat types present in a period, with display labels.
func RecordBookStatOptions(e dynasty.Export, period string) []RecordBookStatOption {
	seen := make(map[string]struct{})
	for _, r := range e.RecordBook {
		if r.Period == period && r.StatType != "" {
			seen[r.StatType] = struct{}{}
		}
	}
	types := make([]string, 0, len(seen))
	for t := range seen {
		types = append(types, t)
	}
	sort.Strings(types)
	out := make([]RecordBookStatOption, 0, len(types))
	for _, t := range types {
		out = append(out, RecordBookStatOption{Value: t, Label: RecordBookStatTypeLabel(t)})
	}
	return out
}
