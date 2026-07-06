package dynasty

// TableHeader holds metadata parsed from a table's binary header.
type TableHeader struct {
	TableID  uint32
	UniqueID uint32
	Name     string
	Offset   int
	Length   int
	Marker   TableMarker

	IsArray         bool
	RecordCount     uint32
	NextRecordToUse uint32
	RecordSize      int
	NumMembers      uint32
	HeaderSize      int
	Table1StartIndex int
	Table2StartIndex int
	Table3StartIndex int
	OffsetStart      int
	Table1Length    uint32
	Table2Length    uint32
	Table3Length    uint32
	TableTotalLength uint32
	HasSecondTable  bool
	HasThirdTable   bool
}

// ActiveRecordCount returns the number of allocated rows in the table.
func (h TableHeader) ActiveRecordCount() uint32 {
	if h.NextRecordToUse > 0 && h.NextRecordToUse <= h.RecordCount {
		return h.NextRecordToUse
	}
	return h.RecordCount
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

// AllocatedRecordCount returns rows allocated in the table header.
func (t *Table) AllocatedRecordCount() uint32 {
	return t.Header.RecordCount
}

// ActiveRecordCount returns rows currently in use per the table header.
func (t *Table) ActiveRecordCount() uint32 {
	return t.Header.ActiveRecordCount()
}
