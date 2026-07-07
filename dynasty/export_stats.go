package dynasty

// gameStatsIndex groups per-game stat rows by SeasonGame index and links each
// stat row back to the player that owns it.
//
// Game-stat rows (GameOffensiveStats/GameDefensiveStats) carry no player
// reference of their own — the ownership lives on the Player side. Every Player
// has a GameStats reference into a "GameStats[]" array store whose elements point
// at that player's individual game-stat rows. offOwners/defOwners invert those
// arrays so each stat row can be attributed to its player.
type gameStatsIndex struct {
	offense   [][]Record
	defense   [][]Record
	offOwners map[int]int // GameOffensiveStats row index -> Player row index
	defOwners map[int]int // GameDefensiveStats row index -> Player row index
	players   []Record    // Player rows indexed by row number
}

func (f *File) buildGameStatsIndex(gameCount int) (gameStatsIndex, error) {
	idx := gameStatsIndex{
		offense: make([][]Record, gameCount),
		defense: make([][]Record, gameCount),
	}
	if gameCount == 0 {
		return idx, nil
	}

	// Only rows that reference the current SeasonGame table belong to this
	// season. Stat rows for other seasons / cut players use a local reference
	// (TableID 0) whose row index no longer maps to a valid game, so require the
	// reference to target the SeasonGame table explicitly.
	var seasonGameID uint32
	if sg, ok := f.PrimaryTableByName("SeasonGame"); ok {
		seasonGameID = sg.Header.TableID
	}
	belongsToSeason := func(ref *RecordReference) bool {
		return ref != nil && seasonGameID != 0 && ref.TableID == seasonGameID
	}

	if off, ok := f.PrimaryTableByName("GameOffensiveStats"); ok {
		if err := off.ReadRecords(); err != nil {
			return idx, err
		}
		for _, record := range off.Records {
			ref, ok := record.Get("SeasonGame")
			if !ok || !belongsToSeason(ref.Reference) {
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
			if !ok || !belongsToSeason(ref.Reference) {
				continue
			}
			gameIdx, ok := GameIndexFromStatReference(ref.Reference, gameCount)
			if !ok {
				continue
			}
			idx.defense[gameIdx] = append(idx.defense[gameIdx], record)
		}
	}

	f.attachGameStatOwners(&idx)

	return idx, nil
}

// attachGameStatOwners walks each Player's GameStats[] array store and records,
// for every referenced game-stat row, which player owns it. Row references are
// direct record indices (RowNumber == Record.Index) into the stat tables.
func (f *File) attachGameStatOwners(idx *gameStatsIndex) {
	playerTable, ok := f.PrimaryTableByName("Player")
	if !ok || playerTable == nil {
		return
	}
	if err := playerTable.ReadRecords(); err != nil {
		return
	}
	idx.players = playerTable.Records

	off, offOK := f.PrimaryTableByName("GameOffensiveStats")
	def, defOK := f.PrimaryTableByName("GameDefensiveStats")
	if !offOK && !defOK {
		return
	}
	var offID, defID uint32
	if offOK {
		offID = off.Header.TableID
	}
	if defOK {
		defID = def.Header.TableID
	}

	idx.offOwners = make(map[int]int)
	idx.defOwners = make(map[int]int)

	for _, player := range playerTable.Records {
		gs, ok := player.Get("GameStats")
		if !ok || gs.Reference == nil {
			continue
		}
		if gs.Reference.TableID == 0 && gs.Reference.RowNumber == 0 {
			continue
		}
		arrTable, ok := f.GetTableByID(gs.Reference.TableID)
		if !ok || arrTable == nil {
			continue
		}
		if err := arrTable.ReadRecords(); err != nil {
			continue
		}
		rowIdx := int(gs.Reference.RowNumber)
		if rowIdx < 0 || rowIdx >= len(arrTable.Records) {
			continue
		}
		for _, value := range arrTable.Records[rowIdx].Fields {
			if value.Reference == nil {
				continue
			}
			row := int(value.Reference.RowNumber)
			switch {
			case offOK && value.Reference.TableID == offID:
				if row >= 0 && row < len(off.Records) {
					if _, exists := idx.offOwners[row]; !exists {
						idx.offOwners[row] = player.Index
					}
				}
			case defOK && value.Reference.TableID == defID:
				if row >= 0 && row < len(def.Records) {
					if _, exists := idx.defOwners[row]; !exists {
						idx.defOwners[row] = player.Index
					}
				}
			}
		}
	}
}

func buildTeamStatsExport(record Record) *TeamStatsExport {
	if len(record.Fields) == 0 {
		return nil
	}
	stats := &TeamStatsExport{}
	setTeamStat := func(dst **int, name string) {
		if v, ok := gameStatIntOK(record, name); ok {
			*dst = &v
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
	return stats
}

func buildOffensiveGameStatsExport(record Record) *OffensiveGameStatsExport {
	if len(record.Fields) == 0 {
		return nil
	}
	stats := &OffensiveGameStatsExport{}
	set := func(dst **int, name string) {
		if v, ok := gameStatIntOK(record, name); ok {
			*dst = &v
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
	return stats
}

func buildDefensiveGameStatsExport(record Record) *DefensiveGameStatsExport {
	if len(record.Fields) == 0 {
		return nil
	}
	stats := &DefensiveGameStatsExport{}
	set := func(dst **int, name string) {
		if v, ok := gameStatIntOK(record, name); ok {
			*dst = &v
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
	return stats
}

func buildPlayerGameStatsExports(gameIdx int, idx gameStatsIndex) []PlayerGameStatsExport {
	offRows := idx.offense[gameIdx]
	defRows := idx.defense[gameIdx]
	if len(offRows) == 0 && len(defRows) == 0 {
		return nil
	}

	out := make([]PlayerGameStatsExport, 0, len(offRows)+len(defRows))
	byPlayer := make(map[int]int, len(offRows)+len(defRows))

	// entryFor returns the index into out for a stat line. Owned lines are keyed
	// by player index so a player's offense and defense merge into one entry;
	// ownerless lines each get their own entry so distinct rows never collapse.
	entryFor := func(playerIdx int, owned bool) int {
		if owned {
			if pos, exists := byPlayer[playerIdx]; exists {
				return pos
			}
		}
		entry := PlayerGameStatsExport{}
		if owned {
			id := playerIdx
			entry.PlayerID = &id
			if playerIdx >= 0 && playerIdx < len(idx.players) {
				entry.Player = buildStatPlayerIdentity(idx.players[playerIdx])
			}
		}
		out = append(out, entry)
		pos := len(out) - 1
		if owned {
			byPlayer[playerIdx] = pos
		}
		return pos
	}

	for _, record := range offRows {
		offense := buildOffensiveGameStatsExport(record)
		if offense == nil {
			continue
		}
		playerIdx, owned := idx.offOwners[record.Index]
		out[entryFor(playerIdx, owned)].Offense = offense
	}
	for _, record := range defRows {
		defense := buildDefensiveGameStatsExport(record)
		if defense == nil {
			continue
		}
		playerIdx, owned := idx.defOwners[record.Index]
		out[entryFor(playerIdx, owned)].Defense = defense
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
