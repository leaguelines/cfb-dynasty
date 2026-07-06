package dynasty

import "sort"

// buildRosterExports groups rostered players by team index.
func (f *File) buildRosterExports() ([]RosterExport, error) {
	playerTable, ok := f.PrimaryTableByName("Player")
	if !ok {
		return nil, nil
	}
	if err := playerTable.ReadRecords(); err != nil {
		return nil, err
	}

	teamNames := f.teamNameByIndex()

	rosters := make(map[int][]PlayerExport)
	for _, record := range playerTable.Records {
		teamIndex, ok := intFieldOK(record, "TeamIndex")
		if !ok || teamIndex <= 0 {
			continue
		}
		player := buildPlayerExport(record)
		if player == nil {
			continue
		}
		rosters[teamIndex] = append(rosters[teamIndex], *player)
	}

	exports := make([]RosterExport, 0, len(rosters))
	for teamID, players := range rosters {
		exports = append(exports, RosterExport{
			TeamID:   teamID,
			TeamName: teamNames[teamID],
			Players:  players,
		})
	}
	sortRosterExports(exports)
	return exports, nil
}

func (f *File) teamNameByIndex() map[int]string {
	teamTable, ok := f.PrimaryTableByName("Team")
	if !ok {
		return nil
	}
	if err := teamTable.ReadRecords(); err != nil {
		return nil
	}
	names := make(map[int]string, teamTable.ActiveRecordCount())
	for _, record := range teamTable.Records {
		if name := bestTeamName(record); name != "" {
			names[record.Index] = name
		}
	}
	return names
}

func sortRosterExports(exports []RosterExport) {
	sort.Slice(exports, func(i, j int) bool {
		return exports[i].TeamID < exports[j].TeamID
	})
}
