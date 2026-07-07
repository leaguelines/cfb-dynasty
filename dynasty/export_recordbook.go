package dynasty

// recordBookCategories are the stat-category fields on a PlayerStatRecords
// struct. Each points at the record-holder row (rank 1) for that category; a
// league board continues with ranks 2..N in the following rows.
var recordBookCategories = []string{
	"PassYards", "PassTDS", "RushYards", "RushTDS",
	"ReceivingYards", "ReceivingCatches", "ReceivingTDS",
	"DefensiveInts", "DefensiveSacks",
}

// recordBookPeriod pairs a StatRecords reference field with its period label.
type recordBookPeriod struct {
	field  string
	period string
}

// League rows expose period record books under Player*-prefixed fields, while
// Conference and Team rows use the shorter field names.
var leaguePeriodFields = []recordBookPeriod{
	{"PlayerCareerStatRecords", "career"},
	{"PlayerSeasonStatRecords", "season"},
	{"PlayerGameStatRecords", "game"},
}

var groupPeriodFields = []recordBookPeriod{
	{"CareerStatRecords", "career"},
	{"SeasonStatRecords", "season"},
	{"GameStatRecords", "game"},
}

// maxRecordBookRank caps how many ranked rows a single category block may yield.
// Conference/team boards resolve to a single holder; league boards store a
// ranked top-N (typically 10). The cap only guards against malformed data.
const maxRecordBookRank = 50

// maxRecordStatValue is the schema's upper bound for a PlayerStatRecord's
// statValue. Career records (e.g. ~19k career passing yards) far exceed the
// per-season sanity cap used elsewhere, so the record book uses this wider
// bound to avoid dropping legitimate career leaders.
const maxRecordStatValue = 24000

// recordStatValue reads a PlayerStatRecord statValue, rejecting the unset
// sentinel (<= 0) and values beyond the schema maximum.
func recordStatValue(record Record) (int, bool) {
	v, ok := intFieldOK(record, "statValue")
	if !ok || v <= 0 || v > maxRecordStatValue {
		return 0, false
	}
	return v, true
}

// buildRecordBookExports assembles the full stat record book across all scopes
// (league, every conference, every team) and periods (career, season, game).
func (f *File) buildRecordBookExports() ([]RecordBookEntry, error) {
	var out []RecordBookEntry

	if league, ok := f.PrimaryTableByName("League"); ok {
		if err := league.ReadRecords(); err == nil {
			for _, row := range league.Records {
				if !recordIsActive(row, league) {
					continue
				}
				out = append(out, f.recordBookForRow(row, "league", "", leaguePeriodFields)...)
			}
		}
	}

	if conf, ok := f.PrimaryTableByName("Conference"); ok {
		if err := conf.ReadRecords(); err == nil {
			for _, row := range conf.Records {
				if !recordIsActive(row, conf) {
					continue
				}
				name := stringField(row, "Name")
				if name == "" {
					continue
				}
				out = append(out, f.recordBookForRow(row, "conference", name, groupPeriodFields)...)
			}
		}
	}

	if team, ok := f.PrimaryTableByName("Team"); ok {
		if err := team.ReadRecords(); err == nil {
			for _, row := range team.Records {
				if !recordIsActive(row, team) {
					continue
				}
				name := bestTeamName(row)
				if !isOfficialTeamName(name) {
					continue
				}
				out = append(out, f.recordBookForRow(row, "team", name, groupPeriodFields)...)
			}
		}
	}

	return out, nil
}

// recordBookForRow expands every period record book referenced by a scope row.
func (f *File) recordBookForRow(row Record, scope, scopeName string, periods []recordBookPeriod) []RecordBookEntry {
	var out []RecordBookEntry
	for _, pf := range periods {
		v, ok := row.Get(pf.field)
		if !ok || v.Reference == nil {
			continue
		}
		if v.Reference.TableID == 0 && v.Reference.RowNumber == 0 {
			continue
		}
		out = append(out, f.recordBookEntriesFromStruct(v.Reference, scope, scopeName, pf.period)...)
	}
	return out
}

// recordBookEntriesFromStruct resolves a PlayerStatRecords struct and expands
// each stat-category block into ranked record entries.
func (f *File) recordBookEntriesFromStruct(ref *RecordReference, scope, scopeName, period string) []RecordBookEntry {
	structRec, ok := f.RecordByReference("PlayerStatRecords", ref)
	if !ok {
		return nil
	}

	var out []RecordBookEntry
	for _, category := range recordBookCategories {
		cv, ok := structRec.Get(category)
		if !ok || cv.Reference == nil {
			continue
		}
		if cv.Reference.TableID == 0 && cv.Reference.RowNumber == 0 {
			continue
		}
		prTable, ok := f.GetTableByID(cv.Reference.TableID)
		if !ok || prTable == nil {
			continue
		}
		if err := prTable.ReadRecords(); err != nil {
			continue
		}
		start := int(cv.Reference.RowNumber)
		if start < 0 || start >= len(prTable.Records) {
			continue
		}
		baseType := stringField(prTable.Records[start], "statType")
		if baseType == "" {
			continue
		}

		// A category block is a contiguous run of rows sharing the same
		// statType. League boards fill 2..N ranks; conference/team boards stop
		// at one row because the next row is a different category. Expansion is
		// anchored on an authoritative struct reference, so it is bounded by the
		// statType and a valid statValue rather than the table's active count
		// (record tables reference valid rows beyond NextRecordToUse). Empty
		// trailing slots decode to the default statType with a zero value and
		// are rejected by recordStatValue, ending the block.
		rank := 0
		for i := start; i < len(prTable.Records) && rank < maxRecordBookRank; i++ {
			rec := prTable.Records[i]
			if stringField(rec, "statType") != baseType {
				break
			}
			value, ok := recordStatValue(rec)
			if !ok {
				break
			}
			rank++
			entry := RecordBookEntry{
				Scope:     scope,
				ScopeName: scopeName,
				Period:    period,
				StatType:  baseType,
				Rank:      rank,
				StatValue: value,
				FirstName: stringField(rec, "firstName"),
				LastName:  stringField(rec, "lastName"),
				Position:  statRecordPosition(rec, baseType),
				TeamName:  f.recordHolderTeam(rec),
			}
			if year, ok := intFieldOK(rec, "calendarYear"); ok && year > 1900 && year < 2100 {
				entry.CalendarYear = &year
			}
			out = append(out, entry)
		}
	}
	return out
}

// recordHolderTeam returns the record holder's team name. It trusts the stored
// teamName first and otherwise resolves only an explicit Team-table reference.
//
// League boards store a full identity (teamName + an explicit TeamRef) for the
// rank-1 holder only; ranks 2..N leave teamName empty and carry an
// uninitialized local TeamRef (TableID 0) that resolves by row index to an
// unrelated team (e.g. Sam Bradford -> "Hawai'i"). Such refs are ignored so a
// missing team stays empty rather than wrong.
func (f *File) recordHolderTeam(record Record) string {
	if name := stringField(record, "teamName"); name != "" {
		return name
	}
	if v, ok := record.Get("TeamRef"); ok && v.Reference != nil && v.Reference.TableID != 0 {
		if team, ok := f.RecordByReference("Team", v.Reference); ok {
			return bestTeamName(team)
		}
	}
	return ""
}
