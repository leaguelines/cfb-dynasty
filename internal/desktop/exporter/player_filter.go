package exporter

import (
	"strings"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

// IsExportablePlayer reports whether a player row should appear in rosters and
// exports. Stale table slots often decode with identical first/last placeholder
// names; FCS pool rows are excluded upstream when schema mapping is correct.
func IsExportablePlayer(p *dynasty.PlayerExport) bool {
	if p == nil {
		return false
	}
	if p.FirstName == "" && p.LastName == "" {
		return false
	}
	// Inactive player slots reuse a single name for both fields (e.g. "Omar Omar").
	if p.FirstName != "" && p.FirstName == p.LastName {
		return false
	}
	return true
}

// IsExportableTeamName reports whether a team name should appear in team-scoped views.
func IsExportableTeamName(name string) bool {
	name = strings.TrimSpace(name)
	if name == "" {
		return false
	}
	if strings.HasPrefix(name, "FCS ") {
		return false
	}
	return true
}

// IsExportableRosterTeam reports whether a team roster should be shown. Generic
// FCS placeholder teams are hidden from the roster browser.
func IsExportableRosterTeam(r *dynasty.RosterExport) bool {
	if r == nil {
		return false
	}
	return IsExportableTeamName(r.TeamName)
}

// IsExportableLeavingPlayer reports whether a leaving-player row should be shown.
func IsExportableLeavingPlayer(p *dynasty.LeavingPlayerExport) bool {
	if p == nil {
		return false
	}
	if p.FirstName == "" && p.LastName == "" {
		return false
	}
	if p.FirstName != "" && p.FirstName == p.LastName {
		return false
	}
	if !IsExportableTeamName(p.TeamName) {
		return false
	}
	return true
}

// FilterRosters removes FCS teams and non-exportable players from every roster.
func FilterRosters(e *dynasty.Export) {
	if e == nil {
		return
	}
	out := make([]dynasty.RosterExport, 0, len(e.Rosters))
	for i := range e.Rosters {
		r := &e.Rosters[i]
		if !IsExportableRosterTeam(r) {
			continue
		}
		players := make([]dynasty.PlayerExport, 0, len(r.Players))
		for j := range r.Players {
			if IsExportablePlayer(&r.Players[j]) {
				players = append(players, r.Players[j])
			}
		}
		if len(players) == 0 {
			continue
		}
		r2 := *r
		r2.Players = players
		out = append(out, r2)
	}
	e.Rosters = out
}

// FilterLeavingPlayers removes placeholder rows and players without a real team.
func FilterLeavingPlayers(e *dynasty.Export) {
	if e == nil {
		return
	}
	out := make([]dynasty.LeavingPlayerExport, 0, len(e.LeavingPlayers))
	for i := range e.LeavingPlayers {
		if IsExportableLeavingPlayer(&e.LeavingPlayers[i]) {
			out = append(out, e.LeavingPlayers[i])
		}
	}
	e.LeavingPlayers = out
}

// PlayerLookup holds a player and where they were found.
type PlayerLookup struct {
	Player   dynasty.PlayerExport
	TeamID   int
	TeamName string
	Recruit  *dynasty.RecruitExport // set when the player is a recruit not on a roster
}

// FindPlayer locates a player id across rosters and recruits.
func FindPlayer(e dynasty.Export, playerID int) (PlayerLookup, bool) {
	for _, r := range e.Rosters {
		if !IsExportableRosterTeam(&r) {
			continue
		}
		for _, p := range r.Players {
			if p.ID == playerID && IsExportablePlayer(&p) {
				return PlayerLookup{
					Player:   p,
					TeamID:   r.TeamID,
					TeamName: r.TeamName,
				}, true
			}
		}
	}
	for i := range e.Recruits {
		rec := &e.Recruits[i]
		if rec.Player == nil || rec.Player.ID != playerID {
			continue
		}
		if !IsExportablePlayer(rec.Player) {
			continue
		}
		return PlayerLookup{
			Player: *rec.Player,
			Recruit: rec,
		}, true
	}
	return PlayerLookup{}, false
}
