package dynasty

// TableHeader holds metadata parsed from a table's binary header.
type TableHeader struct {
	TableID  uint32
	UniqueID uint32
	Name     string
	Offset   int
	Length   int
	Marker   TableMarker
}

// Table represents one SPBF/ASTO/SPEX table inside an unpacked save.
type Table struct {
	Header  TableHeader
	Schema  *TableSchema
	Records []Record
	Data    []byte
	Index   int
}

// Name returns the table name from the header or schema.
func (t *Table) Name() string {
	if t.Header.Name != "" {
		return t.Header.Name
	}
	if t.Schema != nil {
		return t.Schema.Name
	}
	return ""
}

// RecordCount returns the number of parsed records.
func (t *Table) RecordCount() int {
	return len(t.Records)
}
