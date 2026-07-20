package exporter

import (
	"sort"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

// SchoolClassView is the recruiting class page for one school.
type SchoolClassView struct {
	TeamID     int
	Name       string
	SchoolLogo string
	SchoolHref string
	JSONHref   string
	CSVHref    string
	Rows       []RecruitListRow
	Total      int
}

// SchoolClassForTeam builds the recruiting class list for a school.
func SchoolClassForTeam(e dynasty.Export, teamID int) (SchoolClassView, bool) {
	team, ok := teamByID(e.Teams, teamID)
	if !ok {
		return SchoolClassView{}, false
	}
	targets := BuildRecruitTargetByRecruitID(e)
	shortNames := BuildTeamShortNameMap(e.Teams)
	recruits := recruitsForSchool(e.Recruits, targets, teamID)
	rows := make([]RecruitListRow, 0, len(recruits))
	for i := range recruits {
		rows = append(rows, recruitListRowFromExport(&recruits[i], shortNames, targets[recruits[i].ID]))
	}
	name := team.DisplayName
	if name == "" {
		name = team.LongName
	}
	return SchoolClassView{
		TeamID:     teamID,
		Name:       name,
		SchoolLogo: TeamLogoForTeam(team),
		SchoolHref: "/schools/" + fmtInt(teamID),
		JSONHref:   "/schools/" + fmtInt(teamID) + "/class/download?format=json",
		CSVHref:    "/schools/" + fmtInt(teamID) + "/class/download?format=csv",
		Rows:       rows,
		Total:      len(rows),
	}, true
}

// SchoolClassHeader returns the CSV header for a school's recruiting class.
func SchoolClassHeader() []string {
	return []string{"rank", "position", "isAth", "name", "rating", "stage", "class", "playerId"}
}

// SchoolClassRows returns CSV rows for a school's recruiting class.
func SchoolClassRows(view SchoolClassView) [][]string {
	rows := make([][]string, 0, len(view.Rows))
	for _, r := range view.Rows {
		rating := ""
		if r.StarRating > 0 {
			rating = fmtInt(r.StarRating)
		}
		rows = append(rows, []string{
			r.NationalRank,
			r.Position,
			fmtBool(r.IsAth),
			r.DisplayName,
			rating,
			r.Stage,
			r.Class,
			fmtInt(r.ID),
		})
	}
	return rows
}

func recruitsForSchool(recruits []dynasty.RecruitExport, targets map[int]dynasty.RecruitingTargetExport, teamID int) []dynasty.RecruitExport {
	out := make([]dynasty.RecruitExport, 0)
	for i := range recruits {
		r := recruits[i]
		target := targets[r.ID]
		if RecruitBelongsToSchool(&r, &target, teamID) {
			out = append(out, r)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		return compareRecruits(out[i], out[j], RecruitSortNationalRank)
	})
	return out
}

func teamByID(teams []dynasty.TeamExport, teamID int) (dynasty.TeamExport, bool) {
	for _, t := range teams {
		if t.ID == teamID {
			return t, true
		}
	}
	return dynasty.TeamExport{}, false
}
