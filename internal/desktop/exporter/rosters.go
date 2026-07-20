package exporter

import (
	"sort"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

// RosterTeamSummary is one team row on the roster index page.
type RosterTeamSummary struct {
	ID             int
	Name           string
	SchoolLogo     string
	Conference     string
	ConferenceLogo string
	Ratings        TeamRatings
	PlayerCount    int
}

// RosterPositions returns distinct position codes present on the roster, in tab order.
func RosterPositions(players []dynasty.PlayerExport) []string {
	seen := make(map[string]struct{})
	for i := range players {
		p := &players[i]
		if !IsExportablePlayer(p) || p.Position == "" {
			continue
		}
		seen[p.Position] = struct{}{}
	}
	return sortPositionKeys(seen)
}

// RosterTeams returns roster summaries sorted by team name.
func RosterTeams(e dynasty.Export) []RosterTeamSummary {
	conf := teamConferences(e)
	shortNames := BuildTeamShortNameMap(e.Teams)
	out := make([]RosterTeamSummary, 0, len(e.Rosters))
	for _, r := range e.Rosters {
		if !IsExportableRosterTeam(&r) {
			continue
		}
		name := r.TeamName
		if name == "" {
			name = "Team " + fmtInt(r.TeamID)
		}
		out = append(out, RosterTeamSummary{
			ID:             r.TeamID,
			Name:           name,
			SchoolLogo:     TeamLogoForTeamID(shortNames, r.TeamID),
			Conference:     conf[r.TeamID],
			ConferenceLogo: ConferenceLogoPath(conf[r.TeamID]),
			Ratings:        TeamRatingsForTeamID(e, r.TeamID),
			PlayerCount:    len(r.Players),
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		return out[i].ID < out[j].ID
	})
	return out
}

// FindRoster returns the roster for a team id.
func FindRoster(e dynasty.Export, teamID int) (*dynasty.RosterExport, bool) {
	for i := range e.Rosters {
		if e.Rosters[i].TeamID == teamID {
			return &e.Rosters[i], true
		}
	}
	return nil, false
}

// RosterPositionHeader is the player table column header for a position tab.
func RosterPositionHeader() []string {
	return []string{
		"jersey", "firstName", "lastName", "overall", "starRating", "class",
		"archetype", "age", "height", "weight", "homeTown", "homeState", "playerId",
	}
}

// RosterPositionRows returns table rows for players at the given position.
func RosterPositionRows(players []dynasty.PlayerExport, position string) [][]string {
	var rows [][]string
	for i := range players {
		p := &players[i]
		if !IsExportablePlayer(p) || p.Position != position {
			continue
		}
		v := playerVals(p)
		rows = append(rows, []string{
			v.Jersey, v.FirstName, v.LastName, v.Overall, v.StarRating, v.SchoolYear,
			v.ArchetypeLabel, v.Age, v.Height, v.Weight, v.HomeTown, v.HomeState, v.ID,
		})
	}
	return rows
}

// ValidRosterPosition reports whether position appears on the roster.
func ValidRosterPosition(players []dynasty.PlayerExport, position string) bool {
	for i := range players {
		p := &players[i]
		if IsExportablePlayer(p) && p.Position == position {
			return true
		}
	}
	return false
}

func teamConferences(e dynasty.Export) map[int]string {
	m := make(map[int]string, len(e.Teams))
	for _, t := range e.Teams {
		m[t.ID] = t.Conference
	}
	return m
}
