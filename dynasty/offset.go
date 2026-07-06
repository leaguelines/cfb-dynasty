package dynasty

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/leaguelines/cfb-dynasty/internal/bitview"
)

// OffsetEntry maps a schema field to a bit layout inside a record.
type OffsetEntry struct {
	Index              int
	OriginalIndex      int
	Name               string
	Type               string
	Offset             int
	Length             int
	IndexOffset        int
	IsReference        bool
	ValueInSecondTable bool
	ValueInThirdTable  bool
	IsSigned           bool
	MinValue           int
	MaxValue           int
	MaxLength          int
	Final              bool
	Const              bool
	Enum               *EnumSchema
}

func buildOffsetTable(data []byte, schema *TableSchema, header TableHeader) ([]OffsetEntry, error) {
	if schema == nil {
		return nil, fmt.Errorf("cfb-dynasty: table %q: schema required for offset table", header.Name)
	}

	currentIndex := header.OffsetStart
	if currentIndex == 0 {
		currentIndex = tableOffsetListStart
	}

	entries := make([]OffsetEntry, 0, len(schema.Attributes))
	seenIndexOffsets := make(map[int]struct{})

	for i, attr := range schema.Attributes {
		if currentIndex+4 > len(data) {
			return nil, fmt.Errorf("cfb-dynasty: table %q: offset list overrun", header.Name)
		}
		indexOffset := int(byteArrayToLong(data[currentIndex : currentIndex+4]))
		_, repeated := seenIndexOffsets[indexOffset]

		minValue, _ := strconv.Atoi(attr.MinValue)
		maxValue, _ := strconv.Atoi(attr.MaxValue)
		maxLength, _ := strconv.Atoi(attr.MaxLength)

		fieldType := attr.Type
		if minValue < 0 || maxValue < 0 {
			fieldType = "s_" + fieldType
		}

		entry := OffsetEntry{
			Index:              i,
			OriginalIndex:      attr.Index,
			Name:               attr.Name,
			Type:               fieldType,
			IsReference:        isReferenceType(attr),
			ValueInSecondTable: header.HasSecondTable && attr.Type == "string",
			ValueInThirdTable:  header.HasThirdTable && attr.Type == "binaryblob",
			IsSigned:           minValue < 0 || maxValue < 0,
			MinValue:           minValue,
			MaxValue:           maxValue,
			MaxLength:          maxLength,
			Final:              attr.Final || repeated,
			Const:              attr.Const,
			Enum:               attr.Enum,
			IndexOffset:        indexOffset,
		}
		entries = append(entries, entry)

		if !isSkippedOffset(entry) {
			seenIndexOffsets[indexOffset] = struct{}{}
		}
		currentIndex += 4
	}

	sortOffsetEntriesByIndexOffset(entries)
	assignOffsetLengths(entries, header.RecordSize)
	assignBitOffsets(entries, header.RecordSize)
	filtered := filterActiveOffsets(entries)
	sortOffsetEntriesByOffset(filtered)
	return filtered, nil
}

func isReferenceType(attr FieldSchema) bool {
	if attr.Enum != nil {
		return false
	}
	if attr.Type == "record" {
		return true
	}
	if attr.Type == "" {
		return false
	}
	first := attr.Type[0]
	return (first >= 'A' && first <= 'Z') || strings.Contains(attr.Type, "[]")
}

func isSkippedOffset(entry OffsetEntry) bool {
	if entry.Final || entry.Const {
		return true
	}
	if strings.Contains(entry.Type, "()") {
		return true
	}
	return entry.Type == "ITransaction_Sleep"
}

func assignOffsetLengths(entries []OffsetEntry, recordSize int) {
	for i := range entries {
		cur := &entries[i]
		if isSkippedOffset(*cur) {
			continue
		}
		next := nextActiveOffset(entries, i+1)
		if next != nil {
			cur.Length = next.IndexOffset - cur.IndexOffset
		} else {
			cur.Length = recordSize*8 - cur.IndexOffset
		}
		if cur.Length > 32 {
			cur.Length = 32
		}
	}
}

func nextActiveOffset(entries []OffsetEntry, start int) *OffsetEntry {
	for i := start; i < len(entries); i++ {
		if isSkippedOffset(entries[i]) {
			continue
		}
		return &entries[i]
	}
	return nil
}

func assignBitOffsets(entries []OffsetEntry, recordSize int) {
	active := make([]*OffsetEntry, 0, len(entries))
	for i := range entries {
		if !isSkippedOffset(entries[i]) {
			active = append(active, &entries[i])
		}
	}

	currentOffsetIndex := 0
	for bit := 0; bit < recordSize*8; bit += 32 {
		var chunk []*OffsetEntry
		offsetLength := bit % 32
		for currentOffsetIndex < len(active) && offsetLength < 32 {
			entry := active[currentOffsetIndex]
			offsetLength += entry.Length
			chunk = append(chunk, entry)
			currentOffsetIndex++
		}
		if len(chunk) == 0 {
			continue
		}
		first := chunk[0]
		last := chunk[len(chunk)-1]
		last.Offset = first.IndexOffset
		for i := len(chunk) - 2; i >= 0; i-- {
			next := chunk[i+1]
			chunk[i].Offset = next.Offset + next.Length
		}
	}
}

func filterActiveOffsets(entries []OffsetEntry) []OffsetEntry {
	out := make([]OffsetEntry, 0, len(entries))
	for _, entry := range entries {
		if !isSkippedOffset(entry) {
			out = append(out, entry)
		}
	}
	return out
}

func sortOffsetEntriesByIndexOffset(entries []OffsetEntry) {
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j-1].IndexOffset > entries[j].IndexOffset; j-- {
			entries[j-1], entries[j] = entries[j], entries[j-1]
		}
	}
}

func sortOffsetEntriesByOffset(entries []OffsetEntry) {
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j-1].Offset > entries[j].Offset; j-- {
			entries[j-1], entries[j] = entries[j], entries[j-1]
		}
	}
}

func byteArrayToLong(data []byte) uint32 {
	if len(data) < 4 {
		return 0
	}
	return uint32(data[3]) | uint32(data[2])<<8 | uint32(data[1])<<16 | uint32(data[0])<<24
}

func decodeFieldValue(data, table2 []byte, entry OffsetEntry, offsets []OffsetEntry) FieldValue {
	switch {
	case entry.ValueInSecondTable:
		return decodeStringField(data, table2, entry, offsets)
	case entry.Enum != nil:
		return decodeEnumField(data, entry)
	case entry.IsReference:
		return decodeReferenceField(data, entry)
	default:
		return decodePrimitiveField(data, entry)
	}
}

func decodePrimitiveField(data []byte, entry OffsetEntry) FieldValue {
	switch entry.Type {
	case "s_int":
		v := int64(bitview.GetBits(data, entry.Offset, entry.Length)) + int64(entry.MinValue)
		return FieldValue{Raw: v, Int: v}
	case "int":
		v := int64(bitview.GetBits(data, entry.Offset, entry.Length))
		return FieldValue{Raw: v, Int: v}
	case "bool":
		v := bitview.GetBits(data, entry.Offset+entry.Length-1, 1) == 1
		return FieldValue{Raw: v, Bool: v}
	case "float":
		v := float64(bitview.GetFloat32(data, entry.Offset))
		return FieldValue{Raw: v, Float: v}
	default:
		v := bitview.GetBits(data, entry.Offset, entry.Length)
		return FieldValue{Raw: v, Int: int64(v)}
	}
}

func decodeReferenceField(data []byte, entry OffsetEntry) FieldValue {
	tableID := bitview.GetBits(data, entry.Offset, 15)
	rowNumber := bitview.GetBits(data, entry.Offset+15, 17)
	ref := &RecordReference{
		TableID:   tableID,
		RowNumber: rowNumber,
	}
	return FieldValue{Raw: ref, Reference: ref}
}

func decodeStringField(data, table2 []byte, entry OffsetEntry, offsets []OffsetEntry) FieldValue {
	index := resolveTable2StringIndex(data, entry, offsets)
	if index < 0 || index >= len(table2) {
		return FieldValue{}
	}
	end := index + entry.MaxLength
	if end > len(table2) {
		end = len(table2)
	}
	raw := table2[index:end]
	value := stringBeforeNull(raw)
	value = strings.Trim(value, "\x00 ")
	if hash := strings.IndexByte(value, '#'); hash >= 0 {
		value = strings.TrimSpace(value[:hash])
	}
	return FieldValue{Raw: value, String: value}
}

func stringBeforeNull(data []byte) string {
	if i := bytes.IndexByte(data, 0); i >= 0 {
		return string(data[:i])
	}
	return string(data)
}

func resolveTable2StringIndex(data []byte, entry OffsetEntry, offsets []OffsetEntry) int {
	index := int(bitview.GetUint32AtBit(data, entry.Offset))
	if index != table2StringIndexSentinel {
		return index
	}
	// CFB stores the school-name table2 index in DIV_SLOTNUMBER when LongName is indirect.
	if entry.Name == "LongName" {
		if div := findOffsetEntryByName(offsets, "DIV_SLOTNUMBER"); div != nil {
			return int(bitview.GetUint32AtBit(data, div.Offset))
		}
	}
	return index
}

func findOffsetEntryByName(offsets []OffsetEntry, name string) *OffsetEntry {
	for i := range offsets {
		if offsets[i].Name == name {
			return &offsets[i]
		}
	}
	return nil
}

func decodeEnumField(data []byte, entry OffsetEntry) FieldValue {
	bits := bitview.GetBits(data, entry.Offset, entry.Length)
	if entry.Enum != nil {
		if name := entry.Enum.NameForBits(bits, entry.Length); name != "" {
			return FieldValue{Raw: name, String: name, Int: int64(bits)}
		}
	}
	return FieldValue{Raw: int(bits), Int: int64(bits)}
}
