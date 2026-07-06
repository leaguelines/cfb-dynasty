package dynasty

// playerRecordByReference resolves a Player row from a reference without stale fallbacks.
func (f *File) playerRecordByReference(ref *RecordReference) (Record, int, bool) {
	if ref == nil || (ref.TableID == 0 && ref.RowNumber == 0) {
		return Record{}, 0, false
	}

	var table *Table
	var ok bool
	if ref.TableID != 0 {
		table, ok = f.GetTableByID(ref.TableID)
	}
	if !ok || table == nil {
		table, ok = f.PrimaryTableByName("Player")
	}
	if !ok || table == nil {
		return Record{}, 0, false
	}
	if err := table.ReadRecords(); err != nil {
		return Record{}, 0, false
	}

	var candidates []int
	if ref.TableID != 0 {
		candidates = []int{int(ref.RowNumber)}
	} else {
		candidates = referenceRowCandidates(ref.RowNumber, len(table.Records))
	}
	for _, idx := range candidates {
		if idx < 0 || idx >= len(table.Records) {
			continue
		}
		player := table.Records[idx]
		if stringField(player, "FirstName") == "" && stringField(player, "LastName") == "" {
			continue
		}
		return player, player.Index, true
	}
	return Record{}, 0, false
}

// playerRecordFromField resolves a Player row from a reference field, with row-index fallback.
func (f *File) playerRecordFromField(record Record, field string) (Record, int, bool) {
	if value, ok := record.Get(field); ok && value.Reference != nil {
		if player, id, ok := f.playerRecordByReference(value.Reference); ok {
			return player, id, true
		}
	}

	playerTable, ok := f.PrimaryTableByName("Player")
	if !ok {
		return Record{}, 0, false
	}
	if err := playerTable.ReadRecords(); err != nil {
		return Record{}, 0, false
	}
	if record.Index < 0 || record.Index >= len(playerTable.Records) {
		return Record{}, 0, false
	}
	player := playerTable.Records[record.Index]
	if stringField(player, "FirstName") == "" && stringField(player, "LastName") == "" {
		return Record{}, 0, false
	}
	return player, player.Index, true
}

func (f *File) teamNameFromField(record Record, field string) string {
	value, ok := record.Get(field)
	if !ok || value.Reference == nil {
		return ""
	}
	if team, ok := f.RecordByReference("Team", value.Reference); ok {
		return bestTeamName(team)
	}
	return ""
}
