package dynasty

import (
	"fmt"
	"strconv"
	"strings"
)

func intField(record Record, name string) int {
	v, ok := intFieldOK(record, name)
	if !ok {
		return 0
	}
	return v
}

func intFieldOK(record Record, name string) (int, bool) {
	value, ok := record.Get(name)
	if !ok {
		return 0, false
	}
	return int(value.Int), true
}

func statIntOK(record Record, name string) (int, bool) {
	value, ok := record.Get(name)
	if !ok {
		return 0, false
	}
	v := int(value.Int)
	if v <= 0 {
		return 0, false
	}
	return v, true
}

func stringField(record Record, name string) string {
	value, ok := record.Get(name)
	if !ok {
		return ""
	}
	if value.String != "" {
		return value.String
	}
	if value.Reference != nil {
		return ""
	}
	switch raw := value.Raw.(type) {
	case string:
		return raw
	case nil:
		return ""
	default:
		return fmt.Sprint(raw)
	}
}

func isOfficialTeamName(name string) bool {
	if name == "" {
		return false
	}
	if strings.HasPrefix(name, "#") {
		return false
	}
	if len(name) == 1 {
		return false
	}
	return true
}

func isUsableTeamShortName(short, long string) bool {
	if !isOfficialTeamName(short) {
		return false
	}
	// CFB saves currently reuse a stale ShortName decode for many schools.
	if short == "AIRFOR" && long != "Air Force" {
		return false
	}
	return true
}

func isEmptyReference(record Record, field string) bool {
	value, ok := record.Get(field)
	if !ok || value.Reference == nil {
		return true
	}
	return value.Reference.TableID == 0 && value.Reference.RowNumber == 0
}

func schemaVersionCopy(v SchemaVersion) *SchemaVersion {
	copy := v
	return &copy
}

func referenceID(ref *RecordReference) string {
	if ref == nil {
		return ""
	}
	return strconv.FormatUint(uint64(ref.TableID), 10) + ":" + strconv.FormatUint(uint64(ref.RowNumber), 10)
}

// sensibleSeasonCount filters decoded ints that are outside plausible season totals.
func sensibleSeasonCount(record Record, name string) (int, bool) {
	v, ok := intFieldOK(record, name)
	if !ok || v < 0 || v > 20 {
		return 0, false
	}
	return v, true
}

// sensibleRank filters poll and ranking fields to plausible values.
func sensibleRank(record Record, name string) *int {
	v, ok := intFieldOK(record, name)
	if !ok || v <= 0 || v > 200 {
		return nil
	}
	return &v
}

func sumSeasonCounts(record Record, fields ...string) (int, bool) {
	total := 0
	found := false
	for _, field := range fields {
		if v, ok := sensibleSeasonCount(record, field); ok {
			total += v
			found = true
		}
	}
	return total, found
}

// seasonStatIntOK filters unset EA sentinel values from season stat integers.
func seasonStatIntOK(record Record, name string) (int, bool) {
	value, ok := record.Get(name)
	if !ok {
		return 0, false
	}
	v := int(value.Int)
	if v <= 0 {
		return 0, false
	}
	if v >= 10000 {
		return 0, false
	}
	return v, true
}

// careerStatIntOK filters career totals that are within plausible bounds.
func careerStatIntOK(record Record, name string) (int, bool) {
	v, ok := intFieldOK(record, name)
	if !ok || v <= 0 || v > 500 {
		return 0, false
	}
	return v, true
}
