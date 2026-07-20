package exporter

import (
	"sort"
	"strings"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

// Recruiting sort keys.
const (
	RecruitingSortNationalRank = "nationalRank"
	RecruitingSortPositionRank = "positionRank"
	RecruitingSortStateRank    = "stateRank"
	RecruitingSortOverall      = "overall"
	RecruitingSortName         = "name"
	RecruitingSortInfluence    = "influence"
	RecruitingSortNIL          = "nilOffer"
)

// RecruitingUIColumn describes one recruiting board column.
type RecruitingUIColumn struct {
	Label   string
	SortKey string
}

// RecruitingUIColumns is the recruiting browser table layout.
var RecruitingUIColumns = []RecruitingUIColumn{
	{Label: "firstName", SortKey: RecruitingSortName},
	{Label: "lastName", SortKey: RecruitingSortName},
	{Label: "position"},
	{Label: "isAth"},
	{Label: "nationalRank", SortKey: RecruitingSortNationalRank},
	{Label: "positionRank", SortKey: RecruitingSortPositionRank},
	{Label: "stateRank", SortKey: RecruitingSortStateRank},
	{Label: "overall", SortKey: RecruitingSortOverall},
	{Label: "topSchool"},
	{Label: "influence", SortKey: RecruitingSortInfluence},
	{Label: "schools"},
	{Label: "pitches"},
	{Label: "scholarshipStatus"},
	{Label: "nilOffer", SortKey: RecruitingSortNIL},
	{Label: "favorite"},
	{Label: "playerId"},
}

// RecruitingEntry joins pursuit state with recruit identity.
type RecruitingEntry struct {
	Target dynasty.RecruitingTargetExport
	Recruit dynasty.RecruitExport
	HasRecruit bool
}

// BuildRecruitTargetByRecruitID indexes recruiting pursuit rows by recruit ID.
func BuildRecruitTargetByRecruitID(e dynasty.Export) map[int]dynasty.RecruitingTargetExport {
	m := make(map[int]dynasty.RecruitingTargetExport, len(e.Recruiting))
	for _, t := range e.Recruiting {
		m[t.RecruitID] = t
	}
	return m
}

// RecruitingEntries builds joined recruiting rows for the UI.
func RecruitingEntries(e dynasty.Export) []RecruitingEntry {
	byID := recruitByID(e)
	out := make([]RecruitingEntry, 0, len(e.Recruiting))
	for _, t := range e.Recruiting {
		entry := RecruitingEntry{Target: t}
		if r, ok := byID[t.RecruitID]; ok {
			entry.Recruit = r
			entry.HasRecruit = true
		}
		out = append(out, entry)
	}
	return out
}

// RecruitingUIHeader returns flat header labels.
func RecruitingUIHeader() []string {
	h := make([]string, len(RecruitingUIColumns))
	for i, c := range RecruitingUIColumns {
		h[i] = c.Label
	}
	return h
}

// RecruitingUIRows builds table rows for the recruiting browser.
func RecruitingUIRows(entries []RecruitingEntry) [][]string {
	rows := make([][]string, 0, len(entries))
	for _, e := range entries {
		p := playerVals(e.Recruit.Player)
		var teamName, influence string
		if e.Target.TopSchool != nil {
			teamName = e.Target.TopSchool.TeamName
			influence = fmtInt(e.Target.TopSchool.Influence)
		}
		rows = append(rows, []string{
			p.FirstName, p.LastName, p.Position, p.IsAth,
			fmtIntPtr(e.Recruit.NationalRank), fmtIntPtr(e.Recruit.PositionRank), fmtIntPtr(e.Recruit.StateRank),
			p.Overall, teamName, influence,
			fmtInt(len(e.Target.SchoolInterest)), fmtInt(len(e.Target.ActivePitches)),
			e.Target.ScholarshipStatus, fmtIntPtr(e.Target.CurrentNILOffer), fmtBool(e.Target.IsFavorite),
			p.ID,
		})
	}
	return rows
}

// SortRecruitingEntries sorts entries in place.
func SortRecruitingEntries(entries []RecruitingEntry, key string, desc bool) {
	if key == "" {
		key = RecruitingSortNationalRank
	}
	sort.SliceStable(entries, func(i, j int) bool {
		less := compareRecruiting(entries[i], entries[j], key)
		if desc {
			return !less
		}
		return less
	})
}

// RecruitingSortDesc reports default descending for a recruiting sort key.
func RecruitingSortDesc(key string) bool {
	switch key {
	case RecruitingSortOverall, RecruitingSortInfluence, RecruitingSortNIL:
		return true
	default:
		return false
	}
}

func recruitByID(e dynasty.Export) map[int]dynasty.RecruitExport {
	m := make(map[int]dynasty.RecruitExport, len(e.Recruits))
	for _, r := range e.Recruits {
		m[r.ID] = r
	}
	return m
}

func compareRecruiting(a, b RecruitingEntry, key string) bool {
	switch key {
	case RecruitingSortNationalRank:
		return rankOrLast(a.Recruit.NationalRank) < rankOrLast(b.Recruit.NationalRank)
	case RecruitingSortPositionRank:
		return rankOrLast(a.Recruit.PositionRank) < rankOrLast(b.Recruit.PositionRank)
	case RecruitingSortStateRank:
		return rankOrLast(a.Recruit.StateRank) < rankOrLast(b.Recruit.StateRank)
	case RecruitingSortOverall:
		return playerOverall(a.Recruit.Player) < playerOverall(b.Recruit.Player)
	case RecruitingSortName:
		an := strings.ToLower(recruitName(a.Recruit))
		bn := strings.ToLower(recruitName(b.Recruit))
		if an != bn {
			return an < bn
		}
		return strings.ToLower(recruitNameLast(a.Recruit)) < strings.ToLower(recruitNameLast(b.Recruit))
	case RecruitingSortInfluence:
		return recruitingInfluence(a) < recruitingInfluence(b)
	case RecruitingSortNIL:
		return intOrLast(a.Target.CurrentNILOffer) < intOrLast(b.Target.CurrentNILOffer)
	default:
		return rankOrLast(a.Recruit.NationalRank) < rankOrLast(b.Recruit.NationalRank)
	}
}

func recruitingInfluence(e RecruitingEntry) int {
	if e.Target.TopSchool == nil {
		return -1
	}
	return e.Target.TopSchool.Influence
}
