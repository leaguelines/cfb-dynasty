package dynasty

import "sort"

// buildRosterExports groups rostered players by stable team ID (Team.TeamIndex).
// Player.TeamIndex stores a Team table row number; we remap through teamMaps.
func (f *File) buildRosterExports() ([]RosterExport, error) {
	playerTable, ok := f.PrimaryTableByName("Player")
	if !ok {
		return nil, nil
	}
	if err := playerTable.ReadRecords(); err != nil {
		return nil, err
	}

	teams := f.teamMaps()

	rosters := make(map[int][]PlayerExport)
	for _, record := range playerTable.Records {
		row, ok := intFieldOK(record, "TeamIndex")
		if !ok {
			continue
		}
		teamID, ok := teams.exportID(row)
		if !ok {
			continue
		}
		player := buildPlayerExport(record)
		if player == nil {
			continue
		}
		applyCanonicalTeamIndex(player, teams)
		rosters[teamID] = append(rosters[teamID], *player)
	}

	exports := make([]RosterExport, 0, len(rosters))
	for teamID, players := range rosters {
		exports = append(exports, RosterExport{
			TeamID:   teamID,
			TeamName: teams.nameFromID(teamID),
			Players:  players,
		})
	}
	sortRosterExports(exports)
	return exports, nil
}

func sortRosterExports(exports []RosterExport) {
	sort.Slice(exports, func(i, j int) bool {
		return exports[i].TeamID < exports[j].TeamID
	})
}
