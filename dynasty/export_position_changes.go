package dynasty

import "sort"

func (f *File) buildPositionChangeExports() ([]PositionChangeExport, error) {
	table, ok := f.PrimaryTableByName("PlayerPositionChangeHistoryEntry")
	if !ok {
		return nil, nil
	}
	if err := table.ReadRecords(); err != nil {
		return nil, err
	}
	teams := f.teamMaps()
	exports := make([]PositionChangeExport, 0, table.ActiveRecordCount())
	for _, record := range table.Records {
		oldPos := normalizeEnum(stringField(record, "OldPosition"))
		newPos := normalizeEnum(stringField(record, "NewPosition"))
		oldType := normalizeEnum(stringField(record, "OldPlayerType"))
		newType := normalizeEnum(stringField(record, "NewPlayerType"))
		if oldPos == newPos && oldType == newType {
			continue
		}
		export := PositionChangeExport{
			ID:           record.Index,
			OldPosition:  oldPos,
			NewPosition:  newPos,
			OldArchetype: oldType,
			NewArchetype: newType,
			SeasonStage:  normalizeEnum(stringField(record, "SeasonStage")),
		}
		setOptionalPositiveInt(record, "SeasonYear", &export.SeasonYear)
		setOptionalPositiveInt(record, "SeasonWeek", &export.SeasonWeek)
		if player, id, ok := f.playerRecordFromField(record, "Player"); ok {
			export.PlayerID = &id
			_ = player
		}
		export.OldTeamID, _ = teamRefExport(f, record, "OldTeam", teams)
		export.NewTeamID, _ = teamRefExport(f, record, "NewTeam", teams)
		exports = append(exports, export)
	}
	sort.Slice(exports, func(i, j int) bool { return exports[i].ID < exports[j].ID })
	return exports, nil
}
