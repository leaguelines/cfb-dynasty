package dynasty

import (
	"fmt"
	"sort"

	"github.com/leaguelines/cfb-dynasty/internal/binary"
)

const tableTrailingBytes = 8

// discoverTables scans unpacked save bytes and returns parsed tables.
func discoverTables(unpacked []byte) ([]Table, error) {
	starts, markers := discoverTableStarts(unpacked)
	if len(starts) == 0 {
		return nil, nil
	}

	tables := make([]Table, 0, len(starts))
	for i, start := range starts {
		end := len(unpacked) - tableTrailingBytes
		if i+1 < len(starts) {
			end = starts[i+1]
		}
		if start >= end {
			continue
		}

		data := unpacked[start:end]
		header, err := parseTableHeader(data, start, markers[start])
		if err != nil {
			continue
		}

		tables = append(tables, Table{
			Header: header,
			Data:   append([]byte(nil), data...),
			Index:  len(tables),
		})
	}
	return tables, nil
}

func discoverTableStarts(unpacked []byte) ([]int, map[int]TableMarker) {
	seen := make(map[int]TableMarker)
	for _, marker := range AllTableMarkers {
		for _, pos := range binary.FindAll(unpacked, []byte(marker)) {
			start := pos - TableMarkerOffset
			if start < 0 || start+TableMarkerOffset+len(marker) > len(unpacked) {
				continue
			}
			if TableMarker(unpacked[start+TableMarkerOffset:start+TableMarkerOffset+4]) != marker {
				continue
			}
			if _, ok := seen[start]; ok {
				continue
			}
			seen[start] = marker
		}
	}

	starts := make([]int, 0, len(seen))
	for start := range seen {
		starts = append(starts, start)
	}
	sort.Ints(starts)
	return starts, seen
}

func attachTableSchemas(tables []Table, schema *Schema) {
	if schema == nil {
		return
	}
	for i := range tables {
		if tableSchema, ok := schema.Table(tables[i].Header.Name); ok {
			tables[i].Schema = tableSchema
		}
	}
}

// validateTableSchema reports whether a table header matches its schema member count.
func validateTableSchema(header TableHeader, schema *TableSchema) error {
	if schema == nil || schema.NumMembers == 0 {
		return nil
	}
	if int(header.NumMembers) != schema.NumMembers {
		return fmt.Errorf("cfb-dynasty: table %q: header members=%d schema members=%d",
			header.Name, header.NumMembers, schema.NumMembers)
	}
	return nil
}
