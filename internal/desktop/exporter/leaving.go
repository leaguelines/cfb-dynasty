package exporter

import (
	"sort"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

// LeavingTeamSummary is one team row on the leaving players index page.
type LeavingTeamSummary struct {
	ID          int
	Name           string
	SchoolLogo     string
	Conference     string
	ConferenceLogo string
	PlayerCount    int
}

// LeavingTeams returns teams with leaving players, sorted by name.
func LeavingTeams(e dynasty.Export) []LeavingTeamSummary {
	conf := teamConferences(e)
	shortNames := BuildTeamShortNameMap(e.Teams)
	counts := make(map[int]int)
	names := make(map[int]string)
	for _, p := range e.LeavingPlayers {
		if !IsExportableLeavingPlayer(&p) || p.TeamID == nil {
			continue
		}
		id := *p.TeamID
		counts[id]++
		if p.TeamName != "" {
			names[id] = p.TeamName
		}
	}
	out := make([]LeavingTeamSummary, 0, len(counts))
	for id, n := range counts {
		name := names[id]
		if name == "" {
			name = "Team " + fmtInt(id)
		}
		out = append(out, LeavingTeamSummary{
			ID:             id,
			Name:           name,
			SchoolLogo:     TeamLogoForTeamID(shortNames, id),
			Conference:     conf[id],
			ConferenceLogo: ConferenceLogoPath(conf[id]),
			PlayerCount:    n,
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

// LeavingPlayersForTeam returns leaving players for a team id.
func LeavingPlayersForTeam(e dynasty.Export, teamID int) []dynasty.LeavingPlayerExport {
	var out []dynasty.LeavingPlayerExport
	for _, p := range e.LeavingPlayers {
		if !IsExportableLeavingPlayer(&p) || p.TeamID == nil || *p.TeamID != teamID {
			continue
		}
		out = append(out, p)
	}
	sortLeavingPlayers(out)
	return out
}

// LeavingLeaveTypes returns distinct leave types for a team's leaving players.
func LeavingLeaveTypes(players []dynasty.LeavingPlayerExport) []string {
	seen := make(map[string]struct{})
	for _, p := range players {
		if !IsExportableLeavingPlayer(&p) {
			continue
		}
		t := p.LeaveType
		if t == "" {
			t = "Unknown"
		}
		seen[t] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for t := range seen {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

// LeavingHeader is the table column header for leaving player rows.
func LeavingHeader() []string {
	return []string{
		"firstName", "lastName", "position", "leaveType", "leaveStatus",
		"draftClassPosition", "projectRound", "persuadeAttempts", "playerId",
	}
}

// LeavingCollectionRows returns all leaving player rows for CSV/JSON table export.
func LeavingCollectionRows(players []dynasty.LeavingPlayerExport) [][]string {
	rows := make([][]string, 0, len(players))
	for _, p := range players {
		if !IsExportableLeavingPlayer(&p) {
			continue
		}
		rows = append(rows, []string{
			fmtInt(p.ID), fmtInt(p.PlayerID), p.FirstName, p.LastName, p.Position,
			fmtIntPtr(p.TeamID), p.TeamName, p.LeaveType, p.LeaveStatus, p.DraftClassPosition,
			fmtIntPtr(p.ProjectRound), fmtIntPtr(p.PersuadeAttempts),
		})
	}
	return rows
}

// LeavingRows returns table rows for leaving players, optionally filtered by leave type.
func LeavingRows(players []dynasty.LeavingPlayerExport, leaveType string) [][]string {
	var rows [][]string
	for _, p := range players {
		if !IsExportableLeavingPlayer(&p) {
			continue
		}
		lt := p.LeaveType
		if lt == "" {
			lt = "Unknown"
		}
		if leaveType != "" && lt != leaveType {
			continue
		}
		rows = append(rows, []string{
			p.FirstName, p.LastName, p.Position, p.LeaveType, p.LeaveStatus,
			p.DraftClassPosition, fmtIntPtr(p.ProjectRound), fmtIntPtr(p.PersuadeAttempts),
			fmtInt(p.PlayerID),
		})
	}
	return rows
}

// ValidLeavingLeaveType reports whether leaveType appears in the list.
func ValidLeavingLeaveType(players []dynasty.LeavingPlayerExport, leaveType string) bool {
	for _, p := range players {
		if !IsExportableLeavingPlayer(&p) {
			continue
		}
		lt := p.LeaveType
		if lt == "" {
			lt = "Unknown"
		}
		if lt == leaveType {
			return true
		}
	}
	return false
}

func sortLeavingPlayers(players []dynasty.LeavingPlayerExport) {
	sort.SliceStable(players, func(i, j int) bool {
		if players[i].LeaveType != players[j].LeaveType {
			return players[i].LeaveType < players[j].LeaveType
		}
		if players[i].Position != players[j].Position {
			return players[i].Position < players[j].Position
		}
		if players[i].LastName != players[j].LastName {
			return players[i].LastName < players[j].LastName
		}
		return players[i].FirstName < players[j].FirstName
	})
}

// LeavingTeamName resolves the display name for a team id from leaving player data.
func LeavingTeamName(e dynasty.Export, teamID int) string {
	for _, p := range e.LeavingPlayers {
		if p.TeamID != nil && *p.TeamID == teamID && p.TeamName != "" {
			return p.TeamName
		}
	}
	return "Team " + fmtInt(teamID)
}
