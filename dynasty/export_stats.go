package dynasty

// gameStatsIndex groups per-game stat rows by SeasonGame index.
type gameStatsIndex struct {
	offense [][]Record
	defense [][]Record
}

func (f *File) buildGameStatsIndex(gameCount int) (gameStatsIndex, error) {
	idx := gameStatsIndex{
		offense: make([][]Record, gameCount),
		defense: make([][]Record, gameCount),
	}
	if gameCount == 0 {
		return idx, nil
	}

	if off, ok := f.PrimaryTableByName("GameOffensiveStats"); ok {
		if err := off.ReadRecords(); err != nil {
			return idx, err
		}
		for _, record := range off.Records {
			ref, ok := record.Get("SeasonGame")
			if !ok || ref.Reference == nil {
				continue
			}
			gameIdx, ok := GameIndexFromStatReference(ref.Reference, gameCount)
			if !ok {
				continue
			}
			idx.offense[gameIdx] = append(idx.offense[gameIdx], record)
		}
	}

	if def, ok := f.PrimaryTableByName("GameDefensiveStats"); ok {
		if err := def.ReadRecords(); err != nil {
			return idx, err
		}
		for _, record := range def.Records {
			ref, ok := record.Get("SeasonGame")
			if !ok || ref.Reference == nil {
				continue
			}
			gameIdx, ok := GameIndexFromStatReference(ref.Reference, gameCount)
			if !ok {
				continue
			}
			idx.defense[gameIdx] = append(idx.defense[gameIdx], record)
		}
	}

	return idx, nil
}

func buildTeamStatsExport(record Record) *TeamStatsExport {
	stats := &TeamStatsExport{}
	hasData := false
	setTeamStat := func(dst **int, name string) {
		if v, ok := statIntOK(record, name); ok {
			*dst = &v
			hasData = true
		}
	}
	setTeamStat(&stats.Wins, "WINS")
	setTeamStat(&stats.Losses, "LOSSES")
	setTeamStat(&stats.TotalYards, "TOTALYARDS")
	setTeamStat(&stats.PassYards, "OFFPASSYARDS")
	setTeamStat(&stats.RushYards, "OFFRUSHYARDS")
	setTeamStat(&stats.PassAttempts, "PASSATTEMPTS")
	setTeamStat(&stats.PassCompletions, "PASSCOMPLETIONS")
	setTeamStat(&stats.PassTDs, "PASSTDS")
	setTeamStat(&stats.PassInts, "PASSINTS")
	setTeamStat(&stats.RushAttempts, "RUSHATTEMPTS")
	setTeamStat(&stats.RushTDs, "RUSHTDS")
	setTeamStat(&stats.FirstDowns, "FIRSTDOWNS")
	setTeamStat(&stats.Turnovers, "GIVEAWAYS")
	setTeamStat(&stats.Sacks, "SACKS")
	if !hasData {
		return nil
	}
	return stats
}

func buildOffensiveGameStatsExport(record Record) *OffensiveGameStatsExport {
	stats := &OffensiveGameStatsExport{}
	hasData := false
	set := func(dst **int, name string) {
		if v, ok := statIntOK(record, name); ok {
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
	if !hasData {
		return nil
	}
	return stats
}

func buildDefensiveGameStatsExport(record Record) *DefensiveGameStatsExport {
	stats := &DefensiveGameStatsExport{}
	hasData := false
	set := func(dst **int, name string) {
		if v, ok := statIntOK(record, name); ok {
			*dst = &v
			hasData = true
		}
	}
	set(&stats.GameRating, "GAMERATING")
	set(&stats.Tackles, "DEFTACKLES")
	set(&stats.AssistTackles, "ASSDEFTACKLES")
	set(&stats.Sacks, "DLINESACKS")
	set(&stats.Ints, "DSECINTS")
	set(&stats.ForcedFumbles, "DLINEFORCEDFUMBLES")
	set(&stats.FumbleRecover, "DLINEFUMBLERECOVERIES")
	set(&stats.PassDeflections, "DEFPASSDEFLECTIONS")
	if !hasData {
		return nil
	}
	return stats
}

func buildPlayerGameStatsExports(gameIdx int, idx gameStatsIndex) []PlayerGameStatsExport {
	offRows := idx.offense[gameIdx]
	defRows := idx.defense[gameIdx]
	if len(offRows) == 0 && len(defRows) == 0 {
		return nil
	}

	slotCount := len(offRows)
	if len(defRows) > slotCount {
		slotCount = len(defRows)
	}
	out := make([]PlayerGameStatsExport, 0, slotCount)
	for slot := 0; slot < slotCount; slot++ {
		entry := PlayerGameStatsExport{Slot: slot}
		if slot < len(offRows) {
			entry.Offense = buildOffensiveGameStatsExport(offRows[slot])
		}
		if slot < len(defRows) {
			entry.Defense = buildDefensiveGameStatsExport(defRows[slot])
		}
		if entry.Offense == nil && entry.Defense == nil {
			continue
		}
		out = append(out, entry)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func attachTeamGameStats(f *File, game *GameExport, record Record) {
	if cache, ok := record.Get("HomeTeamStatCache"); ok && cache.Reference != nil {
		if row, ok := f.RecordByReference("TeamStats", cache.Reference); ok {
			game.HomeTeamStats = buildTeamStatsExport(row)
		}
	}
	if cache, ok := record.Get("AwayTeamStatCache"); ok && cache.Reference != nil {
		if row, ok := f.RecordByReference("TeamStats", cache.Reference); ok {
			game.AwayTeamStats = buildTeamStatsExport(row)
		}
	}
}
