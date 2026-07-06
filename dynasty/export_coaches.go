package dynasty

// buildCoachExports decodes coaching staff rows.
func (f *File) buildCoachExports() ([]CoachExport, error) {
	coachTable, ok := f.PrimaryTableByName("Coach")
	if !ok {
		return nil, nil
	}
	if err := coachTable.ReadRecords(); err != nil {
		return nil, err
	}

	teamNames := f.teamNameByIndex()

	exports := make([]CoachExport, 0, coachTable.ActiveRecordCount())
	for _, record := range coachTable.Records {
		firstName := stringField(record, "FirstName")
		lastName := stringField(record, "LastName")
		if firstName == "" && lastName == "" {
			continue
		}

		export := CoachExport{
			ID:               record.Index,
			FirstName:        firstName,
			LastName:         lastName,
			Position:         stringField(record, "Position"),
			ContractStatus:   stringField(record, "ContractStatus"),
			OffensiveScheme:  stringField(record, "OffensiveScheme"),
			DefensiveScheme:  stringField(record, "DefensiveScheme"),
			TeamPhilosophy:   stringField(record, "TeamPhilosophy"),
			HomeTown:         stringField(record, "HomeTown"),
			HomeState:        stringField(record, "HomeState"),
			IsUserControlled: boolField(record, "IsUserControlled"),
		}

		if teamID, ok := intFieldOK(record, "TeamIndex"); ok && teamID > 0 {
			export.TeamID = &teamID
			export.TeamName = teamNames[teamID]
		}

		setOptionalPositiveInt(record, "Age", &export.Age)
		setOptionalPositiveInt(record, "Level", &export.Level)
		setOptionalPositiveInt(record, "ContractSalary", &export.ContractSalary)
		setOptionalPositiveInt(record, "ContractYearsRemaining", &export.ContractYearsRemaining)
		setOptionalPositiveInt(record, "ContractLength", &export.ContractLength)
		setOptionalPositiveInt(record, "SeasonsWithTeam", &export.SeasonsWithTeam)

		export.Career = buildCoachCareerStatsExport(f, record)
		exports = append(exports, export)
	}
	return exports, nil
}

func buildCoachCareerStatsExport(f *File, coach Record) *CoachCareerStatsExport {
	if ref, ok := coach.Get("CareerStats"); ok && ref.Reference != nil {
		if row, ok := f.RecordByReference("CareerCoachStats", ref.Reference); ok {
			return coachCareerStatsFromRecord(row)
		}
	}
	return nil
}

func coachCareerStatsFromRecord(record Record) *CoachCareerStatsExport {
	stats := &CoachCareerStatsExport{}
	hasData := false
	set := func(dst **int, name string) {
		if v, ok := careerStatIntOK(record, name); ok {
			*dst = &v
			hasData = true
		}
	}

	set(&stats.Wins, "Wins")
	set(&stats.Losses, "Losses")
	set(&stats.WinsAtCurrentSchool, "WinsAtCurrentSchool")
	set(&stats.LossesAtCurrentSchool, "LossesAtCurrentSchool")
	set(&stats.BowlWins, "BowlWins")
	set(&stats.BowlLosses, "BowlLosses")
	set(&stats.NCWins, "NCWins")
	set(&stats.NCLosses, "NCLosses")
	set(&stats.PlayoffWins, "PlayoffWins")
	set(&stats.PlayoffLosses, "PlayoffLosses")
	set(&stats.ConfChampWins, "ConfChampWins")
	set(&stats.ConfChampLosses, "ConfChampLosses")
	set(&stats.Top25Wins, "Top25Wins")
	set(&stats.RivalWins, "RivalWins")

	if !hasData {
		return nil
	}
	return stats
}
