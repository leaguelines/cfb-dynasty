package dynasty

// careerStatsIndex links career stat rows back to the players that own them.
// Career offensive rows are referenced directly from Player.CareerStats; defensive
// rows and prior-season offensive rows are reached through the SeasonStats[] /
// GameStats[] array stores used for season stats.
type careerStatsIndex struct {
	offByPlayer       map[int]int
	defByPlayer       map[int]int
	seasonOffByPlayer map[int][]int
	seasonDefByPlayer map[int][]int
	special           map[int]*SpecialTeamsStatsExport
	offTable          *Table
	defTable          *Table
	seasonOffTable    *Table
	seasonDefTable    *Table
}

func (f *File) buildCareerStatsIndex() (careerStatsIndex, error) {
	idx := careerStatsIndex{
		offByPlayer:       make(map[int]int),
		defByPlayer:       make(map[int]int),
		seasonOffByPlayer: make(map[int][]int),
		seasonDefByPlayer: make(map[int][]int),
		special:           f.buildCareerSpecialTeamsByPlayer(),
	}

	off, offOK := f.PrimaryTableByName("CareerOffensiveStats")
	def, defOK := f.PrimaryTableByName("CareerDefensiveStats")
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
	seasonOff, seasonOffOK := f.PrimaryTableByName("SeasonOffensiveStats")
	seasonDef, seasonDefOK := f.PrimaryTableByName("SeasonDefensiveStats")
	if seasonOffOK {
		if err := seasonOff.ReadRecords(); err != nil {
			return idx, err
		}
		idx.seasonOffTable = seasonOff
	}
	if seasonDefOK {
		if err := seasonDef.ReadRecords(); err != nil {
			return idx, err
		}
		idx.seasonDefTable = seasonDef
	}

	playerTable, ok := f.PrimaryTableByName("Player")
	if !ok || playerTable == nil {
		return idx, nil
	}
	if err := playerTable.ReadRecords(); err != nil {
		return idx, err
	}

	var offID, defID, seasonOffID, seasonDefID uint32
	if idx.offTable != nil {
		offID = idx.offTable.Header.TableID
	}
	if idx.defTable != nil {
		defID = idx.defTable.Header.TableID
	}
	if idx.seasonOffTable != nil {
		seasonOffID = idx.seasonOffTable.Header.TableID
	}
	if idx.seasonDefTable != nil {
		seasonDefID = idx.seasonDefTable.Header.TableID
	}

	for _, player := range playerTable.Records {
		if ref, ok := player.Get("CareerStats"); ok && !isNilReference(ref.Reference) {
			if ref.Reference.TableID == offID || ref.Reference.TableID == 0 {
				row := int(ref.Reference.RowNumber)
				if idx.offTable != nil && row >= 0 && row < len(idx.offTable.Records) {
					idx.offByPlayer[player.Index] = row
				}
			}
		}

		visit := func(ref *RecordReference) {
			if isNilReference(ref) {
				return
			}
			row := int(ref.RowNumber)
			switch ref.TableID {
			case offID:
				if idx.offTable != nil && row >= 0 && row < len(idx.offTable.Records) {
					if _, exists := idx.offByPlayer[player.Index]; !exists {
						idx.offByPlayer[player.Index] = row
					}
				}
			case defID:
				if idx.defTable != nil && row >= 0 && row < len(idx.defTable.Records) {
					if _, exists := idx.defByPlayer[player.Index]; !exists {
						idx.defByPlayer[player.Index] = row
					}
				}
			case seasonOffID:
				if idx.seasonOffTable != nil && row >= 0 && row < len(idx.seasonOffTable.Records) {
					idx.seasonOffByPlayer[player.Index] = appendUniqueInt(idx.seasonOffByPlayer[player.Index], row)
				}
			case seasonDefID:
				if idx.seasonDefTable != nil && row >= 0 && row < len(idx.seasonDefTable.Records) {
					idx.seasonDefByPlayer[player.Index] = appendUniqueInt(idx.seasonDefByPlayer[player.Index], row)
				}
			}
		}

		for _, value := range player.Fields {
			if value.Reference == nil {
				continue
			}
			if value.Reference.TableID == offID || value.Reference.TableID == defID ||
				value.Reference.TableID == seasonOffID || value.Reference.TableID == seasonDefID {
				visit(value.Reference)
				continue
			}
			arrTable, ok := f.GetTableByID(value.Reference.TableID)
			if !ok || arrTable == nil {
				continue
			}
			if err := arrTable.ReadRecords(); err != nil {
				continue
			}
			rowIdx := int(value.Reference.RowNumber)
			if rowIdx < 0 || rowIdx >= len(arrTable.Records) {
				continue
			}
			for _, slot := range arrTable.Records[rowIdx].Fields {
				if slot.Reference != nil {
					visit(slot.Reference)
				}
			}
		}
	}

	return idx, nil
}

// buildCareerSpecialTeamsByPlayer resolves career kick/punt return totals per player.
func (f *File) buildCareerSpecialTeamsByPlayer() map[int]*SpecialTeamsStatsExport {
	playerTable, ok := f.PrimaryTableByName("Player")
	if !ok || playerTable == nil {
		return nil
	}
	if err := playerTable.ReadRecords(); err != nil {
		return nil
	}
	kpOff, kpOffOK := f.PrimaryTableByName("CareerOffensiveKPReturnStats")
	kpDef, kpDefOK := f.PrimaryTableByName("CareerDefensiveKPReturnStats")
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
		visit := func(ref *RecordReference) {
			if isNilReference(ref) {
				return
			}
			var tbl *Table
			switch ref.TableID {
			case offID:
				tbl = kpOff
			case defID:
				tbl = kpDef
			default:
				return
			}
			row := int(ref.RowNumber)
			if row < 0 || row >= len(tbl.Records) {
				return
			}
			special := buildSpecialTeamsStatsExport(tbl.Records[row])
			if special == nil {
				return
			}
			out[player.Index] = mergeSpecialTeams(out[player.Index], special)
		}

		for _, value := range player.Fields {
			if value.Reference == nil {
				continue
			}
			arrTable, ok := f.GetTableByID(value.Reference.TableID)
			if !ok || arrTable == nil {
				continue
			}
			if err := arrTable.ReadRecords(); err != nil {
				continue
			}
			rowIdx := int(value.Reference.RowNumber)
			if rowIdx < 0 || rowIdx >= len(arrTable.Records) {
				continue
			}
			for _, slot := range arrTable.Records[rowIdx].Fields {
				if slot.Reference != nil {
					visit(slot.Reference)
				}
			}
		}
	}
	return out
}

func playerCareerStatsFromIndex(playerIdx int, idx careerStatsIndex) *PlayerCareerStatsExport {
	var careerOff, careerDef Record
	if row, ok := idx.offByPlayer[playerIdx]; ok && idx.offTable != nil && row < len(idx.offTable.Records) {
		careerOff = idx.offTable.Records[row]
	}
	if row, ok := idx.defByPlayer[playerIdx]; ok && idx.defTable != nil && row < len(idx.defTable.Records) {
		careerDef = idx.defTable.Records[row]
	}

	offRecord, defRecord := selectCareerStatRecords(careerOff, careerDef, idx, playerIdx)
	return playerCareerStatsFromRecords(offRecord, defRecord, idx.special[playerIdx])
}

// selectCareerStatRecords picks career offensive/defensive rows, preferring the
// CareerStats row when it already reflects every linked season and otherwise
// summing prior seasons reached through SeasonStats[] slots.
func selectCareerStatRecords(careerOff, careerDef Record, idx careerStatsIndex, playerIdx int) (Record, Record) {
	offRecord := careerOff
	if len(offRecord.Fields) == 0 || careerSeasonRowsExceed(careerOff, idx.seasonOffByPlayer[playerIdx], idx.seasonOffTable) {
		if merged := mergeSeasonStatRecords(idx.seasonOffByPlayer[playerIdx], idx.seasonOffTable); len(merged.Fields) > 0 {
			offRecord = merged
		}
	}

	defRecord := careerDef
	if len(defRecord.Fields) == 0 || careerSeasonRowsExceed(careerDef, idx.seasonDefByPlayer[playerIdx], idx.seasonDefTable) {
		if merged := mergeSeasonStatRecords(idx.seasonDefByPlayer[playerIdx], idx.seasonDefTable); len(merged.Fields) > 0 {
			defRecord = merged
		}
	}
	return offRecord, defRecord
}

func careerSeasonRowsExceed(career Record, seasonRows []int, table *Table) bool {
	if len(career.Fields) == 0 || table == nil || len(seasonRows) <= 1 {
		return false
	}
	careerGP, careerOK := intFieldOK(career, "GAMESPLAYED")
	if !careerOK {
		return false
	}
	seasonGP := 0
	for _, row := range seasonRows {
		if row < 0 || row >= len(table.Records) {
			continue
		}
		if gp, ok := intFieldOK(table.Records[row], "GAMESPLAYED"); ok {
			seasonGP += gp
		}
	}
	return seasonGP > careerGP
}

func appendUniqueInt(list []int, value int) []int {
	for _, existing := range list {
		if existing == value {
			return list
		}
	}
	return append(list, value)
}
