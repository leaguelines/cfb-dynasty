package dynasty

// buildPlayerSeasonStatsExports assembles per-player season stat lines.
func (f *File) buildPlayerSeasonStatsExports() ([]PlayerSeasonStatsExport, error) {
	offTable, offOK := f.PrimaryTableByName("SeasonOffensiveStats")
	defTable, defOK := f.PrimaryTableByName("SeasonDefensiveStats")
	if !offOK && !defOK {
		return nil, nil
	}

	playerTable, _ := f.PrimaryTableByName("Player")
	if playerTable != nil {
		_ = playerTable.ReadRecords()
	}
	if offOK {
		_ = offTable.ReadRecords()
	}
	if defOK {
		_ = defTable.ReadRecords()
	}

	rowCount := 0
	if offOK && len(offTable.Records) > rowCount {
		rowCount = len(offTable.Records)
	}
	if defOK && len(defTable.Records) > rowCount {
		rowCount = len(defTable.Records)
	}
	if rowCount == 0 {
		return nil, nil
	}

	exports := make([]PlayerSeasonStatsExport, 0, rowCount/4)
	for idx := 0; idx < rowCount; idx++ {
		var offRecord, defRecord Record
		if offOK && idx < len(offTable.Records) {
			offRecord = offTable.Records[idx]
		}
		if defOK && idx < len(defTable.Records) {
			defRecord = defTable.Records[idx]
		}

		offense := buildSeasonOffensiveStatsExport(offRecord)
		defense := buildSeasonDefensiveStatsExport(defRecord)
		if offense == nil && defense == nil {
			continue
		}

		export := PlayerSeasonStatsExport{
			PlayerID: idx,
			Offense:  offense,
			Defense:  defense,
		}
		applySeasonPlayerMeta(&export, offRecord, defRecord)
		if playerTable != nil && idx < len(playerTable.Records) {
			applyPlayerIdentityToSeasonStats(&export, playerTable.Records[idx])
		}
		exports = append(exports, export)
	}
	return exports, nil
}

func applySeasonPlayerMeta(export *PlayerSeasonStatsExport, offRecord, defRecord Record) {
	record := offRecord
	if len(record.Fields) == 0 {
		record = defRecord
	}
	if year, ok := seasonStatIntOK(record, "SEAS_YEAR"); ok {
		export.SeasonYear = &year
	}
	if gp, ok := sensibleSeasonCount(record, "GAMESPLAYED"); ok {
		export.GamesPlayed = &gp
	}
	if gs, ok := sensibleSeasonCount(record, "GAMESSTARTED"); ok {
		export.GamesStarted = &gs
	}
}

func applyPlayerIdentityToSeasonStats(export *PlayerSeasonStatsExport, player Record) {
	export.FirstName = stringField(player, "FirstName")
	export.LastName = stringField(player, "LastName")
	if teamID, ok := intFieldOK(player, "TeamIndex"); ok && teamID > 0 {
		export.TeamID = &teamID
	}
}

func buildSeasonOffensiveStatsExport(record Record) *SeasonOffensiveStatsExport {
	if len(record.Fields) == 0 {
		return nil
	}
	stats := &SeasonOffensiveStatsExport{}
	hasData := false
	set := func(dst **int, name string) {
		if v, ok := seasonStatIntOK(record, name); ok {
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

func buildSeasonDefensiveStatsExport(record Record) *SeasonDefensiveStatsExport {
	if len(record.Fields) == 0 {
		return nil
	}
	stats := &SeasonDefensiveStatsExport{}
	hasData := false
	set := func(dst **int, name string) {
		if v, ok := seasonStatIntOK(record, name); ok {
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

// buildTeamSeasonStatsExports assembles per-team season stat lines.
func (f *File) buildTeamSeasonStatsExports() ([]TeamSeasonStatsExport, error) {
	teamTable, ok := f.PrimaryTableByName("Team")
	if !ok {
		return nil, nil
	}
	if err := teamTable.ReadRecords(); err != nil {
		return nil, err
	}

	statsTable, ok := f.PrimaryTableByName("TeamStats")
	if !ok {
		return nil, nil
	}
	if err := statsTable.ReadRecords(); err != nil {
		return nil, err
	}

	exports := make([]TeamSeasonStatsExport, 0, teamTable.ActiveRecordCount())
	for _, team := range teamTable.Records {
		longName := stringField(team, "LongName")
		if !isOfficialTeamName(longName) {
			continue
		}
		if team.Index >= len(statsTable.Records) {
			continue
		}
		stats := buildTeamSeasonStatsExport(statsTable.Records[team.Index])
		if stats == nil {
			continue
		}
		exports = append(exports, TeamSeasonStatsExport{
			TeamID:   team.Index,
			TeamName: longName,
			Stats:    stats,
		})
	}
	return exports, nil
}

func buildTeamSeasonStatsExport(record Record) *TeamStatsExport {
	stats := &TeamStatsExport{}
	hasData := false
	setSeason := func(dst **int, name string) {
		if v, ok := sensibleSeasonCount(record, name); ok {
			*dst = &v
			hasData = true
		}
	}
	setStat := func(dst **int, name string) {
		if v, ok := seasonStatIntOK(record, name); ok {
			*dst = &v
			hasData = true
		}
	}

	setSeason(&stats.Wins, "WINS")
	setSeason(&stats.Losses, "LOSSES")
	if ties, ok := sensibleSeasonCount(record, "TIES"); ok && ties > 0 {
		stats.Ties = &ties
		hasData = true
	}
	setStat(&stats.TotalYards, "OFFYARDS")
	if stats.TotalYards == nil {
		setStat(&stats.TotalYards, "TOTALYARDS")
	}
	setStat(&stats.PassYards, "OFFPASSYARDS")
	setStat(&stats.RushYards, "OFFRUSHYARDS")
	setStat(&stats.PassAttempts, "PASSATTEMPTS")
	setStat(&stats.PassCompletions, "PASSCOMPLETIONS")
	setStat(&stats.PassTDs, "PASSTDS")
	setStat(&stats.PassInts, "PASSINTS")
	setStat(&stats.RushAttempts, "RUSHATTEMPTS")
	setStat(&stats.RushTDs, "RUSHTDS")
	setStat(&stats.FirstDowns, "FIRSTDOWNS")
	setStat(&stats.Turnovers, "GIVEAWAYS")
	setStat(&stats.Giveaways, "GIVEAWAYS")
	setStat(&stats.Takeaways, "TAKEAWAYS")
	setStat(&stats.Sacks, "SACKS")
	setStat(&stats.DefPassYards, "DEFPASSYARDS")
	setStat(&stats.DefRushYards, "DEFRUSHYARDS")

	if !hasData {
		return nil
	}
	return stats
}
