package dynasty

// isNilReference reports whether a save reference is unset (TableID 0, RowNumber 0).
func isNilReference(ref *RecordReference) bool {
	return ref == nil || (ref.TableID == 0 && ref.RowNumber == 0)
}

// RecordByReference loads a row from the table implied by ref and the schema field type name.
// When TableID is 0 the game stores a row index into the named table (CFB/Madden local ref).
func (f *File) RecordByReference(tableName string, ref *RecordReference) (Record, bool) {
	if ref == nil || tableName == "" {
		return Record{}, false
	}

	var table *Table
	var ok bool
	if ref.TableID != 0 {
		table, ok = f.GetTableByID(ref.TableID)
	} else {
		table, ok = f.PrimaryTableByName(tableName)
	}
	if !ok || table == nil {
		return Record{}, false
	}
	if err := table.ReadRecords(); err != nil {
		return Record{}, false
	}

	for _, idx := range referenceRowCandidates(ref.RowNumber, len(table.Records)) {
		return table.Records[idx], true
	}
	return Record{}, false
}

func referenceRowCandidates(rowNumber uint32, recordCount int) []int {
	if recordCount == 0 {
		return nil
	}
	rn := int(rowNumber)
	candidates := []int{rn, rn - 1, rn + 1}
	seen := make(map[int]struct{}, len(candidates))
	var out []int
	for _, idx := range candidates {
		if idx < 0 || idx >= recordCount {
			continue
		}
		if _, dup := seen[idx]; dup {
			continue
		}
		seen[idx] = struct{}{}
		out = append(out, idx)
	}
	return out
}

// GameIndexFromStatReference maps a stat row's SeasonGame reference to a
// SeasonGame row index. The reference RowNumber is a direct 0-based record index
// into the SeasonGame table (RowNumber == Record.Index), so out-of-range values
// (e.g. stale references to games from other seasons) are rejected rather than
// wrapped into a valid game.
func GameIndexFromStatReference(ref *RecordReference, gameCount int) (int, bool) {
	if ref == nil || gameCount <= 0 {
		return 0, false
	}
	idx := int(ref.RowNumber)
	if idx < 0 || idx >= gameCount {
		return 0, false
	}
	return idx, true
}
