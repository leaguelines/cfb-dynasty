package dynasty

// RecordReference points at a row in another table.
type RecordReference struct {
	TableID    uint32 `json:"tableId"`
	RowNumber  uint32 `json:"rowNumber"`
	TableName  string `json:"tableName,omitempty"`
}

// Record is one row in a table with named field values.
type Record struct {
	Index  int
	Fields map[string]FieldValue
}

// Get returns a field value by name.
func (r Record) Get(name string) (FieldValue, bool) {
	if r.Fields == nil {
		return FieldValue{}, false
	}
	v, ok := r.Fields[name]
	return v, ok
}

// FieldValue holds a decoded column value.
type FieldValue struct {
	Raw       any
	String    string
	Int       int64
	Float     float64
	Bool      bool
	Reference *RecordReference
}

// IsEmpty reports whether the field has no decoded value.
func (v FieldValue) IsEmpty() bool {
	return v.Raw == nil && v.String == "" && v.Int == 0 && !v.Bool && v.Reference == nil
}
