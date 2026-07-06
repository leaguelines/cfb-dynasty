package dynasty

// buildLeavingPlayerExports decodes players in the offseason exit pipeline.
func (f *File) buildLeavingPlayerExports() ([]LeavingPlayerExport, error) {
	leaveTable, ok := f.PrimaryTableByName("LeavingPlayer")
	if !ok {
		return nil, nil
	}
	if err := leaveTable.ReadRecords(); err != nil {
		return nil, err
	}

	teamNames := f.teamNameByIndex()

	exports := make([]LeavingPlayerExport, 0, leaveTable.ActiveRecordCount())
	for _, record := range leaveTable.Records {
		export := LeavingPlayerExport{
			ID:                 record.Index,
			LeaveType:          stringField(record, "LeaveType"),
			LeaveStatus:        stringField(record, "LeaveStatus"),
			DraftClassPosition: stringField(record, "DraftClassPosition"),
		}
		if export.LeaveType == "" && export.LeaveStatus == "" {
			continue
		}
		if export.LeaveStatus == "Unknown" {
			export.LeaveStatus = ""
		}

		setOptionalPositiveInt(record, "ProjectRound", &export.ProjectRound)
		setOptionalPositiveInt(record, "PersuadeAttempts", &export.PersuadeAttempts)

		if player, playerID, ok := f.playerRecordFromField(record, "Player"); ok {
			export.PlayerID = playerID
			export.FirstName = stringField(player, "FirstName")
			export.LastName = stringField(player, "LastName")
			export.Position = stringField(player, "Position")
			if teamID, ok := intFieldOK(player, "TeamIndex"); ok && teamID > 0 {
				export.TeamID = &teamID
				export.TeamName = teamNames[teamID]
			}
		}

		if export.FirstName == "" && export.LastName == "" && export.LeaveType == "" {
			continue
		}
		exports = append(exports, export)
	}
	return exports, nil
}
