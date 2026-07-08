package dynasty

import "sort"

func (f *File) buildRivalryExports() ([]RivalryExport, error) {
	table, ok := f.PrimaryTableByName("Rivalry")
	if !ok {
		return nil, nil
	}
	if err := table.ReadRecords(); err != nil {
		return nil, err
	}
	teams := f.teamMaps()
	exports := make([]RivalryExport, 0, table.ActiveRecordCount())
	for _, record := range table.Records {
		export := RivalryExport{
			ID:   record.Index,
			Name: stringField(record, "Name"),
		}
		export.Team1ID, export.Team1Name = teamRefExport(f, record, "Team1", teams)
		export.Team2ID, export.Team2Name = teamRefExport(f, record, "Team2", teams)
		setOptionalPositiveInt(record, "Team1Wins", &export.Team1Wins)
		setOptionalPositiveInt(record, "Team2Wins", &export.Team2Wins)
		setOptionalPositiveInt(record, "StreakLength", &export.StreakLength)
		setOptionalPositiveInt(record, "Team1LastScore", &export.Team1LastScore)
		setOptionalPositiveInt(record, "Team2LastScore", &export.Team2LastScore)
		if streakTeam, ok := intFieldOK(record, "StreakTeam"); ok {
			if id, ok := teams.exportID(streakTeam); ok {
				export.StreakTeamID = &id
			}
		}
		if export.Name == "" && export.Team1Name == "" && export.Team2Name == "" {
			continue
		}
		exports = append(exports, export)
	}
	sort.Slice(exports, func(i, j int) bool { return exports[i].ID < exports[j].ID })
	return exports, nil
}

func teamRefExport(f *File, record Record, field string, teams teamIndexMaps) (*int, string) {
	value, ok := record.Get(field)
	if !ok || value.Reference == nil {
		if row, ok := intFieldOK(record, field); ok {
			if id, ok := teams.exportID(row); ok {
				return &id, teams.nameFromID(id)
			}
		}
		return nil, ""
	}
	if team, ok := f.RecordByReference("Team", value.Reference); ok {
		if id, ok := teamIDFromRecord(team); ok {
			return &id, bestTeamName(team)
		}
	}
	return nil, ""
}
