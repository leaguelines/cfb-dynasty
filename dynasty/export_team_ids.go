package dynasty

// teamIndexMaps remaps Team table row numbers onto the stable Team.TeamIndex
// IDs that consumers expect.
//
// The Team table is not ordered by TeamIndex: newer programs sit at low row
// numbers while keeping high TeamIndex values (e.g. Appalachian State is row 3
// with TeamIndex 125). Exporting record.Index as the team id therefore
// mislabels nearly every team. Player and coach TeamIndex fields already store
// the stable Team.TeamIndex ID (Oklahoma = 69); do not run those through
// exportID, which treats its argument as a table row and misassigns players
// when row numbers collide with another team's canonical id (row 69 is Middle
// Tennessee, canonical 69 is Oklahoma). TeamIndex 255 marks FCS placeholders.
type teamIndexMaps struct {
	nameByRow map[int]string // Team table row -> display name
	nameByID  map[int]string // Team.TeamIndex -> display name
	idByRow   map[int]int    // Team table row -> Team.TeamIndex
}

const fcsTeamIndexSentinel = 255

// teamMaps builds the row/ID lookup tables used across exports. Returns a
// zero-value map that safely no-ops when the Team table is missing.
func (f *File) teamMaps() teamIndexMaps {
	m := teamIndexMaps{
		nameByRow: make(map[int]string),
		nameByID:  make(map[int]string),
		idByRow:   make(map[int]int),
	}
	teamTable, ok := f.PrimaryTableByName("Team")
	if !ok {
		return m
	}
	if err := teamTable.ReadRecords(); err != nil {
		return m
	}
	for _, record := range teamTable.Records {
		name := bestTeamName(record)
		if name == "" {
			continue
		}
		m.nameByRow[record.Index] = name

		id, ok := intFieldOK(record, "TeamIndex")
		if !ok || id == fcsTeamIndexSentinel {
			continue
		}
		m.idByRow[record.Index] = id
		m.nameByID[id] = name
	}
	return m
}

// teamNameByIndex returns Team.TeamIndex -> name. Prefer teamMaps when both the
// ID and the stored row need translating.
func (f *File) teamNameByIndex() map[int]string {
	return f.teamMaps().nameByID
}

// exportID converts a Team table row index into the stable Team.TeamIndex ID
// used in JSON exports. Use for Team record.Index lookups only — not for
// Player/Coach TeamIndex fields (see playerTeamID).
func (m teamIndexMaps) exportID(row int) (int, bool) {
	id, ok := m.idByRow[row]
	return id, ok
}

// playerTeamID resolves TeamIndex values stored on Player, Coach, and similar
// roster records. Those fields already encode the stable Team.TeamIndex ID.
func (m teamIndexMaps) playerTeamID(stored int) (int, bool) {
	if stored == fcsTeamIndexSentinel {
		return 0, false
	}
	if _, ok := m.nameByID[stored]; ok {
		return stored, true
	}
	return 0, false
}

// nameFromRow resolves a save-stored team row index to a display name.
func (m teamIndexMaps) nameFromRow(row int) string {
	return m.nameByRow[row]
}

// nameFromID resolves an exported Team.TeamIndex to a display name.
func (m teamIndexMaps) nameFromID(id int) string {
	return m.nameByID[id]
}

// teamIDFromRecord returns the stable Team.TeamIndex for a Team table row.
func teamIDFromRecord(record Record) (int, bool) {
	id, ok := intFieldOK(record, "TeamIndex")
	if !ok || id == fcsTeamIndexSentinel {
		return 0, false
	}
	return id, true
}
