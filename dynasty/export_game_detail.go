package dynasty

import "strconv"

type gameDetailIndex struct {
	kickingByGame map[int][]Record
	olineByGame   map[int][]Record
	scoringID     uint32
}

func (f *File) buildGameDetailIndex(gameCount int) gameDetailIndex {
	idx := gameDetailIndex{
		kickingByGame: make(map[int][]Record),
		olineByGame:   make(map[int][]Record),
	}

	sg, sgOK := f.PrimaryTableByName("SeasonGame")
	if !sgOK {
		return idx
	}
	sgID := sg.Header.TableID

	if kick, ok := f.PrimaryTableByName("GameKickingStats"); ok {
		_ = kick.ReadRecords()
		for _, record := range kick.Records {
			if gameIdx, ok := gameIndexFromSeasonGameRef(record, "SeasonGame", sgID, gameCount); ok {
				idx.kickingByGame[gameIdx] = append(idx.kickingByGame[gameIdx], record)
			}
		}
	}

	if oline, ok := f.PrimaryTableByName("GameOLineStats"); ok {
		_ = oline.ReadRecords()
		for _, record := range oline.Records {
			if gameIdx, ok := gameIndexFromSeasonGameRef(record, "SeasonGame", sgID, gameCount); ok {
				idx.olineByGame[gameIdx] = append(idx.olineByGame[gameIdx], record)
			}
		}
	}

	if scoring, ok := f.PrimaryTableByName("ScoringSummary"); ok {
		idx.scoringID = scoring.Header.TableID
	}

	return idx
}

func gameIndexFromSeasonGameRef(record Record, field string, seasonGameID uint32, gameCount int) (int, bool) {
	ref, ok := record.Get(field)
	if !ok || ref.Reference == nil {
		return 0, false
	}
	if ref.Reference.TableID != 0 && ref.Reference.TableID != seasonGameID {
		return 0, false
	}
	return GameIndexFromStatReference(ref.Reference, gameCount)
}

func applyGameDetails(game *GameExport, record Record, f *File, idx gameDetailIndex) {
	if game == nil {
		return
	}

	game.HomeQuarterScores = quarterScores(record, "Home")
	game.AwayQuarterScores = quarterScores(record, "Away")
	setOptionalPositiveInt(record, "HomeScoreOT", &game.HomeScoreOT)
	setOptionalPositiveInt(record, "AwayScoreOT", &game.AwayScoreOT)

	setOptionalPositiveInt(record, "Attendance", &game.Attendance)
	setOptionalPositiveInt(record, "Temperature", &game.Temperature)
	game.Weather = normalizeEnum(stringField(record, "Weather"))
	game.Wind = normalizeEnum(stringField(record, "Wind"))
	setOptionalPositiveInt(record, "WindSpeed", &game.WindSpeed)
	game.IsSimmed = boolField(record, "IsSimmed")
	game.IsOvertime = boolField(record, "IsOvertimeGame")

	if ref, ok := record.Get("BowlGame"); ok && ref.Reference != nil {
		if bowl, ok := f.RecordByReference("BowlGame", ref.Reference); ok {
			id := bowl.Index
			game.BowlGameID = &id
			game.BowlGameName = stringField(bowl, "Name")
		}
	}

	if plays := scoringPlaysForGame(f, record, idx.scoringID); len(plays) > 0 {
		game.ScoringPlays = plays
	}
	if stats := kickingStatsForGame(idx.kickingByGame[int(record.Index)]); len(stats) > 0 {
		game.KickingStats = stats
	}
	if stats := olineStatsForGame(idx.olineByGame[int(record.Index)]); len(stats) > 0 {
		game.OLineStats = stats
	}
}

func quarterScores(record Record, side string) []int {
	var scores []int
	for q := 1; q <= 4; q++ {
		if v, ok := gameStatIntOK(record, fmtQuarterField(side, q)); ok {
			scores = append(scores, v)
		}
	}
	return scores
}

func fmtQuarterField(side string, quarter int) string {
	return side + "ScoreQuarter" + strconv.Itoa(quarter)
}

func scoringPlaysForGame(f *File, game Record, scoringTableID uint32) []ScoringPlayExport {
	value, ok := game.Get("ScoringSummaries")
	if !ok || value.Reference == nil {
		return nil
	}

	var plays []ScoringPlayExport
	for _, ref := range f.arrayStoreMemberRefs(value.Reference) {
		if scoringTableID != 0 && ref.TableID != 0 && ref.TableID != scoringTableID {
			continue
		}
		row, ok := f.RecordByReference("ScoringSummary", ref)
		if !ok {
			continue
		}
		play := ScoringPlayExport{
			Quarter:    int(row.Fields["Quarter"].Int),
			Conversion: normalizeEnum(stringField(row, "Conversion")),
		}
		if v, ok := gameStatIntOK(row, "HomeCurrentScore"); ok {
			play.HomeScore = &v
		}
		if v, ok := gameStatIntOK(row, "AwayCurrentScore"); ok {
			play.AwayScore = &v
		}
		if v, ok := gameStatIntOK(row, "HomePreviousScore"); ok {
			play.HomePreviousScore = &v
		}
		if v, ok := gameStatIntOK(row, "AwayPreviousScore"); ok {
			play.AwayPreviousScore = &v
		}
		if v, ok := gameStatIntOK(row, "TimeStampInSec"); ok {
			play.TimeStampSec = &v
		}
		plays = append(plays, play)
	}
	return plays
}

func kickingStatsForGame(records []Record) []GameKickingStatsExport {
	if len(records) == 0 {
		return nil
	}
	exports := make([]GameKickingStatsExport, 0, len(records))
	for _, record := range records {
		export := GameKickingStatsExport{}
		setGameStatInt(record, "KICKFGMADE", &export.FieldGoalsMade)
		setGameStatInt(record, "KICKFGATTEMPTS", &export.FieldGoalsAttempted)
		setGameStatInt(record, "KICKFGLONGEST", &export.FieldGoalLongest)
		setGameStatInt(record, "KICKEPMADE", &export.ExtraPointsMade)
		setGameStatInt(record, "KICKEPATTEMPTS", &export.ExtraPointsAttempted)
		setGameStatInt(record, "PUNTATTEMPTS", &export.PuntAttempts)
		setGameStatInt(record, "PUNTYARDS", &export.PuntYards)
		setGameStatInt(record, "PUNTLONGEST", &export.PuntLongest)
		setGameStatInt(record, "PUNTIN20", &export.PuntsInside20)
		if export.FieldGoalsMade == nil && export.PuntAttempts == nil {
			continue
		}
		exports = append(exports, export)
	}
	return exports
}

func olineStatsForGame(records []Record) []GameOLineStatsExport {
	if len(records) == 0 {
		return nil
	}
	exports := make([]GameOLineStatsExport, 0, len(records))
	for _, record := range records {
		export := GameOLineStatsExport{}
		setGameStatInt(record, "OLINEPANCAKES", &export.Pancakes)
		setGameStatInt(record, "OLINESACKSALLOWED", &export.SacksAllowed)
		if export.Pancakes == nil && export.SacksAllowed == nil {
			continue
		}
		exports = append(exports, export)
	}
	return exports
}

func setGameStatInt(record Record, name string, dst **int) {
	if v, ok := gameStatIntOK(record, name); ok {
		*dst = &v
	}
}
