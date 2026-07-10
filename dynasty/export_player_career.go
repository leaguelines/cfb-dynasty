package dynasty

// playerCareerStatsFromRecords builds career summary and production from linked rows.
func playerCareerStatsFromRecords(offRecord, defRecord Record, special *SpecialTeamsStatsExport) *PlayerCareerStatsExport {
	stats := &PlayerCareerStatsExport{}
	hasSummary := false
	summaryRecord := offRecord
	if len(summaryRecord.Fields) == 0 {
		summaryRecord = defRecord
	}
	set := func(dst **int, name string) {
		if v, ok := careerStatIntOK(summaryRecord, name); ok {
			*dst = &v
			hasSummary = true
		}
	}
	set(&stats.GamesPlayed, "GAMESPLAYED")
	set(&stats.GamesStarted, "GAMESSTARTED")
	set(&stats.DownsPlayed, "DOWNSPLAYED")
	set(&stats.GameRating, "GAMERATING")

	if offense := buildCareerOffensiveStatsExport(offRecord); offense != nil {
		stats.Offense = offense
	}
	if defense := buildCareerDefensiveStatsExport(defRecord); defense != nil {
		stats.Defense = defense
	}
	if special != nil {
		stats.SpecialTeams = special
	}

	if !hasSummary && stats.Offense == nil && stats.Defense == nil && stats.SpecialTeams == nil {
		return nil
	}
	return stats
}

func (f *File) playerCareerStatsExport(record Record, idx careerStatsIndex) *PlayerCareerStatsExport {
	return playerCareerStatsFromIndex(record.Index, idx)
}

func buildCareerOffensiveStatsExport(record Record) *SeasonOffensiveStatsExport {
	if len(record.Fields) == 0 {
		return nil
	}
	stats := &SeasonOffensiveStatsExport{}
	hasData := false
	set := func(dst **int, name string) {
		if v, ok := careerProductionIntOK(record, name); ok {
			*dst = &v
			hasData = true
		}
	}
	set(&stats.GameRating, "GAMERATING")
	set(&stats.PassYards, "PASSYARDS")
	set(&stats.PassTDs, "PASSTDS")
	set(&stats.PassAttempts, "PASSATTEMPTS")
	set(&stats.PassCompletions, "PASSCOMPLETED")
	set(&stats.PassInts, "PASSINTS")
	set(&stats.RushYards, "RUSHYARDS")
	set(&stats.RushTDs, "RUSHTDS")
	set(&stats.RushAttempts, "RUSHATTEMPTS")
	set(&stats.RecYards, "RECEIVEYARDS")
	set(&stats.RecTDs, "RECEIVETDS")
	set(&stats.Receptions, "RECEIVECATCHES")
	set(&stats.FirstDowns, "FIRSTDOWNS")
	if !hasData {
		return nil
	}
	return stats
}

func buildCareerDefensiveStatsExport(record Record) *SeasonDefensiveStatsExport {
	if len(record.Fields) == 0 {
		return nil
	}
	stats := &SeasonDefensiveStatsExport{}
	hasData := false
	set := func(dst **int, name string) {
		if v, ok := careerProductionIntOK(record, name); ok {
			*dst = &v
			hasData = true
		}
	}
	set(&stats.GameRating, "GAMERATING")
	set(&stats.Tackles, "DEFTACKLES")
	set(&stats.AssistTackles, "ASSDEFTACKLES")
	set(&stats.Sacks, "DLINESACKS")
	set(&stats.TacklesForLoss, "DEFTACKLESFORLOSS")
	set(&stats.Ints, "DSECINTS")
	set(&stats.ForcedFumbles, "DLINEFORCEDFUMBLES")
	set(&stats.FumbleRecoveries, "DLINEFUMBLERECOVERIES")
	set(&stats.PassDeflections, "DEFPASSDEFLECTIONS")
	if !hasData {
		return nil
	}
	return stats
}

// mergeSeasonStatRecords sums counting stats across linked season rows.
func mergeSeasonStatRecords(rows []int, table *Table) Record {
	if table == nil || len(rows) == 0 {
		return Record{}
	}
	merged := Record{Fields: make(map[string]FieldValue)}
	for _, row := range rows {
		if row < 0 || row >= len(table.Records) {
			continue
		}
		for name, value := range table.Records[row].Fields {
			if value.Reference != nil || value.String != "" || value.Bool {
				continue
			}
			if value.Int == 0 {
				continue
			}
			existing := merged.Fields[name]
			existing.Int += value.Int
			merged.Fields[name] = existing
		}
	}
	return merged
}
