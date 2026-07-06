package dynasty

// buildTeamExports decodes Team rows into normalized exports.
func (f *File) buildTeamExports() ([]TeamExport, error) {
	teamTable, ok := f.PrimaryTableByName("Team")
	if !ok {
		return nil, nil
	}
	if err := teamTable.ReadRecords(); err != nil {
		return nil, err
	}

	conferenceByTeam := f.buildConferenceByTeamIndex()

	exports := make([]TeamExport, 0, teamTable.ActiveRecordCount())
	for _, record := range teamTable.Records {
		longName := stringField(record, "LongName")
		if !isOfficialTeamName(longName) {
			continue
		}

		export := TeamExport{
			ID:          record.Index,
			LongName:    longName,
			DisplayName: stringField(record, "DisplayName"),
		}
		if short := stringField(record, "ShortName"); isUsableTeamShortName(short, longName) {
			export.ShortName = short
		}
		if conf, ok := conferenceByTeam[record.Index]; ok {
			export.Conference = conf
		}

		if wins, ok := sumSeasonCounts(record, "HomeWin", "RoadWin"); ok {
			export.OverallWins = &wins
		}
		if losses, ok := sumSeasonCounts(record, "HomeLoss", "RoadLoss"); ok {
			export.OverallLosses = &losses
		}
		if wins, ok := sensibleSeasonCount(record, "ConfWin"); ok {
			export.ConferenceWins = &wins
		}
		if losses, ok := sensibleSeasonCount(record, "ConfLoss"); ok {
			export.ConferenceLosses = &losses
		}

		export.CoachesPollRank = sensibleRank(record, "CoachesPoll_CurrentRank")
		export.MediaPollRank = sensibleRank(record, "MediaPoll_CurrentRank")
		export.CFPPollRank = sensibleRank(record, "CFPPoll_CurrentRank")
		export.OffensiveRank = sensibleRank(record, "OffensiveRank")
		export.DefensiveRank = sensibleRank(record, "DefensiveRank")
		export.PrestigeRank = sensibleRank(record, "PrestigeRank")

		exports = append(exports, export)
	}
	return exports, nil
}

// buildConferenceByTeamIndex maps team row index to conference name when available.
func (f *File) buildConferenceByTeamIndex() map[int]string {
	confTable, ok := f.PrimaryTableByName("Conference")
	if !ok {
		return nil
	}
	if err := confTable.ReadRecords(); err != nil {
		return nil
	}

	teamTable, ok := f.PrimaryTableByName("Team")
	if !ok {
		return nil
	}
	if err := teamTable.ReadRecords(); err != nil {
		return nil
	}

	teamID := teamTable.Header.TableID
	out := make(map[int]string)
	for _, confRecord := range confTable.Records {
		confName := stringField(confRecord, "Name")
		if confName == "" {
			continue
		}
		slots, ok := confRecord.Get("TeamSlots")
		if !ok || slots.Reference == nil {
			continue
		}
		slotTable, ok := f.GetTableByID(slots.Reference.TableID)
		if !ok || slotTable == nil {
			continue
		}
		if err := slotTable.ReadRecords(); err != nil {
			continue
		}
		for _, idx := range referenceRowCandidates(slots.Reference.RowNumber, len(slotTable.Records)) {
			slotRecord := slotTable.Records[idx]
			for _, value := range slotRecord.Fields {
				if value.Reference == nil {
					continue
				}
				if value.Reference.TableID != 0 && value.Reference.TableID != teamID {
					continue
				}
				for _, teamIdx := range referenceRowCandidates(value.Reference.RowNumber, len(teamTable.Records)) {
					if _, exists := out[teamIdx]; !exists {
						out[teamIdx] = confName
					}
				}
			}
		}
	}
	return out
}
