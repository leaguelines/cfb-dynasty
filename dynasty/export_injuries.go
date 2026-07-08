package dynasty

// buildInjuryExports decodes active injury rows.
func (f *File) buildInjuryExports() ([]InjuryExport, error) {
	injuryTable, ok := f.PrimaryTableByName("Injury")
	if !ok {
		return nil, nil
	}
	if err := injuryTable.ReadRecords(); err != nil {
		return nil, err
	}

	teams := f.teamMaps()

	exports := make([]InjuryExport, 0, injuryTable.ActiveRecordCount())
	for _, record := range injuryTable.Records {
		injuryType := stringField(record, "Type")
		if injuryType == "" {
			continue
		}

		export := InjuryExport{
			ID:       record.Index,
			Type:     injuryType,
			Severity: stringField(record, "Severity"),
		}
		setOptionalPositiveInt(record, "MinDuration", &export.MinWeeks)
		setOptionalPositiveInt(record, "MaxDuration", &export.MaxWeeks)

		if player, playerID, ok := f.playerRecordFromField(record, "Player"); ok {
			export.PlayerID = playerID
			export.FirstName = stringField(player, "FirstName")
			export.LastName = stringField(player, "LastName")
			if row, ok := intFieldOK(player, "TeamIndex"); ok {
				if teamID, ok := teams.exportID(row); ok {
					export.TeamID = &teamID
					export.TeamName = teams.nameFromID(teamID)
				}
			}
		}
		if export.TeamName == "" {
			if row, ok := intFieldOK(record, "GameTeam"); ok {
				if teamID, ok := teams.exportID(row); ok {
					export.TeamID = &teamID
					export.TeamName = teams.nameFromID(teamID)
				}
			}
		}

		exports = append(exports, export)
	}
	return exports, nil
}
