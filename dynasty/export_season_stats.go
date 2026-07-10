package dynasty

import (
	"sort"
	"strconv"
)

// seasonStatsIndex links season stat rows back to the players that own them.
// Season stat rows carry no player reference — ownership lives on the Player
// side via the SeasonStats[] array store, same pattern as game stats.
type seasonStatsIndex struct {
	offOwners       map[int]int   // SeasonOffensiveStats row index -> Player row index
	defOwners       map[int]int   // SeasonDefensiveStats row index -> Player row index
	offRowsByPlayer map[int][]int // Player row index -> SeasonOffensiveStats row indices
	defRowsByPlayer map[int][]int // Player row index -> SeasonDefensiveStats row indices
	offSlotByRow    map[int]int   // SeasonOffensiveStats row index -> SeasonStats[] slot
	defSlotByRow    map[int]int   // SeasonDefensiveStats row index -> SeasonStats[] slot
	players         map[int]Record
	teams           teamIndexMaps
	offTable        *Table
	defTable        *Table
}

func (f *File) buildSeasonStatsIndex() (seasonStatsIndex, error) {
	idx := seasonStatsIndex{
		offOwners:       make(map[int]int),
		defOwners:       make(map[int]int),
		offRowsByPlayer: make(map[int][]int),
		defRowsByPlayer: make(map[int][]int),
		offSlotByRow:    make(map[int]int),
		defSlotByRow:    make(map[int]int),
		players:         make(map[int]Record),
		teams:           f.teamMaps(),
	}

	off, offOK := f.PrimaryTableByName("SeasonOffensiveStats")
	def, defOK := f.PrimaryTableByName("SeasonDefensiveStats")
	if !offOK && !defOK {
		return idx, nil
	}
	if offOK {
		if err := off.ReadRecords(); err != nil {
			return idx, err
		}
		idx.offTable = off
	}
	if defOK {
		if err := def.ReadRecords(); err != nil {
			return idx, err
		}
		idx.defTable = def
	}

	playerTable, ok := f.PrimaryTableByName("Player")
	if !ok || playerTable == nil {
		return idx, nil
	}
	if err := playerTable.ReadRecords(); err != nil {
		return idx, err
	}
	for _, player := range playerTable.Records {
		idx.players[player.Index] = player
	}

	f.attachSeasonStatOwners(&idx)
	return idx, nil
}

// attachSeasonStatOwners walks each Player's SeasonStats[] array store and
// records which player owns each season offensive/defensive stat row.
func (f *File) attachSeasonStatOwners(idx *seasonStatsIndex) {
	if idx == nil {
		return
	}
	var offID, defID uint32
	if idx.offTable != nil {
		offID = idx.offTable.Header.TableID
	}
	if idx.defTable != nil {
		defID = idx.defTable.Header.TableID
	}
	if offID == 0 && defID == 0 {
		return
	}

	playerTable, ok := f.PrimaryTableByName("Player")
	if !ok || playerTable == nil {
		return
	}
	if err := playerTable.ReadRecords(); err != nil {
		return
	}

	for _, player := range playerTable.Records {
		ss, ok := player.Get("SeasonStats")
		if !ok || isNilReference(ss.Reference) {
			continue
		}
		arrTable, ok := f.GetTableByID(ss.Reference.TableID)
		if !ok || arrTable == nil {
			continue
		}
		if err := arrTable.ReadRecords(); err != nil {
			continue
		}
		rowIdx := int(ss.Reference.RowNumber)
		if rowIdx < 0 || rowIdx >= len(arrTable.Records) {
			continue
		}
		for slotName, value := range arrTable.Records[rowIdx].Fields {
			if value.Reference == nil {
				continue
			}
			row := int(value.Reference.RowNumber)
			slot := seasonStatsSlotIndex(slotName)
			switch value.Reference.TableID {
			case offID:
				if idx.offTable == nil || row < 0 || row >= len(idx.offTable.Records) {
					continue
				}
				if _, exists := idx.offOwners[row]; !exists {
					idx.offOwners[row] = player.Index
				}
				idx.offRowsByPlayer[player.Index] = appendUniqueInt(idx.offRowsByPlayer[player.Index], row)
				if _, exists := idx.offSlotByRow[row]; !exists {
					idx.offSlotByRow[row] = slot
				}
			case defID:
				if idx.defTable == nil || row < 0 || row >= len(idx.defTable.Records) {
					continue
				}
				if _, exists := idx.defOwners[row]; !exists {
					idx.defOwners[row] = player.Index
				}
				idx.defRowsByPlayer[player.Index] = appendUniqueInt(idx.defRowsByPlayer[player.Index], row)
				if _, exists := idx.defSlotByRow[row]; !exists {
					idx.defSlotByRow[row] = slot
				}
			}
		}
	}
}

// buildPlayerSeasonStatsExports assembles per-player season stat lines.
func (f *File) buildPlayerSeasonStatsExports() ([]PlayerSeasonStatsExport, error) {
	idx, err := f.buildSeasonStatsIndex()
	if err != nil {
		return nil, err
	}
	if len(idx.offRowsByPlayer) == 0 && len(idx.defRowsByPlayer) == 0 {
		return nil, nil
	}

	specialByPlayer := f.buildSeasonSpecialTeamsByPlayer()

	playersWithStats := make(map[int]struct{}, len(idx.offRowsByPlayer)+len(idx.defRowsByPlayer))
	for playerIdx := range idx.offRowsByPlayer {
		playersWithStats[playerIdx] = struct{}{}
	}
	for playerIdx := range idx.defRowsByPlayer {
		playersWithStats[playerIdx] = struct{}{}
	}
	for playerIdx := range specialByPlayer {
		playersWithStats[playerIdx] = struct{}{}
	}

	exports := make([]PlayerSeasonStatsExport, 0, len(playersWithStats))
	for playerIdx := range playersWithStats {
		offRows := idx.offRowsByPlayer[playerIdx]
		defRows := idx.defRowsByPlayer[playerIdx]
		seasonKeys := seasonStatRowKeys(offRows, defRows, idx.offSlotByRow, idx.defSlotByRow)
		if len(seasonKeys) == 0 {
			seasonKeys = []seasonStatKey{{slot: -1}}
		}

		for _, key := range seasonKeys {
			var offRecord, defRecord Record
			if key.offRow >= 0 && idx.offTable != nil && key.offRow < len(idx.offTable.Records) {
				offRecord = idx.offTable.Records[key.offRow]
			}
			if key.defRow >= 0 && idx.defTable != nil && key.defRow < len(idx.defTable.Records) {
				defRecord = idx.defTable.Records[key.defRow]
			}

			offense := buildSeasonOffensiveStatsExport(offRecord)
			defense := buildSeasonDefensiveStatsExport(defRecord)
			if offense == nil && defense == nil {
				if _, ok := specialByPlayer[playerIdx]; !ok {
					continue
				}
			}

			export := PlayerSeasonStatsExport{
				PlayerID: playerIdx,
				Offense:  offense,
				Defense:  defense,
			}
			if key.slot >= 0 {
				slot := key.slot
				export.SeasonSlot = &slot
			}
			applySeasonPlayerMeta(&export, offRecord, defRecord)
			if player, ok := idx.players[playerIdx]; ok {
				applyPlayerIdentityToSeasonStats(&export, player, idx.teams)
			}
			if special, ok := specialByPlayer[playerIdx]; ok && len(seasonKeys) == 1 {
				export.SpecialTeams = special
			}
			exports = append(exports, export)
		}
	}

	sort.Slice(exports, func(i, j int) bool {
		if exports[i].PlayerID != exports[j].PlayerID {
			return exports[i].PlayerID < exports[j].PlayerID
		}
		return seasonSlotValue(exports[i].SeasonSlot) < seasonSlotValue(exports[j].SeasonSlot)
	})
	return exports, nil
}

type seasonStatKey struct {
	slot   int
	offRow int
	defRow int
}

func seasonStatRowKeys(offRows, defRows []int, offSlotByRow, defSlotByRow map[int]int) []seasonStatKey {
	bySlot := make(map[int]seasonStatKey)
	for _, row := range offRows {
		slot := offSlotByRow[row]
		key := bySlot[slot]
		key.slot = slot
		key.offRow = row
		bySlot[slot] = key
	}
	for _, row := range defRows {
		slot := defSlotByRow[row]
		key := bySlot[slot]
		key.slot = slot
		key.defRow = row
		bySlot[slot] = key
	}
	if len(bySlot) == 0 {
		return nil
	}
	slots := make([]int, 0, len(bySlot))
	for slot := range bySlot {
		slots = append(slots, slot)
	}
	sort.Ints(slots)
	out := make([]seasonStatKey, 0, len(slots))
	for _, slot := range slots {
		out = append(out, bySlot[slot])
	}
	return out
}

func seasonStatsSlotIndex(slotName string) int {
	slot, err := strconv.Atoi(slotName)
	if err != nil {
		return 0
	}
	return slot
}

func seasonSlotValue(slot *int) int {
	if slot == nil {
		return -1
	}
	return *slot
}

// buildSeasonSpecialTeamsByPlayer walks each Player's SeasonStats[] array store
// and resolves kick/punt return season totals, keyed by player row index. Return
// stats live in dedicated KPReturn tables that are not parallel-indexed with the
// Player table, so they must be linked through the season stat array.
func (f *File) buildSeasonSpecialTeamsByPlayer() map[int]*SpecialTeamsStatsExport {
	playerTable, ok := f.PrimaryTableByName("Player")
	if !ok || playerTable == nil {
		return nil
	}
	if err := playerTable.ReadRecords(); err != nil {
		return nil
	}
	kpOff, kpOffOK := f.PrimaryTableByName("SeasonOffensiveKPReturnStats")
	kpDef, kpDefOK := f.PrimaryTableByName("SeasonDefensiveKPReturnStats")
	if !kpOffOK && !kpDefOK {
		return nil
	}
	var offID, defID uint32
	if kpOffOK {
		_ = kpOff.ReadRecords()
		offID = kpOff.Header.TableID
	}
	if kpDefOK {
		_ = kpDef.ReadRecords()
		defID = kpDef.Header.TableID
	}

	out := make(map[int]*SpecialTeamsStatsExport)
	for _, player := range playerTable.Records {
		ss, ok := player.Get("SeasonStats")
		if !ok || isNilReference(ss.Reference) {
			continue
		}
		arrTable, ok := f.GetTableByID(ss.Reference.TableID)
		if !ok || arrTable == nil {
			continue
		}
		if err := arrTable.ReadRecords(); err != nil {
			continue
		}
		rowIdx := int(ss.Reference.RowNumber)
		if rowIdx < 0 || rowIdx >= len(arrTable.Records) {
			continue
		}
		for _, value := range arrTable.Records[rowIdx].Fields {
			if value.Reference == nil {
				continue
			}
			var tbl *Table
			switch {
			case kpOffOK && value.Reference.TableID == offID:
				tbl = kpOff
			case kpDefOK && value.Reference.TableID == defID:
				tbl = kpDef
			default:
				continue
			}
			row := int(value.Reference.RowNumber)
			if row < 0 || row >= len(tbl.Records) {
				continue
			}
			special := buildSpecialTeamsStatsExport(tbl.Records[row])
			if special == nil {
				continue
			}
			out[player.Index] = mergeSpecialTeams(out[player.Index], special)
		}
	}
	return out
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

func applyPlayerIdentityToSeasonStats(export *PlayerSeasonStatsExport, player Record, teams teamIndexMaps) {
	export.FirstName = stringField(player, "FirstName")
	export.LastName = stringField(player, "LastName")
	if stored, ok := intFieldOK(player, "TeamIndex"); ok {
		if teamID, ok := teams.playerTeamID(stored); ok {
			export.TeamID = &teamID
		}
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
		teamID, ok := teamIDFromRecord(team)
		if !ok {
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
			TeamID:   teamID,
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
	setStat(&stats.KickReturnYards, "KICKRETURNYARDS")
	setStat(&stats.PuntReturnYards, "PUNTRETURNYARDS")
	setStat(&stats.PuntYards, "PUNTYARDS")
	setStat(&stats.SpecialTeamYards, "SPECIALTEAMYARDS")

	if !hasData {
		return nil
	}
	return stats
}
