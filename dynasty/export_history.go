package dynasty

// buildPlayerAwardExports decodes player award rows.
func (f *File) buildPlayerAwardExports() ([]PlayerAwardExport, error) {
	awardTable, ok := f.PrimaryTableByName("PlayerAward")
	if !ok {
		return nil, nil
	}
	if err := awardTable.ReadRecords(); err != nil {
		return nil, err
	}

	exports := make([]PlayerAwardExport, 0, awardTable.ActiveRecordCount())
	for _, record := range awardTable.Records {
		if !recordIsActive(record, awardTable) {
			continue
		}
		awardType := stringField(record, "AwardType")
		if awardType == "" {
			continue
		}

		export := PlayerAwardExport{
			ID:        record.Index,
			AwardType: awardType,
			Period:    stringField(record, "Period"),
			Position:  stringField(record, "Position"),
		}
		setOptionalPositiveInt(record, "PeriodIndex", &export.PeriodIndex)
		setOptionalPositiveInt(record, "AwardScore", &export.AwardScore)

		if player, playerID, ok := f.playerRecordFromField(record, "Player"); ok {
			export.PlayerID = playerID
			export.FirstName = stringField(player, "FirstName")
			export.LastName = stringField(player, "LastName")
		}
		export.TeamName = f.teamNameFromField(record, "Team")
		if export.TeamName == "" {
			if player, _, ok := f.playerRecordFromField(record, "Player"); ok {
				if row, ok := intFieldOK(player, "TeamIndex"); ok {
					teams := f.teamMaps()
					if teamID, ok := teams.exportID(row); ok {
						export.TeamID = &teamID
						export.TeamName = teams.nameFromID(teamID)
					}
				}
			}
		} else if team, ok := f.teamRecordFromName(export.TeamName); ok {
			if teamID, ok := teamIDFromRecord(team); ok {
				export.TeamID = &teamID
			}
		}

		exports = append(exports, export)
	}
	return exports, nil
}

func (f *File) teamRecordFromName(name string) (Record, bool) {
	teamTable, ok := f.PrimaryTableByName("Team")
	if !ok {
		return Record{}, false
	}
	if err := teamTable.ReadRecords(); err != nil {
		return Record{}, false
	}
	for _, record := range teamTable.Records {
		if bestTeamName(record) == name {
			return record, true
		}
	}
	return Record{}, false
}

// buildLeagueHistoryAwardExports decodes league historical award rows.
func (f *File) buildLeagueHistoryAwardExports() ([]LeagueHistoryAwardExport, error) {
	awardTable, ok := f.PrimaryTableByName("LeagueHistoryAward")
	if !ok {
		return nil, nil
	}
	if err := awardTable.ReadRecords(); err != nil {
		return nil, err
	}

	exports := make([]LeagueHistoryAwardExport, 0, awardTable.ActiveRecordCount())
	for _, record := range awardTable.Records {
		if !recordIsActive(record, awardTable) {
			continue
		}
		export := LeagueHistoryAwardExport{
			ID:              record.Index,
			FirstName:       stringField(record, "firstName"),
			LastName:        stringField(record, "lastName"),
			Position:        stringField(record, "Position"),
			AwardType:       stringField(record, "AwardType"),
			TeamDisplayName: stringField(record, "TeamDisplayName"),
			TeamName:        f.teamNameFromField(record, "TeamIdentity"),
		}
		if export.TeamName == "" {
			export.TeamName = export.TeamDisplayName
		}
		if export.AwardType == "" && export.FirstName == "" && export.LastName == "" {
			continue
		}
		exports = append(exports, export)
	}
	return exports, nil
}

// buildConferenceChampionExports decodes conference championship history.
func (f *File) buildConferenceChampionExports() ([]ConferenceChampionExport, error) {
	champTable, ok := f.PrimaryTableByName("LeagueHistoryConferenceChampion")
	if !ok {
		return nil, nil
	}
	if err := champTable.ReadRecords(); err != nil {
		return nil, err
	}

	exports := make([]ConferenceChampionExport, 0, champTable.ActiveRecordCount())
	for _, record := range champTable.Records {
		if !recordIsActive(record, champTable) {
			continue
		}
		winner := stringField(record, "WinningTeamName")
		if winner == "" {
			winner = f.teamNameFromField(record, "WinningTeamIdentity")
		}
		loser := stringField(record, "LosingTeamName")
		if loser == "" {
			loser = f.teamNameFromField(record, "LosingTeamIdentity")
		}
		confName := stringField(record, "ConferenceName")
		if confName == "" {
			if conf, ok := record.Get("ConferenceIdentity"); ok && conf.Reference != nil {
				if confRecord, ok := f.RecordByReference("Conference", conf.Reference); ok {
					confName = stringField(confRecord, "Name")
				}
			}
		}
		if winner == "" && loser == "" && confName == "" {
			continue
		}

		export := ConferenceChampionExport{
			ID:                record.Index,
			ConferenceName:    confName,
			WinningTeamName:   winner,
			LosingTeamName:    loser,
			WinningCoachFirst: stringField(record, "WinningCoachFirstName"),
			WinningCoachLast:  stringField(record, "WinningCoachLastName"),
		}
		setOptionalPositiveInt(record, "WinningTeamScore", &export.WinningScore)
		setOptionalPositiveInt(record, "LosingTeamScore", &export.LosingScore)
		setOptionalPositiveInt(record, "WinningTeamRank", &export.WinningTeamRank)
		setOptionalPositiveInt(record, "LosingTeamRank", &export.LosingTeamRank)
		exports = append(exports, export)
	}
	return exports, nil
}
