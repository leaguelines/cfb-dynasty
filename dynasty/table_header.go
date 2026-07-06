package dynasty

import (
	"encoding/binary"
	"fmt"
)

const (
	tableBinaryHeaderStart = 0x80
	tableOffsetListStart   = 0xe8
)

// parseTableHeader reads FranTk table metadata from a table blob.
// Layout matches Madden M20/M24 (CFB 27 uses the same header geometry).
func parseTableHeader(data []byte, offset int, marker TableMarker) (TableHeader, error) {
	if len(data) < TableMarkerOffset+4 {
		return TableHeader{}, fmt.Errorf("cfb-dynasty: table %d: header too short (%d bytes)", offset, len(data))
	}

	name := readCString(data, 0)
	if name == "" {
		return TableHeader{}, fmt.Errorf("cfb-dynasty: table @0x%x: empty name", offset)
	}

	isArray := len(name) >= 2 && name[len(name)-2] == '[' && name[len(name)-1] == ']'

	headerStart := tableBinaryHeaderStart
	if len(data) < headerStart+40 {
		return TableHeader{}, fmt.Errorf("cfb-dynasty: table %q @0x%x: binary header too short", name, offset)
	}

	tableStoreLength := binary.BigEndian.Uint32(data[headerStart+36:])
	headerOffset := headerStart + 40
	if tableStoreLength > 0 {
		headerOffset += int(tableStoreLength)
	}
	if len(data) < headerOffset+64 {
		return TableHeader{}, fmt.Errorf("cfb-dynasty: table %q @0x%x: variable header too short", name, offset)
	}

	recordCount := binary.BigEndian.Uint32(data[headerOffset+8:])
	table1Length := binary.BigEndian.Uint32(data[headerOffset+16:])
	table2Length := binary.BigEndian.Uint32(data[headerOffset+20:])
	table3Length := binary.BigEndian.Uint32(data[headerOffset+24:])
	tableTotalLength := binary.BigEndian.Uint32(data[headerOffset+40:])
	recordWords := binary.BigEndian.Uint32(data[headerOffset+44:])
	indexEntries := binary.BigEndian.Uint32(data[headerOffset+52:])
	nextRecordToUse := binary.BigEndian.Uint32(data[headerOffset+60:])

	recordSize := int(recordWords) * 4
	offsetStart := tableOffsetListStart + int(tableStoreLength)
	headerSize := offsetStart
	table1StartIndex := 0
	table2StartIndex := 0

	if !isArray {
		headerSize += int(indexEntries) * 4
		table1StartIndex = headerSize
		table2StartIndex = table1StartIndex + int(recordCount)*recordSize
	} else {
		table1StartIndex = headerSize + int(recordCount)*4
		table2StartIndex = table1StartIndex + int(recordCount)*recordSize
	}

	return TableHeader{
		TableID:          binary.BigEndian.Uint32(data[headerStart:]),
		UniqueID:         binary.BigEndian.Uint32(data[headerStart+4:]),
		Name:             name,
		Offset:           offset,
		Length:           len(data),
		Marker:           marker,
		IsArray:          isArray,
		RecordCount:      recordCount,
		NextRecordToUse:  nextRecordToUse,
		RecordSize:       recordSize,
		NumMembers:       indexEntries,
		HeaderSize:       headerSize,
		Table1StartIndex: table1StartIndex,
		Table2StartIndex: table2StartIndex,
		Table3StartIndex: table2StartIndex + int(table2Length),
		OffsetStart:      offsetStart,
		Table1Length:     table1Length,
		Table2Length:     table2Length,
		Table3Length:     table3Length,
		TableTotalLength: tableTotalLength,
		HasSecondTable:   tableTotalLength > table1Length,
		HasThirdTable:    table3Length > 0,
	}, nil
}

func readCString(data []byte, start int) string {
	if start >= len(data) {
		return ""
	}
	end := start
	for end < len(data) && data[end] != 0 {
		end++
	}
	return string(data[start:end])
}
