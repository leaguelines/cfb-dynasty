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

		teamID, ok := teamIDFromRecord(record)
		if !ok {
			continue
		}

		export := TeamExport{
			ID:          teamID,
			LongName:    longName,
			DisplayName: stringField(record, "DisplayName"),
		}
		if short := stringField(record, "ShortName"); isOfficialTeamName(short) {
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
		setOptionalPositiveInt(record, "TeamPrestige", &export.TeamPrestige)

		export.OffensiveScheme = normalizeEnum(stringField(record, "CurrentOffensiveScheme"))
		export.DefensiveScheme = normalizeEnum(stringField(record, "CurrentDefensiveScheme"))
		export.Philosophy = normalizeEnum(stringField(record, "Philosophy"))
		export.PlayoffStatus = normalizeEnum(stringField(record, "PlayoffStatus"))
		export.PlayoffRoundReached = normalizeEnum(stringField(record, "PlayoffRoundReached"))
		setOptionalPositiveInt(record, "CoachesPoll_CurrentPoints", &export.CoachesPollPoints)
		setOptionalPositiveInt(record, "MediaPoll_CurrentPoints", &export.MediaPollPoints)
		if points, ok := intFieldOK(record, "CFPPoll_CurrentPoints"); ok && points > 0 {
			export.CFPPollPoints = &points
		}
		export.StaffIDs = teamStaffIDsFromRecord(record)

		exports = append(exports, export)
	}
	return exports, nil
}

func teamStaffIDsFromRecord(record Record) *TeamStaffIDsExport {
	staff := &TeamStaffIDsExport{}
	hasData := false
	setCoachID := func(field string, dst **int) {
		value, ok := record.Get(field)
		if !ok || value.Reference == nil {
			return
		}
		id := int(value.Reference.RowNumber)
		if id <= 0 {
			return
		}
		*dst = &id
		hasData = true
	}
	setCoachID("HeadCoach", &staff.HeadCoachID)
	setCoachID("OffensiveCoordinator", &staff.OffensiveCoordinatorID)
	setCoachID("DefensiveCoordinator", &staff.DefensiveCoordinatorID)
	if !hasData {
		return nil
	}
	return staff
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
		// The conference's TeamSlots reference resolves to exactly one row of the
		// Team[] array store; each populated element is a direct Team row index.
		// Row numbers map 1:1 to record indices here, so resolve exactly rather
		// than with fuzzy neighbor candidates (which would cross-assign teams).
		slotIdx := int(slots.Reference.RowNumber)
		if slotIdx < 0 || slotIdx >= len(slotTable.Records) {
			continue
		}
		for _, value := range slotTable.Records[slotIdx].Fields {
			if value.Reference == nil {
				continue
			}
			// Only real Team references count; empty array slots decode to a
			// zeroed reference (TableID 0, RowNumber 0) that must be skipped so
			// team index 0 is not spuriously assigned to a conference.
			if value.Reference.TableID != teamID {
				continue
			}
			teamIdx := int(value.Reference.RowNumber)
			if teamIdx < 0 || teamIdx >= len(teamTable.Records) {
				continue
			}
			if _, exists := out[teamIdx]; !exists {
				out[teamIdx] = confName
			}
		}
	}
	return out
}
