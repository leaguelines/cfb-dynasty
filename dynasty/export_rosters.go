package dynasty

import "sort"

// buildRosterExports groups rostered players by stable team ID (Team.TeamIndex).
// Player.TeamIndex already stores that ID; resolve with playerTeamID.
func (f *File) buildRosterExports() ([]RosterExport, error) {
	playerTable, ok := f.PrimaryTableByName("Player")
	if !ok {
		return nil, nil
	}
	if err := playerTable.ReadRecords(); err != nil {
		return nil, err
	}

	teams := f.teamMaps()
	alternatePositions, err := f.playerAlternatePositions()
	if err != nil {
		return nil, err
	}
	careerIdx, err := f.buildCareerStatsIndex()
	if err != nil {
		return nil, err
	}

	rosters := make(map[int][]PlayerExport)
	for _, record := range playerTable.Records {
		stored, ok := intFieldOK(record, "TeamIndex")
		if !ok {
			continue
		}
		teamID, ok := teams.playerTeamID(stored)
		if !ok {
			continue
		}
		player := f.buildPlayerExport(record, teams, careerIdx)
		if player == nil {
			continue
		}
		if alt, ok := alternatePositions[record.Index]; ok {
			applyPlayerAth(player, alt[0], alt[1])
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
